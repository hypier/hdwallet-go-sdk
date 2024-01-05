package trx

import (
	"hypier.fun/hdwallet/hdwallet-go-sdk/core/base"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
)

type Chain struct {
	client *Client
}

func NewChain() *Chain {
	return &Chain{
		client: nil,
	}
}

func (c *Chain) Client() (*Client, error) {
	if c.client == nil {
		return nil, log.WithError(utils.ErrClientNotInitialized)
	}

	return c.client, nil
}

func (c *Chain) CreateRemoteClientWithTimeout(rpcUrl, apiKey string) (*Chain, error) {
	if c.client != nil && c.client.rpcUrl == rpcUrl {
		return c, nil
	}

	client, err := NewClient(rpcUrl, apiKey)
	if err != nil {
		return nil, log.WithError(err, "CreateRemoteClient failed")
	}

	c.client = client
	return c, nil
}

func (c *Chain) MainToken() base.Token {
	return &Token{chain: c}
}
