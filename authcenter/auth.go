// Copyright 2017, 2018 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package authcenter

import (
	"sort"
	"time"
)

var Default Auth

func init() {
	Default = NewDNSAuth(DefaultMainDomains, 1*time.Minute)
}

func InitWithDomains(domains []string) {
	if Default != nil {
		Default.StopUpdates()
	}
	Default = NewDNSAuth(domains, 1*time.Minute)
}

type Auth interface {
	Entries() map[string]Entry
	HasPermissions(key string, perms ...Permission) bool
	AllPermissions(key string) []Permission
	StopUpdates()
}

type Permission string

const (
	RecordWritePermission Permission = "write"
)

type Entry struct {
	Key         string
	Permissions []Permission
}

func (e *Entry) AllPermissions() []Permission {
	return e.Permissions
}

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

type Permissions []Permission

func (s Permissions) Len() int           { return len(s) }
func (s Permissions) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Permissions) Less(i, j int) bool { return s[i] < s[j] }
func (s Permissions) Search(x Permission) int {
	return sort.Search(len(s), func(i int) bool { return s[i] >= x })
}
