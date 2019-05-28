// Copyright 2017-2019 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"

	cli "github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"

	"github.com/AtlantPlatform/atlant-go/fs"
)

var app = cli.App("atlant-keygen", "Generates various keys for ATLANT Node.")

func main() {
	app.Command("net", "Generate a new private network key", netKeyCmd)
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func netKeyCmd(c *cli.Cmd) {
	c.Action = func() {
		data, err := fs.NewPrivateKey()
		if err != nil {
			log.Fatalln("failed to generate IPFS key:", err)
		}
		parts := strings.Split(string(data), "\n")
		fmt.Println(parts[len(parts)-1])
	}
}
