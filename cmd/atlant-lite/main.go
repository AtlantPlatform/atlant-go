// Copyright 2017-2021 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD-3-Clause "New" or "Revised"
// License (BSD-3-Clause) that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/AtlantPlatform/atlant-go/client"
	cli "github.com/jawher/mow.cli"
)

var app = cli.App("atlant-lite", getBanner()+"\nA lightweight ATLANT node client.")

var nodeAddr = app.StringOpt("A addr", "testnet", "Full node address (ex. localhost:33780)")

func init() {
	// log.SetFlags(log.Lshortfile | log.LstdFlags)
}

var defaultLogLevel = "4"

// logLevel is set in main func
var logLevel *string

func toNatural(s string, defaults uint64) int {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		// defaults in case of incorrect or empty "" value
		return int(defaults)
	} else if i < 0 {
		// not defaults, because nobody expects +100 while specifying -100
		return 0
	}
	return int(i)
}

func main() {

	logLevel = app.String(cli.StringOpt{
		Name:   "l log-level",
		Desc:   "Logging verbosity (0 = minimum, 1...4, 5 = debug).",
		EnvVar: "ANC_LOG_LEVEL",
		Value:  defaultLogLevel,
	})
	log.SetLevel(log.Level(toNatural(*logLevel, 4)))

	app.Command("ping", "Ping node and get its ID", cmdPing)
	app.Command("version", "Get node version", cmdVersion)
	app.Command("put", "Put an object into the store", cmdPutObject)
	app.Command("get", "Get object contents from the store", cmdGetContents)
	app.Command("meta", "Get object meta data from the store", cmdGetMeta)
	app.Command("delete", "Delete object from a store by its ID", cmdDeleteObject)
	app.Command("versions", "List all object versions", cmdListVersions)
	app.Command("ls", "List all objects and sub-directories in a prefix", cmdListObjects)
	if err := app.Run(os.Args); err != nil {
		log.Fatalln("[ERR]", err)
	}
}

func cmdPing(c *cli.Cmd) {
	c.Action = func() {
		cli := getClient()
		ctx := context.Background()
		ver, err := cli.Ping(ctx)
		if err != nil {
			log.Fatalln("[ERR]", err)
		}
		fmt.Println(ver)
	}
}

func cmdVersion(c *cli.Cmd) {
	c.Action = func() {
		cli := getClient()
		ctx := context.Background()
		ver, err := cli.Version(ctx)
		if err != nil {
			log.Fatalln("[ERR]", err)
		}
		fmt.Println(ver)
	}
}

func cmdPutObject(c *cli.Cmd) {
	src := c.StringArg("SRC", "", "Source file path on the disk")
	dst := c.StringArg("DST", "", "Destination object path in the store")
	meta := c.StringOpt("M meta", "", "User meta to keep with object")
	c.Spec = "[-M] SRC DST"
	c.Action = func() {
		f, err := os.Open(*src)
		if err != nil {
			log.Fatalln("[ERR]", err)
		}
		defer f.Close()
		fileInfo, err := f.Stat()
		if err != nil {
			log.Fatalln("[ERR]", err)
		}

		cli := getClient()
		ctx := context.Background()
		meta, err := cli.PutObject(ctx, *dst, &client.PutObjectInput{
			Body:     f,
			Size:     fileInfo.Size(),
			UserMeta: *meta,
		})
		if err != nil {
			log.Fatalln("[ERR]", err)
		}
		fmt.Println(jsonPrint(meta))
	}
}

func cmdGetContents(c *cli.Cmd) {
	path := c.StringArg("PATH", "", "Object path in the store")
	version := c.StringOpt("V version", "", "Specify the exact version of the object")
	c.Spec = "[-V] PATH"
	c.Action = func() {
		cli := getClient()
		ctx := context.Background()
		data, err := cli.GetContents(ctx, *path, *version)
		if err != nil {
			log.Fatalln("[ERR]", err)
		}
		io.Copy(os.Stdout, bytes.NewReader(data))
	}
}

