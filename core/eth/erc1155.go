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

type Erc1155Token struct {
	*Token
	contractAddress common.Address
}

func NewErc1155Token(chain *Chain, contractAddress common.Address) *Erc1155Token {
	return &Erc1155Token{
		Token:           NewToken(chain),
		contractAddress: contractAddress,
	}
}

func (e *Erc1155Token) ContractAddress() common.Address {
	return e.contractAddress
}

func (e *Erc1155Token) TokenInfo() (*base.TokenInfo, error) {
	client, err := e.chain.Client()
	if err != nil {
		return nil, err
	}

	erc20 := NewErc1155Contract(e.contractAddress, client.RPCClient())

	name, err := erc20.Name()
	if err != nil {
		return nil, err
	}

	symbol, err := erc20.Symbol()
	if err != nil {
		return nil, err
	}

	e.Info.Name = name
	e.Info.Symbol = symbol
	e.Info.Decimal = int16(1)

	return e.Info, nil
}

func (e *Erc1155Token) BalanceOfAddress(address common.Address, tokenId *big.Int) (*base.Balance, error) {
	client, err := e.chain.Client()
	if err != nil {
		return nil, err
	}

	erc1155 := NewErc1155Contract(e.contractAddress, client.RPCClient())

	of, err := erc1155.BalanceOf(address, tokenId)
	if err != nil {
		return nil, err
	}

	return &base.Balance{
		Total: utils.NewOptAmount(of.String(), e.Info.Decimal),
	}, nil
}

func (e *Erc1155Token) EstimateGasLimit(from common.Address, gasPrice *big.Int, input []byte) (uint64, error) {

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

func (e *Erc1155Token) SafeTransferFrom(caller Account, from common.Address, to common.Address, tokenId *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	tx := NewTransaction(caller.Address(), to, nil, e.chain)

	err := tx.BuildTransfer()
	if err != nil {
		return nil, log.WithError(err)
	}

	client, err := e.chain.Client()
	if err != nil {
		return nil, log.WithError(err)
	}
	//预估gas费
	erc1155 := NewErc1155Contract(e.contractAddress, client.RPCClient())
	//得到input
	input, err := erc1155.abi.Pack("safeTransferFrom", from, to, tokenId, amount, data)
	if err != nil {
		return nil, log.WithError(err)
	}
	//预估gas费
	tx.GasLimit, err = e.EstimateGasLimit(caller.Address(), tx.GasPrice, input)
	if err != nil {
		tx.GasLimit = uint64(DefaultContractGasLimit)
	}
	transfer, err := erc1155.SafeTransferFrom(tx.ToTransactOpts(caller.privateKeyECDSA), from, to, tokenId, amount, data)
	if err != nil {
		return nil, log.WithError(err)
	}

	return transfer, nil
}

func (e *Erc1155Token) SafeBatchTransferFrom(caller Account, from common.Address, to common.Address, tokenIds []*big.Int, amounts []*big.Int, data []byte) (*types.Transaction, error) {
	tx := NewTransaction(caller.Address(), to, nil, e.chain)

	err := tx.BuildTransfer()
	if err != nil {
		return nil, log.WithError(err)
	}

	client, err := e.chain.Client()
	if err != nil {
		return nil, log.WithError(err)
	}
	erc1155 := NewErc1155Contract(e.contractAddress, client.RPCClient())
	//得到input
	input, err := erc1155.abi.Pack("safeBatchTransferFrom", from, to, tokenIds, amounts, data)
	if err != nil {
		return nil, log.WithError(err)
	}
	//预估gas费
	tx.GasLimit, err = e.EstimateGasLimit(caller.Address(), tx.GasPrice, input)
	if err != nil {
		tx.GasLimit = uint64(DefaultContractGasLimit)
	}
	transfer, err := erc1155.SafeBatchTransferFrom(tx.ToTransactOpts(caller.privateKeyECDSA), from, to, tokenIds, amounts, data)
	if err != nil {
		return nil, log.WithError(err)
	}

	return transfer, nil
}

func (e *Erc1155Token) SetApprovalForAll(caller Account, operator common.Address, approved bool) (*types.Transaction, error) {

	tx := NewTransaction(caller.Address(), e.contractAddress, nil, e.chain)

	err := tx.BuildTransfer()
	if err != nil {
		return nil, log.WithError(err)
	}

	client, err := e.chain.Client()
	if err != nil {
		return nil, log.WithError(err)
	}
	erc1155 := NewErc1155Contract(e.contractAddress, client.RPCClient())
	//得到input
	input, err := erc1155.abi.Pack("setApprovalForAll", operator, approved)
	if err != nil {
		return nil, log.WithError(err)
	}
	//预估gas费
	tx.GasLimit, err = e.EstimateGasLimit(caller.Address(), tx.GasPrice, input)
	if err != nil {
		tx.GasLimit = uint64(DefaultContractGasLimit)
	}
	all, err := erc1155.SetApprovalForAll(tx.ToTransactOpts(caller.privateKeyECDSA), operator, approved)
	if err != nil {
		return nil, log.WithError(err)
	}

	return all, nil
}
