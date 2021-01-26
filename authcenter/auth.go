// Copyright 2017-21 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD 3-Clause "New" or "Revised"
// License (BSD 3) that can be found in the LICENSE file.

package authcenter

import (
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
)

// Default is that interface that
var Default Auth

// commented as it is not used
// func init() {
// 	log.WithField("domains", DefaultMainDomains).Debugln("auth.init")
// 	Default = NewDNSAuth(DefaultMainDomains, 1*time.Minute)
// }

// InitWithDomains - initialize with DNS checks of certain domains
func InitWithDomains(domains []string) {
	log.WithField("domains", domains).Debugln("auth.InitWithDomains")
	if Default != nil {
		Default.StopUpdates()
	}
	Default = NewDNSAuth(domains, 1*time.Minute)
}

// InitWithURLs - initialize with URL checks
func InitWithURLs(urls []string) {
	log.WithField("urls", urls).Debugln("auth.InitWithURLs")
	if Default != nil {
		Default.StopUpdates()
	}
	Default = NewURLAuth(urls, 1*time.Minute)
}

// Auth is an interface for checking permissions
type Auth interface {
	Entries() map[string]Entry
	HasPermissions(key string, perms ...Permission) bool
	AllPermissions(key string) []Permission
	StopUpdates()
}

// Permission is a text word
type Permission string

const (
	// RecordWritePermission is a permission to write
	RecordWritePermission Permission = "write"
	// RecordSyncPermission is a permission to sync
	RecordSyncPermission Permission = "sync"
)

// Entry is a node permissions record
type Entry struct {
	Key         string
	Permissions []Permission
}

// AllPermissions returns list of permissions for the node
func (e *Entry) AllPermissions() []Permission {
	return e.Permissions
}

// HasPermissions checks node permission presence
func (e *Entry) HasPermissions(perms ...Permission) bool {
	if e == nil {
		return false
	}
	for _, p := range perms {
		i := Permissions(e.Permissions).Search(p)
		if i < len(e.Permissions) && e.Permissions[i] == p {
			continue
		}
		// a permission is missing
		return false
	}
	return true
}

// Permissions is the collection of Permission
type Permissions []Permission

func (s Permissions) Len() int           { return len(s) }
func (s Permissions) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Permissions) Less(i, j int) bool { return s[i] < s[j] }

// Search allows to search the collection of permission for specific one
func (s Permissions) Search(x Permission) int {
	return sort.Search(len(s), func(i int) bool { return s[i] >= x })
}
