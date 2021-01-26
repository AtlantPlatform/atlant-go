// Copyright 2017-21 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD 3-Clause "New" or "Revised"
// License (BSD 3) that can be found in the LICENSE file.

package fs

import (
	"strconv"

	config "github.com/ipfs/go-ipfs-config"
	log "github.com/sirupsen/logrus"
)

// PlanetaryCache interface for cache
type PlanetaryCache interface{}

type ipfsOptions struct {
	StoreEnabled   bool
	RelayEnabled   bool
	PubSubEnabled  bool
	NetworkProfile NetworkProfile
	BootstrapPeers []config.BootstrapPeer
	ListenHost     string
	ListenPort     int
	Cache          PlanetaryCache
}

// IpfsOpt handler for options
type IpfsOpt func(o *ipfsOptions)

func defaultIpfsOptions() *ipfsOptions {
	return &ipfsOptions{
		StoreEnabled:   true,
		RelayEnabled:   false,
		PubSubEnabled:  true,
		NetworkProfile: NetworkDefault,
		BootstrapPeers: []config.BootstrapPeer{},
		ListenHost:     "0.0.0.0",
		ListenPort:     33770,
	}
}

// UseStoreOpt handler for StoreEnabled IPFS config option
func UseStoreOpt(v bool) IpfsOpt {
	return func(o *ipfsOptions) {
		o.StoreEnabled = v
	}
}

// UseCacheOpt handler for Cache IPFS config option
func UseCacheOpt(cache PlanetaryCache) IpfsOpt {
	return func(o *ipfsOptions) {
		o.Cache = cache
	}
}

// UseRelayOpt handler for RelayEnabled IPFS config option
func UseRelayOpt(v bool) IpfsOpt {
	return func(o *ipfsOptions) {
		o.RelayEnabled = v
	}
}

// UsePubSubOpt handler for PubSubEnabled IPFS config option
func UsePubSubOpt(v bool) IpfsOpt {
	return func(o *ipfsOptions) {
		o.PubSubEnabled = v
	}
}

// UseNetworkProfileOpt - handler to set NetworkProfile IPFS config option
func UseNetworkProfileOpt(profile NetworkProfile) IpfsOpt {
	return func(o *ipfsOptions) {
		switch profile {
		case NetworkDefault, NetworkServer, NetworkTest, NetworkNoModify:
			o.NetworkProfile = profile
		default:
			log.Warnln("unknown network profile:", profile)
		}
	}
}

// UseBootstrapPeersOpt - handler to set BootstrapPeers IPFS config option
func UseBootstrapPeersOpt(peers []string) IpfsOpt {
	return func(o *ipfsOptions) {
		usePeers := make([]config.BootstrapPeer, 0, len(peers))
		for _, addr := range peers {
			if peer, err := config.ParseBootstrapPeer(addr); err != nil {
				log.Warnf("failed to parse bootstrap addr %s: %v", addr, err)
			} else {
				usePeers = append(usePeers, peer)
			}
		}
		if len(usePeers) > 0 {
			o.BootstrapPeers = usePeers
		} else if len(peers) == 0 {
			o.BootstrapPeers = nil
		} else if len(peers) > 0 && len(usePeers) == 0 {
			log.Warnln("using default bootstrap peers, since all specified failed to parse")
		}
	}
}

// ListenHostOpt handler to set ListenHost IPFS config option
func ListenHostOpt(v string) IpfsOpt {
	return func(o *ipfsOptions) {
		if len(v) > 0 {
			o.ListenHost = v
		}
	}
}

// ListenPortOpt handler to set ListenPort IPFS config option
func ListenPortOpt(v string) IpfsOpt {
	return func(o *ipfsOptions) {
		if len(v) == 0 {
			return
		} else if port, err := strconv.Atoi(v); err != nil {
			log.Warningf("failed to parse port option: %v", err)
		} else if port <= 1024 || port > 65000 {
			log.Warningf("ignoring listening TCP port that is out of range: %v", v)
		} else {
			o.ListenPort = port
		}
	}
}

// NetworkProfile string that stores network settings
type NetworkProfile string

const (
	// NetworkDefault restores default network settings. Agressively discovers private IPs in local network.
	//
	// Activates `default-networking` and `local-discovery` profiles for IPFS.
	NetworkDefault NetworkProfile = "default"

	// NetworkServer is recommended for nodes with public IPv4 address (servers, VPSes, etc.),
	// disables host and content discovery in local networks. Use if the provider warns about DDoS from your node.
	//
	// Activates `default-networking` and `server` profiles for IPFS.
	NetworkServer NetworkProfile = "server"

	// NetworkTest reduces external interference, useful for running ipfs in test environments.
	// Note that with these settings node won't be able to talk to the rest of the
	// network without manual bootstrap.
	//
	// Activates `test` profile for IPFS.
	NetworkTest NetworkProfile = "test"

	// NetworkNoModify skips settings network profile for existing IPFS repos.
	NetworkNoModify NetworkProfile = "no-modify"
)
