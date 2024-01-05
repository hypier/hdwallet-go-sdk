package btc

import "hypier.fun/hdwallet/hdwallet-go-sdk/utils"

type Coin struct {
}

func (c *Coin) CoinType() uint32 {
	return utils.BTC
}

func (c *Coin) Symbol() string {
	return "BTC"
}

func (c *Coin) Name() string {
	return "Bitcoin"
}
