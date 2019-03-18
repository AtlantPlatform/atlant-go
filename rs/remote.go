// Copyright 2017-2019 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package rs

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	capn "github.com/glycerine/go-capnproto"
	log "github.com/sirupsen/logrus"

	"github.com/AtlantPlatform/atlant-go/proto"
)

type nodeState int

const (
	stateAlive     nodeState = 1
	stateTimeout   nodeState = 2
	stateCancelled nodeState = 3
	stateError     nodeState = 4
)

func (r *recordStore) aliveNodes(ctx context.Context, nodeIDs []string) []string {
	alive := make(map[string]struct{}, 100)
	wg := new(sync.WaitGroup)
	ctx, cancelFn := context.WithTimeout(ctx, 15*time.Second)
	defer cancelFn()
	for _, nodeID := range nodeIDs {
		wg.Add(1)
		go func(nodeID string) {
			defer wg.Done()
			r.outboundWork()
			if state := r.pingNode(ctx, nodeID); state == stateAlive {
				alive[nodeID] = struct{}{}
			}
		}(nodeID)
	}
	wg.Wait()
	aliveList := make([]string, 0, len(alive))
	for id := range alive {
		aliveList = append(aliveList, id)
	}
	return aliveList
}

func (r *recordStore) pingNode(ctx context.Context, nodeID string) nodeState {
	u := fmt.Sprintf("http://%s/private/v1/ping", nodeID)
	req, _ := http.NewRequest("GET", u, nil)
	req = req.WithContext(ctx)
	resp, err := r.fs.Client().Do(req)
	if err != nil {
		// log.Debugln("pingNode:", nodeID, err)
		select {
		case <-ctx.Done():
			if ctx.Err() == context.Canceled {
				return stateCancelled
			}
			return stateTimeout
		default:
			return stateError
		}
	}
	if resp.StatusCode != http.StatusOK {
		return stateError
	}
	return stateAlive
}

func (r *recordStore) getNodeRecords(ctx context.Context, nodeID string, rC chan<- *proto.Record) error {
	u := fmt.Sprintf("http://%s/private/v1/records", nodeID)
	req, _ := http.NewRequest("GET", u, nil)
	req = req.WithContext(ctx)
	resp, err := r.fs.Client().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		if len(body) == 0 {
			return fmt.Errorf("error %d: %s", resp.StatusCode, resp.Status)
		}
		return fmt.Errorf("%s", string(body))
	}

	for {
		seg, err := capn.ReadFromStream(resp.Body, nil)
		if err == io.EOF {
			return nil
		} else if err != nil {
			err = fmt.Errorf("failed to read segment: %v", err)
			return err
		}
		r := proto.ReadRootRecord(seg)
		rC <- &r
	}
	return nil
}

func (r *recordStore) collectRecords(ctx context.Context, peers []string, rC chan<- *proto.Record) {
	defer close(rC)
	log.Debugln("collecting records from:", peers)

	wg := new(sync.WaitGroup)
	for _, nodeID := range peers {
		wg.Add(1)
		go func(nodeID string) {
			defer wg.Done()
			r.outboundWork()
			if err := r.getNodeRecords(ctx, nodeID, rC); err != nil {
				log.WithField("nodeID", nodeID).Warningf("failed to get node records: %v", err)
			}
		}(nodeID)
	}
	wg.Wait()
}
