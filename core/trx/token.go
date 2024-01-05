package trx

import (
	"github.com/shopspring/decimal"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/base"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
)

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
		Name:    "TRX",
		Symbol:  "TRX",
		Decimal: 6,
	}
	base.AddToken(t.CoinType(), "", t.Info)
	return t.Info, nil
}

// BalanceOfAddress 查询余额
func (t *Token) BalanceOfAddress(address string) (*base.Balance, error) {
	rpcClient := t.chain.client.RPCClient()
	ac, err := rpcClient.GetAccount(address)
	if err != nil {
		return nil, err
	}
	balance := decimal.New(ac.Balance, int32(t.GetDecimal()))
	return &base.Balance{
		Total:  utils.NewOptAmount(balance.String(), t.GetDecimal()),
		Usable: utils.NewOptAmount(balance.String(), t.GetDecimal()),
	}, nil
}

// Transfer 转账
func (t *Token) Transfer(from *Account, to string, value int64) (string, error) {
	chain := NewChain()
	client, err := chain.Client()
	if err != nil {
		return "", err
	}
	txExt, err := client.RPCClient().Transfer(from.Address(), to, value)
	if err != nil {
		return "", err
	}
	transaction, err := NewTransaction(txExt)
	if err != nil {
		return "", err
	}
	//签名
	err = transaction.Sign(from.privateKeyECDSA)
	if err != nil {
		return "", err
	}
	//发送
	return transaction.Send(client.RPCClient())
}
