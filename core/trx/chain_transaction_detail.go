package trx

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	tron_comm "github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/base"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"math/big"
	"strconv"
)

func (c *Chain) FetchTransactionDetail(hash string) (*base.TransactionDetail, error) {
	tx, err := c.client.RPCClient().GetTransactionByID(hash)
	if err != nil {
		return nil, log.WithError(err, "GetTransactionByID failed")
	}
	info, err := c.client.RPCClient().GetTransactionInfoByID(hash)
	if err != nil {
		return nil, log.WithError(err, "GetTransactionInfoByID failed")
	}
	detail := &base.TransactionDetail{
		GasFee:          big.NewInt(info.Fee),
		GasUsed:         uint64(info.Receipt.EnergyUsageTotal),
		Hash:            common.Bytes2Hex(info.Id),
		ContractAddress: tron_comm.EncodeCheck(info.ContractAddress),
		BlockNumber:     big.NewInt(info.BlockNumber),
		Time:            info.BlockTimeStamp,
	}
	//交易状态处理
	switch tx.Ret[0].ContractRet.Number() {
	case 0: //交易中
		detail.Status = base.TransactionStatusPending
		break
	case 1: //交易成功
		detail.Status = base.TransactionStatusSuccess
	default: //交易失败
		detail.Status = base.TransactionStatusFailure
	}
	switch tx.RawData.Contract[0].Type {
	case core.Transaction_Contract_TriggerSmartContract: //USDT
		tsc := core.TriggerSmartContract{}
		err := tx.RawData.Contract[0].Parameter.UnmarshalTo(&tsc)
		if err != nil {
			return nil, log.WithError(err, "UnmarshalTo failed")
		}
		detail.Form = tron_comm.EncodeCheck(tsc.OwnerAddress)
		data := tsc.Data[8:]
		To := data[8:28]
		detail.To = tron_comm.EncodeCheck(To)
		Amount := hex.EncodeToString(data[32:])
		Value, err := strconv.ParseUint(Amount, 16, len(Amount))
		if err != nil {
			return nil, log.WithError(err, "ParseUint failed")
		}
		detail.Value = big.NewInt(int64(Value))
		break
	case core.Transaction_Contract_TransferContract: //主币
		tsc := core.TransferContract{}
		err := tx.RawData.Contract[0].Parameter.UnmarshalTo(&tsc)
		if err != nil {
			return nil, log.WithError(err, "UnmarshalTo failed")
		}
		detail.Value = big.NewInt(tsc.Amount)
		detail.To = tron_comm.EncodeCheck(tsc.ToAddress)
		detail.Form = tron_comm.EncodeCheck(tsc.OwnerAddress)
		break

	}
	return detail, nil
}

func (c *Chain) FetchTransactionStatus(hash string) (base.TransactionStatus, error) {
	tx, err := c.client.RPCClient().GetTransactionByID(hash)
	if err != nil {
		return base.TransactionStatusNone, log.WithError(err, "GetTransactionByID failed")
	}
	switch tx.Ret[0].ContractRet.Number() {
	case 0: //交易中
		return base.TransactionStatusPending, nil
	case 1: //交易成功
		return base.TransactionStatusSuccess, nil
	default: //交易失败
		return base.TransactionStatusFailure, nil
	}
}
