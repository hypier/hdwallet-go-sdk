package trx

import (
	"fmt"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/base"
	"testing"
)

func TestToken_BalanceOfAddress(t1 *testing.T) {
	// todo 差rpcUrl和apiKey
	//chain, _ := NewChain().CreateRemoteClientWithTimeout("https://trx.getblock.io/9b015f60-6a39-4624-b60d-a9954de4f61e/mainnet/", "")
	// 暂定使用默认查询主链节点
	chain, _ := NewChain().CreateRemoteClientWithTimeout("", "")
	type fields struct {
		Coin  Coin
		Info  *base.TokenInfo
		chain *Chain
	}
	type args struct {
		address string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestToken_BalanceOfAddress_查询余额",
			fields: fields{
				Coin: Coin{},
				Info: &base.TokenInfo{
					Name:    "TRX",
					Symbol:  "TRX",
					Decimal: 6,
				},
				chain: chain,
			},
			args: args{
				address: "TC4HdkqWtq8Y1AsNQiQ9GM4fRkoZcNjKZm",
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Token{
				Coin:  tt.fields.Coin,
				Info:  tt.fields.Info,
				chain: tt.fields.chain,
			}
			got, err := t.BalanceOfAddress(tt.args.address)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(got.Total.AmountString())
			fmt.Println(got.Usable.AmountString())
		})
	}

}

func TestToken_Transfer(t1 *testing.T) {
	// 暂定使用默认查询主链节点
	chain, _ := NewChain().CreateRemoteClientWithTimeout("", "")
	// TODO 通过私钥还原账户，测试的时候需要把私钥填写进去，提交的时候去掉私钥，避免自己的私钥泄露
	account, _ := NewAccountWithPrivateKey("")
	fmt.Println("account--------->", account.Address())
	type fields struct {
		Coin  Coin
		Info  *base.TokenInfo
		chain *Chain
	}
	type args struct {
		from  *Account
		to    string
		value int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestToken_Transfer_转账",
			fields: fields{
				Coin: Coin{},
				Info: &base.TokenInfo{
					Name:    "TRX",
					Symbol:  "TRX",
					Decimal: 6,
				},
				chain: chain,
			},
			args: args{
				from:  account,
				to:    "TL158aZdqVfxpm1bLPuPkV48bhaDHFUaRh",
				value: int64(100000),
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Token{
				Coin:  tt.fields.Coin,
				Info:  tt.fields.Info,
				chain: tt.fields.chain,
			}
			got, err := t.Transfer(tt.args.from, tt.args.to, tt.args.value)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(got)
		})
	}
}
