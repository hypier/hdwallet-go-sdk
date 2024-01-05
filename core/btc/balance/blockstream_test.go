package balance

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/base"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/btc"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"testing"
)

func TestStreamToken_GetBalance1(t1 *testing.T) {
	param, _ := utils.GetBtcChainParams(utils.BtcChainTestNet3)
	address, _ := btcutil.DecodeAddress("2MzQfDPhMpCHpuGcKLwMtBNWJXpXismGLfi", param)
	type fields struct {
		Coin  btc.Coin
		Info  *base.TokenInfo
		chain *btc.Chain
	}
	type args struct {
		address btcutil.Address
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestStreamToken测试查询余额",
			fields: fields{
				Coin:  btc.Coin{},
				chain: btc.NewChain(),
				Info: &base.TokenInfo{
					Name:    "btc",
					Symbol:  "btc",
					Decimal: 8,
				},
			},
			args: args{
				address: address,
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &StreamToken{
				Coin:  tt.fields.Coin,
				Info:  tt.fields.Info,
				chain: tt.fields.chain,
			}
			got, err := t.GetBalance(tt.args.address)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("got", got.Usable.AmountString())
			fmt.Println("got", got.Total.AmountString())
		})
	}
}
