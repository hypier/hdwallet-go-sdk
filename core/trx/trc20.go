package trx

import (
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/base"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"math/big"
)

type Trc20Token struct {
	*Token
	contractAddress string
}

func NewTrc20Token(chain *Chain, contractAddress string) *Trc20Token {
	return &Trc20Token{
		Token:           NewToken(chain),
		contractAddress: contractAddress,
	}
}

func (e *Trc20Token) ContractAddress() string {
	return e.contractAddress
}

func (e *Trc20Token) GetDecimal() int16 {
	if e.Info.Decimal == 0 {
		e.Info, _ = e.TokenInfo()
	}

	return e.Info.Decimal
}

func (e *Trc20Token) TokenInfo() (*base.TokenInfo, error) {
	token := base.GetToken(e.CoinType(), e.contractAddress)
	if token != nil {
		return token, nil
	}

	client, err := e.chain.Client()
	if err != nil {
		return nil, log.WithError(err)
	}

	name, err := client.RPCClient().TRC20GetName(e.contractAddress)
	if err != nil {
		return nil, log.WithError(err)
	}

	symbol, err := client.RPCClient().TRC20GetSymbol(e.contractAddress)
	if err != nil {
		return nil, log.WithError(err)
	}

	decimal, err := client.RPCClient().TRC20GetDecimals(e.contractAddress)
	if err != nil {
		return nil, log.WithError(err)
	}

	e.Info.Name = name
	e.Info.Symbol = symbol
	e.Info.Decimal = int16(decimal.Int64())

	base.AddToken(e.CoinType(), e.contractAddress, e.Info)

	return e.Info, nil
}

// BalanceOfAddress 查询余额
func (e *Trc20Token) BalanceOfAddress(address string) (*base.Balance, error) {
	client, err := e.chain.Client()
	if err != nil {
		return nil, log.WithError(err)
	}

	of, err := client.RPCClient().TRC20ContractBalance(address, e.contractAddress)
	if err != nil {
		return nil, log.WithError(err)
	}

	return &base.Balance{
		Total: utils.NewOptAmount(of.String(), e.GetDecimal()),
	}, nil
}

// Transfer 转账
func (e *Trc20Token) Transfer(from *Account, to string, value *big.Int) (string, error) {
	chain := NewChain()
	client, err := chain.Client()
	if err != nil {
		return "", err
	}
	txExt, err := client.RPCClient().TRC20Send(from.Address(), to, e.ContractAddress(), value, 10000)
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

// Approve 授权
func (e *Trc20Token) Approve(from *Account, to string, value *big.Int) (string, error) {
	chain := NewChain()
	client, err := chain.Client()
	if err != nil {
		return "", err
	}
	txExt, err := client.RPCClient().TRC20Approve(from.Address(), to, e.ContractAddress(), value, 10000)
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
	send, err := transaction.Send(client.RPCClient())
	if err != nil {
		return "", err
	}
	return send, nil
}
