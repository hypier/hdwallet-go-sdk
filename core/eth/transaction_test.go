package eth

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"math/big"
	"testing"
)

var (
	from   = common.HexToAddress("0xC7bFc5CB33ADA7566217ebbF15B5Dc25f5e609D7")
	pri, _ = crypto.HexToECDSA("1032adbf75a73f959d30dcae3e35a2c12252daac44abf8d9d2e21b29754db496")
	to     = common.HexToAddress("0x80b8ddcFAeAC83ba89f4A256929d13b311E3A974")
	url    = "https://ethereum-goerli.publicnode.com"
)

func TestNewTransaction(t *testing.T) {
	// 创建chain
	ch := NewChain()
	// 创建链接
	ch.CreateRemoteClient("https://ethereum-goerli.publicnode.com")
	type args struct {
		from  common.Address
		to    common.Address
		value *big.Int
		chain *Chain
	}
	tests := []struct {
		name string
		args args
		want *Transaction
	}{
		{
			name: "TestNewTransaction",
			args: args{
				from:  from,
				to:    to,
				value: big.NewInt(100),
				chain: ch,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTransaction(tt.args.from, tt.args.to, tt.args.value, tt.args.chain)
			t.Log("transaction:", got)
		})
	}
}

func TestTransaction_ToTransactOpts(t *testing.T) {
	// 创建chain
	ch := NewChain()
	// 创建链接
	ch.CreateRemoteClient("https://ethereum-goerli.publicnode.com")
	type fields struct {
		From      common.Address
		To        common.Address
		Data      []byte
		Value     *big.Int
		GasPrice  *big.Int
		GasFeeCap *big.Int
		GasTipCap *big.Int
		BaseFee   *big.Int
		Nonce     *big.Int
		GasLimit  uint64
		chain     *Chain
		ctx       context.Context
	}
	type args struct {
		privateKeyCDSA *ecdsa.PrivateKey
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *bind.TransactOpts
	}{
		{
			name: "TestTransaction_ToTransactOpts",
			fields: fields{
				From:      from,
				To:        to,
				Data:      nil,
				Value:     big.NewInt(100),
				GasPrice:  big.NewInt(100),
				GasFeeCap: big.NewInt(100),
				GasTipCap: big.NewInt(100),
				BaseFee:   big.NewInt(100),
				Nonce:     big.NewInt(100),
				GasLimit:  100,
				chain:     ch,
				ctx:       context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Transaction{
				From:      tt.fields.From,
				To:        tt.fields.To,
				Data:      tt.fields.Data,
				Value:     tt.fields.Value,
				GasPrice:  tt.fields.GasPrice,
				GasFeeCap: tt.fields.GasFeeCap,
				GasTipCap: tt.fields.GasTipCap,
				BaseFee:   tt.fields.BaseFee,
				Nonce:     tt.fields.Nonce,
				GasLimit:  tt.fields.GasLimit,
				chain:     tt.fields.chain,
				ctx:       tt.fields.ctx,
			}
			got := tx.ToTransactOpts(tt.args.privateKeyCDSA)
			t.Log("TransactOpts:", got)
		})
	}
}

