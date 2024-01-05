package eth

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/base"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"math/big"
	"strings"
)

// FetchTransactionDetail GetTransactionInfo 获取交易明细
func (c *Chain) FetchTransactionDetail(hash string) (*base.TransactionDetail, error) {
	if len(hash) == 0 {
		return nil, fmt.Errorf("无效的交易Hash")
	}
	client, err := c.Client()
	if err != nil {
		return nil, log.WithError(err, "Client failed")
	}
	detail := &base.TransactionDetail{
		Hash: hash,
	}
	//根据hash获取交易对象
	transaction, isPending, err := client.RPCClient().TransactionByHash(context.Background(), common.HexToHash(detail.Hash))
	if err != nil {
		return nil, log.WithError(err, "TransactionByHash failed")
	}
	//组装基础交易对象存在的信息
	//主币种
	if len(transaction.Data()) == 0 {
		//转移目标
		detail.To = transaction.To().String()
		detail.Value = transaction.Value()
		detail.TokenType = 0
	} else {
		//如果不是主币 To地址则是合约地址
		detail.ContractAddress = transaction.To().String()
		//得到to地址
		method := getToAndValue(transaction)
		detail.To = method.Recipient.String()
		detail.Value = method.Amount
		detail.TokenType = 1
		detail.MethodName = method.Name
	}
	switch detail.TokenType {
	case 0: //主币种
		amount := utils.NewOptAmount(detail.Value.String(), NewToken(c).GetDecimal())
		detail.Amount = amount.AmountString()
		break
	case 1: //其他合约
		decimal, err := NewErc20Token(c, common.HexToAddress(detail.ContractAddress)).GetDecimal()
		if err != nil { //其他精度默认为0 也就是原样返回
			decimal = 0
		}
		amount := utils.NewOptAmount(detail.Value.String(), decimal)
		detail.Amount = amount.AmountString()
		break
	}
	detail.GasLimit = transaction.Gas()
	detail.GasPrice = transaction.GasPrice()
	detail.Data = transaction.Data()
	detail.Type = transaction.Type()
	detail.Time = transaction.Time().UnixMilli()
	//解码交易对象 得到form
	sender, err := decodeSigner(transaction)
	if err != nil {
		return nil, log.WithError(err, "decodeSigner failed")
	}
	//form
	detail.Form = sender.String()
	//如果在等待中
	if isPending {
		detail.Status = base.TransactionStatusPending
		return detail, nil
	}
	//获取交易明细
	receipt, err := client.RPCClient().TransactionReceipt(context.Background(), common.HexToHash(detail.Hash))
	if err != nil {
		return nil, log.WithError(err, "TransactionReceipt failed")
	}
	//组装明细对象中的数据
	detail.GasUsed = receipt.GasUsed
	detail.EffectiveGasPrice = receipt.EffectiveGasPrice
	detail.Logs = receipt.Logs
	detail.BlockHash = receipt.BlockHash.String()
	detail.BlockNumber = receipt.BlockNumber
	detail.TransactionIndex = receipt.TransactionIndex
	//失败了
	if receipt.Status == 0 {
		detail.Status = base.TransactionStatusFailure
		//获取失败的详情
		_, err := client.RPCClient().CallContract(context.Background(), ethereum.CallMsg{
			From:       sender,
			To:         transaction.To(),
			Data:       transaction.Data(),
			Gas:        transaction.Gas(),
			GasPrice:   transaction.GasPrice(),
			GasFeeCap:  transaction.GasFeeCap(),
			GasTipCap:  transaction.GasTipCap(),
			Value:      transaction.Value(),
			AccessList: transaction.AccessList(),
		}, receipt.BlockNumber)
		if err != nil {
			detail.Err = err.Error()
		}
	} else {
		detail.Status = base.TransactionStatusSuccess
	}
	detail.GasFee = big.NewInt(0).Mul(receipt.EffectiveGasPrice, big.NewInt(0).SetUint64(receipt.GasUsed))
	return detail, nil
}

type Method struct {
	Recipient common.Address
	Amount    *big.Int
	Name      string
}

func getToAndValue(tx *types.Transaction) Method {
	resultMethod := Method{}
	parsed, err := abi.JSON(strings.NewReader(ERC20InterfaceABI))
	if err != nil {
		return resultMethod
	}
	method, err := parsed.MethodById(tx.Data())
	//ERC20解析失败
	if err != nil {
		parsed, err = abi.JSON(strings.NewReader(ERC721InterfaceABI))
		if err != nil {
			return resultMethod
		}
		method, err = parsed.MethodById(tx.Data())
		//解析ERC721失败
		if err != nil {
			parsed, err = abi.JSON(strings.NewReader(ERC1155InterfaceABI))
			if err != nil {
				return resultMethod
			}
			method, err = parsed.MethodById(tx.Data())
			if err != nil {
				return resultMethod
			}
		}
	}
	resultMethod.Name = method.Name
	toIndex := -1
	valueIndex := -1
	for i, param := range method.Inputs {
		if param.Name == "recipient" {
			toIndex = i
		} else if param.Name == "amount" {
			valueIndex = i
		}
	}
	params, err := method.Inputs.Unpack(tx.Data()[4:])
	if err != nil || toIndex == -1 || valueIndex == -1 {
		return resultMethod
	}
	resultMethod.Amount = params[valueIndex].(*big.Int)
	resultMethod.Recipient = params[toIndex].(common.Address)
	return resultMethod
}

// decodeSigner 解码交易信息 得到form
func decodeSigner(txn *types.Transaction) (common.Address, error) {
	var signer types.Signer
	switch {
	case txn.Type() == types.AccessListTxType:
		signer = types.NewEIP2930Signer(txn.ChainId())
	case txn.Type() == types.DynamicFeeTxType:
		signer = types.NewLondonSigner(txn.ChainId())
	default:
		signer = types.NewEIP155Signer(txn.ChainId())
	}
	return types.Sender(signer, txn)
}

func (c *Chain) FetchTransactionStatus(hash string) (base.TransactionStatus, error) {
	_, isPending, err := c.client.RPCClient().TransactionByHash(context.Background(), common.HexToHash(hash))
	if err != nil {
		return base.TransactionStatusNone, log.WithError(err, "TransactionByHash failed")
	}
	if isPending {
		return base.TransactionStatusPending, nil
	}
	//获取交易明细
	receipt, err := c.client.RPCClient().TransactionReceipt(context.Background(), common.HexToHash(hash))
	if err != nil {
		return base.TransactionStatusNone, log.WithError(err, "TransactionReceipt failed")
	}
	if receipt.Status == 0 {
		return base.TransactionStatusFailure, nil
	}
	return base.TransactionStatusSuccess, nil
}
