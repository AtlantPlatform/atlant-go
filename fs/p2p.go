// Copyright 2017-2019 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package fs

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/p2p"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
	ma "github.com/multiformats/go-multiaddr"
)

const streamProtoName = "/p2p/atlant"

var (
	// ErrStreamDisabled is thrown if libp2pStreamMounting is disabled
	ErrStreamDisabled = errors.New("libp2pStreamMounting is disabled")
	// ErrListenerRegistered is thrown if listener is already registered
	ErrListenerRegistered = errors.New("listener is already registered")
	// ErrListenerClosed is thrown if listener is closed
	ErrListenerClosed = errors.New("listener is closed")
)

// PlanetaryListener is interface for P2P listening
type PlanetaryListener interface {
	Listen(addr string) error
	IsRunning() bool
	Close() error
}

type p2pListener struct {
	mux  *sync.Mutex
	node *core.IpfsNode
}

func newListener(n *core.IpfsNode) *p2pListener {
	return &p2pListener{
		node: n,
		mux:  new(sync.Mutex),
	}
}

func (l *p2pListener) Listen(addr string) error {
	mlAddr, err := ma.NewMultiaddr(addr)
	if err != nil {
		err = fmt.Errorf("failed to assemble multiaddress: %v", err)
		return err
	}
	l.mux.Lock()
	defer l.mux.Unlock()

	log.WithFields(log.Fields{
		"multiaddress": mlAddr.String(),
		"protocol":     streamProtoName,
	}).Debugln("p2pListener started, p2p ForwardRemote")

	_, errForward := l.node.P2P.ForwardRemote(l.node.Context(), protocol.ID(streamProtoName), mlAddr, true)
	if errForward != nil {
		errForward = fmt.Errorf("failed to init P2P listener: %v", err)
		return errForward
	}
	return nil
}

func (l *p2pListener) Close() error {
	log.Infoln("p2pListener is closing")
	l.node.Close()
	return nil
}

func (l *p2pListener) IsRunning() bool {
	if l.node.P2P == nil {
		return false
	}
	// not very sure that this is correct return
	return len(l.node.P2P.ListenersP2P.Listeners) > 0
}

// PlanetaryClient http client for PlanetaryListener
type PlanetaryClient interface {
	// Do performs a HTTP request over the pipe to PlanetaryListener, e.g.
	// GET http://14V8BYb2dEc3wEwLZroaaTDhoW9TjAMXBnH8BBHj8e5ZEF4hB/private/v1/ping
	Do(req *http.Request) (*http.Response, error)
	Close()
}

type p2pClient struct {
	node *core.IpfsNode
	doWG *sync.WaitGroup
	cli  *http.Client

	remoteMap map[string]*p2p.Listeners
	remoteMux *sync.RWMutex
}

func newClient(n *core.IpfsNode) *p2pClient {
	return &p2pClient{
		node: n,
		doWG: new(sync.WaitGroup),
		cli:  &http.Client{},
	}
}

func (c *p2pClient) dial(ctx context.Context, nodeID string) (p2p.Listener, error) {
	id, err := peer.IDB58Decode(nodeID)
	if err != nil {
		err = fmt.Errorf("failed to parse remote ID: %v", err)
		return nil, err
	}
	bindAddr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/0")
	if err != nil {
		err = fmt.Errorf("failed to assemble multiaddress: %v", err)
		return nil, err
	}
	remote, err := c.node.P2P.ForwardLocal(ctx, id, streamProtoName, bindAddr)
	if err != nil {
		err = fmt.Errorf("failed to dial remote P2P listener: %v", err)
		return nil, err
	}
	return remote, nil
}

func (c *p2pClient) Do(req *http.Request) (*http.Response, error) {
	var nodeID string
	if _, err := peer.IDB58Decode(req.URL.Host); err != nil {
		err = fmt.Errorf("failed to parse nodeID from URL: %v", err)
		return nil, err
	}
	nodeID = req.URL.Host

	c.doWG.Add(1)
	defer c.doWG.Done()

	log.WithFields(log.Fields{
		"url":        req.URL,
		"remoteAddr": req.RemoteAddr,
	}).Debugln("Do.Request")

	remote, err := c.dial(req.Context(), nodeID)
	if err != nil {
		err = fmt.Errorf("dial error: %v", err)
		return nil, err
	}

	host, _ := remote.ListenAddress().ValueForProtocol(ma.P_IP4)
	port, _ := remote.ListenAddress().ValueForProtocol(ma.P_TCP)
	req.URL.Scheme = "http"
	req.URL.Host = fmt.Sprintf("%s:%s", host, port)
	return c.cli.Do(req)
}

func (c *p2pClient) Close() {
	c.doWG.Wait()
	c.remoteMux.Lock()
	defer c.remoteMux.Unlock()
	for id, r := range c.remoteMap {
		r.Close(func(listener p2p.Listener) bool {
			log.WithField("listener", listener).Infoln("p2pClient closed")
			return true
		})
		delete(c.remoteMap, id)
	}
}