func cmdGetMeta(c *cli.Cmd) {
	path := c.StringArg("PATH", "", "Object path in the store")
	version := c.StringOpt("V version", "", "Specify the exact version of the object")
	c.Spec = "[-V] PATH"
	c.Action = func() {
		cli := getClient()
		ctx := context.Background()
		meta, err := cli.GetMeta(ctx, *path, *version)
		if err != nil {
			log.Fatalln("[ERR]", err)
		}
		fmt.Println(jsonPrint(meta))
	}
}

func cmdDeleteObject(c *cli.Cmd) {
	id := c.StringArg("ID", "", "Object ID in the store")
	c.Spec = "ID"
	c.Action = func() {
		if strings.Contains(*id, "/") {
			log.Fatalln("[ERR]", "ID can't contain slashes. Are you trying to use path instead?")
		}
		cli := getClient()
		ctx := context.Background()
		err := cli.DeleteObject(ctx, *id)
		if err != nil {
			log.Fatalln("[ERR]", err)
		}
		fmt.Println("Deleted")
	}
}

func cmdListVersions(c *cli.Cmd) {
	path := c.StringArg("PATH", "", "Object path in the store")
	c.Spec = "PATH"
	c.Action = func() {
		cli := getClient()
		ctx := context.Background()
		id, versions, err := cli.ListVersions(ctx, *path)
		if err != nil {
			log.Fatalln("[ERR]", err)
		}
		fmt.Println("ID:", id)
		fmt.Println("Versions:", jsonPrint(versions))
	}
}

func cmdListObjects(c *cli.Cmd) {
	prefix := c.StringArg("PREFIX", "/", "Path prefix in the store")
	c.Spec = "[PREFIX]"
	c.Action = func() {
		cli := getClient()
		ctx := context.Background()
		dirs, files, err := cli.ListObjects(ctx, *prefix)
		if err != nil {
			log.Fatalln("[ERR]", err)
		}
		fmt.Println("Dirs:", jsonPrint(dirs))
		fmt.Println("Files:", jsonPrint(files))
	}
}

func getClient() client.Client {
	var urlPrefix string
	switch *nodeAddr {
	case "testnet", "testnet1":
		urlPrefix = "http://node-dev1.atlant.io:33780"
	case "testnet2":
		urlPrefix = "http://node-dev2.atlant.io:33780"
	case "testnet3":
		urlPrefix = "http://node-dev3.atlant.io:33780"
	case "testnet4":
		urlPrefix = "http://node-dev4.atlant.io:33780"
	default:
		urlPrefix = "http://" + *nodeAddr
	}
	return client.New(urlPrefix)
}

func getBanner() string {
	var text string

	text += "                     :-------.                    \n"
	text += "                   `ymyoooooo+.                   \n"
	text += "                  `hmmmhooooooo-                  \n"
	text += "                 .dmmmmmdooooooo:                 \n"
	text += "                 `ymmmmmhsssssss-                 \n"
	text += "               ---:hmmmhsssyooo/---`              \n"
	text += "              /mdoooymysssdmyoooooo+.             \n"
	text += "             ommmdsooooooymmmyoooooo+.            \n"
	text += "            smmmmmmsooooooymmmhooooooo-           \n"
	text += "            +mmmmmdssssssshmmmhsssssss.           \n"
	text += "         .---ommmdsssysosyhhhysssssss:---         \n"
	text += "        -dmsoosmdssshmhooooooossssssooooo/`       \n"
	text += "       :mmmmsoooooosmmmdooooooosdddyoooooo+`      \n"
	text += "      +mmmmmmyooooooodmmdoooooooymmmyoooooo+.     \n"
	text += "      :mmmmmmyssssssymmmdssssssshmmmysssssso`     \n"
	text += "       -dmmmyssssssommmhsssssssymmmyssssss+`      \n"
	text += "        .hdsssssss/ /mhsssssss.`hmyssssss/        \n"
	text += "         `-.......   ........`  `........ 		\n"
	return text
}

func jsonPrint(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "\t")
	return string(data)
}
