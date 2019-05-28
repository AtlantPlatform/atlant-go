// Copyright 2017-2019 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const someValidNodeID = "14V8BbA8ip" // yes, this is not real node.

func TestMemoryStorage(t *testing.T) {
	var storage Storage
	storage = NewMemoryStorage()
	if storage.String() != "" {
		t.Error("TestMemoryStorage: Expected empty in the beginning")
	}
	storage.Set(Entry{someValidNodeID, make([]Permission, 0)})
	got := storage.String()
	expected := someValidNodeID + ":\n"
	if got != expected {
		t.Errorf("TestMemoryStorage: got '%v', expected '%v'", got, expected)
	}
	storage.Reset()
	if storage.String() != "" {
		t.Error("TestMemoryStorage: Expected empty in the end")
	}
}

func TestDiskStorage(t *testing.T) {
	f, err := ioutil.TempFile("", "atlant-auth-temp-*.txt")
	f.Close()

	if err != nil {
		t.Error("TestDiskStorage: ", err)
	}
	var storage Storage
	storage = NewDiskStorage(f.Name())

	if storage.String() != "" {
		t.Error("TestDiskStorage: Expected empty in the beginning")
	}
	storage.Set(Entry{someValidNodeID, make([]Permission, 0)})
	got := strings.TrimSpace(storage.String())
	expected := someValidNodeID + ":"
	if got != expected {
		t.Errorf("TestDiskStorage: got '%v', expected '%v'", got, expected)
	}

	// checking the file on disk
	fi, err := os.Stat(f.Name())
	if err != nil {
		t.Error("TestDiskStorage: checking file ", err)
	}
	if fi.Size() == 0 {
		t.Errorf("TestDiskStorage: check file size got '%v', expected greater than zero", fi.Size())
	}

	storage.Reset()
	if storage.String() != "" {
		t.Error("TestDiskStorage: Expected empty in the end")
	}
}
