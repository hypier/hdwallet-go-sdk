package ens

import (
	"encoding/hex"
	"hypier.fun/hdwallet/hdwallet-go-sdk/config"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var client, _ = ethclient.Dial("https://mainnet.infura.io/v3/831a5442dc2e4536a9f8dee4ea1707a6")

var lixClient, _ = ethclient.Dial("https://test.lixb.io")

func TestResolveEmpty(t *testing.T) {
	_, err := Resolve(client, "")
	assert.NotNil(t, err, "Resolved empty name")
}

func TestResolveZero(t *testing.T) {
	_, err := Resolve(client, "0")
	assert.NotNil(t, err, "Resolved empty name")
}

func TestResolveImName(t *testing.T) {
	registry := config.Ens[2888].Registry
	at, err := NewRegistryAt(lixClient, registry)
	if err != nil {
		t.Error(err)
		return
	}
	resolver, err := at.Resolver("shaw.byte.im")
	if err != nil {
		t.Error(err)
		return
	}
	address, err := resolver.Address()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("解析结果", address.String())
	expected := "0x352f0d4b871650c899497B1B8d5741e320ce977e"
	require.Nil(t, err, "Error resolving name")
	assert.Equal(t, expected, address.String(), "Did not receive expected result")
}

func TestReverserSolveImAddress(t *testing.T) {
	registry := config.Ens[2888].Registry
	resolverFor, err := NewReverseResolverForA(lixClient, registry, common.HexToAddress("0x352f0d4b871650c899497B1B8d5741e320ce977e"))
	if err != nil {
		t.Error(err)
		return
	}
	name, err := resolverFor.Name(common.HexToAddress("0x352f0d4b871650c899497B1B8d5741e320ce977e"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("解析结果", name)
	expected := "shaw.byte.im"
	require.Nil(t, err, "Error resolving name")
	assert.Equal(t, expected, name, "Did not receive expected result")
}
func TestResolveNotPresent(t *testing.T) {
	_, err := Resolve(client, "sirnotappearinginthisregistry.eth")
	require.NotNil(t, err, "Resolved name that does not exist")
	assert.Equal(t, "unregistered name", err.Error(), "Unexpected error")
}

// func TestResolveNoResolver(t *testing.T) {
// 	_, err := Resolve(client, "noresolver.eth")
// 	require.NotNil(t, err, "Resolved name without a resolver")
// 	assert.Equal(t, "no resolver", err.Error(), "Unexpected error")
// }

func TestResolveBadResolver(t *testing.T) {
	_, err := Resolve(client, "resolvestozero.eth")
	require.NotNil(t, err, "Resolved name with a bad resolver")
	assert.Equal(t, "no address", err.Error(), "Unexpected error")
}

func TestResolveTestEnsTest(t *testing.T) {
	expected := "ed96dd3be847b387217ef9de5b20d8392a6cdf40"
	actual, err := Resolve(client, "test.enstest.eth")
	require.Nil(t, err, "Error resolving name")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

func TestResolveResolverEth(t *testing.T) {
	expected := "4976fb03c32e5b8cfe2b6ccb31c09ba78ebaba41"
	actual, err := Resolve(client, "resolver.eth")
	require.Nil(t, err, "Error resolving name")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

func TestResolveEthereum(t *testing.T) {
	expected := "de0b295669a9fd93d5f28d9ec85e40f4cb697bae"
	actual, err := Resolve(client, "ethereum.eth")
	require.Nil(t, err, "Error resolving name")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

func TestResolveAddress(t *testing.T) {
	expected := "b8c2c29ee19d8307cb7255e1cd9cbde883a267d5"
	actual, err := Resolve(client, "0xb8c2C29ee19D8307cb7255e1Cd9CbDE883A267d5")
	require.Nil(t, err, "Error resolving address")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

func TestResolveShortAddress(t *testing.T) {
	expected := "0000000000000000000000000000000000000001"
	actual, err := Resolve(client, "0x1")
	require.Nil(t, err, "Error resolving address")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

func TestResolveHexString(t *testing.T) {
	_, err := Resolve(client, "0xe32c6d1a964749b6de2130e20daed821a45b9e7261118801ff5319d0ffc6b54a")
	assert.NotNil(t, err, "Resolved too-long hex string")
}

func TestReverseResolveTestEnsTest(t *testing.T) {
	expected := "nick.eth"
	address := common.HexToAddress("b8c2C29ee19D8307cb7255e1Cd9CbDE883A267d5")
	actual, err := ReverseResolve(client, address)
	require.Nil(t, err, "Error resolving address")
	assert.Equal(t, expected, actual, "Did not receive expected result")
}
