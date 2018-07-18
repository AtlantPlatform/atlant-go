// Copyright 2017, 2018 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package authcenter

import (
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var DefaultMainDomains = []string{
	"node-main.atlant-dev.io",
	"node-main.atlant.io",
}

var DefaultTestDomains = []string{
	"node-test.atlant-dev.io",
	"node-test.atlant.io",
}

func NewDNSAuth(domains []string, dur time.Duration) Auth {
	d := &dnsAuth{
		mux:     new(sync.RWMutex),
		dur:     dur,
		domains: domains,
		entries: make(map[string][]Entry),

		stopC: make(chan struct{}),
	}
	go d.refresh()
	return d
}

type dnsAuth struct {
	mux     *sync.RWMutex
	dur     time.Duration
	domains []string
	entries map[string][]Entry

	stopC chan struct{}
}

func (d *dnsAuth) refresh() {
	sync := func() error {
		seen := make(map[string]struct{})
		promoted := make(map[string]int)
		checkDomain := func(domain string) {
			if _, ok := seen[domain]; ok {
				return
			}
			labels, err := net.LookupTXT(domain)
			if err != nil {
				if strings.Contains(err.Error(), "no such host") {
					return
				}
				log.WithField("domain", domain).Infoln("failed to fetch TXT records:", err)
				return
			}
			seen[domain] = struct{}{}
			for _, label := range labels {
				key, tags, ok := parseLabel(label)
				if !ok {
					log.WithField("domain", domain).Infoln("malformed label on auth domain:", label)
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
						log.WithField("domain", domain).Infoln("unknown permission tag:", tag)
					}
				}
				sort.Sort(Permissions(entry.Permissions))
				d.entries[domain] = append(d.entries[domain], entry)
			}
		}

		d.mux.Lock()
		defer d.mux.Unlock()
		d.entries = make(map[string][]Entry, len(d.entries))
		for _, domain := range d.domains {
			checkDomain(domain)
		}
		for domain, n := range promoted {
			if _, ok := seen[domain]; ok {
				// already seen that domain
				continue
			} else if shouldCare := checkRatio(n, len(seen)); !shouldCare {
				// should not care for promotions without majority
				continue
			}
			d.domains = append(d.domains, domain)
			checkDomain(domain)
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
				log.Warningln("DNS auth sync failed: %v", err)
				t.Reset(time.Minute)
				continue
			}
			t.Reset(d.dur)
		}
	}
}

// checkRatio returns true if majority of total promotes a domain.
func checkRatio(n, total int) bool {
	switch {
	case total >= 0 && total <= 2:
		return n >= 1
	case total == 3:
		return n >= 2
	default:
		return n >= 3
	}
}

func parseLabel(label string) (key string, tags []string, ok bool) {
	parts := strings.Split(label, ":")
	if len(parts) != 2 {
		return "", nil, false
	}
	key = strings.TrimSpace(parts[0])
	tagsRaw := strings.Split(parts[1], ",")
	for _, tag := range tagsRaw {
		tags = append(tags, strings.TrimSpace(tag))
	}
	return key, tags, true
}

func (d *dnsAuth) StopUpdates() {
	close(d.stopC)
}

func (d *dnsAuth) AllPermissions(key string) []Permission {
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

func (d *dnsAuth) HasPermissions(key string, perms ...Permission) bool {
	d.mux.RLock()
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

func (d *dnsAuth) Entries() map[string]Entry {
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
