// Copyright 2017-2021 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD-3-Clause "New" or "Revised"
// License (BSD-3-Clause) that can be found in the LICENSE file.

package main

import "testing"

const validNodeID = "14V8BbA8ip" // yes, this is not real node.

func TestEntryToString(t *testing.T) {
	entry := Entry{
		validNodeID,
		[]Permission{RecordWritePermission, RecordSyncPermission},
	}
	got := entry.String()
	expected := validNodeID + ":write,sync"

	if got != expected {
		t.Errorf("Entry.String got '%v', expected '%v'", got, expected)
	}
}

func TestEntryFromString(t *testing.T) {
	expected := validNodeID + ":write,sync"
	entry, err := NewEntryFromString(expected)
	if err != nil {
		t.Error(err)
	}
	got := entry.String() // getting string again
	if got != expected {
		t.Errorf("NewEntryFromString: got '%v' expected '%v'", got, expected)
	}
}

func TestWriteEntryFromString(t *testing.T) {
	expected := validNodeID + ":write"
	entry, err := NewEntryFromString(expected)
	if err != nil {
		t.Error(err)
	}
	got := entry.String() // getting string again
	if got != expected {
		t.Errorf("TestWriteEntryFromString: got '%v' expected '%v'", got, expected)
	}
}
