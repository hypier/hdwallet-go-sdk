package eth

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"

	//_ "hypier.fun/hdwallet/hdwallet-go-sdk/core/base"
	"testing"
)

func TestToken_BalanceOfAddress(t *testing.T) {
	chain := NewChain()
	token := NewToken(chain)
	token.TokenInfo()
	// 正面测试用例：确保给定地址的余额正确返回
	t.Run("PositiveTest - goerli", func(t *testing.T) {
		testcase := CaseAccountGoerli()

		client, err := chain.CreateRemoteClient(testcase.url)
		assert.NoError(t, err)
		assert.NotNilf(t, client, "client should not be nil")

		balance, err := token.BalanceOfAddress(testcase.from.String())
		assert.NoError(t, err)
		assert.NotNilf(t, balance, "balance should not be nil")

		fmt.Printf("balance: %v\n", balance.Total.AmountString())
	})

}

func TestToken_Transfer(t *testing.T) {
	chain := NewChain()
	token := NewToken(chain)
	token.TokenInfo()
	// 正向测试
	t.Run("transfer success - goerli", func(t *testing.T) {
		testcase := CaseAccountGoerli()

		chain.CreateRemoteClient(testcase.url)
		// 模拟账户
		from, _ := NewAccountWithPrivateKey(testcase.privateKeyHex)
		to := testcase.to
		value, _ := utils.ParseAmount("0.00001", token.Info.Decimal)

		transfer, err := token.Transfer(*from, to, value.BigInt())
		assert.Nilf(t, err, "transfer should be nil")
		assert.NotNilf(t, transfer.Hash(), "transfer should not be nil")

		fmt.Printf("transfer hash: %x\n", transfer.Hash())
	})

	// 负向测试
	t.Run("transfer failed - Insufficient balance", func(t *testing.T) {
		testcase := CaseAccountGoerliNew()

		chain.CreateRemoteClient(testcase.url)

		from, _ := NewAccountWithPrivateKey(testcase.privateKeyHex)
		to := testcase.to
		value, _ := utils.ParseAmount("10000000", token.Info.Decimal)

		transfer, err := token.Transfer(*from, to, value.BigInt())
		assert.Error(t, err, "transfer should be error")
		assert.Nilf(t, transfer, "transfer should not be nil")
	})
}

func TestToken_EstimateGasLimit(t *testing.T) {
	chain := NewChain()
	token := NewToken(chain)
	token.TokenInfo()
	// 正向测试
	t.Run("transfer success - goerli", func(t *testing.T) {
		testcase := CaseAccountGoerli()

		chain.CreateRemoteClient(testcase.url)

		// 模拟账户
		from := testcase.from
		to := testcase.to
		value, _ := utils.ParseAmount("0.00001", token.Info.Decimal)

		gasLimit, err := token.EstimateGasLimit(from, to, value.BigInt())
		assert.Nilf(t, err, "estimateGasLimit should be nil")
		assert.NotNilf(t, gasLimit, "estimateGasLimit should not be nil")

		fmt.Printf("estimateGasLimit gas: %v\n", gasLimit)
	})

	// 负向测试
	t.Run("transfer failed - Insufficient balance", func(t *testing.T) {
		testcase := CaseAccountGoerliNew()

		chain.CreateRemoteClient(testcase.url)

		from := testcase.from
		to := testcase.to
		value, _ := utils.ParseAmount("0.00001", token.Info.Decimal)

		_, err := token.EstimateGasLimit(from, to, value.BigInt())
		assert.NotNilf(t, err, "estimateGasLimit should be nil, but got %v", err)

		fmt.Printf("estimateGasLimit error: %v\n", err)

	})
}
