package trx

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/tyler-smith/go-bip39"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
)

type Account struct {
	Coin
	privateKeyECDSA *ecdsa.PrivateKey
}

// NewAccount 使用助记词创建账户
func NewAccount(mnemonic string) (*Account, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, log.WithError(err, "bip39.NewSeedWithErrorChecking failed")
	}

	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, log.WithError(err, "hdkeychain.NewMaster failed")
	}

	path, err := accounts.ParseDerivationPath("m/44'/195'/0'/0/0")
	if err != nil {
		return nil, log.WithError(err, "accounts.ParseDerivationPath failed")
	}

	key := masterKey
	for _, n := range path {
		key, err = key.Derive(n)
		if err != nil {
			return nil, log.WithError(err, "key.Derive failed")
		}
	}

	privateKey, err := key.ECPrivKey()
	if err != nil {
		return nil, log.WithError(err, "key.ECPrivKey failed")
	}

	privateKeyECDSA := privateKey.ToECDSA()

	return &Account{
		privateKeyECDSA: privateKeyECDSA,
	}, nil
}

func NewAccountWithPrivateKey(privateKey string) (*Account, error) {
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, log.WithError(err, "crypto.HexToECDSA failed")
	}

	return &Account{
		privateKeyECDSA: privateKeyECDSA,
	}, nil
}

func (a *Account) PrivateKey() []byte {
	return crypto.FromECDSA(a.privateKeyECDSA)
}
func (a *Account) PrivateKeyHex() string {
	return hex.EncodeToString(crypto.FromECDSA(a.privateKeyECDSA))
}

func (a *Account) PublicKey() []byte {
	return crypto.FromECDSAPub(&a.privateKeyECDSA.PublicKey)
}

func (a *Account) PublicKeyHex() string {
	ecdsaPub := crypto.FromECDSAPub(&a.privateKeyECDSA.PublicKey)
	return hex.EncodeToString(ecdsaPub)
}

func (a *Account) Address() string {
	return (address.PubkeyToAddress(a.privateKeyECDSA.PublicKey)).String()
}
