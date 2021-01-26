// Copyright 2017-2021 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD-3-Clause "New" or "Revised"
// License (BSD-3-Clause) that can be found in the LICENSE file.

package main

import (
	"strings"
)

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

// ToString - represents Entry as one line string
func (e *Entry) String() string {
	str := make([]string, len(e.Permissions))
	for index, perm := range e.Permissions {
		str[index] = string(perm)
	}
	return e.Key + ":" + strings.Join(str, ",")
}

// NewEntryFromString - contructs an Entry from string
func NewEntryFromString(input string) (Entry, error) {
	keySlice := strings.Split(input, ":")
	permSlice := make([]string, 0)
	if len(keySlice) > 1 {
		permSlice = strings.Split(keySlice[1], ",")
	}

	perms := make([]Permission, len(permSlice))
	for index, perm := range permSlice {
		perms[index] = Permission(perm)
	}
	return Entry{keySlice[0], perms}, nil
}
