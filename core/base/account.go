package base

type Account interface {
	// PrivateKey 私钥
	PrivateKey() []byte

	// PrivateKeyHex 私钥
	PrivateKeyHex() string

	// PublicKey 公钥
	PublicKey() []byte

	// PublicKeyHex 公钥
	PublicKeyHex() string
}
