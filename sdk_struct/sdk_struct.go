package sdk_struct

import (
	"math/rand"
)

// BuildTime 编译时间 会在编译时注入
var BuildTime = ""

// GitId 编译会注入
var GitId = ""

const Tag = "GoSdkBtc"

var okLinkApikey = map[int]string{
	0: "b8989144-3a62-494d-9206-0f79d8244b55",
}

var apiKey = ""

func GetOkLinkApiKey() string {
	if len(apiKey) <= 0 {
		apiKey = okLinkApikey[rand.Intn(len(okLinkApikey))]
	}
	return apiKey
}

func SetApiKey(data string) {
	apiKey = data
}
