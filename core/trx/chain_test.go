package trx

import (
	"fmt"
	"hypier.fun/hdwallet/hdwallet-go-sdk/config"
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
func TestChain_FetchTransactionDetail(t *testing.T) {
	chain := NewChain()
	_, err := chain.CreateRemoteClientWithTimeout("", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	//主币交易 b541de25851a7ca70a1e4d6bbe9c8c006142c60d05cd341c0c7dc475f2ded75b
	//自定义代币交易 4dad085811268d8f97315e42782fb8ef8da39db9986f4b0b7a77e27085d96463
	//usdt 35f0bbda2d0093f27c42e4f03acf9f01058c32d6f34326fd01703de7efbd21a7
	detail, err := chain.FetchTransactionDetail("b541de25851a7ca70a1e4d6bbe9c8c006142c60d05cd341c0c7dc475f2ded75b")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(detail)
}
