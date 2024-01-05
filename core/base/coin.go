package base

type Coin interface {
	CoinType() uint32
	Symbol() string
	Name() string
}
