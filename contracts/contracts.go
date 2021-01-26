// Copyright 2017-2021 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD 3-Clause "New" or "Revised"
// License (BSD 3) that can be found in the LICENSE file.

package contracts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/serialx/hashring"
	log "github.com/sirupsen/logrus"

	"github.com/AtlantPlatform/atlant-go/rs"
	"github.com/AtlantPlatform/ethfw"
)

var (
	ErrNoAddress       = errors.New("no contract address specified")
	ErrNoABI           = errors.New("no contract abi specified")
	ErrNodeUnavailable = errors.New("no geth node available")
)

const (
	TokenETH string = "eth"
	TokenATL string = "atl"
	TokenPTO string = "pto"
)

type ContractConfig struct {
	Address string `json:"address"`
	ABI     []byte `json:"abi"`
}

type Manager interface {
	TokenManager(typ, name string) (TokenManager, error)
	KYCManager() (KYCManager, error)
}

type TokenManager interface {
	AccountBalance(address string) (float64, error)
}

type KYCStatus string

const (
	StatusUnknown   KYCStatus = "unknown"
	StatusApproved  KYCStatus = "approved"
	StatusSuspended KYCStatus = "suspended"
)

type KYCManager interface {
	AccountStatus(account string) (KYCStatus, error)
}

var DefaultTestNodes = []string{
	"http://node-dev1.atlant.io:8545",
	"http://node-dev2.atlant.io:8545",
	"http://node-dev3.atlant.io:8545",
}

var DefaultMainNodes = []string{}

func NewManager(session string, store rs.PlanetaryRecordStore, testnet bool) Manager {
	m := &manager{
		store:   store,
		session: session,
		ringMux: new(sync.RWMutex),
		fails:   make(map[string]int),
	}
	if testnet {
		m.ring = hashring.New(DefaultTestNodes)
	} else {
		m.ring = hashring.New(DefaultMainNodes)
	}
	return m
}

type manager struct {
	session string
	store   rs.PlanetaryRecordStore
	ring    *hashring.HashRing
	ringMux *sync.RWMutex
	fails   map[string]int
}

func (m *manager) getClient() (cli *rpc.Client, addr string, ok bool) {
	for {
		m.ringMux.RLock()
		addr, ok = m.ring.GetNode(m.session)
		m.ringMux.RUnlock()
		if !ok {
			log.Warningln("no available geth nodes in pool, all dead x_X")
			return nil, "", false
		}
		newCli, err := rpc.DialHTTP(addr)
		if err == nil {
			cli = newCli
			break
		}
		log.Warningf("failed to connect to geth node: %v", err)
		m.failNode(addr)
		time.Sleep(3 * time.Second)
	}

	return cli, addr, ok
}

func (m *manager) failNode(addr string) {
	m.ringMux.Lock()
	defer m.ringMux.Unlock()
	if m.fails[addr] < 0 {
		// node been removed
		return
	}
	m.fails[addr]++
	if m.fails[addr] < 3 {
		return
	}
	m.fails[addr] = -1
	m.ring = m.ring.RemoveNode(addr)
	log.Warningf("geth node %s has been removed from pool and will be checked again in 5min", addr)
	go func() {
		// schedule a revival
		time.Sleep(5 * time.Minute)
		m.reviveNode(addr)
	}()
}

func (m *manager) reviveNode(addr string) {
	m.ringMux.Lock()
	defer m.ringMux.Unlock()
	if m.fails[addr] >= 0 {
		// node been restored
		return
	}
	log.Warningf("geth node %s has been added back into pool", addr)
	m.ring = m.ring.AddNode(addr)
	m.fails[addr] = 0
}

func (m *manager) TokenManager(typ, name string) (TokenManager, error) {
	switch typ {
	case TokenETH:
		return m.bindETH(), nil
	case TokenATL:
		ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
		r, err := m.store.ReadRecord(ctx, fmt.Sprintf("/configs/atl/atl.json"))
		cancelFn()
		if err != nil {
			err = fmt.Errorf("failed to read contract config: %v", err)
			return nil, err
		}
		buf, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		var cfg ContractConfig
		if err := json.Unmarshal(buf, &cfg); err != nil {
			err = fmt.Errorf("failed to unmarshal contract config: %v", err)
			return nil, err
		}
		return m.bindATL(cfg.Address, cfg.ABI)
	case TokenPTO:
		ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
		r, err := m.store.ReadRecord(ctx, fmt.Sprintf("/configs/pto/%s.json", name))
		cancelFn()
		if err != nil {
			err = fmt.Errorf("failed to read contract config: %v", err)
			return nil, err
		}
		buf, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		var cfg ContractConfig
		if err := json.Unmarshal(buf, &cfg); err != nil {
			err = fmt.Errorf("failed to unmarshal contract config: %v", err)
			return nil, err
		}
		return m.bindPTO(cfg.Address, cfg.ABI)
	default:
		err := fmt.Errorf("unknown token: %s %s", typ, name)
		return nil, err
	}
}

func (m *manager) KYCManager() (KYCManager, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	r, err := m.store.ReadRecord(ctx, fmt.Sprintf("/configs/kyc/kyc.json"))
	cancelFn()
	if err != nil {
		err = fmt.Errorf("failed to read contract config: %v", err)
		return nil, err
	}
	buf, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	var cfg ContractConfig
	if err := json.Unmarshal(buf, &cfg); err != nil {
		err = fmt.Errorf("failed to unmarshal contract config: %v", err)
		return nil, err
	}
	return m.bindKYC(cfg.Address, cfg.ABI)
}

type baseContract struct {
	m        *manager
	contract *ethfw.BoundContract
}
