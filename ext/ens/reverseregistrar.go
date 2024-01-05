package ens

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/wealdtech/go-ens/v3/contracts/reverseregistrar"
)

// ReverseRegistrar is the structure for the reverse registrar
type ReverseRegistrar struct {
	Contract     *reverseregistrar.Contract
	ContractAddr common.Address
}

// NewReverseRegistrar obtains the reverse registrar
func NewReverseRegistrar(backend bind.ContractBackend) (*ReverseRegistrar, error) {
	registry, err := NewRegistry(backend)
	if err != nil {
		return nil, err
	}

	// Obtain the registry address from the registrar
	address, err := registry.Owner("addr.reverse")
	if err != nil {
		return nil, err
	}
	if address == UnknownAddress {
		return nil, errors.New("no registrar for that network")
	}
	return NewReverseRegistrarAt(backend, address)
}

// NewReverseRegistrarAt obtains the reverse registrar at a given address
func NewReverseRegistrarAt(backend bind.ContractBackend, address common.Address) (*ReverseRegistrar, error) {
	contract, err := reverseregistrar.NewContract(address, backend)
	if err != nil {
		return nil, err
	}
	return &ReverseRegistrar{
		Contract:     contract,
		ContractAddr: address,
	}, nil
}

// SetName sets the name
func (r *ReverseRegistrar) SetName(opts *bind.TransactOpts, name string) (tx *types.Transaction, err error) {
	return r.Contract.SetName(opts, name)
}

// DefaultResolverAddress obtains the default resolver address
func (r *ReverseRegistrar) DefaultResolverAddress() (common.Address, error) {
	return r.Contract.DefaultResolver(nil)
}
