package balance

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/pkg/errors"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/base"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/btc"
	"hypier.fun/hdwallet/hdwallet-go-sdk/sdk_struct"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"math"
	"sort"
	"strconv"
)

type OkLinkToken struct {
	btc.Coin
	Info  *base.TokenInfo
	chain *btc.Chain
}

// MainOKLinkUTXO 未花费的账本
type MainOKLinkUTXO struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Page      string `json:"page"`
		Limit     string `json:"limit"`
		TotalPage string `json:"totalPage"`
		UTXOList  []struct {
			TxId          string `json:"txid"`
			UnspentAmount string `json:"unspentAmount"`
			Index         string `json:"index"`
		} `json:"utxoList"`
	} `json:"data"`
}

// MainOKLinkBalance 余额
type MainOKLinkBalance struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Balance string `json:"balance"`
	} `json:"data"`
}

// MainOKLinkFee 旷工费
type MainOKLinkFee struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		BestTransactionFee string `json:"bestTransactionFee"`
	} `json:"data"`
}

func (t *OkLinkToken) GetDecimal() int16 {
	if t.Info.Decimal == 0 {
		t.Info, _ = t.TokenInfo()
	}

	return t.Info.Decimal
}

func (t *OkLinkToken) TokenInfo() (*base.TokenInfo, error) {
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

// GetBtcUnspent 获取未花费的比特币交易输出
func (t *OkLinkToken) GetBtcUnspent(current btcutil.Address, amount uint64) ([]btc.BtcUnspent, error) {
	data := make([]btc.BtcUnspent, 0)
	script, err := txscript.PayToAddrScript(current)
	if err != nil {
		return nil, err
	}
	scriptStr := fmt.Sprintf("%x", script)
	total := uint64(0)
	index := 1
	for total < amount {
		log.Infof("start begin %d", index)
		unSpent, er := getUnSpentTxByMain(current, index)
		if er != nil {
			return nil, er
		}
		if len(unSpent.Data) > 0 {
			list := unSpent.Data[0].UTXOList
			if len(list) == 0 {
				break
			}
			sort.Slice(list, func(i, j int) bool {
				return list[i].UnspentAmount < list[j].UnspentAmount
			})
			for i := range list {
				item := list[i]
				vout, err := strconv.ParseUint(item.Index, 10, 32)
				if err != nil {

					return nil, err
				}
				amountValue, err := strconv.ParseFloat(item.UnspentAmount, 64)
				if err != nil {
					return nil, err
				}
				newAmount, err := btcutil.NewAmount(amountValue)
				if err != nil {
					return nil, err
				}
				unspent := btc.BtcUnspent{TxID: item.TxId, Vout: uint32(vout), ScriptPubKey: scriptStr, Amount: amountValue, Value: uint64(newAmount.ToUnit(btcutil.AmountSatoshi))}
				data = append(data, unspent)
				total += unspent.Value
			}
			index++
		}
	}
	return data, err
}

// getUnSpentTxByMain 获取未花费的比特币交易输出
func getUnSpentTxByMain(current btcutil.Address, page int) (MainOKLinkUTXO, error) {
	url := fmt.Sprintf("https://www.oklink.com/api/v5/explorer/address/utxo?chainShortName=BTC&address=%s&page=%s&limit=20", current.String(), strconv.Itoa(page))
	var data MainOKLinkUTXO
	request, err := utils.DoGetAndHeader(url, 3, map[string]string{"Ok-Access-Key": sdk_struct.GetOkLinkApiKey()})
	if err != nil {
		log.Error(err)
		return data, err
	}
	if err = json.Unmarshal(request, &data); err != nil {
		log.Error(err)
		return data, err
	}
	return data, err
}

// GetBalance 获取余额
func (t *OkLinkToken) GetBalance(address btcutil.Address) (*base.Balance, error) {
	url := fmt.Sprintf("https://www.oklink.com/api/v5/explorer/address/address-summary?chainShortName=BTC&address=%s", address.String())
	request, err := utils.DoGetAndHeader(url, 3, map[string]string{"Ok-Access-Key": sdk_struct.GetOkLinkApiKey()})
	if err != nil {
		return base.EmptyBalance(), err
	}
	var data MainOKLinkBalance
	if err = json.Unmarshal(request, &data); err != nil {
		return base.EmptyBalance(), err
	}
	if len(data.Data) > 0 {
		amount, _ := utils.ParseAmount(data.Data[0].Balance, t.GetDecimal())
		return &base.Balance{
			Total:  amount,
			Usable: amount,
		}, nil
	}
	return base.EmptyBalance(), nil
}

// PushTx 广播交易
func (t *OkLinkToken) PushTx(signedTx string, transaction *btc.Transaction) (string, error) {
	client, err := rpcclient.New(&rpcclient.ConnConfig{
		Host:              "btc.getblock.io/fdc930ba-aa6a-4a45-8ce9-506eb03f2e98/mainnet/",
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
func (t *OkLinkToken) GetGasFee() (uint64, error) {
	var data MainOKLinkFee
	request, err := utils.DoGetAndHeader("https://www.oklink.com/api/v5/explorer/blockchain/fee?chainShortName=BTC",
		3, map[string]string{"Ok-Access-Key": sdk_struct.GetOkLinkApiKey()})

	if err = json.Unmarshal(request, &data); err != nil {
		log.Error(err)
		return 0, err
	}
	if data.Code == "0" {
		if len(data.Data) > 0 {
			fee, err := strconv.ParseUint(data.Data[0].BestTransactionFee, 10, 32)
			if err != nil {
				return 0, err
			}
			log.Info(fee)
			return uint64(math.Ceil(float64(fee) * 1.5)), nil
		} else {
			return 0, errors.New("Data length is 0")
		}
	} else {
		return 0, errors.New(data.Msg)
	}
}
