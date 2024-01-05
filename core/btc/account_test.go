package btc

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"testing"
)

func TestNewAccountWithPrivateKey(t *testing.T) {
	key, err := NewAccountWithPrivateKey("cTFhdQbsU1xQfziSbM3FYz21a1NX6ukms12w5B3jTq1ZSXDQZqVN", utils.BtcChainTestNet3)
	if err != nil {
		return
	}
	fmt.Println(key)
}

func TestNewAccount(t *testing.T) {
	key, err := NewAccount("dwarf unaware dragon car curve stage include output picture organ skin talk", utils.BtcChainTestNet3)
	if err != nil {
		return
	}
	fmt.Println(key)
}

func TestAccount_ChainParams(t *testing.T) {
	account, err := NewAccountWithPrivateKey("cTFhdQbsU1xQfziSbM3FYz21a1NX6ukms12w5B3jTq1ZSXDQZqVN", utils.BtcChainTestNet3)
	if err != nil {
		return
	}
	fmt.Println(account.ChainParams())
}

func TestAccount_ComingTaprootAddress(t *testing.T) {
	account, err := NewAccountWithPrivateKey("cTFhdQbsU1xQfziSbM3FYz21a1NX6ukms12w5B3jTq1ZSXDQZqVN", utils.BtcChainTestNet3)
	if err != nil {
		return
	}
	fmt.Println(account.ComingTaprootAddress())
}

func TestAccount_Address(t *testing.T) {
	account, err := NewAccountWithPrivateKey("cTFhdQbsU1xQfziSbM3FYz21a1NX6ukms12w5B3jTq1ZSXDQZqVN", utils.BtcChainTestNet3)
	if err != nil {
		return
	}
	fmt.Println(account.Address())
}

func TestAccount_GetKey(t *testing.T) {
	account, err := NewAccountWithPrivateKey("cTFhdQbsU1xQfziSbM3FYz21a1NX6ukms12w5B3jTq1ZSXDQZqVN", utils.BtcChainTestNet3)
	if err != nil {
		return
	}
	param, _ := utils.GetBtcChainParams(utils.BtcChainTestNet3)
	from, _ := btcutil.DecodeAddress("2MuKWyXzED48Rag2WrrLC97BgtCuteUzLDS", param)
	fmt.Println(account.GetKey(from))
}

func TestAccount_LegacyAddress(t *testing.T) {
	account, err := NewAccountWithPrivateKey("cTFhdQbsU1xQfziSbM3FYz21a1NX6ukms12w5B3jTq1ZSXDQZqVN", utils.BtcChainTestNet3)
	if err != nil {
		return
	}
	fmt.Println(account.LegacyAddress())
}

func TestAccount_NativeSegwitAddress(t *testing.T) {
	account, err := NewAccountWithPrivateKey("cTFhdQbsU1xQfziSbM3FYz21a1NX6ukms12w5B3jTq1ZSXDQZqVN", utils.BtcChainTestNet3)
	if err != nil {
		return
	}
	fmt.Println(account.NativeSegwitAddress())
}

func TestAccount_NestedSegwitAddress(t *testing.T) {
	account, err := NewAccountWithPrivateKey("cTFhdQbsU1xQfziSbM3FYz21a1NX6ukms12w5B3jTq1ZSXDQZqVN", utils.BtcChainTestNet3)
	if err != nil {
		return
	}
	fmt.Println(account.NestedSegwitAddress())
}

func TestAccount_PrivateKey(t *testing.T) {
	account, err := NewAccountWithPrivateKey("cTFhdQbsU1xQfziSbM3FYz21a1NX6ukms12w5B3jTq1ZSXDQZqVN", utils.BtcChainTestNet3)
	if err != nil {
		return
	}
	fmt.Println(account.PrivateKey())
}

func TestAccount_PrivateKeyHex(t *testing.T) {
	account, err := NewAccountWithPrivateKey("cTFhdQbsU1xQfziSbM3FYz21a1NX6ukms12w5B3jTq1ZSXDQZqVN", utils.BtcChainTestNet3)
	if err != nil {
		return
	}
	fmt.Println(account.PrivateKeyHex())
}

func TestAccount_PublicKey(t *testing.T) {
	account, err := NewAccountWithPrivateKey("cTFhdQbsU1xQfziSbM3FYz21a1NX6ukms12w5B3jTq1ZSXDQZqVN", utils.BtcChainTestNet3)
	if err != nil {
		return
	}
	fmt.Println(account.PublicKey())
}

func TestAccount_PublicKeyHex(t *testing.T) {
	account, err := NewAccountWithPrivateKey("cTFhdQbsU1xQfziSbM3FYz21a1NX6ukms12w5B3jTq1ZSXDQZqVN", utils.BtcChainTestNet3)
	if err != nil {
		return
	}
	fmt.Println(account.PublicKeyHex())
}

func TestAccount_TaprootAddress(t *testing.T) {
	account, err := NewAccountWithPrivateKey("cTFhdQbsU1xQfziSbM3FYz21a1NX6ukms12w5B3jTq1ZSXDQZqVN", utils.BtcChainTestNet3)
	if err != nil {
		return
	}
	fmt.Println(account.TaprootAddress())
}
