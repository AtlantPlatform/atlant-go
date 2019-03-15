// Copyright 2017, 2018 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package fs

import (
	"context"
	"errors"
	"io"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/xlab/catcher"

	cid "github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/core"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
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
	node    *core.IpfsNode
	subs    []*pubsub.Subscription
	subsMux *sync.RWMutex
}

func newIpfsPubSub(node *core.IpfsNode) *ipfsPubSub {
	return &ipfsPubSub{
		node:    node,
		subsMux: new(sync.RWMutex),
	}
}

func (p *ipfsPubSub) Publish(topic string, data []byte) error {
	if p.node == nil || p.node.PubSub == nil {
		return ErrNoPubSub
	}
	return p.node.PubSub.Publish(topic, data)
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
	if p.node == nil || p.node.PubSub == nil {
		return ErrNoPubSub
	}
	for _, topic := range topics {
		sub, err := p.node.PubSub.Subscribe(topic)
		if err != nil {
			return err
		}
		p.subsMux.Lock()
		p.subs = append(p.subs, sub)
		p.subsMux.Unlock()
		go func(sub *pubsub.Subscription) {
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
