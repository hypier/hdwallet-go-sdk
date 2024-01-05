package eth

import (
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/base"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"math/big"
)

type Chain struct {
	client  *Client
	chainId *big.Int
}

func NewChain() *Chain {
	return &Chain{
		client: nil,
	}
}

func (c *Chain) ChainId() *big.Int {
	return c.chainId
}

func (c *Chain) Client() (*Client, error) {
	if c.client == nil {
		return nil, log.WithError(utils.ErrClientNotInitialized)
	}

	return c.client, nil
}

func (c *Chain) CreateRemoteClient(rpcUrl string) (*Chain, error) {
	return c.CreateRemoteClientWithTimeout(rpcUrl, 60000)
}

func (c *Chain) CreateRemoteClientWithTimeout(rpcUrl string, timeout int64) (*Chain, error) {
	if c.client != nil && c.client.rpcUrl == rpcUrl {
		return c, nil
	}

	client, err := NewClient(rpcUrl, timeout)
	if err != nil {
		return nil, log.WithError(err, "CreateRemoteClient failed")
	}

	c.client = client
	c.chainId = client.chainId
	return c, nil
}

func (c *Chain) MainToken() base.Token {
	return &Token{chain: c}
}
