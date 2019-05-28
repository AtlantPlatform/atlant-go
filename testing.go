// Copyright 2017-2019 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

//+build testing

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/AtlantPlatform/atlant-go/fs"
	"github.com/AtlantPlatform/atlant-go/logging"
	"github.com/AtlantPlatform/atlant-go/rs"
	"github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
	"github.com/xlab/closer"
)

func init() {
	// this init will be called only within builds that have `-tags testing`
	defaultLogLevel = "5"

	testingCommands = []testingCmd{{
		Name: "test-ipfs-list",
		Desc: "Test for IPFS API: List objects method",
		Init: testIpfsList,
	}, {
		Name: "test-update-rs",
		Desc: "Test for rs update",
		Init: testUpdate,
	}}
}

func testUpdate(c *cli.Cmd) {
	meta := logging.WithFn()
	path := c.StringOpt("P path", "test.txt", "Specify object path.")
	fsDir = c.StringOpt("F fsdir", "var/fs1", "")
	stateDir = c.StringOpt("S statedir", "var/state1", "")
	content := c.StringArg("CONTENT", "hello world!", "Content to store.")
	c.Action = func() {
		printInfo(meta)
		runWithPlanetaryContext(func(ctx PlanetaryContext) {
			log.Println("Node ID:", ctx.NodeID())
			log.Println("Session ID:", ctx.SessionID())

			store, err := rs.NewPlanetaryRecordStore(ctx.NodeID(), ctx.FileStore(), ctx.StateStore())
			if err != nil {
				log.Fatalln("store:", err)
				return
			}
			body := func(str *string) io.ReadCloser {
				buf := bytes.NewBuffer(nil)
				buf.WriteString(*str)
				return fakeCloser{buf}
			}(content)
			r, err := store.UpdateRecord(ctx, *path, body)
			if err != nil {
				log.Fatalln("update:", err)
				return
			}
			fmt.Println(r.Object.Version, "vers")
			fmt.Println(r.Object.VersionOffset, "offset")
			fmt.Println(r.Object.ID, "id")
			fmt.Println(r.Object.Path, "path")
			fmt.Println(r.Object.VersionPrevious, "prev ver")
			fmt.Println("meta vers:", r.Object.Meta().Version())
			fmt.Println("meta prev ver:", r.Object.Meta().VersionPrevious())
			fmt.Println(err, "err")

			if err := ctx.FileStore().PinNewest(r.Object, 3); err != nil {
				log.Fatalln("range pin", err)
				return
			}
			if err := store.Close(); err != nil {
				log.Fatalln("closing store:", err)
				return
			}
		})
	}
}

func testIpfsList(c *cli.Cmd) {
	meta := logging.WithFn()
	id := c.StringOpt("I id", "", "Specify object ID.")
	path := c.StringOpt("P path", "hello.txt", "Specify object path.")
	version := c.StringOpt("V version", "", "Specify object version (IPFS CID).")
	prevVersion := c.StringOpt("P prev-version", "", "Specify object previous version (IPFS CID).")
	versionOffset := c.IntOpt("O offset", 0, "Specify version offset.")
	c.Action = func() {
		printInfo(meta)
		runWithPlanetaryContext(func(ctx PlanetaryContext) {
			log.Println("IPFS identity:", ctx.FileStore().NodeID())
			ref := fs.ObjectRef{
				ID:              *id,
				Path:            *path,
				Version:         *version,
				VersionOffset:   *versionOffset,
				VersionPrevious: *prevVersion,
			}
			list, err := ctx.FileStore().ListObjects(context.Background(), ref)
			if err != nil {
				log.Errorln(err)
				return
			}
			v, _ := json.MarshalIndent(list, "", "\t")
			fmt.Printf("List: %v\n", string(v))
		})
	}
}

func printInfo(meta log.Fields) {
	log.WithFields(meta).Println("---TEST START---")
	closer.Bind(func() {
		log.WithFields(meta).Println("---TEST ENDED---")
	})
}
