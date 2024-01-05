package config

import (
	"github.com/ethereum/go-ethereum/common"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"sync"
)

var (
	once sync.Once
	Base *BaseConfig
	Ens  map[int]EnsConfig //chainId
)

type BaseConfig struct {
	BaseDir      string
	LogSwitch    string
	KeyStorePath string
	LogFile      string
	Platform     string
	DeviceType   string
	BtcParam     string
}

func init() {
	Ens = make(map[int]EnsConfig)
	Ens[2888] = EnsConfig{
		Registry: common.HexToAddress("0x78F0390de2eF54a8B6De164D0E62E340455c57dc"),
	}
}

type EnsConfig struct {
	Registry common.Address //注册器合约地址
}

// InitConfig 初始化方法
func InitConfig(baseConfig *BaseConfig) string {
	//初始化日志
	once.Do(func() {
		Base = &BaseConfig{
			KeyStorePath: baseConfig.BaseDir + "/key",
			LogFile:      baseConfig.BaseDir + "/logs",
			BtcParam:     baseConfig.BtcParam,
		}
		log.InitLog(baseConfig.BaseDir+"/logs", baseConfig.LogSwitch)
	})
	return baseConfig.BaseDir
}
