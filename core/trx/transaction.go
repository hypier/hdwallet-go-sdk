package trx

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/golang/protobuf/proto"
)

type Transaction struct {
	tx   *core.Transaction
	txId string
}

func NewTransaction(txExt *api.TransactionExtention) (*Transaction, error) {
	return &Transaction{tx: txExt.Transaction, txId: hex.EncodeToString(txExt.Txid)}, nil
}

func (t *Transaction) TxHash() ([]byte, error) {
	rawData, err := proto.Marshal(t.tx.GetRawData())
	if err != nil {
		return nil, err
	}
	txHash := sha256.Sum256(rawData)
	return txHash[:], nil
}

func (t *Transaction) Sign(privateKeyCDSA *ecdsa.PrivateKey) error {
	txHash, err := t.TxHash()
	if err != nil {
		return err
	}
	signature, err := crypto.Sign(txHash, privateKeyCDSA)
	if err != nil {
		return err
	}
	t.tx.Signature = append(t.tx.Signature, signature)
	return nil
}

func (t *Transaction) Send(client *client.GrpcClient) (string, error) {
	result, err := client.Broadcast(t.tx)
	if err != nil {
		return "", err
	}
	if result.Code != api.Return_SUCCESS {
		return "", fmt.Errorf("send transaction fail: %s, %s", result.Code.String(), string(result.GetMessage()))
	}

	h, _ := t.TxHash()
	return common.BytesToHexString(h), nil
}

func SignAndSendTx(privateKeyCDSA *ecdsa.PrivateKey, client *client.GrpcClient, txExt *api.TransactionExtention) (string, error) {
	tx, err := NewTransaction(txExt)
	if err != nil {
		return "", err
	}
	err = tx.Sign(privateKeyCDSA)
	if err != nil {
		return "", err
	}
	return tx.Send(client)
}
