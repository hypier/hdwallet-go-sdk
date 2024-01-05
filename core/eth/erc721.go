package eth

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/base"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"math/big"
)

type Erc721Token struct {
	*Token
	contractAddress common.Address
}

func NewErc721Token(chain *Chain, contractAddress common.Address) *Erc721Token {
	return &Erc721Token{
		Token:           NewToken(chain),
		contractAddress: contractAddress,
	}
}

func (e *Erc721Token) ContractAddress() common.Address {
	return e.contractAddress
}

func (e *Erc721Token) TokenInfo() (*base.TokenInfo, error) {
	client, err := e.chain.Client()
	if err != nil {
		return nil, log.WithError(err)
	}

	erc721 := NewErc721Contract(e.contractAddress, client.RPCClient())

	name, err := erc721.Name()
	if err != nil {
		return nil, log.WithError(err)
	}

	symbol, err := erc721.Symbol()
	if err != nil {
		return nil, log.WithError(err)
	}

	e.Info.Name = name
	e.Info.Symbol = symbol

	return e.Info, nil
}

func (e *Erc721Token) BalanceOfAddress(address string) (*base.Balance, error) {
	client, err := e.chain.Client()
	if err != nil {
		return nil, log.WithError(err)
	}

	erc721 := NewErc721Contract(e.contractAddress, client.RPCClient())

	of, err := erc721.BalanceOf(common.HexToAddress(address))
	if err != nil {
		return nil, log.WithError(err)
	}

	return &base.Balance{
		Total: utils.NewOptAmount(of.String(), e.Info.Decimal),
	}, nil
}

func (e *Erc721Token) Transfer(from Account, to common.Address, value *big.Int) (*types.Transaction, error) {
	if value == nil {
		return nil, log.WithError(utils.ErrInvalidValue)
	}
	//gas消耗需要 + value  这里如果把tokenId当成value放入交易中 会出现gas费不足的情况
	tx := NewTransaction(from.Address(), to, nil, e.chain)
	err := tx.BuildTransfer()
	if err != nil {
		return nil, log.WithError(err)
	}

	//预估gas费
	tx.GasLimit, err = e.estimateGasLimit(from.Address(), to, tx.GasPrice, value, "transferFrom")
	if err != nil {
		tx.GasLimit = uint64(DefaultContractGasLimit)
	}

	client, err := e.chain.Client()
	if err != nil {
		return nil, log.WithError(err)
	}
	transfer, err := NewErc721Contract(e.contractAddress, client.RPCClient()).
		TransferFrom(tx.ToTransactOpts(from.privateKeyECDSA), from.Address(), to, value)
	if err != nil {
		return nil, log.WithError(err)
	}
	return transfer, nil
}

func (e *Erc721Token) EstimateGasLimit(from, to common.Address, amount *big.Int, method string) (uint64, error) {
	tx := NewTransaction(from, to, nil, e.chain)
	err := tx.BuildTransfer()
	if err != nil {
		return 0, log.WithError(err)
	}

	return e.estimateGasLimit(from, to, tx.GasPrice, amount, method)
}

func (e *Erc721Token) estimateGasLimit(from, to common.Address, gasPrice, amount *big.Int, method string) (uint64, error) {
	if amount == nil {
		return 0, log.WithError(utils.ErrInvalidValue)
	}

	client, err := e.chain.Client()
	if err != nil {
		return uint64(DefaultContractGasLimit), log.WithError(err)
	}

	//得到input
	erc721 := NewErc721Contract(e.contractAddress, client.RPCClient())
	input, err := erc721.abi.Pack(method, from, to, amount)
	if err != nil {
		return uint64(DefaultContractGasLimit), log.WithError(err)
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

// Approve 授权转账功能
// e 		包含了ERC20合约地址 花的钱
// from 	操作者	付钱的人
// spender 	授权目标	花钱的人
// tokens	授权金额	花多少钱
func (e *Erc721Token) Approve(from *Account, spender string, tokens *big.Int) (*types.Transaction, error) {
	client, err := e.chain.Client()
	if err != nil {
		return nil, log.WithError(err)
	}
	tx := NewTransaction(from.Address(), common.HexToAddress(spender), nil, e.chain)
	erc721 := NewErc721Contract(e.contractAddress, client.RPCClient())
	return erc721.Approve(tx.ToTransactOpts(from.privateKeyECDSA), common.HexToAddress(spender), tokens)
}

// ApprovalForAll 全部授权
func (e *Erc721Token) ApprovalForAll(from *Account, spender string, flag bool) (*types.Transaction, error) {
	client, err := e.chain.Client()
	if err != nil {
		return nil, log.WithError(err)
	}
	tx := NewTransaction(from.Address(), common.HexToAddress(spender), nil, e.chain)
	erc721 := NewErc721Contract(e.contractAddress, client.RPCClient())
	return erc721.SetApprovalForAll(tx.ToTransactOpts(from.privateKeyECDSA), common.HexToAddress(spender), flag)
}
