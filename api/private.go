// Copyright 2017, 2018 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package api

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/AtlantPlatform/atlant-go/rs"
)

type PrivateServer struct {
	mux http.Handler
}

func NewPrivateServer() *PrivateServer {
	return &PrivateServer{}
}

// Listen starts a TCP listener, for private server it is advised to use a randomly
// assinged port, e.g. "127.0.0.1:0". Returns the final address or an error if any.
func (p *PrivateServer) Listen(addr string) (string, error) {
	l, err := net.Listen("tcp4", addr)
	if err != nil {
		return "", err
	}
	log.Debugln("PrivateServer listen on", l.Addr().String())
	// start a HTTP server using node's private listener
	go http.Serve(l, p.mux)
	return l.Addr().String(), nil
}

func (p *PrivateServer) RouteAPI(ctx APIContext) {
	r := gin.Default()
	r.GET("/private/v1/ping", p.PingHandler(ctx))
	r.GET("/private/v1/records", p.RecordsHandler(ctx))
	r.POST("/private/v1/announce", p.AnnounceHandler(ctx))
	p.mux = r
}

func (p *PrivateServer) PingHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, ctx.NodeID())
	}
}

func (p *PrivateServer) RecordsHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := ctx.RecordStore().ExportRecords(ctx, c.Writer); err != nil {
			c.AbortWithStatus(500)
		}
		c.Status(200)
	}
}

func (p *PrivateServer) AnnounceHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var event *rs.EventAnnounce
		if err := c.BindJSON(&event); err != nil {
			return
		}
		ctx.RecordStore().ReceiveEventAnnounce(event)
		c.Status(200)
	}
}
