package trx

import (
	"fmt"
	"math/big"
	"testing"
)

func TestTrc20Token_BalanceOfAddress(t *testing.T) {
	// 暂定使用默认查询主链节点
	chain, _ := NewChain().CreateRemoteClientWithTimeout("", "")
	token := NewToken(chain)

	//token := NewToken(NewChain())
	type fields struct {
		Token           *Token
		contractAddress string
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
			name: "TestTrc20Token_BalanceOfAddress",
			fields: fields{
				Token: token,
				// todo 需要更换为正式合约地址
				contractAddress: "0x0000000000000000000000000000000000000000",
			},
			args: args{
				address: "TL158aZdqVfxpm1bLPuPkV48bhaDHFUaRh",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Trc20Token{
				Token:           tt.fields.Token,
				contractAddress: tt.fields.contractAddress,
			}
			got, err := e.BalanceOfAddress(tt.args.address)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(got)
		})
	}
}

func TestTrc20Token_Transfer(t *testing.T) {
	type fields struct {
		Token           *Token
		contractAddress string
	}
	type args struct {
		from  *Account
		to    string
		value *big.Int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "",
			fields: fields{
				Token:           NewToken(NewChain()),
				contractAddress: "0x0000000000000000000000000000000000000000",
			},
			args: args{
				from:  &Account{},
				to:    "TL158aZdqVfxpm1bLPuPkV48bhaDHFUaRh",
				value: big.NewInt(100),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Trc20Token{
				Token:           tt.fields.Token,
				contractAddress: tt.fields.contractAddress,
			}
			got, err := e.Transfer(tt.args.from, tt.args.to, tt.args.value)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(got)
		})
	}
}
