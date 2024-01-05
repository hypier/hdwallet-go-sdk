package eth

import "hypier.fun/hdwallet/hdwallet-go-sdk/utils"

type Coin struct {
}

func (c *Coin) CoinType() uint32 {
	return utils.ETH
}

func (c *Coin) Symbol() string {
	return "ETH"
}

func (c *Coin) Name() string {
	return "Ether"
}
