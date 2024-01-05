package eth

import (
	"fmt"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/base"
	"math/big"
	"testing"
)

func TestChain_FetchTransactionDetail(t *testing.T) {
	type fields struct {
		client  *Client
		chainId *big.Int
	}
	type args struct {
		hash string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantDetail *base.TransactionDetail
		wantErr    bool
	}{
		{
			name:   "TestTransaction_GetTransactionInfo",
			fields: fields{},
			args: args{
				//hash: "0xcb8fa953786eb8af07543f85fd169ec97d704e870841d61f06d22368c8ead432",//主币转移
				//hash: "0x2db3456b5f258f62dfc3e4a859ce2b8145d883441b7904d77217951ca30c325c", //ERC20转移
				hash: "0xd33357a1d0e2b35898cd5b25b31076cc40047ec84273e06647a87bac16bbd480", //ERC721转移
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Chain{
				client:  tt.fields.client,
				chainId: tt.fields.chainId,
			}
			//c.CreateRemoteClient("https://ethereum-goerli.publicnode.com")
			c.CreateRemoteClient("https://goerli.blockpi.network/v1/rpc/public")
			gotDetail, err := c.FetchTransactionDetail(tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchTransactionDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log("Detail:", gotDetail.Form, gotDetail.To, gotDetail.ContractAddress, gotDetail.Value, gotDetail.Amount, gotDetail.MethodName, gotDetail.TokenType)
		})
	}
}

func TestChain_FetchTransactionStatus(t *testing.T) {
	c := NewChain()
	c.CreateRemoteClient("https://ethereum-goerli.publicnode.com")
	status, err := c.FetchTransactionStatus("0x8d397a473a7c3af41b35a54c6f32723f3aa177f9f72d0789fd9ea2db4050f8c5")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(status)
}
