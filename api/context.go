// Copyright 2017, 2018 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package api

import (
	"context"

	"github.com/AtlantPlatform/atlant-go/contracts"
	"github.com/AtlantPlatform/atlant-go/fs"
	"github.com/AtlantPlatform/atlant-go/rs"
	"github.com/AtlantPlatform/atlant-go/state"
)

type APIContext struct {
	context.Context
}

func NewContext(ctx context.Context, r rs.PlanetaryRecordStore, mgr contracts.Manager, ethAddr, logDir string) APIContext {
	ctx = context.WithValue(ctx, "rs", r)
	ctx = context.WithValue(ctx, "eth_addr", ethAddr)
	ctx = context.WithValue(ctx, "contracts", mgr)
	ctx = context.WithValue(ctx, "log_dir", logDir)
	return APIContext{ctx}
}

func (c APIContext) NodeID() string {
	return c.Value("node_id").(string)
}

func (c APIContext) SessionID() string {
	return c.Value("session_id").(string)
}

func (c APIContext) Version() string {
	return c.Value("ver").(string)
}

func (c APIContext) LogDir() string {
	return c.Value("log_dir").(string)
}

func (c APIContext) RecordStore() rs.PlanetaryRecordStore {
	return c.Value("rs").(rs.PlanetaryRecordStore)
}

func (c APIContext) FileStore() fs.PlanetaryFileStore {
	return c.Value("fs").(fs.PlanetaryFileStore)
}

func (c APIContext) StateStore() state.IndexedStore {
	return c.Value("ss").(state.IndexedStore)
}

func (c APIContext) ContractsManager() contracts.Manager {
	v := c.Value("contracts")
	if v == nil {
		return nil
	}
	return v.(contracts.Manager)
}

func (c APIContext) ETHAddr() string {
	v := c.Value("eth_addr")
	if v == nil {
		return ""
	}
	return v.(string)
}

func (c APIContext) Env() string {
	return c.Value("env").(string)
}
