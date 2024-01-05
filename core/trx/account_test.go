package trx

import (
	"crypto/ecdsa"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"hypier.fun/hdwallet/hdwallet-go-sdk/config"
	"testing"
)

var (
	mnemonic = "dwarf unaware dragon car curve stage include output picture organ skin talk"
	private  = "a19a3f51386f2325312cd4c4a81b27aef389e3376efcf8440f681f480b4020ff"
	addr     = "TU7YLwySaoPJCDuhF2tZcsxJVDX8umJuXe"
)

func init() {
	config.InitConfig(&config.BaseConfig{
		BaseDir:    "..",
		LogSwitch:  "CONSOLE_FILE",
		Platform:   "ALL",
		DeviceType: "UNKNOWN",
	})
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
				mnemonic: mnemonic,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccount(tt.args.mnemonic)
			if err != nil {
				t.Error(err)
				panic(err)
			}
			t.Log("privateKey:", got.PrivateKeyHex())
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
			name: "Test NewAccountWithPrivateKey",
			args: args{
				privateKey: private,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountWithPrivateKey(tt.args.privateKey)
			if err != nil {
				t.Error(err)
				panic(err)
			}
			t.Log("address:", got.Address())
		})
	}
}

func TestAccount_Address(t *testing.T) {
	got, _ := NewAccountWithPrivateKey(private)
	type fields struct {
		Coin            Coin
		privateKeyECDSA *ecdsa.PrivateKey
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test Account Address",
			fields: fields{
				Coin:            got.Coin,
				privateKeyECDSA: got.privateKeyECDSA,
			},
			want: addr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				Coin:            tt.fields.Coin,
				privateKeyECDSA: tt.fields.privateKeyECDSA,
			}
			if got := a.Address(); got != tt.want {
				t.Errorf("Address() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_PrivateKey(t *testing.T) {
	got, _ := NewAccountWithPrivateKey(private)
	type fields struct {
		Coin            Coin
		privateKeyECDSA *ecdsa.PrivateKey
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Test Account PrivateKey",
			fields: fields{
				Coin:            got.Coin,
				privateKeyECDSA: got.privateKeyECDSA,
			},
			want: []byte(private),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				Coin:            tt.fields.Coin,
				privateKeyECDSA: tt.fields.privateKeyECDSA,
			}
			got := a.PrivateKey()
			t.Log("privateKey:", common.Bytes2Hex(got))
		})
	}
}

func TestAccount_PrivateKeyHex(t *testing.T) {
	got, _ := NewAccountWithPrivateKey(private)
	type fields struct {
		Coin            Coin
		privateKeyECDSA *ecdsa.PrivateKey
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test Account PrivateKeyHex",
			fields: fields{
				Coin:            got.Coin,
				privateKeyECDSA: got.privateKeyECDSA,
			},
			want: private,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				Coin:            tt.fields.Coin,
				privateKeyECDSA: tt.fields.privateKeyECDSA,
			}
			if got := a.PrivateKeyHex(); got != tt.want {
				t.Errorf("PrivateKeyHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_PublicKey(t *testing.T) {
	got, _ := NewAccountWithPrivateKey(private)
	type fields struct {
		Coin            Coin
		privateKeyECDSA *ecdsa.PrivateKey
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "Test Account PublicKey",
			fields: fields{
				Coin:            got.Coin,
				privateKeyECDSA: got.privateKeyECDSA,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				Coin:            tt.fields.Coin,
				privateKeyECDSA: tt.fields.privateKeyECDSA,
			}
			key := a.PublicKey()
			t.Log("publicKey:", common.Bytes2Hex(key))
		})
	}
}

func TestAccount_PublicKeyHex(t *testing.T) {
	got, _ := NewAccountWithPrivateKey(private)
	type fields struct {
		Coin            Coin
		privateKeyECDSA *ecdsa.PrivateKey
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test Account PublicKeyHex",
			fields: fields{
				Coin:            got.Coin,
				privateKeyECDSA: got.privateKeyECDSA,
			},
			want: "04a271d9a3257b1bd570db65e77668a39eb3c164ed8e7860ce422728f514e04a02c9106a983f2c800e66dc560039b27c41311e524bb51d1d727d9787ae30d5a631",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				Coin:            tt.fields.Coin,
				privateKeyECDSA: tt.fields.privateKeyECDSA,
			}
			if got := a.PublicKeyHex(); got != tt.want {
				t.Errorf("PublicKeyHex() = %v, want %v", got, tt.want)
			}
		})
	}
}
