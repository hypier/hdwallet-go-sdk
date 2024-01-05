package btc

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/base"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
)

func (c *Chain) FetchTransactionDetail(txHash string) (*base.TransactionDetail, error) {
	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return nil, log.WithError(err, "NewHashFromStr failed")
	}

	rawResult, err := c.client.rpcClient.GetRawTransactionVerbose(hash)
	if err != nil {
		return nil, log.WithError(err, "GetRawTransactionVerbose failed")
	}
	//	Hex           string `json:"hex"`
	//	Txid          string `json:"txid"`
	//	Hash          string `json:"hash,omitempty"`
	//	Size          int32  `json:"size,omitempty"`
	//	Vsize         int32  `json:"vsize,omitempty"`
	//	Weight        int32  `json:"weight,omitempty"`
	//	Version       uint32 `json:"version"`
	//	LockTime      uint32 `json:"locktime"`
	//	Vin           []Vin  `json:"vin"`
	//	Vout          []Vout `json:"vout"`
	//	BlockHash     string `json:"blockhash,omitempty"`
	//	Confirmations uint64 `json:"confirmations,omitempty"`
	//	Time          int64  `json:"time,omitempty"`
	//	Blocktime     int64  `json:"blocktime,omitempty"`
	//&{
	//Hex 02000000000101847b31f2d6bcf706b7f747aa22bf35d7e16a96ab8b6e2c22a5f271e9e47486340000000017160014b48208851445d7cd8753e32d25bd9f5bdda43f69feffffff03f4a487000000000017a914724d60b094581aec5756c33ec54ca9eb1850fd1b870000000000000000166a146f6d6e69000000008000050d000000001dcd65001c0200000000000017a91474d9b52e9e2457a737e233d336577b2c8b52188a8702473044022027cec650607f72f36a138f0d799dff0fa1b428b701583c4732a0fa3239e7c94802201051a5bae6b3d4bdcceaf4666b8b3fedb96ab27a1351b5eb77470ab04d294aed01210277e3a69e37416feb1896ce38b992a90a15db5e0f07f4bfd576436a0586427b4dff352600
	//Hash 604520d6133dbacc55d15ea76d42797e88a0cc384153d3eb6524da90dbcc33f6
	//Txid b1db4a20e9c35e8ef64232d99c5ee941461ef1748a62fd1360364d5e2df3c416
	//278		Size 大小
	//197		Vsize 交易费用?
	//785	  Weight 重量
	//2	  	  Version 版本
	//2504191 LockTime 锁定时间
	//[
	//	{
	//		348674e4e971f2a5222c6e8bab966ae1d735bf22aa47f7b706f7bcd6f2317b84
	//		0
	//		0x14000076140
	//		4294967294
	//		[
	//			3044022027cec650607f72f36a138f0d799dff0fa1b428b701583c4732a0fa3239e7c94802201051a5bae6b3d4bdcceaf4666b8b3fedb96ab27a1351b5eb77470ab04d294aed01
	//			0277e3a69e37416feb1896ce38b992a90a15db5e0f07f4bfd576436a0586427b4d
	//		]
	//	}
	//]
	//[
	//	{
	//		0.08889588	第一个人接收的
	//		0
	//		{
	//			OP_HASH160
	//			724d60b094581aec5756c33ec54ca9eb1850fd1b
	//			OP_EQUAL
	//			a914724d60b094581aec5756c33ec54ca9eb1850fd1b87
	//			0
	//			scripthash
	//			[]
	//		}
	//	}
	//	{
	//		0
	//		1
	//		{
	//			OP_RETURN
	//			6f6d6e69000000008000050d000000001dcd6500
	//			6a146f6d6e69000000008000050d000000001dcd6500
	//			0
	//			nulldata
	//			[]
	//		}
	//	}
	//	{
	//		5.4e-06  第三个人接收的
	//		2
	//		{
	//			OP_HASH160
	//			74d9b52e9e2457a737e233d336577b2c8b52188a
	//			OP_EQUAL
	//			a91474d9b52e9e2457a737e233d336577b2c8b52188a87
	//			0
	//			scripthash
	//			[]
	//		}
	//	}
	//]
	//00000000000000122988c78057b5739633c1e71e31299ab1dbe0857a93bc3319
	//7
	//1695177573
	//1695177573
	//}
	status := base.TransactionStatusPending
	if rawResult.Confirmations > 0 {
		status = base.TransactionStatusSuccess
	}
	return &base.TransactionDetail{
		Hash:   txHash,
		Status: status,
		Time:   rawResult.Time,
	}, nil
}

func (c *Chain) FetchTransactionStatus(hash string) (base.TransactionStatus, error) {
	detail, err := c.FetchTransactionDetail(hash)
	if err != nil {
		return base.TransactionStatusNone, log.WithError(err, "FetchTransactionDetail failed")
	}
	return detail.Status, nil
}
