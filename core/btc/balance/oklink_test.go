package balance

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"testing"
)

func TestOkLinkToken_GetBalance(t1 *testing.T) {
	param, _ := utils.GetBtcChainParams(utils.BtcChainTestNet3)
	address, _ := btcutil.DecodeAddress("2MzQfDPhMpCHpuGcKLwMtBNWJXpXismGLfi", param)
	type args struct {
		address btcutil.Address
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestOkLinkToken测试查询余额",
			args: args{
				address: address,
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			x := &OkLinkToken{}
			got, err := x.GetBalance(tt.args.address)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("got", got.Usable.AmountString())
			fmt.Println("got", got.Total.AmountString())
		})
	}
}
