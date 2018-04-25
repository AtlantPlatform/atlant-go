package fs

import (
	"context"
	"errors"
	"io"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/xlab/catcher"

	"github.com/AtlantPlatform/go-ipfs/core"
	cid "github.com/AtlantPlatform/go-ipfs/go-cid"
	floodsub "github.com/AtlantPlatform/go-ipfs/go-libp2p-floodsub"
)

var (
	ErrNoPubSub = errors.New("IPFS pubsub is not initialized")
	ErrSubStop  = errors.New("stop subscription")
)

type PlanetaryPubSub interface {
	Publish(topic string, data []byte) error
	Subscribe(fn MessagePeekFunc, topics ...string) error
	Close() error
}

type ipfsPubSub struct {
	node  *core.IpfsNode
	flood *floodsub.PubSub

	subs    []*floodsub.Subscription
	subsMux *sync.RWMutex
}

func newIpfsPubSub(node *core.IpfsNode) *ipfsPubSub {
	return &ipfsPubSub{
		node:    node,
		flood:   node.Floodsub,
		subsMux: new(sync.RWMutex),
	}
}

func (p *ipfsPubSub) Publish(topic string, data []byte) error {
	if p.node == nil || p.flood == nil {
		return ErrNoPubSub
	}
	return p.flood.Publish(topic, data)
}

// TODO(max):
// RegisterTopicValidator
// WithValidatorConcurrency
// WithValidatorTimeout
// etc

type MessagePeekFunc func(m *Message) error

type Message struct {
	From     string   `json:"from,omitempty"`
	Data     []byte   `json:"data,omitempty"`
	Seqno    []byte   `json:"seqno,omitempty"`
	TopicIDs []string `json:"topicIDs,omitempty"`
}

func (p *ipfsPubSub) Subscribe(fn MessagePeekFunc, topics ...string) error {
	if p.node == nil || p.flood == nil {
		return ErrNoPubSub
	}
	for _, topic := range topics {
		sub, err := p.flood.Subscribe(topic)
		if err != nil {
			return err
		}
		p.subsMux.Lock()
		p.subs = append(p.subs, sub)
		p.subsMux.Unlock()
		go func(sub *floodsub.Subscription) {
			defer catcher.Catch(catcher.RecvLog(true))
			for {
				msg, err := sub.Next(context.Background())
				if err == io.EOF || err == context.Canceled {
					return
				}
				var fromID string
				if id, err := cid.Parse(msg.From); err == nil {
					fromID = id.String()
				}
				m := &Message{
					From:     fromID,
					Data:     msg.Data,
					Seqno:    msg.Seqno,
					TopicIDs: msg.TopicIDs,
				}
				if err := fn(m); err == ErrSubStop {
					return
				} else if err != nil {
					log.Warningf("MessagePeekFunc error: %v", err)
				}
			}
		}(sub)
	}
	return nil
}

func (p *ipfsPubSub) Close() error {
	if p == nil {
		return nil
	}
	p.subsMux.Lock()
	defer p.subsMux.Unlock()
	for _, sub := range p.subs {
		if sub != nil {
			sub.Cancel()
		}
	}
	p.subs = nil
	return nil
}
