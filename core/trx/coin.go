package trx

import "hypier.fun/hdwallet/hdwallet-go-sdk/utils"

type Coin struct {
}

func (c *Coin) CoinType() uint32 {
	return utils.TRX
}

func (c *Coin) Symbol() string {
	return "TRX"
}

func (c *Coin) Name() string {
	return "Tron"
}
