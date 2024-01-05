package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
)

func TestErc721Token_Approve(t *testing.T) {

	// 创建chain
	ch := NewChain()
	// 创建链接
	ch.CreateRemoteClient("https://test.lixb.io")
	// 得到token
	tok := NewToken(ch)
	// 先根据私钥恢复出Account信息
	acc, _ := NewAccountWithPrivateKey("915360ebde333f1db34447c7132b06491aa0142387656bdab449ada74578e837")
	// 得到钱包地址：0x61AA16822d6CCC7DBD6cA8d28272773bE7AB08aD
	type fields struct {
		Token           *Token
		contractAddress common.Address
	}
	type args struct {
		from    *Account
		spender string // 授权目标	花钱的人
		tokens  *big.Int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "测试Approve",
			fields: fields{
				Token:           tok,
				contractAddress: common.HexToAddress("0x46834493377772db72Ae802470200Fc711178B45"), // 根据钱包地址得到common.Address
			},
			args: args{
				from:    acc,
				spender: "0x1800999d037FA401286A9977316555728b2808dA", // 授权钱包地址
				tokens:  big.NewInt(513348752459698707),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Erc721Token{
				Token:           tt.fields.Token,
				contractAddress: tt.fields.contractAddress,
			}
			got, _ := e.Approve(tt.args.from, tt.args.spender, tt.args.tokens)
			fmt.Println("got结果：", got)
			fmt.Println("got结果-hash：", got.Hash().String())
		})
	}
}

func TestErc721Token_ApproveAll(t *testing.T) {
	// 创建chain
	ch := NewChain()
	// 创建链接
	ch.CreateRemoteClient("https://test.lixb.io")
	// 得到token
	tok := NewToken(ch)
	// 先根据私钥恢复出Account信息
	acc, _ := NewAccountWithPrivateKey("915360ebde333f1db34447c7132b06491aa0142387656bdab449ada74578e837")
	// 得到钱包地址：0x61AA16822d6CCC7DBD6cA8d28272773bE7AB08aD
	type fields struct {
		Token           *Token
		contractAddress common.Address
	}
	type args struct {
		from    *Account
		spender string
		flag    bool // 是否授权全部 true 是  false 否
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "测试ApproveAll",
			fields: fields{
				Token:           tok,
				contractAddress: common.HexToAddress("0x46834493377772db72Ae802470200Fc711178B45"), // 合约地址
			},
			args: args{
				from:    acc,
				spender: "0x1800999d037FA401286A9977316555728b2808dA", // 授权钱包地址
				flag:    true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Erc721Token{
				Token:           tt.fields.Token,
				contractAddress: tt.fields.contractAddress,
			}
			got, _ := e.ApprovalForAll(tt.args.from, tt.args.spender, tt.args.flag)
			fmt.Println("got结果：", got)
			fmt.Println("got结果-hash：", got.Hash().String())
		})
	}
}

// 测试通过
func TestErc721Token_Transfer(t *testing.T) {
	// 创建chain
	ch := NewChain()
	// 创建链接
	ch.CreateRemoteClient("https://ethereum-goerli.publicnode.com")
	// 得到token
	tok := NewToken(ch)

	// 通过私钥得到账户
	acc, _ := NewAccountWithPrivateKey("915360ebde333f1db34447c7132b06491aa0142387656bdab449ada74578e837")
	fmt.Println("通过私钥还原得到钱包地址(公钥):", acc.Address().String())
	type fields struct {
		Token           *Token
		contractAddress common.Address
	}
	type args struct {
		from  Account
		to    common.Address
		value *big.Int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "测试Erc721Token_Transfer",
			fields: fields{
				Token:           tok,
				contractAddress: common.HexToAddress("0xA0B3Fe2d3465974D5DbeA36783f15ae624C8107A"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("本次操作合约地址:", tt.fields.contractAddress)
			e := &Erc721Token{
				Token:           tt.fields.Token,
				contractAddress: tt.fields.contractAddress,
			}
			got, err := e.Transfer(*acc, common.HexToAddress("0x61AA16822d6CCC7DBD6cA8d28272773bE7AB08aD"), big.NewInt(513348752514224658))
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("交易返回结果:", got.Hash())
		})
	}
}
