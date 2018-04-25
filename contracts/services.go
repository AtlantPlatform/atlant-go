package contracts

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"

	"github.com/AtlantPlatform/ethfw/sol"
)

type kycManager struct {
	baseContract
}

func (m *manager) bindKYC(address string, abi []byte) (KYCManager, error) {
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
	return &kycManager{
		baseContract: baseContract{
			contract: boundContract,
			m:        m,
		},
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
