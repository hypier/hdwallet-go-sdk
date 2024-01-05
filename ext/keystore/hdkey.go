package keystore

import (
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/btc"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/eth"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/trx"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
)

type Key struct {
	//别名
	Alias string
	// 账户的扩展ID
	KeyID string
	//更新时间
	Time uint64 `json:"time"`
	//助记词，加密保存
	Mnemonic string
	// 账户
	Accounts Accounts
}
type Accounts map[uint32]*Account

type Account struct {
	Address    string
	PrivateKey string
}

type cipherParamsJSON struct {
	IV string `json:"iv"`
}

type CryptoJSON struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"ciphertext"`
	CipherParams cipherParamsJSON       `json:"cipherparams"`
	KDF          string                 `json:"kdf"`
	KDFParams    map[string]interface{} `json:"kdfparams"`
	MAC          string                 `json:"mac"`
}

type EncryptedAccountJSON struct {
	Address string     `json:"address"`
	Crypto  CryptoJSON `json:"crypto"`
}

type EncryptedKeyJSON struct {
	//别名
	Alias string `json:"alias"`
	// 账户的扩展ID
	KeyID string `json:"key_id"`
	//更新时间
	Time uint64 `json:"time"`
	//助记词，加密保存
	Mnemonic *CryptoJSON `json:"mnemonic_encrypted,omitempty"`
	// 账户
	Accounts map[uint32]EncryptedAccountJSON `json:"accounts_encrypted"`
	// 版本
	Version string `json:"version"`
}

func NewHDKeyWithMnemonic(alias, mnemonic string) (*Key, error) {
	if mnemonic == "" {
		return nil, utils.ErrInvalidMnemonicPhrase

	}

	return &Key{
		Alias:    alias,
		KeyID:    alias,
		Mnemonic: mnemonic,
	}, nil
}

func NewHDKeyWithPrivateKey(codeType uint32, chainId int, alias, privateKey string) (*Key, error) {
	if privateKey == "" {
		return nil, utils.ErrInvalidPrivateKey
	}
	var address string
	var privateKeyHex string
	var coinType uint32
	if codeType == utils.ETH {
		a, err := eth.NewAccountWithPrivateKey(privateKey)
		if err != nil {
			return nil, err
		}
		address = a.Address().String()
		privateKeyHex = a.PrivateKeyHex()
		coinType = a.CoinType()
	} else if codeType == utils.BTC {
		a, err := btc.NewAccountWithPrivateKey(privateKey, chainId)
		if err != nil {
			return nil, err
		}
		address, err = a.Address()
		if err != nil {
			return nil, err
		}
		privateKeyHex = a.PrivateKeyHex()
		coinType = a.CoinType()

	} else if codeType == utils.TRX {
		a, err := trx.NewAccountWithPrivateKey(privateKey)
		if err != nil {
			return nil, err
		}
		address = a.Address()
		privateKeyHex = a.PrivateKeyHex()
		coinType = a.CoinType()
	}

	account := &Account{
		Address:    address,
		PrivateKey: privateKeyHex,
	}

	accounts := make(Accounts)
	accounts[coinType] = account

	return &Key{
		Alias:    alias,
		KeyID:    alias,
		Accounts: accounts,
	}, nil
}

func (k *Key) FileName() string {
	return k.KeyID
}

func (k *Key) VerifyKey() string {
	return k.KeyID
}
