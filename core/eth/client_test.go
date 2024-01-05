package eth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient(t *testing.T) {

	testcase := ClientCase2()

	// Positive test case - valid URL and timeout 60000秒
	t.Run("ValidURLAndTimeout", func(t *testing.T) {
		c, err := NewClient(testcase.url, testcase.timeout)
		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, testcase.chainId, c.chainId)
		assert.Equal(t, testcase.url, c.rpcUrl)
		assert.Equal(t, testcase.timeout, c.timeout.Milliseconds())
	})

	// Positive test case - valid URL and default timeout 60秒
	t.Run("ValidURLAndDefaultTimeout", func(t *testing.T) {
		c, err := NewClient(testcase.url, 0)
		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, testcase.chainId, c.chainId)
		assert.Equal(t, testcase.url, c.rpcUrl)
		assert.Equal(t, testcase.timeout, c.timeout.Milliseconds())
	})

	// Negative test case - invalid URL
	t.Run("InvalidURL", func(t *testing.T) {
		c, err := NewClient("invalid-url", 500)
		assert.Error(t, err)
		assert.Nil(t, c)
	})

	// Negative test case - RPC dial error
	t.Run("RPCDialError", func(t *testing.T) {
		c, err := NewClient("http://example.com", 500)
		assert.Error(t, err)
		assert.Nil(t, c)
	})

	// Negative test case - ChainID error
	t.Run("ChainIDError", func(t *testing.T) {
		c, err := NewClient("http://example.com", 500)
		assert.Error(t, err)
		assert.Nil(t, c)
	})

}

func TestGetClient(t *testing.T) {
	// Positive test case
	t.Run("Positive", func(t *testing.T) {
		url := ClientCase1().url
		client := &Client{}
		addClient(url, client)

		got := getClient(url)
		assert.Equal(t, client, got)
	})

	// Negative test case
	t.Run("Negative", func(t *testing.T) {
		nonExistentURL := "https://nonexistent.com"
		got := getClient(nonExistentURL)
		assert.Nil(t, got)
	})
}

func TestAddClient(t *testing.T) {
	// Positive test case 正向
	t.Run("Positive", func(t *testing.T) {
		url := ClientCase2().url
		client := &Client{}
		addClient(url, client)

		got, _ := cachedClient.Load(url)
		assert.Equal(t, client, got)
	})

	// Negative test case 反向
	t.Run("Negative", func(t *testing.T) {
		duplicateURL := ClientCase2().url
		anotherClient := &Client{}
		addClient(duplicateURL, anotherClient)

		got, _ := cachedClient.Load(duplicateURL)
		assert.Equal(t, anotherClient, got)
	})
}
