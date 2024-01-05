package trx

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	type args struct {
		url    string
		apiKey string
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{
			name: "TestNewClient",
			args: args{
				url:    "",
				apiKey: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.url, tt.args.apiKey)
			if err != nil {
				t.Error(err)
			}

			num, err := got.RPCClient().GetBlockByNum(10000)
			if err != nil {
				return
			}
			t.Log("num:", num)
		})
	}
}
