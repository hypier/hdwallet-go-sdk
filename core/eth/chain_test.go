package eth

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestCreateRemoteClient(t *testing.T) {

	rpcUrl := ClientCase1().url

	// Positive test case 正向
	t.Run("Positive", func(t *testing.T) {
		chain := &Chain{}
		_, err := chain.CreateRemoteClient(rpcUrl)
		assert.NoError(t, err)

		// Verify client and chainId are set correctly
		assert.NotNil(t, chain.client)
		assert.Equal(t, ClientCase1().chainId, chain.chainId)
	})

	// Negative test case - invalid RPC URL 反向
	t.Run("InvalidURL", func(t *testing.T) {
		chain := &Chain{}
		_, err := chain.CreateRemoteClient("")
		assert.Error(t, err)
	})

}

func TestChain_Client(t *testing.T) {
	// Create a mock client
	mockClient := &Client{}

	// Create a chain with the mock client
	chain := &Chain{
		client: mockClient,
	}

	// Positive test case 正向
	t.Run("Positive", func(t *testing.T) {
		client, err := chain.Client()
		assert.Equal(t, mockClient, client)
		assert.Nil(t, err)
	})

	// Negative test case 反向
	t.Run("Negative", func(t *testing.T) {
		// Create a chain without initializing the client
		chain := &Chain{}

		client, err := chain.Client()
		assert.Nil(t, client)
		assert.Error(t, err)
	})
}

func TestChain_ChainId(t *testing.T) {
	// Positive test case 正向
	t.Run("Positive", func(t *testing.T) {
		expected := big.NewInt(123)
		chain := &Chain{chainId: expected}

		actual := chain.ChainId()
		assert.Equal(t, expected, actual)
	})

	// Negative test case 反向
	t.Run("Negative", func(t *testing.T) {
		chain := NewChain()
		actual := chain.ChainId()
		assert.Nil(t, actual)
	})
}
