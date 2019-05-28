// Copyright 2017-2019 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package main

import (
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Auth interface
type Auth interface {
	SetRemoteAddress([]string)
	SetKeys(string, string)
	IsAllowed([]byte) bool
}

// AuthTrustLocal - implementation of Auth interface
type AuthTrustLocal struct {
	Addresses []string
}

// SetRemoteAddress - set list of remote addresses
func (a *AuthTrustLocal) SetRemoteAddress(addr []string) {
	a.Addresses = addr
}

// SetKeys - ignored as irrelevant to this method
func (a *AuthTrustLocal) SetKeys(string, string) {
}

// IsLocal - checks if single proviced IP address is local
// like 127.0.0.1 or 172.x.x.x or 192.168.x.x
func (a *AuthTrustLocal) IsLocal(ip string) bool {

	IP := net.ParseIP(ip)
	if IP == nil {
		log.WithField("ip", ip).Warning("Invalid IP provided, cannot be parsed")
		return false
	}

	var privateIPBlocks []*net.IPNet
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}

	for _, block := range privateIPBlocks {
		if block.Contains(IP) {
			return true
		}
	}
	return false
}

// IsAllowed - checks if all the addresses in the list are in local networks
// Body is not checked. If at least one IP address is not local - it will return false
func (a *AuthTrustLocal) IsAllowed([]byte) bool {
	for _, addr := range a.Addresses {
		if strings.TrimSpace(addr) == "" {
			continue
		}
		if !a.IsLocal(addr) {
			return false
		}
	}
	return true
}

// NewAuthTrustLocal - create Auth to trust only local network
func NewAuthTrustLocal() *AuthTrustLocal {
	return &AuthTrustLocal{}
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - //

// AuthTrustAll - implementation of Auth interface that trusts everything
type AuthTrustAll struct {
}

// SetRemoteAddress - ignored as irrelevant to this method
func (a *AuthTrustAll) SetRemoteAddress([]string) {
}

// SetKeys - ignored as irrelevant to this method
func (a *AuthTrustAll) SetKeys(string, string) {
}

// IsAllowed - always true for this method
func (a *AuthTrustAll) IsAllowed([]byte) bool {
	return true
}

// NewAuthTrustAll - create Auth to trust everybody
func NewAuthTrustAll() *AuthTrustAll {
	return &AuthTrustAll{}
}
