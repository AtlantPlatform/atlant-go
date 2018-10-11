// Copyright 2017, 2018 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package contracts

import (
	"context"
	"math/big"

	"github.com/AtlantPlatform/ethfw"
	"github.com/AtlantPlatform/ethfw/sol"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ethManager struct {
	baseContract
}

func (m *manager) bindETH() TokenManager {
	return &ethManager{
		baseContract: baseContract{
			m: m,
		},
	}
}

func (c *ethManager) AccountBalance(account string) (float64, error) {
	rpc, _, ok := c.m.getClient()
	if !ok {
		return 0, ErrNodeUnavailable
	}
	cli := ethclient.NewClient(rpc)
	bigint, err := cli.BalanceAt(context.TODO(), common.HexToAddress(account), nil)
	if err != nil {
		// m.failNode(addr)
		return 0, err
	}
	wei := ethfw.BigWei(bigint)
	return wei.Ether(), nil
}

type atlManager struct {
	baseContract
	manager Manager
}

func (m *manager) bindATL(address string, abi []byte) (TokenManager, error) {
	if len(address) == 0 {
		return nil, ErrNoAddress
	} else if abi == nil {
		return nil, ErrNoABI
	}
	rpc, _, ok := m.getClient()
	if !ok {
		return nil, ErrNodeUnavailable
	}
	cli := ethclient.NewClient(rpc)
	boundContract, err := ethfw.BindContract(cli, &sol.Contract{
		Address: common.HexToAddress(address),
		ABI:     abi,
	})
	if err != nil {
		return nil, err
	}
	return &atlManager{
		baseContract: baseContract{
			contract: boundContract,
			m:        m,
		},
	}, nil
}

func (c *atlManager) AccountBalance(account string) (float64, error) {
	opts := &bind.CallOpts{
		Context: context.TODO(),
	}
	balance := new(*big.Int)
	err := c.contract.Call(opts, balance, "balanceOf", common.HexToAddress(account))
	if err != nil {
		// c.m.failNode(addr)
		return 0, ErrNodeUnavailable
	}
	wei := ethfw.BigWei(*balance)
	return wei.Tokens(), nil
}

type ptoManager struct {
	baseContract
}

func (m *manager) bindPTO(address string, abi []byte) (TokenManager, error) {
	if len(address) == 0 {
		return nil, ErrNoAddress
	} else if abi == nil {
		return nil, ErrNoABI
	}
	rpc, _, ok := m.getClient()
	if !ok {
		return nil, ErrNodeUnavailable
	}
	cli := ethclient.NewClient(rpc)
	boundContract, err := ethfw.BindContract(cli, &sol.Contract{
		Address: common.HexToAddress(address),
		ABI:     abi,
	})
	if err != nil {
		return nil, err
	}
	return &ptoManager{
		baseContract: baseContract{
			contract: boundContract,
			m:        m,
		},
	}, nil
}

func (c *ptoManager) AccountBalance(account string) (float64, error) {
	opts := &bind.CallOpts{
		Context: context.TODO(),
	}
	balance := new(*big.Int)
	err := c.contract.Call(opts, balance, "balanceOf", common.HexToAddress(account))
	if err != nil {
		// c.m.failNode(addr)
		return 0, ErrNodeUnavailable
	}
	wei := ethfw.BigWei(*balance)
	return wei.Tokens(), nil
}
