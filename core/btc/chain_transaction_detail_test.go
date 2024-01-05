package btc

import (
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"testing"
)

func TestChain_FetchTransactionDetail(t *testing.T) {
	chain := NewChain()
	_, err := chain.CreateRemoteClientWithTimeout("btc.getblock.io/fdc930ba-aa6a-4a45-8ce9-506eb03f2e98/testnet/", "u", "p", utils.BtcChainTestNet3)
	if err != nil {
		t.Error(err)
		return
	}
	detail, err := chain.FetchTransactionDetail("604520d6133dbacc55d15ea76d42797e88a0cc384153d3eb6524da90dbcc33f6")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(detail)
}
