package fs

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/AtlantPlatform/go-ipfs/core"
	peer "github.com/AtlantPlatform/go-ipfs/go-libp2p-peer"
	ma "github.com/AtlantPlatform/go-ipfs/go-multiaddr"
	"github.com/AtlantPlatform/go-ipfs/p2p"
)

const streamProtoName = "/p2p/atlant"

var (
	ErrStreamDisabled     = errors.New("libp2pStreamMounting is disabled")
	ErrListenerRegistered = errors.New("listener is already registered")
	ErrListenerClosed     = errors.New("listener is closed")
)

type PlanetaryListener interface {
	Listen(addr string) error
	IsRunning() bool
	Close() error
}

type p2pListener struct {
	mux  *sync.Mutex
	node *core.IpfsNode
	info *p2p.ListenerInfo
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
	if l.info != nil {
		for _, lstInfo := range l.info.Registry.Listeners {
			if lstInfo.Address.Equal(mlAddr) {
				return ErrListenerRegistered
			}
		}
		if l.info.Running {
			// TODO: allow multiple listeners
			return ErrListenerRegistered
		}
	}
	log.Debugln("p2pListener on", mlAddr.String())
	listenInfo, err := l.node.P2P.NewListener(l.node.Context(), streamProtoName, mlAddr)
	if err != nil {
		err = fmt.Errorf("failed to init P2P listener: %v", err)
		return err

	}
	l.info = listenInfo
	return err
}

func (l *p2pListener) Close() error {
	if !l.info.Running {
		return ErrListenerClosed
	}
	return l.info.Close()
}

func (l *p2pListener) IsRunning() bool {
	if l.info == nil {
		return false
	}
	return l.info.Running
}

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

	remoteMap map[string]*p2p.ListenerInfo
	remoteMux *sync.RWMutex
}

func newClient(n *core.IpfsNode) *p2pClient {
	return &p2pClient{
		node: n,
		doWG: new(sync.WaitGroup),
		cli:  &http.Client{},
	}
}

func (c *p2pClient) dial(ctx context.Context, nodeID string) (*p2p.ListenerInfo, error) {
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
	remote, err := c.node.P2P.Dial(ctx, nil, id, streamProtoName, bindAddr)
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
	} else {
		nodeID = req.URL.Host
	}
	c.doWG.Add(1)
	defer c.doWG.Done()

	remote, err := c.dial(req.Context(), nodeID)
	if err != nil {
		err = fmt.Errorf("dial error: %v", err)
		return nil, err
	}
	host, _ := remote.Address.ValueForProtocol(ma.P_IP4)
	port, _ := remote.Address.ValueForProtocol(ma.P_TCP)
	req.URL.Scheme = "http"
	req.URL.Host = fmt.Sprintf("%s:%s", host, port)
	return c.cli.Do(req)
}

func (c *p2pClient) Close() {
	c.doWG.Wait()
	c.remoteMux.Lock()
	defer c.remoteMux.Unlock()
	for id, r := range c.remoteMap {
		if err := r.Close(); err != nil {
			log.Infoln(err)
		}
		delete(c.remoteMap, id)
	}
}
