package eth

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/base"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"math/big"
)

type IToken interface {
	base.Token
	Transfer(from Account, to common.Address, value *big.Int) (*types.Transaction, error)
}
type Token struct {
	Coin
	Info  *base.TokenInfo
	chain *Chain
}

func NewToken(chain *Chain) *Token {
	return &Token{chain: chain, Info: &base.TokenInfo{}}
}

func (t *Token) Chain() base.Chain {
	return t.chain
}

func (t *Token) GetDecimal() int16 {
	if t.Info.Decimal == 0 {
		t.Info, _ = t.TokenInfo()
	}

	return t.Info.Decimal
}

func (t *Token) TokenInfo() (*base.TokenInfo, error) {
	token := base.GetToken(t.CoinType(), "")
	if token != nil {
		return token, nil
	}

	t.Info = &base.TokenInfo{
		Name:    "ETH",
		Symbol:  "ETH",
		Decimal: 18,
	}

	base.AddToken(t.CoinType(), "", t.Info)

	return t.Info, nil
}

func (t *Token) BalanceOfAddress(address string) (*base.Balance, error) {
	client, err := t.chain.Client()
	if err != nil {
		return nil, log.WithError(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout())
	defer cancel()

	balance, err := client.RPCClient().BalanceAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		return base.EmptyBalance(), log.WithError(err)
	}

	return &base.Balance{
		Total: utils.NewOptAmount(balance.String(), t.GetDecimal()),
	}, nil
}

func (t *Token) EstimateGasLimit(from, to common.Address, amount *big.Int) (uint64, error) {
	tx := NewTransaction(from, to, amount, t.chain)
	err := tx.BuildTransfer()
	if err != nil {
		return 0, log.WithError(err)
	}

	msg := &ethereum.CallMsg{
		From:     from,
		To:       &to,
		Gas:      uint64(DefaultEthGasList),
		GasPrice: tx.GasPrice,
		Value:    amount,
	}

	gasLimit, err := t.chain.EstimateGasLimit(msg)
	if err != nil {
		return 0, log.WithError(err)
	}

	return gasLimit, nil
}

func (t *Token) Transfer(from Account, to common.Address, value *big.Int) (*types.Transaction, error) {
	// 构建交易
	tx := NewTransaction(from.Address(), to, value, t.chain)
	err := tx.BuildTransfer()
	if err != nil {
		return nil, log.WithError(err)
	}

	//预估gasLimit
	msg := &ethereum.CallMsg{
		From:     from.Address(),
		To:       &to,
		Gas:      uint64(DefaultEthGasList),
		GasPrice: tx.GasPrice,
		Value:    value,
	}
	tx.GasLimit, err = t.chain.EstimateGasLimit(msg)
	if err != nil {
		tx.GasLimit = uint64(DefaultEthGasList)
	}

	// 签名
	signTx, err := tx.SignTx(from.privateKeyECDSA)
	if err != nil {
		return nil, log.WithError(err)
	}

	// 执行交易
	client, err := t.chain.Client()
	if err != nil {
		return nil, log.WithError(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout())
	defer cancel()

	err = client.RPCClient().SendTransaction(ctx, signTx)
	if err != nil {
		return nil, log.WithError(err)
	}

	return signTx, nil
}
