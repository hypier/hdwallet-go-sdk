package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"hypier.fun/hdwallet/hdwallet-go-sdk/config"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"math/big"
	"testing"
)

func init() {
	config.InitConfig(&config.BaseConfig{
		BaseDir:    "..",
		LogSwitch:  "CONSOLE_FILE",
		Platform:   "ALL",
		DeviceType: "UNKNOWN",
	})
}

func TestErc20Token_BalanceOfAddress(t *testing.T) {
	chain := NewChain()

	// 正面测试用例：确保给定地址的余额正确返回
	t.Run("PositiveTest - goerli", func(t *testing.T) {
		testcase := CaseAccountGoerli()
		token := NewErc20Token(chain, testcase.erc20)

		client, err := chain.CreateRemoteClient(testcase.url)
		assert.NoError(t, err)
		assert.NotNilf(t, client, "client should not be nil")

		//tokenInfo, err := token.TokenInfo()
		//assert.Nilf(t, err, "token info should be nil")
		//assert.NotNilf(t, tokenInfo, "token info should not be nil")
		//fmt.Printf("token info: %v\n", tokenInfo)

		balance, err := token.BalanceOfAddress(testcase.from.String())
		assert.NoError(t, err)
		assert.NotNilf(t, balance, "balance should not be nil")
		fmt.Printf("balance: %v\n", balance.Total.AmountString())
	})

	t.Run("PositiveTest - goerli", func(t *testing.T) {
		testcase := CaseAccountGoerli()
		token := NewErc20Token(chain, testcase.erc20)

		client, err := chain.CreateRemoteClient(testcase.url)
		assert.NoError(t, err)
		assert.NotNilf(t, client, "client should not be nil")

		//tokenInfo, err := token.TokenInfo()
		//assert.Nilf(t, err, "token info should be nil")
		//assert.NotNilf(t, tokenInfo, "token info should not be nil")
		//fmt.Printf("token info: %v\n", tokenInfo)

		balance, err := token.BalanceOfAddress(testcase.to.String())
		assert.NoError(t, err)
		assert.NotNilf(t, balance, "balance should not be nil")
		fmt.Printf("balance: %v\n", balance.Total.AmountString())
	})

}

func TestErc20Token_Transfer(t *testing.T) {
	chain := NewChain()

	// 正向测试
	t.Run("transfer success - goerli", func(t *testing.T) {
		testcase := CaseAccountGoerli()
		token := NewErc20Token(chain, testcase.erc20)

		chain.CreateRemoteClient(testcase.url)

		//tokenInfo, err := token.TokenInfo()
		//assert.Nilf(t, err, "token info should be nil")
		//assert.NotNilf(t, tokenInfo, "token info should not be nil")
		//fmt.Printf("token info: %v\n", tokenInfo)

		// 模拟账户
		from, _ := NewAccountWithPrivateKey(testcase.privateKeyHex)
		to := common.HexToAddress("0x80b8ddcFAeAC83ba89f4A256929d13b311E3A974")
		//value, _ := utils.ParseAmount("0.00001", token.Info.Decimal)

		transfer, err := token.Transfer(from, to, big.NewInt(1000))
		assert.Nilf(t, err, "transfer should be nil")
		assert.NotNilf(t, transfer.Hash(), "transfer should not be nil")

		fmt.Printf("transfer hash: %x\n", transfer.Hash())
	})

	// 负向测试
	t.Run("transfer failed - Insufficient balance", func(t *testing.T) {
		testcase := CaseAccountGoerliNew()
		token := NewErc20Token(chain, testcase.erc20)

		chain.CreateRemoteClient(testcase.url)

		from, _ := NewAccountWithPrivateKey(testcase.privateKeyHex)
		to := testcase.to
		value, _ := utils.ParseAmount("0.00001", token.Info.Decimal)

		transfer, err := token.Transfer(from, to, value.BigInt())
		assert.NotNilf(t, err, "transfer should be nil, but got %v", err)
		assert.Nilf(t, transfer, "transfer should not be nil")

	})
}

