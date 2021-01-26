// Copyright 2017-21 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD 3-Clause "New" or "Revised"
// License (BSD 3) that can be found in the LICENSE file.

package main

import (
	"bufio"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"

	gin "github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	cli "github.com/jawher/mow.cli"
)

var app = cli.App("atlant-auth", getBanner()+"\nATLANT node authorization center.")

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

type appError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func main() {

	logLevel = app.String(cli.StringOpt{
		Name:   "l log-level",
		Desc:   "Logging verbosity (0 = minimum, 1...4, 5 = debug).",
		EnvVar: "AUTH_LOG_LEVEL",
		Value:  defaultLogLevel,
	})

	app.Before = func() {
		log.SetLevel(log.Level(toNatural(*logLevel, 4)))
		if log.GetLevel() > log.InfoLevel {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		log.Debugf("set app logging to %v", log.GetLevel())
		procs := runtime.GOMAXPROCS(toNatural(*goMaxProcs, 128))
		log.Debugf("set GOMAXPROCS to %d", procs)
	}

	app.Action = func() {

		// authmode should be only from the set that is allowed
		if *authMode != "trust-local" && *authMode != "trust-all" {
			log.Fatalln("Invalid auth mode: trust-local or trust-all expected")
		}
		var auth Auth
		if *authMode == "trust-local" {
			auth = NewAuthTrustLocal()
		} else if *authMode == "trust-all" {
			auth = NewAuthTrustAll()
			// } else if *authMode == "rsa" {
			// auth = NewAuthRSA()
		}

		log.WithFields(log.Fields{
			"mode":    *authMode,
			"address": "http://" + *webListenAddr,
		}).Println("Starting HTTP Server")

		r := gin.Default()
		// if *enableMetrics {
		// 	p := ginprom.New(
		// 		ginprom.Engine(r),
		// 		ginprom.Subsystem("gin"),
		// 		ginprom.Path("/metrics"),
		// 	)
		// 	r.Use(p.Instrument())
		// }
		// public endpoint for retrieving all nodes
		var storage Storage
		if *storagePath != "" {
			storage = NewDiskStorage(*storagePath)
			values, err := storage.GetAll()
			if err != nil {
				log.Warningf("Error reading disk storage %v", err)
			}
			log.WithFields(log.Fields{
				"path":    *storagePath,
				"records": len(values),
			}).Infoln("Using Disk Storage")
			for _, entry := range values {
				log.WithField("value", entry.String()).Infoln("Loaded Entry")
			}
		} else {
			storage = NewMemoryStorage()
			log.Infoln("Using memory storage, no persistance")
		}
		r.HEAD("/", func(c *gin.Context) {
			// endpoint to make sure server is alive
			c.Data(200, "plain/text; charset=utf-8", []byte(""))
		})
		r.GET("/", func(c *gin.Context) {
			records, err := storage.GetAll()
			if err != nil {
				log.WithField("error", err.Error()).Error("GetAll")
				c.Data(500, "text/plain; charset=utf-8", []byte(err.Error()))
			} else {
				log.WithField("records", len(records)).Debugln("GetAll")
				c.Data(200, "text/plain; charset=utf-8", []byte(storage.String()))
			}
		})
		r.POST("/", func(c *gin.Context) {
			body, errBody := c.GetRawData()
			if errBody != nil {
				log.WithField("error", errBody.Error()).Error("c.GetRawData")
				c.Data(400, "text/plain; charset=utf-8", []byte("Invalid Body\n"))
				return
			}
			ip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
			auth.SetRemoteAddress([]string{
				ip,
				c.Request.Header.Get("x-forwarded-for"),
			})
			// in more advanced methods, it would be keys verification
			auth.SetKeys(
				c.Request.Header.Get("x-api-key"),
				c.Request.Header.Get("x-api-secret"),
			)
			// validating permission
			if !auth.IsAllowed(body) {
				c.Data(403, "text/plain; charset=utf-8", []byte("No permission to update\n"))
				return
			}
			// ok. we know this person can writes
			updated := 0
			scanner := bufio.NewScanner(strings.NewReader(string(body)))
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if len(line) > 0 {
					entry, errEntry := NewEntryFromString(line)
					if errEntry != nil {
						log.Warning(errEntry)
					} else {
						storage.Set(entry)
						log.WithField("entry", entry.String()).Info("Saved")
						updated++
					}
				}
			}
			if updated > 0 {
				c.Data(200, "text/plain; charset=utf-8", []byte("Updated: "+strconv.Itoa(updated)+"\n"))
			} else {
				c.Data(400, "text/plain; charset=utf-8", []byte("Nothing Updated\n"))
			}
		})

		r.Run(*webListenAddr)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln("[ERR]", err)
	}
}
