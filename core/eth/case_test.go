package eth

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

type TestAccountCase struct {
	caseName      string
	mnemonic      string
	privateKey    *ecdsa.PrivateKey
	privateKeyHex string
	from          common.Address
	to            common.Address
	url           string
	chainId       *big.Int
	erc20         common.Address
	erc721        common.Address
}

func CaseAccountGoerliNew() *TestAccountCase {

	pri := "fdf55d604973d4afc0781e6baa98046b11cd876b1b0e3e2bfbf21c917ad83520"
	privateKey, _ := crypto.HexToECDSA(pri)
	return &TestAccountCase{
		caseName:      "case1-B5aC",
		mnemonic:      "tattoo season illegal swallow embody hundred face vast moon tent answer trim",
		privateKey:    privateKey,
		privateKeyHex: pri,
		from:          common.HexToAddress("0x986E43FcC2911f2De260d610E38D390A2D0a824F"),
		to:            common.HexToAddress("0xC7bFc5CB33ADA7566217ebbF15B5Dc25f5e609D7"),
		erc20:         common.HexToAddress("0x168ecd12b8C96ed4F0684Ce0D3A3C9a03dAb32C9"),
		url:           "https://eth-goerli.g.alchemy.com/v2/xnUWpWzH39i8uYi2PR7MW-jG-Vnl-Y5u",
		chainId:       big.NewInt(5),
	}
}

func CaseAccountGoerli() *TestAccountCase {

	pri := "1032adbf75a73f959d30dcae3e35a2c12252daac44abf8d9d2e21b29754db496"
	privateKey, _ := crypto.HexToECDSA(pri)
	return &TestAccountCase{
		caseName:      "case2-5c0Bf",
		mnemonic:      "",
		privateKey:    privateKey,
		privateKeyHex: pri,
		from:          common.HexToAddress("0xC7bFc5CB33ADA7566217ebbF15B5Dc25f5e609D7"),
		to:            common.HexToAddress("0x986E43FcC2911f2De260d610E38D390A2D0a824F"),
		erc20:         common.HexToAddress("0x168ecd12b8C96ed4F0684Ce0D3A3C9a03dAb32C9"),
		url:           "https://eth-goerli.g.alchemy.com/v2/xnUWpWzH39i8uYi2PR7MW-jG-Vnl-Y5u",
		chainId:       big.NewInt(5),
	}
}

type TestClientCase struct {
	caseName string
	url      string
	chainId  *big.Int
	timeout  int64
}

func ClientCase1() *TestClientCase {
	return &TestClientCase{
		caseName: "goerli",
		url:      "https://ethereum-goerli.publicnode.com",
		chainId:  big.NewInt(5),
		timeout:  60000,
	}
}

func ClientCase2() *TestClientCase {
	return &TestClientCase{
		caseName: "goerli",
		url:      "https://ethereum-goerli.publicnode.com",
		chainId:  big.NewInt(5),
		timeout:  60000,
	}
}
