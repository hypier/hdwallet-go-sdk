package eth

import (
	"context"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	cachedClient = &sync.Map{}
)

type Client struct {
	rpcClient *ethclient.Client
	rpcUrl    string
	timeout   time.Duration
	chainId   *big.Int

	client *rpc.Client
}

func (c *Client) Client() *rpc.Client {
	return c.client
}

func (c *Client) RPCClient() *ethclient.Client {
	return c.rpcClient
}

func (c *Client) Timeout() time.Duration {
	return c.timeout
}

// NewClient 创建一个新的客户端
func NewClient(url string, timeout int64) (*Client, error) {
	if url == "" {
		return nil, log.WithError(utils.ErrInvalidURL, "NewClient failed")
	}

	c := getClient(url)
	if c != nil {
		return c, nil
	}

	var t time.Duration
	if timeout == 0 {
		t = 60 * time.Second
	} else {
		t = time.Duration(timeout * int64(time.Millisecond))
	}

	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	client, err := rpc.DialContext(ctx, url)
	if err != nil {
		return nil, log.WithError(err, "DialContext failed")
	}

	rpcClient := ethclient.NewClient(client)
	chainId, err1 := rpcClient.ChainID(ctx)

	if err1 != nil {
		return nil, log.WithError(err1, "ChainID failed")
	}

	c2 := &Client{rpcClient: rpcClient, client: client, timeout: t, rpcUrl: url, chainId: chainId}

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