func TestTransaction_ensureGasPrice(t *testing.T) {
	// 创建chain
	ch := NewChain()
	// 创建链接
	ch.CreateRemoteClient("https://ethereum-goerli.publicnode.com")
	type fields struct {
		From      common.Address
		To        common.Address
		Data      []byte
		Value     *big.Int
		GasPrice  *big.Int
		GasFeeCap *big.Int
		GasTipCap *big.Int
		BaseFee   *big.Int
		Nonce     *big.Int
		GasLimit  uint64
		chain     *Chain
		ctx       context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestTransaction_ensureGasPrice",
			fields: fields{
				From:      from,
				To:        to,
				Data:      nil,
				Value:     big.NewInt(100),
				GasPrice:  big.NewInt(100),
				GasFeeCap: big.NewInt(100),
				GasTipCap: big.NewInt(100),
				BaseFee:   big.NewInt(100),
				Nonce:     big.NewInt(100),
				GasLimit:  2100000,
				chain:     ch,
				ctx:       context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Transaction{
				From:      tt.fields.From,
				To:        tt.fields.To,
				Data:      tt.fields.Data,
				Value:     tt.fields.Value,
				GasPrice:  tt.fields.GasPrice,
				GasFeeCap: tt.fields.GasFeeCap,
				GasTipCap: tt.fields.GasTipCap,
				BaseFee:   tt.fields.BaseFee,
				Nonce:     tt.fields.Nonce,
				GasLimit:  tt.fields.GasLimit,
				chain:     tt.fields.chain,
				ctx:       tt.fields.ctx,
			}
			if err := tx.ensureGasPrice(); (err != nil) != tt.wantErr {
				t.Errorf("ensureGasPrice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransaction_getNonce(t *testing.T) {
	// 创建chain
	ch := NewChain()
	// 创建链接
	ch.CreateRemoteClient("https://ethereum-goerli.publicnode.com")
	type fields struct {
		From      common.Address
		To        common.Address
		Data      []byte
		Value     *big.Int
		GasPrice  *big.Int
		GasFeeCap *big.Int
		GasTipCap *big.Int
		BaseFee   *big.Int
		Nonce     *big.Int
		GasLimit  uint64
		chain     *Chain
		ctx       context.Context
	}
	type args struct {
		address common.Address
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *big.Int
		wantErr bool
	}{
		{
			name: "TestTransaction_getNonce",
			fields: fields{
				From:      from,
				To:        to,
				Data:      nil,
				Value:     big.NewInt(100),
				GasPrice:  big.NewInt(100),
				GasFeeCap: big.NewInt(100),
				GasTipCap: big.NewInt(100),
				BaseFee:   big.NewInt(100),
				Nonce:     big.NewInt(100),
				GasLimit:  2100000,
				chain:     ch,
				ctx:       context.Background(),
			},
			args: args{
				address: from,
			},
			want:    big.NewInt(1),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Transaction{
				From:      tt.fields.From,
				To:        tt.fields.To,
				Data:      tt.fields.Data,
				Value:     tt.fields.Value,
				GasPrice:  tt.fields.GasPrice,
				GasFeeCap: tt.fields.GasFeeCap,
				GasTipCap: tt.fields.GasTipCap,
				BaseFee:   tt.fields.BaseFee,
				Nonce:     tt.fields.Nonce,
				GasLimit:  tt.fields.GasLimit,
				chain:     tt.fields.chain,
				ctx:       tt.fields.ctx,
			}
			got, err := tx.getNonce(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNonce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log("Nonce:", got)
		})
	}
}

//func TestTransaction_GetTransactionInfo(t *testing.T) {
//	from := common.HexToAddress("0x61AA16822d6CCC7DBD6cA8d28272773bE7AB08aD")
//	to := common.HexToAddress("0x61AA16822d6CCC7DBD6cA8d28272773bE7AB08aD")
//	value := big.NewInt(1000)
//	chain, err := NewChain().CreateRemoteClient("https://ethereum-goerli.publicnode.com", 1000*60)
//	if err != nil {
//		fmt.Print("构建链对象失败")
//		return
//	}
//	transaction := NewTransaction(from, to, value, chain)
//	type fields struct {
//	}
//	type args struct {
//		hash string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{
//			name:   "TestTransaction_GetTransactionInfo",
//			fields: fields{},
//			args: args{
//				hash: "0xcb8fa953786eb8af07543f85fd169ec97d704e870841d61f06d22368c8ead432",
//			}, wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotDetail, err := transaction.GetTransactionInfo(tt.args.hash)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetTransactionInfo() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if gotDetail == nil {
//				t.Errorf("明细为空")
//				return
//			}
//			fmt.Println("明细", gotDetail)
//		})
//	}
//}

func TestTransaction_SignTx(t *testing.T) {
	testcase := CaseAccountGoerli()
	// 创建chain
	ch := NewChain()
	// 创建链接
	ch.CreateRemoteClient("https://ethereum-goerli.publicnode.com")
	tx := NewTransaction(from, to, big.NewInt(100), ch)
	tx.BuildTransfer()
	key, err := NewAccountWithPrivateKey(testcase.privateKeyHex)
	if err != nil {
		return
	}
	type args struct {
		privateKeyCDSA *ecdsa.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Transaction
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "TestTransaction_SignTx",
			args: args{privateKeyCDSA: key.privateKeyECDSA},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tx.SignTx(tt.args.privateKeyCDSA)
			if err != nil {
				t.Error("signTx error", err)
				panic(err)
			}
			t.Log("Transaction_SignTx:", got)
			t.Log("hash:", got.Hash().String())
		})
	}
}

func TestTransaction_SignRawTx(t *testing.T) {
	testcase := CaseAccountGoerli()
	// 创建chain
	ch := NewChain()
	// 创建链接
	ch.CreateRemoteClient("https://ethereum-goerli.publicnode.com")
	tx := NewTransaction(from, to, big.NewInt(100), ch)
	tx.BuildTransfer()
	key, err := NewAccountWithPrivateKey(testcase.privateKeyHex)
	if err != nil {
		panic(err)
	}
	token := NewToken(ch)
	value, _ := utils.ParseAmount("0.00001", token.Info.Decimal)

	transfer, err := token.Transfer(*key, to, value.BigInt())
	if err != nil {
		panic(err)
	}
	type args struct {
		privateKeyCDSA *ecdsa.PrivateKey
		t              *types.Transaction
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Transaction
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "TestTransaction_SignRawTx",
			args: args{
				privateKeyCDSA: key.privateKeyECDSA,
				t:              transfer,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tx.SignRawTx(tt.args.privateKeyCDSA, tt.args.t)
			if err != nil {
				t.Error("signRawTx error", err)
				panic(err)
			}
			t.Log("Transaction_SignRawTx:", got)
		})
	}
}
