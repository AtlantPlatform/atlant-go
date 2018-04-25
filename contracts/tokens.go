package contracts

import (
	"context"
	"math/big"

	"github.com/AtlantPlatform/ethfw"
	"github.com/AtlantPlatform/ethfw/sol"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
	cli, _, ok := c.m.getClient()
	if !ok {
		return 0, ErrNodeUnavailable
	}
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
	cli, _, ok := m.getClient()
	if !ok {
		return nil, ErrNodeUnavailable
	}
	boundContract, err := cli.BindContract(&sol.Contract{
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
	cli, _, ok := m.getClient()
	if !ok {
		return nil, ErrNodeUnavailable
	}
	boundContract, err := cli.BindContract(&sol.Contract{
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
