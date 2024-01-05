package base

type Chain interface {
	// MainToken 主币种
	MainToken() Token

	//// SendRawTransaction 发送原始签名的交易
	//SendRawTransaction(signedTx string) (string, error)

	// FetchTransactionDetail 获取交易明细
	FetchTransactionDetail(hash string) (*TransactionDetail, error)

	// FetchTransactionStatus 获取交易状态
	FetchTransactionStatus(hash string) (TransactionStatus, error)
}
