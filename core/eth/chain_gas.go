package eth

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"strings"
)

var (
	DefaultContractGasLimit = 63000
	DefaultEthGasList       = 21000
)

func (c *Chain) EstimateGasLimit(msg *ethereum.CallMsg) (uint64, error) {
	client, err := c.Client()
	if err != nil {
		return 0, log.WithError(err, "Client failed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout())
	defer cancel()
	gas, err := client.RPCClient().EstimateGas(ctx, *msg)
	if err != nil {
		//余额不足
		if strings.Contains(err.Error(), "gas required exceeds allowance") {
			return 0, log.WithError(utils.ErrAccountNoAllowance)
		}
		return 0, log.WithError(err, "EstimateGas failed")
	}

	return gas, nil
}
