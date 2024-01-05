package eth

import (
	"github.com/ethereum/go-ethereum"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/base"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Erc20Token struct {
	*Token
	contractAddress common.Address
}

func NewErc20Token(chain *Chain, contractAddress common.Address) *Erc20Token {
	return &Erc20Token{
		Token:           NewToken(chain),
		contractAddress: contractAddress,
	}
}

func (e *Erc20Token) ContractAddress() common.Address {
	return e.contractAddress
}

func (e *Erc20Token) GetDecimal() (int16, error) {
	var err error
	if e.Info.Decimal == 0 {
		e.Info, err = e.TokenInfo()
		if err != nil {
			return 0, err
		}
	}

	return e.Info.Decimal, err
}

func (e *Erc20Token) TokenInfo() (*base.TokenInfo, error) {
	token := base.GetToken(e.CoinType(), e.contractAddress.String())
	if token != nil {
		return token, nil
	}

	client, err := e.chain.Client()
	if err != nil {
		return nil, log.WithError(err)
	}

	erc20 := NewErc20Contract(e.contractAddress, client.RPCClient())

	name, err := erc20.Name()
	if err != nil {
		return nil, log.WithError(err)
	}

	symbol, err := erc20.Symbol()
	if err != nil {
		return nil, log.WithError(err)
	}

	decimal, err := erc20.Decimals()
	if err != nil {
		return nil, log.WithError(err)
	}

	e.Info.Name = name
	e.Info.Symbol = symbol
	e.Info.Decimal = int16(decimal)

	base.AddToken(e.CoinType(), e.contractAddress.String(), e.Info)

	return e.Info, nil
}

func (e *Erc20Token) BalanceOfAddress(address string) (*base.Balance, error) {
	client, err := e.chain.Client()
	if err != nil {
		return nil, log.WithError(err)
	}

	erc20 := NewErc20Contract(e.contractAddress, client.RPCClient())

	of, err := erc20.BalanceOf(common.HexToAddress(address))
	if err != nil {
		return nil, log.WithError(err)
	}
	decimal, err := e.GetDecimal()
	if err != nil {
		return nil, log.WithError(err)
	}
	return &base.Balance{
		Total: utils.NewOptAmount(of.String(), decimal),
	}, nil
}

func (e *Erc20Token) EstimateGasLimit(from, to common.Address, amount *big.Int, method string) (uint64, error) {

	tx := NewTransaction(from, to, nil, e.chain)
	err := tx.BuildTransfer()
	if err != nil {
		return 0, log.WithError(err)
	}

	return e.estimateGasLimit(from, to, tx.GasPrice, amount, method)
}

func (e *Erc20Token) estimateGasLimit(from, to common.Address, gasPrice, amount *big.Int, method string) (uint64, error) {
	if amount == nil {
		return 0, log.WithError(utils.ErrInvalidValue)
	}
	client, err := e.chain.Client()
	if err != nil {
		return 0, log.WithError(err)
	}
	//得到input
	erc20 := NewErc20Contract(e.contractAddress, client.RPCClient())
	input, err := erc20.abi.Pack(method, to, amount)
	if err != nil {
		return 0, log.WithError(err)
	}
	msg := &ethereum.CallMsg{
		From: from,
		To:   &e.contractAddress,
		//Gas:      uint64(DefaultContractGasLimit),
		GasPrice: gasPrice,
		Value:    big.NewInt(0),
		Data:     input,
	}
	gasLimit, err := e.chain.EstimateGasLimit(msg)
	if err != nil {
		return 0, log.WithError(err)
	}
	return gasLimit, nil
}

func (e *Erc20Token) Transfer(from *Account, to common.Address, value *big.Int) (*types.Transaction, error) {
	if value == nil {
		return nil, log.WithError(utils.ErrInvalidValue, "")
	}

	tx := NewTransaction(from.Address(), to, nil, e.chain)
	err := tx.BuildTransfer()
	if err != nil {
		return nil, log.WithError(err)
	}
	tx.GasLimit, err = e.estimateGasLimit(from.Address(), to, tx.GasPrice, value, "transfer")
	if err != nil {
		tx.GasLimit = uint64(DefaultContractGasLimit)
	}
	client, err := e.chain.Client()
	if err != nil {
		return nil, log.WithError(err)
	}
	erc20 := NewErc20Contract(e.contractAddress, client.RPCClient())
	transfer, err := erc20.Transfer(tx.ToTransactOpts(from.privateKeyECDSA), to, value)
	if err != nil {
		return nil, log.WithError(err)
	}

	return transfer, nil
}

// Approve 授权转账功能
// e 		包含了ERC20合约地址 花的钱
// from 	操作者	付钱的人
// spender 	授权目标	花钱的人
// tokens	授权金额	花多少钱
func (e *Erc20Token) Approve(from *Account, spender common.Address, tokens *big.Int) (*types.Transaction, error) {

	tx := NewTransaction(from.Address(), spender, nil, e.chain)
	err := tx.BuildTransfer()
	if err != nil {
		return nil, log.WithError(err)
	}

	client, err := e.chain.Client()
	if err != nil {
		return nil, log.WithError(err)
	}
	erc20 := NewErc20Contract(e.contractAddress, client.RPCClient())

	tx.GasLimit, err = e.estimateGasLimit(from.Address(), spender, tx.GasPrice, tokens, "approve")
	if err != nil {
		return nil, log.WithError(err)
	}
	return erc20.Approve(tx.ToTransactOpts(from.privateKeyECDSA), spender, tokens)
}
