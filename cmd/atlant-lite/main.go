// Copyright 2017-2019 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/AtlantPlatform/atlant-go/client"
	cli "github.com/jawher/mow.cli"
)

var app = cli.App("atlant-lite", getBanner()+"\nA lightweight ATLANT node client.")

var nodeAddr = app.StringOpt("A addr", "testnet", "Full node address (ex. localhost:33780)")

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	app.Command("ping", "Ping node and get its ID", cmdPing)
	app.Command("version", "Get node version", cmdVersion)
	app.Command("put", "Put an object into the store", cmdPutObject)
	app.Command("get", "Get object contents from the store", cmdGetContents)
	app.Command("meta", "Get object meta data from the store", cmdDeleteObject)
	app.Command("delete", "Delete object from a store", cmdPing)
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
		cli := getClient()
		ctx := context.Background()
		meta, err := cli.DeleteObject(ctx, *id)
		if err != nil {
			log.Fatalln("[ERR]", err)
		}
		fmt.Println(jsonPrint(meta))
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
