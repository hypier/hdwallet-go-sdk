package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/net/context"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"math/big"
	"strings"
)

func DeployContract(opts *bind.TransactOpts, backend bind.ContractBackend, jsonAbi string, bytecode string, params ...interface{}) (common.Address, *types.Transaction, error) {
	parsed, err := abi.JSON(strings.NewReader(jsonAbi))
	if err != nil {
		return common.Address{}, nil, log.WithError(err, "abi.JSON failed")
	}
	address, tx, _, err := bind.DeployContract(opts, parsed, common.FromHex(bytecode), backend, params...)
	if err != nil {
		return common.Address{}, nil, log.WithError(err, "DeployContract failed")
	}
	return address, tx, nil
}

// SendTransaction 执行读方法
func SendTransaction(chain *Chain, contract common.Address, jsonAbi string, method string, params ...interface{}) ([]interface{}, error) {

	parsed, err := abi.JSON(strings.NewReader(jsonAbi))
	if err != nil {
		return nil, log.WithError(err, "abi.JSON failed")
	}
	c := bind.NewBoundContract(contract, parsed, chain.client.RPCClient(), chain.client.RPCClient(), chain.client.RPCClient())
	var out []interface{}
	err = c.Call(&bind.CallOpts{}, &out, method, params...)
	if err != nil {
		return nil, log.WithError(err, "Call failed")
	}
	return out, nil
}

// SendRawTransaction 执行写方法
func SendRawTransaction(chain *Chain, privateKey string, contract common.Address, jsonAbi string, method string, params ...interface{}) (*types.Transaction, error) {
	account, err := NewAccountWithPrivateKey(privateKey)
	if err != nil {
		return nil, log.WithError(err)
	}
	tx := *NewTransaction(account.Address(), contract, nil, chain)
	opts := tx.ToTransactOpts(account.privateKeyECDSA)
	parsed, err := abi.JSON(strings.NewReader(jsonAbi))
	if err != nil {
		return nil, log.WithError(err, "abi.JSON failed")
	}
	c := bind.NewBoundContract(contract, parsed, chain.client.RPCClient(), chain.client.RPCClient(), chain.client.RPCClient())
	return c.Transact(opts, method, params...)
}
func EstimateGasLimit(from common.Address, backend bind.ContractBackend, contractAddress common.Address, input []byte) (int64, error) {
	param := TransactBaseParam{
		From: from,
	}
	err := param.EnsureGasPrice(backend)
	if err != nil {
		return 0, log.WithError(err, "EnsureGasPrice failed")
	}

	ethValue := param.EthValue
	if ethValue == nil {
		ethValue = big.NewInt(0)
	}
	msg := ethereum.CallMsg{From: param.From, To: &contractAddress,
		GasPrice:  param.GasPrice,
		GasFeeCap: param.GasFeeCap,
		GasTipCap: param.GasTipCap,
		Value:     ethValue, Data: input}

	gasLimit, err := backend.EstimateGas(context.Background(), msg)
	if err != nil {
		return 0, log.WithError(err, "failed to estimate gas needed: %v")
	}
	return int64(gasLimit), nil
}

type TransactBaseParam struct {
	From      common.Address
	EthValue  *big.Int
	GasPrice  *big.Int
	GasFeeCap *big.Int
	GasTipCap *big.Int
	BaseFee   *big.Int
}

func (_self *TransactBaseParam) EnsureGasPrice(backend bind.ContractBackend) error {
	head, err := backend.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return log.WithError(err, "HeaderByNumber failed")
	}
	_self.BaseFee = head.BaseFee

	if head.BaseFee == nil {
		if _self.GasPrice == nil {
			price, err := backend.SuggestGasPrice(context.Background())
			if err != nil {
				return log.WithError(err, "SuggestGasPrice failed")
			}
			_self.GasPrice = price
		}
	} else {
		if _self.GasTipCap == nil {
			tip, err := backend.SuggestGasTipCap(context.Background())
			if err != nil {
				return log.WithError(err, "SuggestGasTipCap failed")
			}
			_self.GasTipCap = tip
		}
		if _self.GasFeeCap == nil {
			gasFeeCap := new(big.Int).Add(
				_self.GasTipCap,
				new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
			)
			_self.GasFeeCap = gasFeeCap
		}
		if _self.GasFeeCap.Cmp(_self.GasTipCap) < 0 {
			return fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", _self.GasFeeCap, _self.GasTipCap)
		}
	}
	return nil
}
