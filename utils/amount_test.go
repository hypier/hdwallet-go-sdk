package utils

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestOptAmount_String(t *testing.T) {
	type fields struct {
		Value   string
		Decimal int16
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "TestOptAmount_ToAmount_01",
			fields: fields{
				Value:   "123000",
				Decimal: 5,
			},
			want:    "1.23",
			wantErr: false,
		},

		{
			name: "TestOptAmount_ToAmount_02",
			fields: fields{
				Value:   "123000",
				Decimal: 3,
			},
			want:    "123",
			wantErr: false,
		},
		{
			name: "TestOptAmount_ToAmount_03",
			fields: fields{
				Value:   "123",
				Decimal: 4,
			},
			want:    "0.0123",
			wantErr: false,
		},
		{
			name: "TestOptAmount_ToAmount_04",
			fields: fields{
				Value:   "123",
				Decimal: 3,
			},
			want:    "0.123",
			wantErr: false,
		},
		{
			name: "TestOptAmount_ToAmount_05",
			fields: fields{
				Value:   "1230",
				Decimal: 2,
			},
			want:    "12.3",
			wantErr: false,
		},
		{
			name: "TestOptAmount_ToAmount_06",
			fields: fields{
				Value:   "1230",
				Decimal: 1,
			},
			want:    "123",
			wantErr: false,
		},
		{
			name: "TestOptAmount_ToAmount_07",
			fields: fields{
				Value:   "12000000000000000000003",
				Decimal: 22,
			},
			want:    "1.2000000000000000000003",
			wantErr: false,
		},
		{
			name: "TestOptAmount_ToAmount_08",
			fields: fields{
				Value:   "123000",
				Decimal: 0,
			},
			want:    "123000",
			wantErr: false,
		},
		{
			name: "TestOptAmount_ToAmount_09",
			fields: fields{
				Value:   "123000",
				Decimal: 2,
			},
			want:    "1230",
			wantErr: false,
		},
		{
			name: "TestOptAmount_ToAmount_10",
			fields: fields{
				Value:   "100232654000000585230003650000000001",
				Decimal: 18,
			},
			want:    "100232654000000585.230003650000000001",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OptAmount{
				value:   tt.fields.Value,
				decimal: tt.fields.Decimal,
			}
			got := o.AmountString()
			if got != tt.want {
				t.Errorf("ToAmountString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAmount(t *testing.T) {
	type args struct {
		amount  string
		decimal int16
	}
	tests := []struct {
		name    string
		args    args
		want    *OptAmount
		wantErr bool
	}{
		{
			name: "TestParseAmount",
			args: args{
				amount:  "12300036",
				decimal: 2,
			},
			want: &OptAmount{
				value:   "1230003600",
				decimal: 2,
			},
			wantErr: false,
		},
		{
			name: "TestParseAmount",
			args: args{
				amount:  "123000.36",
				decimal: 2,
			},
			want: &OptAmount{
				value:   "12300036",
				decimal: 2,
			},
			wantErr: false,
		},
		{
			name: "TestParseAmount",
			args: args{
				amount:  "12.30005",
				decimal: 3,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "TestParseAmount",
			args: args{
				amount:  "10058.523000365",
				decimal: 18,
			},
			want: &OptAmount{
				value:   "10058523000365000000000",
				decimal: 18,
			},
			wantErr: false,
		},
		{
			name: "TestParseAmount",
			args: args{
				amount:  "10023265400000058.523000365",
				decimal: 18,
			},
			want: &OptAmount{
				value:   "10023265400000058523000365000000000",
				decimal: 18,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAmount(tt.args.amount, tt.args.decimal)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAmount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAmount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOptAmount(t *testing.T) {
	type args struct {
		value   string
		decimal int16
	}
	tests := []struct {
		name string
		args args
		want *OptAmount
	}{
		{
			name: "TestNewOptAmount",
			args: args{
				value:   "123000",
				decimal: 5,
			},
			want: &OptAmount{
				value:   "123000",
				decimal: 5,
			},
		},
		{
			name: "TestNewOptAmount",
			args: args{
				value:   "1230.005",
				decimal: 5,
			},
			want: nil,
		},
		{
			name: "TestNewOptAmount",
			args: args{
				value:   "10023265400000058523000365000000000",
				decimal: 18,
			},
			want: &OptAmount{
				value:   "10023265400000058523000365000000000",
				decimal: 18,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewOptAmount(tt.args.value, tt.args.decimal)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOptAmount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigInt(t *testing.T) {
	type args struct {
		value   string
		decimal int16
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{
			name: "TestNewOptAmount",
			args: args{
				value:   "123000",
				decimal: 5,
			},
			want: big.NewInt(123000),
		},
		{
			name: "TestNewOptAmount",
			args: args{
				value:   "1230.005",
				decimal: 5,
			},
			want: nil,
		},
		{
			name: "TestNewOptAmount",
			args: args{
				value:   "100232654000000585",
				decimal: 18,
			},
			want: big.NewInt(100232654000000585),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := NewOptAmount(tt.args.value, tt.args.decimal)
			got := opt.BigInt()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOptAmount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWei2ethStr(t *testing.T) {
	str := Wei2ethStr(big.NewInt(1001000), 6)
	fmt.Println(str)
}
