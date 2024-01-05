package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"math/big"
)

//type ITransaction interface {
//	EstimateGasLimit(from, to common.Address, amount *big.Int, input []byte) (uint64, error)
//}

type Transaction struct {
	From  common.Address // 发送地址
	To    common.Address // 接收地址
	Data  []byte         // 附加数据
	Value *big.Int       // 转账金额

	GasPrice  *big.Int // 燃料价格
	GasFeeCap *big.Int // 燃料费用上限
	GasTipCap *big.Int // 燃料费用
	BaseFee   *big.Int // 基础费用
	Nonce     *big.Int // 交易的nonce
	GasLimit  uint64   // 燃料限制

	chain *Chain
	ctx   context.Context
}

//func (tx *Transaction) EstimateGasLimit(from, to common.Address, amount *big.Int, input []byte) (uint64, error) {
//	client, err := tx.chain.Client()
//	if err != nil {
//		return 0, err
//	}
//
//	gasLimit, err := client.RPCClient().EstimateGas(context.Background(), ethereum.CallMsg{
//		From:      from,
//		To:        &to,
//		Gas:       uint64(DefaultContractGasLimit),
//		GasPrice:  tx.GasPrice,
//		GasFeeCap: tx.GasFeeCap,
//		GasTipCap: tx.GasTipCap,
//		Value:     amount,
//		Data:      input,
//	})
//	if err != nil {
//		fmt.Println("connect client get estimateGas error", err)
//		return 0, err
//	}
//
//	tx.GasLimit = gasLimit
//
//	return gasLimit, nil
//}

// NewTransaction 创建交易
func NewTransaction(from common.Address, to common.Address, value *big.Int, chain *Chain) *Transaction {

	return &Transaction{
		From:  from,
		To:    to,
		Value: value,
		chain: chain,
		ctx:   context.Background(),
	}
}

// 确保燃料价格的方法
func (tx *Transaction) ensureGasPrice() error {
	client, err := tx.chain.Client()
	if err != nil {
		return log.WithError(err)
	}

	// Get the latest block header
	head, err := client.RPCClient().HeaderByNumber(tx.ctx, nil)
	if err != nil {
		return log.WithError(err)
	}
	tx.BaseFee = head.BaseFee

	// If the base fee is nil, set the gas price if it is also nil
	if head.BaseFee == nil {
		tx.GasPrice, err = client.RPCClient().SuggestGasPrice(tx.ctx)
		if err != nil {
			return log.WithError(err)
		}

	} else {
		// If the gas tip cap is nil, set it from the backend
		// eip1159 gas price
		if tx.GasTipCap == nil {
			tip, err := client.RPCClient().SuggestGasTipCap(tx.ctx)
			if err != nil {
				return log.WithError(err)
			}
			tx.GasTipCap = tip
		}
		// If the gas fee cap is nil, calculate it based on the base fee and set it
		if tx.GasFeeCap == nil {
			gasFeeCap := new(big.Int).Add(
				tx.GasTipCap,
				new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
			)
			tx.GasFeeCap = gasFeeCap
		}
		// Check if the gas fee cap is less than the gas tip cap
		if tx.GasFeeCap.Cmp(tx.GasTipCap) < 0 {
			return fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", tx.GasFeeCap, tx.GasTipCap)
		}
	}
	return nil
}

// getNonce 获取给定地址的nonce。
func (tx *Transaction) getNonce(address common.Address) (*big.Int, error) {
	client, err := tx.chain.Client()
	if err != nil {
		return nil, log.WithError(err)
	}

	ctx, cancel := context.WithTimeout(tx.ctx, client.Timeout())
	defer cancel()

	nonce, err := client.RPCClient().PendingNonceAt(ctx, address)
	if err != nil {
		return nil, log.WithError(err)
	}

	return new(big.Int).SetUint64(nonce), nil

}

// BuildTransfer 构建转账交易
func (tx *Transaction) BuildTransfer() error {
	// 1. 获取nonce
	if tx.Nonce == nil {
		nonce, err := tx.getNonce(tx.From)
		if err != nil {
			return log.WithError(err)
		}
		tx.Nonce = nonce
	}
	// 2. 获取燃料价格
	err := tx.ensureGasPrice()
	if err != nil {
		return log.WithError(err)
	}

	return nil
}
func (tx *Transaction) ToTransactOpts(privateKeyCDSA *ecdsa.PrivateKey) *bind.TransactOpts {
	return &bind.TransactOpts{
		From:      tx.From,
		Nonce:     tx.Nonce,
		Value:     tx.Value,
		GasPrice:  tx.GasPrice,
		GasFeeCap: tx.GasFeeCap,
		GasTipCap: tx.GasTipCap,
		GasLimit:  tx.GasLimit,
		Context:   tx.ctx,

		Signer: func(address common.Address, t *types.Transaction) (*types.Transaction, error) {
			return tx.SignRawTx(privateKeyCDSA, t)
		},
	}
}

// SignTx 对交易进行签名
func (tx *Transaction) SignTx(privateKeyCDSA *ecdsa.PrivateKey) (*types.Transaction, error) {
	var t *types.Transaction
	if tx.GasFeeCap == nil {
		baseTx := &types.LegacyTx{
			Nonce:    tx.Nonce.Uint64(),
			To:       &tx.To,
			GasPrice: tx.GasPrice,
			Gas:      tx.GasLimit,
			Value:    tx.Value,
			Data:     tx.Data,
		}
		t = types.NewTx(baseTx)
	} else {
		baseTx := &types.DynamicFeeTx{
			Nonce:     tx.Nonce.Uint64(),
			To:        &tx.To,
			GasFeeCap: tx.GasFeeCap,
			GasTipCap: tx.GasTipCap,
			Gas:       tx.GasLimit,
			Value:     tx.Value,
			Data:      tx.Data,
		}
		t = types.NewTx(baseTx)
	}
	signedTx, err := types.SignTx(t, types.LatestSignerForChainID(tx.chain.ChainId()), privateKeyCDSA)
	if err != nil {
		return nil, log.WithError(err)
	}
	return signedTx, nil
}
func (tx *Transaction) SignRawTx(privateKeyCDSA *ecdsa.PrivateKey, t *types.Transaction) (*types.Transaction, error) {
	signer := types.LatestSignerForChainID(tx.chain.ChainId())
	signedTx, err := types.SignTx(t, signer, privateKeyCDSA)
	if err != nil {
		return nil, log.WithError(err)
	}

	return signedTx, nil
}
