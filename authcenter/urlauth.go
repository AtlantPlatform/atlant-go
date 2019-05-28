package authcenter

import (
	"bufio"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// NewURLAuth initializes authcenter with urls, that will be requested
func NewURLAuth(urls []string, dur time.Duration) Auth {
	d := &urlAuth{
		mux:     new(sync.RWMutex),
		dur:     dur,
		urls:    urls,
		entries: make(map[string][]Entry),

		stopC: make(chan struct{}),
	}
	go d.refresh()
	return d
}

type urlAuth struct {
	mux     *sync.RWMutex
	dur     time.Duration
	urls    []string
	entries map[string][]Entry

	stopC chan struct{}
}

func (d *urlAuth) refresh() {
	sync := func() error {
		seen := make(map[string]struct{})
		promoted := make(map[string]int)
		checkURL := func(url string) {
			if _, ok := seen[url]; ok {
				return
			}
			log.WithField("url", url).Debugln("urlAuth: Looking up")

			httpResponse, err := http.Get(url)
			if err != nil {
				log.WithField("url", url).Infoln("Failed to fetch records from URL:", err)
				return
			}
			defer httpResponse.Body.Close()
			lineReader := io.LimitReader(httpResponse.Body, 2048)

			seen[url] = struct{}{}
			scanner := bufio.NewScanner(lineReader)
			for i := 0; i < 2048 && scanner.Scan(); i++ {
				// reading it line by line
				label := string(scanner.Bytes())
				key, tags, ok := parseLabel(label)
				if !ok {
					log.WithFields(log.Fields{
						"url":   url,
						"label": label,
					}).Infoln("malformed label on auth url")
					continue
				}
				if key == "promote" {
					seenTags := make(map[string]struct{})
					for _, tag := range tags {
						if _, ok := seenTags[tag]; ok {
							continue
						}
						seenTags[tag] = struct{}{}
						promoted[tag]++
					}
					continue
				}
				entry := Entry{
					Key: key,
				}
				for _, tag := range tags {
					switch p := Permission(tag); p {
					case RecordWritePermission, RecordSyncPermission:
						entry.Permissions = append(entry.Permissions, p)
					default:
						log.WithFields(log.Fields{
							"url": url,
							"tag": tag,
						}).Infoln("Unknown permission tag")
					}
				}
				sort.Sort(Permissions(entry.Permissions))
				d.entries[url] = append(d.entries[url], entry)
			}
		}

		d.mux.Lock()
		defer d.mux.Unlock()
		d.entries = make(map[string][]Entry, len(d.entries))
		for _, url := range d.urls {
			checkURL(url)
		}
		for url, n := range promoted {
			if _, ok := seen[url]; ok {
				// already seen that domain
				continue
			} else if shouldCare := checkRatio(n, len(seen)); !shouldCare {
				// should not care for promotions without majority
				continue
			}
			d.urls = append(d.urls, url)
			checkURL(url)
		}
		return nil
	}
	t := time.NewTimer(time.Millisecond)
	for {
		select {
		case <-d.stopC:
			return
		case <-t.C:
			if err := sync(); err != nil {
				log.Warningf("URL auth sync failed: %v", err)
				t.Reset(time.Minute)
				continue
			}
			t.Reset(d.dur)
		}
	}
}
func (d *urlAuth) StopUpdates() {
	close(d.stopC)
}

func (d *urlAuth) AllPermissions(key string) []Permission {
	var perms []Permission
	d.mux.RLock()
	for _, list := range d.entries {
		for _, e := range list {
			if e.Key != key {
				continue
			}
			perms = append(perms, e.AllPermissions()...)
		}
	}
	d.mux.RUnlock()
	return perms
}

func (d *urlAuth) HasPermissions(key string, perms ...Permission) bool {
	d.mux.RLock()
	log.WithFields(log.Fields{
		"key":   key,
		"perms": perms,
	}).Debugln("urlAuth permissions check")
	for _, list := range d.entries {
		for _, e := range list {
			if e.Key != key {
				continue
			}
			if e.HasPermissions(perms...) {
				d.mux.RUnlock()
				return true
			}
		}
	}
	d.mux.RUnlock()
	return false
}

func (d *urlAuth) Entries() map[string]Entry {
	d.mux.RLock()
	m := make(map[string]Entry, len(d.entries))
	for _, list := range d.entries {
		for _, e := range list {
			m[e.Key] = e
		}
	}
	d.mux.RUnlock()
	return m
}
