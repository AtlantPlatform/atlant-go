package main

import (
	"context"

	"github.com/AtlantPlatform/atlant-go/fs"
	"github.com/AtlantPlatform/atlant-go/proto"
	"github.com/AtlantPlatform/atlant-go/state"
)

type PlanetaryContext struct {
	context.Context
}

func NewPlanetaryContext(ctx context.Context, env, ver string,
	fileStore fs.PlanetaryFileStore, stateStore state.IndexedStore) PlanetaryContext {
	ctx = context.WithValue(ctx, "env", env)
	ctx = context.WithValue(ctx, "ver", ver)
	ctx = context.WithValue(ctx, "node_id", fileStore.NodeID())
	ctx = context.WithValue(ctx, "session_id", proto.NewID())
	ctx = context.WithValue(ctx, "fs", fileStore)
	ctx = context.WithValue(ctx, "ss", stateStore)
	return PlanetaryContext{ctx}
}

func (c PlanetaryContext) FileStore() fs.PlanetaryFileStore {
	return c.Value("fs").(fs.PlanetaryFileStore)
}

func (c PlanetaryContext) StateStore() state.IndexedStore {
	return c.Value("ss").(state.IndexedStore)
}

func (c PlanetaryContext) Env() string {
	return c.Value("env").(string)
}

func (c PlanetaryContext) Ver() string {
	return c.Value("ver").(string)
}

func (c PlanetaryContext) NodeID() string {
	return c.Value("node_id").(string)
}

func (c PlanetaryContext) SessionID() string {
	return c.Value("session_id").(string)
}
