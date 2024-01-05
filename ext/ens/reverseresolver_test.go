package ens_test

import (
	"fmt"
	"hypier.fun/hdwallet/hdwallet-go-sdk/ext/ens"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestReverseResolve(t *testing.T) {

	tests := []struct {
		name    string
		address common.Address
		res     string
		err     string
	}{
		//{
		//	name:    "NoResolver",
		//	address: common.Address{},
		//	err:     "not a resolver",
		//},
		{
			name:    "NoReverseRecord",
			address: common.HexToAddress("0x03C852a3a8E0E2D048d2F86E3CeBb0a63aD58365"),
			err:     "no resolution",
		},
		//{
		//	name:    "Exists",
		//	address: common.HexToAddress("0x809FA673fe2ab515FaA168259cB14E2BeDeBF68e"),
		//	res:     "avsa.eth",
		//},
	}

	client, err := ethclient.Dial("https://mainnet.infura.io/v3/831a5442dc2e4536a9f8dee4ea1707a6")
	require.NoError(t, err)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := ens.ReverseResolve(client, test.address)
			fmt.Println(res, err)
		})
	}
}
