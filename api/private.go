// Copyright 2017-2021 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD-3-Clause "New" or "Revised"
// License (BSD-3-Clause) that can be found in the LICENSE file.

package api

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/AtlantPlatform/atlant-go/rs"
)

// PrivateServer contains http server descriptor
type PrivateServer struct {
	mux http.Handler
}

// NewPrivateServer returns new instance of http server
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

// RouteAPI sets up GIN routes for private API
func (p *PrivateServer) RouteAPI(ctx APIContext) {
	r := gin.Default()
	r.GET("/private/v1/ping", p.PingHandler(ctx))
	r.GET("/private/v1/records", p.RecordsHandler(ctx))
	r.POST("/private/v1/announce", p.AnnounceHandler(ctx))
	p.mux = r
}

// PingHandler returns HTTP response with NodeID
func (p *PrivateServer) PingHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, ctx.NodeID())
	}
}

// RecordsHandler returns HTTP response with records exported
func (p *PrivateServer) RecordsHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := ctx.RecordStore().ExportRecords(ctx, c.Writer); err != nil {
			c.AbortWithStatus(500)
		}
		c.Status(200)
	}
}

// AnnounceHandler endpoint to reace on event announcements
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
