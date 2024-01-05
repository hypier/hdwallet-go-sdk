package btc

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/tyler-smith/go-bip39"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
)

type Account struct {
	Coin
	privateKey *btcec.PrivateKey
	address    *btcutil.AddressPubKey
	chain      *chaincfg.Params
}

// NewAccount 使用助记词创建账户
func NewAccount(mnemonic string, chainId int) (*Account, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, log.WithError(err, "bip39.NewSeedWithErrorChecking failed")
	}

	pri, pub := btcec.PrivKeyFromBytes(seed)
	chain, err := utils.GetBtcChainParams(chainId)
	if err != nil {
		return nil, log.WithError(err, "ChainID failed")
	}
	address, err := btcutil.NewAddressPubKey(pub.SerializeCompressed(), chain)
	if err != nil {
		return nil, log.WithError(err, "NewAddressPubKey failed")
	}

	return &Account{
		privateKey: pri,
		address:    address,
		chain:      chain,
	}, nil
}

func NewAccountWithPrivateKey(privateKey string, chainId int) (*Account, error) {
	var (
		pri     *btcec.PrivateKey
		pubData []byte
	)
	chain, err := utils.GetBtcChainParams(chainId)
	if err != nil {
		return nil, log.WithError(err, "ChainID failed")
	}
	wif, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		seed, err := hex.DecodeString(privateKey)
		if err != nil {
			return nil, log.WithError(err, "hex.DecodeString failed")
		}
		var pub *btcec.PublicKey
		pri, pub = btcec.PrivKeyFromBytes(seed)
		pubData = pub.SerializeCompressed()
	} else {
		if !wif.IsForNet(chain) {
			return nil, log.WithError(fmt.Errorf("the specified chainnet does not match the wif private key"))
		}
		pri = wif.PrivKey
		pubData = wif.SerializePubKey()
	}

	address, err := btcutil.NewAddressPubKey(pubData, chain)
	if err != nil {
		return nil, err
	}
	return &Account{
		privateKey: pri,
		address:    address,
		chain:      chain,
	}, nil
}
func (a *Account) GetKey(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
	return a.privateKey, true, nil
}

func (a *Account) Address() (string, error) {
	return a.address.EncodeAddress(), nil
}

func (a *Account) ChainParams() *chaincfg.Params {
	return a.chain
}

func (a *Account) GetScript(addr btcutil.Address) ([]byte, error) {
	return nil, errors.New("GetScript not supported")
}
func (a *Account) PrivateKey() []byte {
	return a.privateKey.Serialize()
}
func (a *Account) PrivateKeyHex() string {
	return hex.EncodeToString(a.privateKey.Serialize())
}

func (a *Account) PublicKey() []byte {
	return a.address.ScriptAddress()
}

func (a *Account) PublicKeyHex() string {
	return hex.EncodeToString(a.address.ScriptAddress())
}

// NativeSegwitAddress P2WPKH just for m/84'/
func (a *Account) NativeSegwitAddress() (string, error) {
	address, err := btcutil.NewAddressWitnessPubKeyHash(a.address.AddressPubKeyHash().ScriptAddress(), a.chain)
	if err != nil {
		return "", log.WithError(err, "NewAddressWitnessPubKeyHash failed")
	}
	return address.EncodeAddress(), nil
}

// NestedSegwitAddress P2SH-P2WPKH just for m/49'/
func (a *Account) NestedSegwitAddress() (string, error) {
	witAddr, err := btcutil.NewAddressWitnessPubKeyHash(a.address.AddressPubKeyHash().ScriptAddress(), a.chain)
	if err != nil {
		return "", log.WithError(err, "NewAddressWitnessPubKeyHash failed")
	}
	witnessProgram, err := txscript.PayToAddrScript(witAddr)
	if err != nil {
		return "", log.WithError(err, "PayToAddrScript failed")
	}
	address, err := btcutil.NewAddressScriptHash(witnessProgram, a.chain)
	if err != nil {
		return "", log.WithError(err, "NewAddressScriptHash failed")
	}
	return address.EncodeAddress(), nil
}

// TaprootAddress P2TR just for m/86'/
func (a *Account) TaprootAddress() (string, error) {
	tapKey := txscript.ComputeTaprootKeyNoScript(a.address.PubKey())
	address, err := btcutil.NewAddressTaproot(
		schnorr.SerializePubKey(tapKey), a.chain,
	)
	if err != nil {
		return "", log.WithError(err, "NewAddressTaproot failed")
	}
	return address.EncodeAddress(), nil
}

func (a *Account) ComingTaprootAddress() (string, error) {
	taproot, err := btcutil.NewAddressTaproot(a.address.ScriptAddress()[1:33], a.chain)
	if err != nil {
		return "", log.WithError(err, "NewAddressTaproot failed")
	}
	return taproot.EncodeAddress(), nil
}

// LegacyAddress P2PKH just for m/44'/
func (a *Account) LegacyAddress() string {
	return a.address.AddressPubKeyHash().EncodeAddress()
}
