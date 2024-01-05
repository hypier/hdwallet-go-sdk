package btc

import (
	"github.com/btcsuite/btcd/rpcclient"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"strings"
	"sync"
)

var (
	cachedClient = &sync.Map{}
)

type Client struct {
	rpcClient *rpcclient.Client
	rpcUrl    string
}

func (c *Client) RPCClient() *rpcclient.Client {
	return c.rpcClient

}

func NewClient(url, user, pass string, chainId int) (*Client, error) {
	if url == "" {
		return nil, log.WithError(utils.ErrInvalidURL, "NewClient failed")
	}

	if strings.HasPrefix(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
	}

	c := getClient(url)
	if c != nil {
		return c, nil
	}

	params, err := utils.GetBtcChainParams(chainId)
	if err != nil {
		return nil, log.WithError(err, "ChainID failed")
	}

	rpcClient, err := rpcclient.New(&rpcclient.ConnConfig{
		Host:         url,
		User:         user,
		Pass:         pass,
		HTTPPostMode: true,
		DisableTLS:   false,
		Params:       params.Name,
	}, nil)

	if err != nil {
		return nil, log.WithError(err, "ChainID failed")
	}

	c2 := &Client{rpcClient: rpcClient, rpcUrl: url}

	addClient(url, c2)
	return c2, nil
}

func getClient(url string) *Client {
	value, ok := cachedClient.Load(url)
	if ok {
		return value.(*Client)
	}
	return nil
}

func addClient(url string, client *Client) {
	cachedClient.Store(url, client)
}
