// Copyright 2017-2021 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD-3-Clause "New" or "Revised"
// License (BSD-3-Clause) that can be found in the LICENSE file.

package api

import (
	"context"

	"github.com/AtlantPlatform/atlant-go/contracts"
	"github.com/AtlantPlatform/atlant-go/fs"
	"github.com/AtlantPlatform/atlant-go/rs"
	"github.com/AtlantPlatform/atlant-go/state"
)

// APIContext structure to store API Context
type APIContext struct {
	context.Context
}

// NewContext creates object of API Context
func NewContext(ctx context.Context, r rs.PlanetaryRecordStore, mgr contracts.Manager, ethAddr, logDir string) APIContext {
	ctx = context.WithValue(ctx, "rs", r)
	ctx = context.WithValue(ctx, "eth_addr", ethAddr)
	ctx = context.WithValue(ctx, "contracts", mgr)
	ctx = context.WithValue(ctx, "log_dir", logDir)
	return APIContext{ctx}
}

// NodeID returns Node ID
func (c APIContext) NodeID() string {
	return c.Value("node_id").(string)
}

// SessionID returns Session ID
func (c APIContext) SessionID() string {
	return c.Value("session_id").(string)
}

// Version returns Version
func (c APIContext) Version() string {
	return c.Value("ver").(string)
}

// LogDir returns path for logs
func (c APIContext) LogDir() string {
	return c.Value("log_dir").(string)
}

// RecordStore returns PlanetaryRecordStore object
func (c APIContext) RecordStore() rs.PlanetaryRecordStore {
	return c.Value("rs").(rs.PlanetaryRecordStore)
}

// FileStore returns PlanetaryFileStore object
func (c APIContext) FileStore() fs.PlanetaryFileStore {
	return c.Value("fs").(fs.PlanetaryFileStore)
}

// StateStore returns IndexedStore object
func (c APIContext) StateStore() state.IndexedStore {
	return c.Value("ss").(state.IndexedStore)
}

// ContractsManager returns manager for contracts
func (c APIContext) ContractsManager() contracts.Manager {
	v := c.Value("contracts")
	if v == nil {
		return nil
	}
	return v.(contracts.Manager)
}

// ETHAddr returns address of ethereum wallet of this node
func (c APIContext) ETHAddr() string {
	v := c.Value("eth_addr")
	if v == nil {
		return ""
	}
	return v.(string)
}

// Env returns environment
func (c APIContext) Env() string {
	return c.Value("env").(string)
}
