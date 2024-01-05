package btc

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcwallet/wallet/txauthor"
	"github.com/btcsuite/btcwallet/wallet/txsizes"
	"hypier.fun/hdwallet/hdwallet-go-sdk/config"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/base"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"strconv"
	"strings"
)

type Token struct {
	Coin
	Info  *base.TokenInfo
	chain *Chain
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

type NetParams interface {
	// GetBtcUnspent 获取没有花费的账本
	GetBtcUnspent(address btcutil.Address, amount uint64) ([]BtcUnspent, error)
	// GetBalance 获取余额
	GetBalance(address btcutil.Address) (*base.Balance, error)
	// PushTx 广播交易
	PushTx(signedTx string, transaction *Transaction) (string, error)
	// GetGasFee 获取推荐的矿工费
	GetGasFee() (uint64, error)
}

func (t *Token) BalanceOfAddress(address string) (*base.Balance, error) {
	// 默认使用blockstream实现
	url := fmt.Sprintf("https://blockstream.info/testnet/api/address/%s", address)
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

	return nil, nil
}
func NewToken(chain *Chain) *Token {
	return &Token{chain: chain, Info: &base.TokenInfo{}}
}

func (t *Token) Chain() base.Chain {
	return t.chain
}

func (t *Token) GetDecimal() int16 {
	if t.Info.Decimal == 0 {
		t.Info, _ = t.TokenInfo()
	}

	return t.Info.Decimal
}

func (t *Token) TokenInfo() (*base.TokenInfo, error) {
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

type FeeRate struct {
	Low     int64
	Average int64
	High    int64
}

// EstimateGas 获取当前燃气价格
func EstimateGas() (*FeeRate, error) {
	url := "https://mempool-mainnet.coming.chat/api/v1/fees/recommended"

	bs, err := utils.DoGet(url, 3)
	if err != nil {
		return nil, log.WithError(err, "EstimateGas failed")
	}
	respDict := make(map[string]interface{})
	err = json.Unmarshal(bs, &respDict)
	if err != nil {
		return nil, err
	}

	var low, avg, high float64
	var ok bool
	if low, ok = respDict["minimumFee"].(float64); !ok {
		low = 1
	}
	if avg, ok = respDict["halfHourFee"].(float64); !ok {
		avg = low
	}
	if high, ok = respDict["fastestFee"].(float64); !ok {
		high = avg
	}
	return &FeeRate{
		Low:     int64(low),
		Average: int64(avg),
		High:    int64(high),
	}, nil
}

// BtcUnspent 结构体定义了未花费的比特币交易输出
type BtcUnspent struct {
	TxID         string  `json:"txid"`                   // 交易ID
	Vout         uint32  `json:"vout"`                   // 输出索引
	ScriptPubKey string  `json:"scriptPubKey"`           // 输出脚本
	RedeemScript string  `json:"redeemScript,omitempty"` // 兑换脚本，可选
	Amount       float64 `json:"amount"`                 // 输出金额
	Value        uint64  `json:"value"`
}

// TransferParam 转账目标,因为可以一次转多个 所以定义一个结构体来封装
type TransferParam struct {
	To     btcutil.Address `json:"to"`     // 转账目标
	Amount int64           `json:"amount"` // 输出金额
}

// EstimateGasLimit 估计消耗费用
func (t *Token) EstimateGasLimit(from btcutil.Address, params []TransferParam, chainId int) (int, error) {
	// 生成找零脚本
	changeBytes, err := txscript.PayToAddrScript(from)
	if err != nil {
		return 0, log.WithError(err, "PayToAddrScript failed")
	}
	// 创建用于生成找零脚本的ChangeSource对象
	changeSource := txauthor.ChangeSource{
		NewScript: func() ([]byte, error) {
			return changeBytes, nil
		},
		ScriptSize: len(changeBytes),
	}
	chainCfg, err := utils.GetBtcChainParams(chainId)
	if err != nil {
		return 0, log.WithError(err, "ChainID failed")
	}
	txOuts, err := makeTxOutputs(params, chainCfg)
	if err != nil {
		return 0, log.WithError(err, "makeTxOutputs failed")
	}
	estimatedSize := txsizes.EstimateVirtualSize(
		0, 0, 1, 0, txOuts, changeSource.ScriptSize,
	)
	return estimatedSize, err
}
func (t *Token) Transfer(from *Account, to string, value int64) (string, error) {
	// 拿配置
	chainCfg, err := utils.GetBtcChainParam(config.Base.BtcParam)
	if err != nil {
		return "", log.WithError(err, "ChainID failed")
	}
	toAddr, err := btcutil.DecodeAddress(to, chainCfg)
	if err != nil {
		return "", log.WithError(err, "DecodeAddress failed")
	}
	// 构建交易
	var outputs = []TransferParam{{
		To:     toAddr,
		Amount: value,
	}}
	nsa, err := from.NestedSegwitAddress()
	if err != nil {
		return "", log.WithError(err, "TaprootAddress failed")
	}
	fromAddr, err := btcutil.DecodeAddress(nsa, chainCfg)
	if err != nil {
		return "", log.WithError(err, "DecodeAddress failed")
	}
	allBtcUnspent, err := GetBtcUnspent(fromAddr)
	if err != nil {
		return "", log.WithError(err, "GetBtcUnspent failed")
	}
	var btcUnspent = make([]BtcUnspent, 0)
	//拼接
	for i := range allBtcUnspent {
		btcUnspent = append(btcUnspent, allBtcUnspent[i])
	}
	tx, err := NewTransaction(btcUnspent, outputs, fromAddr, 1000, chainCfg)
	if err != nil {
		return "", log.WithError(err, "NewTransaction failed")
	}
	err = tx.SignWithSecretsSource(from)
	hash, err := tx.SendRawTransaction(t.chain.client.rpcClient)
	if err != nil {
		return "", log.WithError(err, "SendRawTransaction failed")
	}
	return hash.String(), nil
}
