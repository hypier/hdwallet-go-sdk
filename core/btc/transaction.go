package btc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/wallet/txauthor"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
)

type Transaction struct {
	txauthor.AuthoredTx                  // 交易对象
	chainParams         *chaincfg.Params // 链参数
	feePerKb            int64            // 每千字节手续费
}
type BlockUnspent struct {
	TxId   string `json:"txid"`
	VOut   uint64 `json:"vout"`
	Value  uint64 `json:"value"`
	Status struct {
		Confirmed bool `json:"confirmed"`
	} `json:"status"`
}

func NewTransaction(unspents []BtcUnspent, params []TransferParam, changeAddress btcutil.Address, feePerKb int64, chainParam *chaincfg.Params) (*Transaction, error) {
	// 检查参数是否正确
	if len(unspents) == 0 || changeAddress == nil || feePerKb <= 0 {
		return nil, errors.New("invalid params")
	}
	// 将每千字节手续费转换为金额
	feeRatePerKb := btcutil.Amount(feePerKb)
	// 生成交易输出
	txOuts, err := makeTxOutputs(params, chainParam)
	if err != nil {
		return nil, log.WithError(err, "makeTxOutputs failed")
	}
	// 生成找零脚本
	changeBytes, err := txscript.PayToAddrScript(changeAddress)
	if err != nil {
		return nil, log.WithError(err, "PayToAddrScript failed")
	}

	// 创建用于生成找零脚本的ChangeSource对象
	changeSource := txauthor.ChangeSource{
		NewScript: func() ([]byte, error) {
			return changeBytes, nil
		},
		ScriptSize: len(changeBytes),
	}

	unsignedTx, err := txauthor.NewUnsignedTransaction(txOuts, feeRatePerKb, makeInputSource(unspents), &changeSource)

	if err != nil {
		return nil, log.WithError(err, "NewUnsignedTransaction failed")
	}
	// 如果存在找零输出，则随机化找零位置
	if unsignedTx.ChangeIndex >= 0 {
		unsignedTx.RandomizeChangePosition()
	}
	// 返回创建的BtcTransaction对象
	return &Transaction{*unsignedTx, chainParam, feePerKb}, nil
}

func (t *Transaction) SignWithSecretsSource(account *Account) error {
	err := t.AddAllInputScripts(account)
	if err != nil {
		return err
	}
	err = validateMsgTx(t.Tx, t.PrevScripts, t.PrevInputValues)
	if err != nil {
		return err
	}
	return nil
}
func (t *Transaction) SendRawTransaction(c *rpcclient.Client) (*chainhash.Hash, error) {
	txHex := ""
	tx := t.Tx
	if tx != nil {
		// Serialize the transaction and convert to hex string.
		buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
		if err := tx.Serialize(buf); err != nil {
			return nil, log.WithError(err, "tx serialize failed")
		}
		txHex = hex.EncodeToString(buf.Bytes())
	}
	//cmd := btcjson.NewSendRawTransactionCmd(txHex, &allowHighFees)
	cmd := btcjson.NewBitcoindSendRawTransactionCmd(txHex, btcutil.SatoshiPerBitcent/20)

	var sendCmd rpcclient.FutureSendRawTransactionResult = c.SendCmd(cmd)
	hash, err := sendCmd.Receive()
	if err != nil {
		return nil, log.WithError(err, "send raw transaction failed")
	}
	return hash, nil
}

// validateMsgTx 私有函数用于验证交易输入脚本
func validateMsgTx(tx *wire.MsgTx, prevScripts [][]byte, inputValues []btcutil.Amount) error {
	inputFetcher, err := txauthor.TXPrevOutFetcher(
		tx, prevScripts, inputValues,
	)
	if err != nil {
		return err
	}
	hashCache := txscript.NewTxSigHashes(tx, inputFetcher)
	for i, prevScript := range prevScripts {
		vm, err := txscript.NewEngine(
			prevScript, tx, i, txscript.StandardVerifyFlags, nil,
			hashCache, int64(inputValues[i]), inputFetcher,
		)
		if err != nil {
			return fmt.Errorf("无法创建脚本引擎: %s", err)
		}
		err = vm.Execute()
		if err != nil {
			return fmt.Errorf("无法验证交易: %s", err)
		}
	}
	return nil
}

// makeInputSource 函数用于生成输入源
func makeInputSource(unspents []BtcUnspent) txauthor.InputSource {
	sz := len(unspents)
	currentTotal := btcutil.Amount(0)
	currentInputs := make([]*wire.TxIn, 0, sz)
	currentInputValues := make([]btcutil.Amount, 0, sz)
	currentScripts := make([][]byte, 0, sz)
	return func(target btcutil.Amount) (btcutil.Amount, []*wire.TxIn, []btcutil.Amount, [][]byte, error) {
		for currentTotal < target && len(unspents) != 0 {
			u := unspents[0]
			unspents = unspents[1:]
			hash, _ := chainhash.NewHashFromStr(u.TxID)
			nextInput := wire.NewTxIn(&wire.OutPoint{
				Hash:  *hash,
				Index: u.Vout,
			}, nil, nil)
			amount, _ := btcutil.NewAmount(u.Amount)
			s, _ := hex.DecodeString(u.ScriptPubKey)
			nextInput.Sequence = uint32(amount)
			currentTotal += amount
			currentInputs = append(currentInputs, nextInput)
			currentInputValues = append(currentInputValues, amount)
			currentScripts = append(currentScripts, s)
		}
		return currentTotal, currentInputs, currentInputValues, currentScripts, nil
	}
}

// makeTxOutputs 函数用于生成交易输出
func makeTxOutputs(params []TransferParam, chainCfg *chaincfg.Params) ([]*wire.TxOut, error) {
	paramLen := len(params)
	if paramLen == 0 {
		return nil, log.WithError(errors.New("output is empty"))
	}
	txOuts := make([]*wire.TxOut, 0, paramLen)

	for i := 0; i < paramLen; i++ {
		param := &params[i]
		if !param.To.IsForNet(chainCfg) {
			return nil, log.WithError(errors.New("invalid address"))
		}
		// 创建一个新的脚本，支付给提供的地址
		pkScript, err := txscript.PayToAddrScript(param.To)
		if err != nil {
			return nil, err
		}
		txOut := &wire.TxOut{
			Value:    param.Amount,
			PkScript: pkScript,
		}
		txOuts = append(txOuts, txOut)
	}
	return txOuts, nil
}
func GetBtcUnspent(current btcutil.Address) ([]BtcUnspent, error) {
	data := make([]BtcUnspent, 0)
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
		unspent := BtcUnspent{TxID: item.TxId, Vout: uint32(item.VOut), ScriptPubKey: scriptStr, Amount: btcutil.Amount(item.Value).ToBTC(), Value: item.Value}
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
