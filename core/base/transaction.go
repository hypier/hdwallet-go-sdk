package base

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type Transaction interface {
	HexTransfer(from Account, to common.Address, value *big.Int) (string, error)
}

type TransactionDetail struct {
	Form              string
	To                string
	Time              int64
	GasLimit          uint64   //最大的limit
	GasUsed           uint64   //当前使用的limit
	GasPrice          *big.Int //最大gas单价
	EffectiveGasPrice *big.Int //当前使用的gas单价
	GasFee            *big.Int //本次交易消耗的gas费
	Value             *big.Int //交易金额
	Amount            string   //格式化之后的金额
	Data              []byte
	Type              uint8 //类型 2-EIP-1559
	TokenType         uint8 // 0-主币种 1-其他
	MethodName        string
	Status            TransactionStatus
	Err               string
	Logs              []*types.Log
	Hash              string
	ContractAddress   string
	BlockHash         string
	BlockNumber       *big.Int
	TransactionIndex  uint
}

type TransactionStatus = int

const (
	TransactionStatusNone    TransactionStatus = 0
	TransactionStatusPending TransactionStatus = 1
	TransactionStatusSuccess TransactionStatus = 2
	TransactionStatusFailure TransactionStatus = 3
)
