package utils

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
)

// zero is deafult of uint32
const (
	Zero      uint32 = 0
	ZeroQuote uint32 = 0x80000000
	BTCToken  uint32 = 0x10000000
	ETHToken  uint32 = 0x20000000
)

// wallet type from bip44
const (
	// BTC https://github.com/satoshilabs/slips/blob/master/slip-0044.md#registered-coin-types
	BTC  = ZeroQuote + 0
	LTC  = ZeroQuote + 2
	DOGE = ZeroQuote + 3
	DASH = ZeroQuote + 5
	DCR  = ZeroQuote + 42
	NEM  = ZeroQuote + 43

	ETH   = ZeroQuote + 60
	ETC   = ZeroQuote + 61
	QTUM  = ZeroQuote + 88
	ATOM  = ZeroQuote + 118
	XMR   = ZeroQuote + 128
	ZCash = ZeroQuote + 133
	XRP   = ZeroQuote + 144
	BCH   = ZeroQuote + 145
	BTM   = ZeroQuote + 153
	HC    = ZeroQuote + 171
	RVN   = ZeroQuote + 175
	XLM   = ZeroQuote + 184

	EOS  = ZeroQuote + 194
	TRX  = ZeroQuote + 195
	ALGO = ZeroQuote + 283
	CKB  = ZeroQuote + 309

	AE    = ZeroQuote + 457
	BNB   = ZeroQuote + 714
	VET   = ZeroQuote + 818
	NEO   = ZeroQuote + 888
	ONT   = ZeroQuote + 1024
	XTZ   = ZeroQuote + 1729
	Libra = ZeroQuote + 9999
	WAVES = ZeroQuote + 5741564
	// USDT btc token
	USDT = BTCToken + 1

	// IOST eth token
	IOST = ETHToken + 1
	USDC = ETHToken + 2

	BtcChainMainNet  = int(wire.MainNet)
	BtcChainTestNet3 = int(wire.TestNet3)
	BtcChainRegtest  = int(wire.TestNet)
	BtcChainSimNet   = int(wire.SimNet)
)

var CoinTypes = map[uint32]uint32{
	USDT: BTC,
	IOST: ETH,
	USDC: ETH,
}

var MainChainId = map[uint32]int{
	ETH: 1,
}

func GetBtcChainParams(chainId int) (*chaincfg.Params, error) {
	switch chainId {
	case BtcChainMainNet:
		return &chaincfg.MainNetParams, nil
	case BtcChainTestNet3:
		return &chaincfg.TestNet3Params, nil
	case BtcChainRegtest:
		return &chaincfg.RegressionNetParams, nil
	case BtcChainSimNet:
		return &chaincfg.SimNetParams, nil
	default:
		return nil, fmt.Errorf("unknown btc chainId: %d", chainId)
	}
}
func GetBtcChainParam(name string) (*chaincfg.Params, error) {
	switch name {
	case wire.MainNet.String():
		return &chaincfg.MainNetParams, nil
	case wire.TestNet3.String():
		return &chaincfg.TestNet3Params, nil
	case wire.TestNet.String():
		return &chaincfg.RegressionNetParams, nil
	case wire.SimNet.String():
		return &chaincfg.SimNetParams, nil
	default:
		return nil, fmt.Errorf("unknown btc name: %s", name)
	}
}
func GetBtcChainId(name string) (int, error) {
	switch name {
	case wire.MainNet.String():
		return BtcChainMainNet, nil
	case wire.TestNet3.String():
		return BtcChainTestNet3, nil
	case wire.TestNet.String():
		return BtcChainRegtest, nil
	case wire.SimNet.String():
		return BtcChainSimNet, nil
	default:
		return -1, fmt.Errorf("unknown btc name: %s", name)
	}
}
