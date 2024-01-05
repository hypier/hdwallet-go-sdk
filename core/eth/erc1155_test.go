package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/base"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"math/big"
	"reflect"
	"testing"
)

const (
	Erc1155ContractAddress = "0x53E31aecD9d607143647B90e4AB4b1AB5F191d46"
	T_From                 = "0xac4C7169A6F069034359e695db932027431866DF"
	T_To                   = "0xD68f8adD960094066E198CDbb3e0f59eBba5c0Bf"
	T_P                    = "3ed6d2390fd5269e06f16ca352e873354efc26aa2e36fa60b1e9d00c2dd8c937"
)

// success
func TestErc1155Token_BalanceOfAddress(t *testing.T) {
	chain := NewChain()
	_, err := chain.CreateRemoteClient("https://ethereum-goerli.publicnode.com")
	if err != nil {
		return
	}
	type fields struct {
		Token           *Token
		contractAddress common.Address
	}
	type args struct {
		address common.Address
		tokenId *big.Int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *base.Balance
		wantErr bool
	}{
		{
			name: "Positive-TestErc1155Token_BalanceOfAddress",
			fields: fields{
				Token: &Token{
					Info: &base.TokenInfo{
						Name:    "ETH",
						Symbol:  "ETH",
						Decimal: 18,
					}, chain: chain,
				},
				contractAddress: common.HexToAddress(Erc1155ContractAddress),
			},
			args: args{
				address: common.HexToAddress(T_From),
				tokenId: big.NewInt(1),
			},
			want: &base.Balance{
				Total:  utils.NewOptAmount("1", 16),
				Usable: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Erc1155Token{
				Token:           tt.fields.Token,
				contractAddress: tt.fields.contractAddress,
			}
			got, err := e.BalanceOfAddress(tt.args.address, tt.args.tokenId)
			assert.NoError(t, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("BalanceOfAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Total.AmountString() != tt.want.Total.AmountString() {
				fmt.Println("BalanceOfAddress() got =", got.Total.AmountString(), " want =", tt.want.Total.AmountString(), " 与期望不符")
			}
			if got.Total.AmountString() == tt.want.Total.AmountString() {
				fmt.Println("BalanceOfAddress() got =", got.Total.AmountString(), " want =", tt.want.Total.AmountString(), " 与期望相符")
			}
		})
	}
}

// success
func TestErc1155Token_ContractAddress(t *testing.T) {
	type fields struct {
		Token           *Token
		contractAddress common.Address
	}
	tests := []struct {
		name   string
		fields fields
		want   common.Address
	}{
		{
			name: "TestErc1155Token_ContractAddress",
			fields: fields{
				Token: &Token{
					Info: &base.TokenInfo{
						Name:    "ETH",
						Symbol:  "ETH",
						Decimal: 18,
					},
				},
				contractAddress: common.HexToAddress(Erc1155ContractAddress),
			},
			want: common.HexToAddress(Erc1155ContractAddress),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Erc1155Token{
				Token:           tt.fields.Token,
				contractAddress: tt.fields.contractAddress,
			}
			got := e.ContractAddress()
			if got.String() != tt.want.String() {
				t.Errorf("ContractAddress() = %v, want %v", got, tt.want)
			}
			if got := e.ContractAddress(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ContractAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

// success
func TestErc1155Token_SafeBatchTransferFrom(t *testing.T) {
	account, err := NewAccountWithPrivateKey(T_P)
	if err != nil {
		return
	}

	chain := NewChain()
	_, err = chain.CreateRemoteClient("https://ethereum-goerli.publicnode.com")
	if err != nil {
		return
	}
	type fields struct {
		Token           *Token
		contractAddress common.Address
	}
	type args struct {
		caller   Account
		from     common.Address
		to       common.Address
		tokenIds []*big.Int
		amounts  []*big.Int
		data     []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Transaction
		wantErr bool
	}{
		{
			name: "TestErc1155Token_SafeBatchTransferFrom",
			fields: fields{
				Token: &Token{
					Info: &base.TokenInfo{
						Name:    "ETH",
						Symbol:  "ETH",
						Decimal: 18,
					},
					chain: chain,
				},
				contractAddress: common.HexToAddress(Erc1155ContractAddress),
			},
			args: args{
				caller:   *account,
				from:     common.HexToAddress(T_From),
				to:       common.HexToAddress(T_To),
				tokenIds: []*big.Int{big.NewInt(2), big.NewInt(3)},
				amounts:  []*big.Int{big.NewInt(1), big.NewInt(1)},
				data:     nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Erc1155Token{
				Token:           tt.fields.Token,
				contractAddress: tt.fields.contractAddress,
			}
			got, err := e.SafeBatchTransferFrom(tt.args.caller, tt.args.from, tt.args.to, tt.args.tokenIds, tt.args.amounts, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeBatchTransferFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log("hash:", got.Hash().String())
		})
	}
}

// success
func TestErc1155Token_SafeTransferFrom(t *testing.T) {
	account, err := NewAccountWithPrivateKey(T_P)
	if err != nil {
		return
	}

	chain := NewChain()
	_, err = chain.CreateRemoteClient("https://ethereum-goerli.publicnode.com")
	if err != nil {
		return
	}
	type fields struct {
		Token           *Token
		contractAddress common.Address
	}
	type args struct {
		caller  Account
		from    common.Address
		to      common.Address
		tokenId *big.Int
		amount  *big.Int
		data    []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Transaction
		wantErr bool
	}{
		{
			name: "TestErc1155Token_SafeTransferFrom",
			fields: fields{
				Token: &Token{
					Info: &base.TokenInfo{
						Name:    "ETH",
						Symbol:  "ETH",
						Decimal: 18,
					},
					chain: chain,
				},
				contractAddress: common.HexToAddress(Erc1155ContractAddress),
			},
			args: args{
				caller:  *account,
				from:    common.HexToAddress(T_From),
				to:      common.HexToAddress(T_To),
				tokenId: big.NewInt(2),
				amount:  big.NewInt(1),
				data:    nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Erc1155Token{
				Token:           tt.fields.Token,
				contractAddress: tt.fields.contractAddress,
			}
			got, err := e.SafeTransferFrom(tt.args.caller, tt.args.from, tt.args.to, tt.args.tokenId, tt.args.amount, tt.args.data)
			fmt.Println(got.GasPrice())
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeTransferFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("hash:", got.Hash().String())
		})
	}
}

// success
func TestErc1155Token_TokenInfo(t *testing.T) {
	chain := NewChain()
	_, err := chain.CreateRemoteClient("https://ethereum-goerli.publicnode.com")
	if err != nil {
		return
	}
	type fields struct {
		Token           *Token
		contractAddress common.Address
	}
	tests := []struct {
		name    string
		fields  fields
		want    *base.TokenInfo
		wantErr bool
	}{
		{
			name: "TestErc1155Token_TokenInfo",
			fields: fields{
				Token: &Token{
					Info: &base.TokenInfo{
						Name:    "ETH",
						Symbol:  "ETH",
						Decimal: 18,
					}, chain: chain,
				},
				contractAddress: common.HexToAddress(Erc1155ContractAddress),
			},
			want: &base.TokenInfo{
				Name:    "cat",
				Symbol:  "cat",
				Decimal: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Erc1155Token{
				Token:           tt.fields.Token,
				contractAddress: tt.fields.contractAddress,
			}
			got, err := e.TokenInfo()
			fmt.Println("Name:", got.Name)
			fmt.Println("Symbol:", got.Symbol)
			fmt.Println("Decimal:", got.Decimal)
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TokenInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// success
func TestErc1155Token_SetApprovalForAll(t *testing.T) {
	account, err := NewAccountWithPrivateKey(T_P)
	if err != nil {
		return
	}

	chain := NewChain()
	_, err = chain.CreateRemoteClient("https://ethereum-goerli.publicnode.com")
	if err != nil {
		return
	}
	type fields struct {
		Token           *Token
		contractAddress common.Address
	}
	type args struct {
		caller   Account
		operator common.Address
		approved bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Transaction
		wantErr bool
	}{
		{
			name: "TestErc1155Token_SetApprovalForAll",
			fields: fields{
				Token: &Token{
					Info: &base.TokenInfo{
						Name:    "ETH",
						Symbol:  "ETH",
						Decimal: 18,
					},
					chain: chain,
				},
				contractAddress: common.HexToAddress(Erc1155ContractAddress),
			},
			args: args{
				caller:   *account,
				operator: common.HexToAddress(T_To),
				approved: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Erc1155Token{
				Token:           tt.fields.Token,
				contractAddress: tt.fields.contractAddress,
			}
			got, err := e.SetApprovalForAll(tt.args.caller, tt.args.operator, tt.args.approved)
			fmt.Println(got.GasPrice())
			fmt.Println(got.Hash().String())
			if (err != nil) != tt.wantErr {
				t.Errorf("SetApprovalForAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
