package eth

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
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

	path, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/0")
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

func (a *Account) Address() common.Address {
	return crypto.PubkeyToAddress(a.privateKeyECDSA.PublicKey)
}

func (a *Account) Sign(input []byte) (string, error) {
	msgHashBytes := crypto.Keccak256Hash(input)
	fmt.Printf("Message Hash: %x\n", msgHashBytes)
	//以太坊签名前缀
	//prefix := []byte(ETH_SIGN_PREFIX)

	/*hash := crypto.Keccak256Hash(
		msgHashBytes.Bytes(),
	)*/
	hashBytes := signHash(msgHashBytes.Bytes())
	fmt.Printf("Sign Hash: %x\n", hashBytes)

	//签名
	signature, err := signMessage(hashBytes, a.privateKeyECDSA)
	if err != nil {
		return "", log.WithError(err, "signMessage failed")
	}

	fmt.Printf("Signature: %x\n", signature)
	return hexutil.Encode(signature), nil
}

func signMessage(sigHash []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {

	// 使用该签名方法生成的签名，需要在最后一位加上27
	signature, err := crypto.Sign(sigHash, privateKey)
	if err != nil {
		return nil, log.WithError(err, "crypto.Sign failed")
	}

	if len(signature) != 65 {
		return nil, log.WithError(errors.New("invalid signature length"), "invalid signature length")
	}
	v := signature[64]
	if v == 0 || v == 1 {
		v += 27
		signature[64] = v
	}
	return signature, nil
}

func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

// 验证签名  message 签名约定的字符串
func (a *Account) VerifySign(signatureHex string, input []byte) (bool, error) {
	// 签名 Hex 转 Bytes
	signature, err := hexutil.Decode(signatureHex)
	if err != nil {
		return false, log.WithError(err, "hexutil.Decode failed")
	}
	keccak256 := crypto.Keccak256(input)
	// msg 转 hash
	newHash := signHash(keccak256)
	fmt.Printf("Sign Hash: %x\n", newHash)
	// 由于前端工具差异, 可能需要重新 Hash
	//var newHash []byte
	if signature[64] > 30 {
		signature[64] -= 31
		//newHash = signHash(hash.Bytes())
	} else {
		signature[64] -= 27
		//newHash = hash.Bytes()
	}

	sigPublicKeyECDSA, err := crypto.SigToPub(newHash, signature)
	if err != nil {
		return false, log.WithError(err, "crypto.SigToPub failed")
	}
	sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
	signatureNoRecoverID := signature[:len(signature)-1]

	recoveredAddr := crypto.PubkeyToAddress(*sigPublicKeyECDSA)
	if a.Address().Hex() != recoveredAddr.Hex() {
		return false, log.WithError(errors.New("address not match"), "address not match")
	}

	verified := crypto.VerifySignature(sigPublicKeyBytes, newHash, signatureNoRecoverID)
	if !verified {
		return false, log.WithError(errors.New("verify signature failed"), "verify signature failed")
	}
	return true, nil
}
