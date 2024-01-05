package trx

import (
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"google.golang.org/grpc"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"sync"
	"time"
)

var (
	cachedClient = &sync.Map{}
)

type Client struct {
	rpcClient *client.GrpcClient
	rpcUrl    string
}

func (c *Client) RPCClient() *client.GrpcClient {
	return c.rpcClient

}

func NewClient(url, apiKey string) (*Client, error) {
	//if url == "" {
	//	return nil, log.WithError(utils.ErrInvalidURL, "NewClient failed")
	//}

	c := getClient(url)
	if c != nil {
		return c, nil
	}

	rpcClient := client.NewGrpcClient(url)
	if apiKey != "" {
		err := rpcClient.SetAPIKey(apiKey)
		if err != nil {
			return nil, log.WithError(err, "SetAPIKey failed")
		}
	}

	rpcClient.SetTimeout(30 * time.Second)

	err := rpcClient.Start(grpc.WithInsecure())
	if err != nil {
		return nil, log.WithError(err, "Start failed")
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
