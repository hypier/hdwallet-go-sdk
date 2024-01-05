package base

import (
	"fmt"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"sync"
)

var (
	cachedToken map[string]*TokenInfo
	cacheMutex  sync.Mutex
)

type TokenInfo struct {
	Name    string
	Symbol  string
	Decimal int16
}

func AddToken(coinType uint32, contractAddress string, info *TokenInfo) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if cachedToken == nil {
		cachedToken = make(map[string]*TokenInfo)
	}

	if contractAddress == "" {
		contractAddress = "0x0"
	}
	log.Debugf("addToken:coinType%d_contractAddress%s", coinType, contractAddress)
	cachedToken[fmt.Sprintf("%d_%s", coinType, contractAddress)] = info
}

func GetToken(coinType uint32, contractAddress string) *TokenInfo {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if contractAddress == "" {
		contractAddress = "0x0"
	}
	log.Debugf("getToken:coinType%d_contractAddress%s", coinType, contractAddress)
	log.Debugf("cachedToken:%v", cachedToken[fmt.Sprintf("%d_%s", coinType, contractAddress)])
	return cachedToken[fmt.Sprintf("%d_%s", coinType, contractAddress)]
}

type Token interface {
	// Chain 链相关
	Chain() Chain

	// TokenInfo token信息
	TokenInfo() (*TokenInfo, error)

	// BalanceOfAddress 查询token余额
	BalanceOfAddress(address string) (*Balance, error)
}
