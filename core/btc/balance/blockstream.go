package balance

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/base"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/btc"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"strconv"
	"strings"
)

type StreamToken struct {
	btc.Coin
	Info  *base.TokenInfo
	chain *btc.Chain
}

// BlockBalance 结构体定义了比特币区块的余额
type BlockBalance struct {
	FinalBalance uint64 `json:"final_balance"`
	Balance      uint64 `json:"balance"`
	ChainStats   struct {
		FundedSum uint64 `json:"funded_txo_sum"`
		SpentSum  uint64 `json:"spent_txo_sum"`
	} `json:"chain_stats"`
}

// BlockUnspent 结构体定义了未花费的块的交易输出
type BlockUnspent struct {
	TxId   string `json:"txid"`
	VOut   uint64 `json:"vout"`
	Value  uint64 `json:"value"`
	Status struct {
		Confirmed bool `json:"confirmed"`
	} `json:"status"`
}

type BlockChair struct {
	Data struct {
		ByteSat uint64 `json:"suggested_transaction_fee_per_byte_sat"`
	} `json:"Data"`
}

func (t *StreamToken) GetDecimal() int16 {
	if t.Info.Decimal == 0 {
		t.Info, _ = t.TokenInfo()
	}

	return t.Info.Decimal
}
func (t *StreamToken) TokenInfo() (*base.TokenInfo, error) {
	token := base.GetToken(t.CoinType(), "")
	if token != nil {
		return token, nil
	}

	t.Info = &base.TokenInfo{
		Name:    "BTC",
		Symbol:  "BTC",
		Decimal: 8,
	}

	base.AddToken(t.CoinType(), "", t.Info)

	return t.Info, nil
}

// GetBalance 获取余额
func (t *StreamToken) GetBalance(address btcutil.Address) (*base.Balance, error) {
	url := fmt.Sprintf("https://blockstream.info/testnet/api/address/%s", address.String())
	bytes, err := utils.DoGet(url, 3)
	if err != nil {
		return base.EmptyBalance(), err
	}
	var data BlockBalance
	if err = json.Unmarshal(bytes, &data); err != nil {
		return base.EmptyBalance(), err
	}
	u := data.ChainStats.FundedSum - data.ChainStats.SpentSum
	if u <= 0 {
		u = 0
		return base.EmptyBalance(), nil
	}
	btcAmount := btcutil.Amount(u).ToBTC()
	amountStr := strconv.FormatFloat(btcAmount, 'f', 8, 64)
	amountStr = strings.TrimRight(amountStr, "0")
	amountStr = strings.TrimRight(amountStr, ".")

	amount, _ := utils.ParseAmount(amountStr, t.GetDecimal())
	return &base.Balance{
		Total:  amount,
		Usable: amount,
	}, nil
}

func (t *StreamToken) GetBtcUnspent(current btcutil.Address, amount uint64) ([]btc.BtcUnspent, error) {
	data := make([]btc.BtcUnspent, 0)
	script, err := txscript.PayToAddrScript(current)
	if err != nil {
		return nil, err
	}
	scriptStr := fmt.Sprintf("%x", script)
	unSpent, err := getUnSpentTx(current)
	if err != nil {
		return nil, err
	}
	for _, item := range unSpent {
		if item.Status.Confirmed == false {
			continue
		}
		unspent := btc.BtcUnspent{TxID: item.TxId, Vout: uint32(item.VOut), ScriptPubKey: scriptStr, Amount: btcutil.Amount(item.Value).ToBTC(), Value: item.Value}
		data = append(data, unspent)
	}
	return data, nil
}

func getUnSpentTx(current btcutil.Address) ([]BlockUnspent, error) {
	url := fmt.Sprintf("https://blockstream.info/testnet/api/address/%s/utxo", current.String())
	var data []BlockUnspent
	request, err := utils.DoGet(url, 3)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if err = json.Unmarshal(request, &data); err != nil {
		log.Error(err)
		return nil, err
	}
	return data, err
}

func (t *StreamToken) PushTx(signedTx string, transaction *btc.Transaction) (string, error) {
	client, err := rpcclient.New(&rpcclient.ConnConfig{
		Host:              "btc.getblock.io/fdc930ba-aa6a-4a45-8ce9-506eb03f2e98/testnet/",
		User:              "u",
		Pass:              "p",
		HTTPPostMode:      true,
		Params:            chaincfg.TestNet3Params.Name,
		DisableTLS:        false,
		EnableBCInfoHacks: true,
	}, nil)

	send, err := transaction.SendRawTransaction(client)
	if err != nil {
		return "", err
	}
	return send.String(), nil
}

// GetGasFee 获取推荐的矿工费
func (t *StreamToken) GetGasFee() (uint64, error) {
	request, err := utils.DoGet("https://api.blockchair.com/bitcoin/testnet/stats", 3)
	if err != nil {
		return 0, err
	}
	var data BlockChair
	if err = json.Unmarshal(request, &data); err != nil {
		log.Error(err)
		return 0, err
	}
	return data.Data.ByteSat * 1000, nil
}
