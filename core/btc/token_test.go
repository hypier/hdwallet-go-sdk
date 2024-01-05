package btc

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"hypier.fun/hdwallet/hdwallet-go-sdk/config"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"math/rand"
	"testing"
)

func init() {
	config.InitConfig(&config.BaseConfig{
		BaseDir:    "..",
		LogSwitch:  "CONSOLE_FILE",
		Platform:   "ALL",
		DeviceType: "UNKNOWN",
		BtcParam:   "TestNet3",
	})
}
func TestEstimateGas(t *testing.T) {
	got, err := EstimateGas()
	if err != nil {
		t.Errorf("EstimateGas() error = %v", err)
		return
	}
	fmt.Println(got)
}

func TestEstimateGasLimit(t *testing.T) {
	param, err := utils.GetBtcChainParams(utils.BtcChainTestNet3)
	if err != nil {
		t.Errorf("EstimateFee() error = %v", err)
		return
	}
	from, _ := btcutil.DecodeAddress("2MuKWyXzED48Rag2WrrLC97BgtCuteUzLDS", param)
	to, _ := btcutil.DecodeAddress("2MzQfDPhMpCHpuGcKLwMtBNWJXpXismGLfi", param)
	var params = []TransferParam{{
		To:     to,
		Amount: int64(rand.Intn(5000) + 1),
	}}
	token := NewToken(NewChain())
	amount, err := token.EstimateGasLimit(from, params, utils.BtcChainTestNet3)
	if err != nil {
		t.Errorf("EstimateFee() error = %v", err)
		return
	}
	fmt.Println(amount)
}

func TestBtcTransaction_sendRawTransaction(t1 *testing.T) {
	param, err := utils.GetBtcChainParams(utils.BtcChainTestNet3)
	b, err := NewAccountWithPrivateKey("cTFhdQbsU1xQfziSbM3FYz21a1NX6ukms12w5B3jTq1ZSXDQZqVN", utils.BtcChainTestNet3)
	if err != nil {
		fmt.Println(err)
		return
	}
	from, _ := btcutil.DecodeAddress("2MuKWyXzED48Rag2WrrLC97BgtCuteUzLDS", param)
	to, _ := btcutil.DecodeAddress("2MzQfDPhMpCHpuGcKLwMtBNWJXpXismGLfi", param)
	randomNumber := rand.Intn(1000) + 1
	allBtcUnspent, err := GetBtcUnspent(from)

	var transaction *Transaction
	var outputs = []TransferParam{{
		To:     to,
		Amount: int64(randomNumber),
	}}
	var btcUnspent = make([]BtcUnspent, 0)
	//拼接
	for i := range allBtcUnspent {
		btcUnspent = append(btcUnspent, allBtcUnspent[i])
	}
	transaction, err = NewTransaction(btcUnspent, outputs, from, 1000, param)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("size:", transaction.Tx.SerializeSize(), transaction.Tx.SerializeSizeStripped())
	err = transaction.SignWithSecretsSource(b)

	fmt.Println(len(transaction.Tx.TxIn), len(transaction.Tx.TxOut))

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("size:", transaction.Tx.SerializeSize(), transaction.Tx.SerializeSizeStripped())
	fmt.Println(transaction.Tx.TxHash())
	chain := NewChain()
	_, err = chain.CreateRemoteClientWithTimeout("https://btc.getblock.io/fdc930ba-aa6a-4a45-8ce9-506eb03f2e98/testnet/", "u", "p", utils.BtcChainTestNet3)
	if err != nil {
		fmt.Println(err)
		return
	}
	hash, err := transaction.SendRawTransaction(chain.client.rpcClient)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hash)
}

func TestToken_Transfer(t1 *testing.T) {
	chainId, err := utils.GetBtcChainId(config.Base.BtcParam)
	if err != nil {
		fmt.Println(err)
		return
	}
	chain := NewChain()
	chain.CreateRemoteClientWithTimeout("https://btc.getblock.io/fdc930ba-aa6a-4a45-8ce9-506eb03f2e98/testnet/", "u", "p", chainId)
	token := NewToken(chain)
	from, err := NewAccountWithPrivateKey("cTFhdQbsU1xQfziSbM3FYz21a1NX6ukms12w5B3jTq1ZSXDQZqVN", chainId)
	if err != nil {
		fmt.Println(err)
		return
	}
	hash, err := token.Transfer(from, "2MzQfDPhMpCHpuGcKLwMtBNWJXpXismGLfi", 100)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hash)
}
