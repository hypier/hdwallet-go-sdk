package keystore

import (
	"hypier.fun/hdwallet/hdwallet-go-sdk/config"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"reflect"
	"testing"
)

func init() {
	config.InitConfig(&config.BaseConfig{
		BaseDir:    "..",
		LogSwitch:  "CONSOLE_FILE",
		Platform:   "ALL",
		DeviceType: "UNKNOWN",
	})
}

var (
	WalletCase1 = WalletCase{
		Mnemonic:    "frequent retreat lawn demise daughter syrup bid lonely decade dwarf wild describe",
		Password:    "test",
		NewPassword: "test1",
		Name:        "eth wallet 1",
		NewName:     "eth wallet 2",
		Address:     "0xEa2a9Ce354F787791597dF0686Ee2EB9716AE293",
		Code:        "815b5d2895644219b0a4ae4b2cf08cfd",
		PrivateKey:  "340ccf973af62ccfa48622b947f89f57793630fe34919ec1a06be7b02273fc46",
	}
)

type WalletCase struct {
	Mnemonic    string
	Password    string
	NewPassword string
	Name        string
	NewName     string
	Address     string
	Code        string
	PrivateKey  string
}

func Test_keyStorePassphrase_StoreKey(t *testing.T) {
	key, err := NewHDKeyWithMnemonic(WalletCase1.Code, "wedding vessel humble pupil gadget fee rotate bomb camp coconut detect wrist")

	if err != nil {
		return
	}

	type fields struct {
		keysDirPath             string
		scryptN                 int
		scryptP                 int
		skipKeyFileVerification bool
	}
	type args struct {
		key  *Key
		auth string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "store key",
			fields: fields{
				keysDirPath:             config.Base.KeyStorePath,
				scryptN:                 LightScryptN,
				scryptP:                 LightScryptP,
				skipKeyFileVerification: false,
			},
			args: args{
				key:  key,
				auth: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ks := StorePassphrase{
				keysDirPath:             tt.fields.keysDirPath,
				scryptN:                 tt.fields.scryptN,
				scryptP:                 tt.fields.scryptP,
				skipKeyFileVerification: tt.fields.skipKeyFileVerification,
			}
			if err := ks.StoreKey(tt.args.key, tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("StoreKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_keyStorePassphrase_StoreKey_PrivateKey(t *testing.T) {
	key, err := NewHDKeyWithPrivateKey(utils.ETH, 1, "0x9F53dd23f9Ff55E6b7a53a03993D843Bfb2f0F5d", "0620786c1a7f2eb938027f3f2e3edd7183ef27ab84231139b2314f36c3a5cf41")
	if err != nil {
		return
	}

	type fields struct {
		keysDirPath             string
		scryptN                 int
		scryptP                 int
		skipKeyFileVerification bool
	}
	type args struct {
		key  *Key
		auth string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "store key",
			fields: fields{
				keysDirPath:             config.Base.KeyStorePath,
				scryptN:                 LightScryptN,
				scryptP:                 LightScryptP,
				skipKeyFileVerification: false,
			},
			args: args{
				key:  key,
				auth: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ks := StorePassphrase{
				keysDirPath:             tt.fields.keysDirPath,
				scryptN:                 tt.fields.scryptN,
				scryptP:                 tt.fields.scryptP,
				skipKeyFileVerification: tt.fields.skipKeyFileVerification,
			}
			if err := ks.StoreKey(tt.args.key, tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("StoreKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_keyStorePassphrase_GetKey(t *testing.T) {
	key, err := NewHDKeyWithMnemonic(WalletCase1.Code, "wedding vessel humble pupil gadget fee rotate bomb camp coconut detect wrist")
	key.Time = 1692584324
	if err != nil {
		return
	}

	type fields struct {
		keysDirPath             string
		scryptN                 int
		scryptP                 int
		skipKeyFileVerification bool
	}
	type args struct {
		key  *Key
		auth string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Key
		wantErr bool
	}{
		{
			name: "get key",
			fields: fields{
				keysDirPath:             config.Base.KeyStorePath,
				scryptN:                 LightScryptN,
				scryptP:                 LightScryptP,
				skipKeyFileVerification: false,
			},
			args: args{
				key: &Key{
					KeyID: WalletCase1.Code,
				},
				auth: "test",
			},
			want:    key,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ks := StorePassphrase{
				keysDirPath:             tt.fields.keysDirPath,
				scryptN:                 tt.fields.scryptN,
				scryptP:                 tt.fields.scryptP,
				skipKeyFileVerification: tt.fields.skipKeyFileVerification,
			}
			got, err := ks.GetKey(tt.args.key, tt.args.auth)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Mnemonic, tt.want.Mnemonic) {
				t.Errorf("GetKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}
