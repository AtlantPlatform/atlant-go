package contracts

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"

	"github.com/AtlantPlatform/ethfw/sol"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/AtlantPlatform/ethfw"
)

type kycManager struct {
	baseContract
	cli ethfw.Client
}

func (m *manager) bindKYC(kycAddress string, abi []byte) (KYCManager, error) {
	if len(kycAddress) == 0 {
		return nil, ErrNoAddress
	} else if abi == nil {
		return nil, ErrNoABI
	}
	cli, _, ok := m.getClient()
	if !ok {
		return nil, ErrNodeUnavailable
	}
	boundContract, err := cli.BindContract(&sol.Contract{
		Address: common.HexToAddress(kycAddress),
		ABI:     abi,
	})
	if err != nil {
		return nil, err
	}
	return &kycManager{
		baseContract: baseContract{
			contract: boundContract,
			m:        m,
		},
		cli: cli,
	}, nil
}

func (k *kycManager) AccountStatus(account string) (KYCStatus, error) {
	// cli, _, ok := k.m.getClient()
	// if !ok {
	// 	// k.m.failNode(addr)
	// 	return StatusUnknown, ErrNodeUnavailable
	// }
	var status uint8
	opts := &bind.CallOpts{
		Context: context.TODO(),
	}
	err := k.contract.Call(opts, &status, "getStatus", common.HexToAddress(account))
	if err != nil {
		// k.m.failNode(addr)
		return StatusUnknown, ErrNodeUnavailable
	}
	switch status {
	case 0:
		return StatusUnknown, nil
	case 1:
		return StatusApproved, nil
	case 2:
		return StatusSuspended, nil
	default:
		log.Warningf("received usupported KYC status: %d", status)
		return StatusUnknown, nil
	}
}


func (k *kycManager) ApproveAddr(account string) (*types.Transaction, error) {
	opts, err := k.cli.TransactOpts(context.Background(), common.HexToAddress(k.m.ownAddr), k.m.pass)
	if err != nil {
		return nil, err
	}
	tx, err := k.contract.Transact(opts, "approveAddr", common.HexToAddress(account))
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (k *kycManager) SuspendAddr(account string) (*types.Transaction, error) {
	opts, err := k.cli.TransactOpts(context.Background(), common.HexToAddress(k.m.ownAddr), k.m.pass)
	if err != nil {
		return nil, err
	}
	tx, err := k.contract.Transact(opts, "suspendAddr", common.HexToAddress(account))
	if err != nil {
		return nil, err
	}
	return tx, nil
}