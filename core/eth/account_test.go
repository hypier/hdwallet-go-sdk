package eth

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestAccount_Address(t *testing.T) {
	type fields struct {
		privateKeyECDSA *ecdsa.PrivateKey
	}
	tests := []struct {
		name   string
		fields fields
		want   common.Address
	}{
		{
			name: CaseAccountGoerliNew().caseName,
			fields: fields{
				privateKeyECDSA: CaseAccountGoerliNew().privateKey,
			},
			want: CaseAccountGoerliNew().from,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				privateKeyECDSA: tt.fields.privateKeyECDSA,
			}
			if got := a.Address(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Address() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_PrivateKey(t *testing.T) {
	type fields struct {
		privateKeyECDSA *ecdsa.PrivateKey
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: CaseAccountGoerliNew().caseName,
			fields: fields{
				privateKeyECDSA: CaseAccountGoerliNew().privateKey,
			},
			want: common.Hex2Bytes(CaseAccountGoerliNew().privateKeyHex),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				privateKeyECDSA: tt.fields.privateKeyECDSA,
			}
			got := a.PrivateKey()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrivateKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_PublicKey(t *testing.T) {
	type fields struct {
		privateKeyECDSA *ecdsa.PrivateKey
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: CaseAccountGoerliNew().caseName,
			fields: fields{
				privateKeyECDSA: CaseAccountGoerliNew().privateKey,
			},
			want: crypto.FromECDSAPub(&CaseAccountGoerliNew().privateKey.PublicKey),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				privateKeyECDSA: tt.fields.privateKeyECDSA,
			}
			if got := a.PublicKey(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PublicKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_PublicKeyHex(t *testing.T) {
	type fields struct {
		privateKeyECDSA *ecdsa.PrivateKey
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: CaseAccountGoerliNew().caseName,
			fields: fields{
				privateKeyECDSA: CaseAccountGoerliNew().privateKey,
			},
			want: hex.EncodeToString(crypto.FromECDSAPub(&CaseAccountGoerliNew().privateKey.PublicKey)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				privateKeyECDSA: tt.fields.privateKeyECDSA,
			}
			if got := a.PublicKeyHex(); got != tt.want {
				t.Errorf("PublicKeyHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAccount(t *testing.T) {
	type args struct {
		mnemonic string
	}
	tests := []struct {
		name    string
		args    args
		want    *Account
		wantErr bool
	}{
		{
			name: "Test NewAccount with valid mnemonic",
			args: args{
				mnemonic: CaseAccountGoerliNew().mnemonic,
			},
			want: &Account{
				privateKeyECDSA: CaseAccountGoerliNew().privateKey,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccount(tt.args.mnemonic)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got.PrivateKeyHex(), tt.want.PrivateKeyHex()) {
				t.Errorf("NewAccount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAccountWithPrivateKey(t *testing.T) {
	type args struct {
		privateKey string
	}
	tests := []struct {
		name    string
		args    args
		want    *Account
		wantErr bool
	}{
		{
			name: "Test NewAccountWithPrivateKey with valid private key",
			args: args{
				privateKey: CaseAccountGoerliNew().privateKeyHex,
			},
			want: &Account{
				privateKeyECDSA: CaseAccountGoerliNew().privateKey,
			},
			wantErr: false,
		},
		{
			name: "Test NewAccountWithPrivateKey with invalid private key",
			args: args{
				privateKey: "invalid private key",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountWithPrivateKey(tt.args.privateKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAccountWithPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAccountWithPrivateKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_PrivateKeyHex(t *testing.T) {
	type fields struct {
		privateKeyECDSA *ecdsa.PrivateKey
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: CaseAccountGoerliNew().caseName,
			fields: fields{
				privateKeyECDSA: CaseAccountGoerliNew().privateKey,
			},
			want: CaseAccountGoerliNew().privateKeyHex,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				privateKeyECDSA: tt.fields.privateKeyECDSA,
			}
			if got := a.PrivateKeyHex(); got != tt.want {
				t.Errorf("PrivateKeyHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_VerifySign(t *testing.T) {
	type fields struct {
		Coin            Coin
		privateKeyECDSA *ecdsa.PrivateKey
	}
	type args struct {
		message      string
		signatureHex string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "TestAccount_VerifySign",
			args: args{
				message:      "hello world",
				signatureHex: "0x56b40874007fb6c84bdcd82f7fe9aec51b7519194407da91a1325f6505d9d41e02c562c6c1c9cbbc5dd272ae5f373e78004f3cca29f1e640bc1bbea8106e30291c",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				privateKeyECDSA: CaseAccountGoerliNew().privateKey,
			}
			got, err := a.VerifySign(tt.args.signatureHex, []byte(tt.args.message))

			if err != nil {
				t.Error("verify sign error,", err)
				return
			}
			t.Log("verify sign result,", got)
		})
	}
}

func TestAccount_Sign(t *testing.T) {

	type fields struct {
		Coin            Coin
		privateKeyECDSA *ecdsa.PrivateKey
	}
	type args struct {
		input []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{ //47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad
			name:    "TestAccount_Sign",
			wantErr: assert.NoError,
			args: args{
				input: []byte("hello world"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				privateKeyECDSA: CaseAccountGoerliNew().privateKey,
			}
			got, err := a.Sign(tt.args.input)
			if err != nil {
				t.Error("sign error,", err)
				return
			}
			t.Log("sign success,", got)
		})
	}
}
