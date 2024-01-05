package ext

import (
	ne "errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"hypier.fun/hdwallet/hdwallet-go-sdk/config"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"testing"
)

func init() {
	config.InitConfig(&config.BaseConfig{
		BaseDir:    "..",
		LogSwitch:  "CONSOLE_FILE",
		Platform:   "ALL",
		DeviceType: "UNKNOWN",
	})
}

func TestWithError(t *testing.T) {

	t.Run("系统错误", func(t *testing.T) {
		err := log.WithError(ne.New("test error"), "test error")
		assert.Errorf(t, err, "test error")
		assert.Equalf(t, 500, err.(*utils.Error).ErrCode, "test error")
	})

	t.Run("Message is null", func(t *testing.T) {
		err := log.WithError(log.WithError(log.WithError(ne.New("test error"))))
		assert.Errorf(t, err, "test error")
		assert.Equalf(t, 500, err.(*utils.Error).ErrCode, "test error")
	})

	t.Run("自定义错误", func(t *testing.T) {
		err := log.WithError(utils.ErrInvalidURL, "test error")
		assert.Errorf(t, err, "test error")
		assert.Equalf(t, 106, err.(*utils.Error).ErrCode, "test error")
	})

	t.Run("嵌套错误", func(t *testing.T) {
		e := log.WithError(utils.ErrInvalidURL, "test error")
		err := log.WithError(e)
		assert.Errorf(t, err, "test error")
		assert.Equalf(t, 106, err.(*utils.Error).ErrCode, "test error")
	})

	t.Run("多层嵌套错误", func(t *testing.T) {
		e := log.WithError(utils.ErrInvalidURL, "test error 1")
		err := log.WithError(log.WithError(log.WithError(e, "test error 2"), "test error 3"), "test error 4")
		assert.Errorf(t, err, "test error")
		assert.Equalf(t, 106, err.(*utils.Error).ErrCode, "test error")
	})

	t.Run("自定义错误", func(t *testing.T) {
		err := log.WithError(errors.Wrap(errors.Wrap(utils.ErrInvalidURL, "test error 1"), "test error 2"), "test error 3")

		assert.Errorf(t, err, "test error")
		assert.Equalf(t, 106, err.(*utils.Error).ErrCode, "test error")
	})
}
