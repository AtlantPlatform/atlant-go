// Copyright 2017-2019 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package main

import (
	"io"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

// Storage interface
type Storage interface {
	Set(e Entry)
	GetAll() (map[string]Entry, error)
	String() string
	Reset()
}

// MemoryStorage - non-persistant storage
type MemoryStorage struct {
	mux     *sync.RWMutex
	entries map[string]Entry
}

// DiskStorage - non-persistant storage
type DiskStorage struct {
	path       string
	memStorage *MemoryStorage
}

// NewMemoryStorage initializes non-persistant storage, running in memory only
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		mux:     new(sync.RWMutex),
		entries: make(map[string]Entry),
	}
}

// Set - set value in the storage. For memory storage this is just updating the key
func (m *MemoryStorage) Set(e Entry) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.entries[e.Key] = e
}

// GetAll - get all memory storage as list
func (m *MemoryStorage) GetAll() (map[string]Entry, error) {
	return m.entries, nil
}

// String - GetAll as text string
func (m *MemoryStorage) String() string {
	var sb strings.Builder
	records, err := m.GetAll()
	if err == nil {
		for _, entry := range records {
			sb.WriteString(entry.String())
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// Reset - resetting
func (m *MemoryStorage) Reset() {
	m.mux.Lock()
	defer m.mux.Unlock()
	for k := range m.entries {
		delete(m.entries, k)
	}
}

// NewDiskStorage initializes dummy persistent file storage
// which is immutable to service restart
func NewDiskStorage(path string) *DiskStorage {
	return &DiskStorage{
		path:       path,
		memStorage: NewMemoryStorage(),
	}
}

func (d *DiskStorage) isError(err error) bool {
	if err != nil {
		log.Error(err.Error())
	}
	return (err != nil)
}

// Open - open storage. Load from disk to memory
func (d *DiskStorage) Open() error {
	var file, err = os.OpenFile(d.path, os.O_RDWR, 0644)
	if d.isError(err) {
		return err
	}
	defer file.Close()

	// read file, line by line
	var text = make([]byte, 2048)
	for {
		_, err = file.Read(text)
		// break if finally arrived at end of file
		if err == io.EOF {
			break
		}
		// break if error occured
		if err != nil && err != io.EOF {
			d.isError(err)
			break
		}
		line := strings.TrimSpace(string(text))
		if len(line) > 0 {
			entry, errEntry := NewEntryFromString(line)
			if !d.isError(errEntry) {
				d.memStorage.Set(entry)
			}
		}
	}
	return nil
}

// Save - save storage to the disk
func (d *DiskStorage) Save() error {

	var file, err = os.OpenFile(d.path, os.O_RDWR, 0644)
	if d.isError(err) {
		return err
	}
	defer file.Close()

	_, errWrite := file.WriteString(d.String())
	return errWrite
}

// Set - set value in the storage. For memory storage this is just updating the key
func (d *DiskStorage) Set(e Entry) {
	d.Open()
	d.memStorage.Set(e)
	d.Save()
}

// GetAll - get all storage as list (from the memory)
func (d *DiskStorage) GetAll() (map[string]Entry, error) {
	return d.memStorage.GetAll()
}

// String - get all storage as string (from the memory)
func (d *DiskStorage) String() string {
	return d.memStorage.String()
}

// Reset - resetting
func (d *DiskStorage) Reset() {
	d.Open()
	d.memStorage.Reset()
	d.Save()
}