func TestErc20Token_EstimateGasLimit(t *testing.T) {
	chain := NewChain()
	// 正向测试
	t.Run("EstimateGasLimit success - goerli", func(t *testing.T) {
		testcase := CaseAccountGoerli()
		token := NewErc20Token(chain, testcase.erc20)

		chain.CreateRemoteClient(testcase.url)

		// 模拟账户
		from := testcase.from
		to := common.HexToAddress("0x80b8ddcFAeAC83ba89f4A256929d13b311E3A974")
		decimal, _ := token.GetDecimal()
		value, err := utils.ParseAmount("0.00001", decimal)
		assert.Nilf(t, err, "parse amount should be nil")

		gasLimit, err := token.EstimateGasLimit(from, to, value.BigInt(), "transfer")
		assert.Nilf(t, err, "estimateGasLimit should be nil")
		assert.NotNilf(t, gasLimit, "estimateGasLimit should not be nil")

		fmt.Printf("estimateGasLimit gas: %v\n", gasLimit)
	})

	// 正向测试
	t.Run("EstimateGasLimit success - goerli", func(t *testing.T) {
		testcase := CaseAccountGoerli()
		token := NewErc20Token(chain, testcase.erc20)

		chain.CreateRemoteClient(testcase.url)

		// 模拟账户
		from := CaseAccountGoerliNew().from
		to := testcase.to

		gasLimit, err := token.EstimateGasLimit(from, to, big.NewInt(10000), "approve")
		assert.Nilf(t, err, "estimateGasLimit should be nil")
		assert.NotNilf(t, gasLimit, "estimateGasLimit should not be nil")

		fmt.Printf("estimateGasLimit gas: %v\n", gasLimit)
	})

	// 负向测试
	t.Run("EstimateGasLimit failed - Insufficient balance", func(t *testing.T) {
		testcase := CaseAccountGoerliNew()
		token := NewErc20Token(chain, testcase.erc20)
		chain.CreateRemoteClient(testcase.url)

		tokenInfo, err := token.TokenInfo()
		assert.Nilf(t, err, "token info should be nil")
		assert.NotNilf(t, tokenInfo, "token info should not be nil")

		from := testcase.from
		to := testcase.to
		value, _ := utils.ParseAmount("0.00001", token.Info.Decimal)

		_, err = token.EstimateGasLimit(from, to, value.BigInt(), "transfer")
		assert.NotNilf(t, err, "estimateGasLimit should be nil, but got %v", err)

		fmt.Printf("estimateGasLimit error: %v\n", err)

	})
}

func TestErc20Token_Approve(t *testing.T) {
	chain := NewChain()

	// 正向测试
	t.Run("Approve success - goerli", func(t *testing.T) {
		testcase := CaseAccountGoerli()
		token := NewErc20Token(chain, testcase.erc20)

		chain.CreateRemoteClient(testcase.url)

		// 模拟账户
		from, _ := NewAccountWithPrivateKey(testcase.privateKeyHex)
		to := common.HexToAddress("0x80b8ddcFAeAC83ba89f4A256929d13b311E3A974")

		approve, err := token.Approve(from, to, big.NewInt(10000))
		assert.Nilf(t, err, "Approve should be nil")
		assert.NotNilf(t, approve, "transfer should not be nil")
		assert.NotNilf(t, approve.Hash(), "transfer should not be nil")

		fmt.Printf("transfer hash: %x\n", approve.Hash())
	})

	// 负向测试
	t.Run("Approve failed - Insufficient balance", func(t *testing.T) {
		testcase := CaseAccountGoerliNew()
		token := NewErc20Token(chain, testcase.erc20)

		chain.CreateRemoteClient(testcase.url)

		from, _ := NewAccountWithPrivateKey(testcase.privateKeyHex)
		to := testcase.to

		transfer, err := token.Approve(from, to, big.NewInt(10000))
		assert.NotNilf(t, err, "transfer should be nil, but got %v", err)
		assert.Nilf(t, transfer, "transfer should not be nil")

	})
}
