package keystore

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"hypier.fun/hdwallet/hdwallet-go-sdk/config"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"

	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

const (
	keyHeaderKDF = "scrypt"

	// StandardScryptN is the N parameter of Scrypt encryption algorithm, using 256MB
	// memory and taking approximately 1s CPU time on a modern processor.
	StandardScryptN = 1 << 18

	// StandardScryptP is the P parameter of Scrypt encryption algorithm, using 256MB
	// memory and taking approximately 1s CPU time on a modern processor.
	StandardScryptP = 1

	// LightScryptN is the N parameter of Scrypt encryption algorithm, using 4MB
	// memory and taking approximately 100ms CPU time on a modern processor.
	LightScryptN = 1 << 12

	// LightScryptP is the P parameter of Scrypt encryption algorithm, using 4MB
	// memory and taking approximately 100ms CPU time on a modern processor.
	LightScryptP = 6

	scryptR     = 8
	scryptDKLen = 32

	version = "1"
)

type StorePassphrase struct {
	keysDirPath string
	scryptN     int
	scryptP     int
	// skipKeyFileVerification disables the security-feature which does
	// reads and decrypts any newly created keyfiles. This should be 'false' in all
	// cases except tests -- setting this to 'true' is not recommended.
	skipKeyFileVerification bool
}

func NewKeyStore() *StorePassphrase {
	return &StorePassphrase{
		keysDirPath:             config.Base.KeyStorePath,
		scryptN:                 LightScryptN,
		scryptP:                 LightScryptP,
		skipKeyFileVerification: false,
	}
}

func (ks StorePassphrase) GetKey(key *Key, auth string) (*Key, error) {
	return getKey(ks.JoinPath(key.FileName()), key.KeyID, auth)
}

func getKey(filename, keyID, auth string) (*Key, error) {
	// Load the key from the keystore and decrypt its contents
	keyJson, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	newKey, err := DecryptKey(keyJson, auth)
	if err != nil {
		return nil, err
	}
	// Make sure we're really operating on the requested key (no swap attacks)
	if keyID != newKey.KeyID {
		return nil, fmt.Errorf("key content mismatch: have account %x, want %x", keyID, newKey.KeyID)
	}
	return newKey, nil
}

//StoreKey generates a key, encrypts with 'auth' and stores in the given directory
//func StoreKey(dir, auth string, scryptN, scryptP int) (accounts.Account, error) {
//	_, a, err := storeNewKey(&keyStorePassphrase{dir, scryptN, scryptP, false}, rand.Reader, auth)
//	return err
//}

func (ks StorePassphrase) StoreKey(key *Key, auth string) error {
	if auth == "" {
		return utils.ErrInvalidPassword
	}

	accountJSONS := make(map[uint32]EncryptedAccountJSON)
	for coin, account := range key.Accounts {
		keyJson, err := EncryptAccount(account, auth, ks.scryptN, ks.scryptP)
		if err != nil {
			return err
		}

		accountJSONS[coin] = keyJson
	}

	encryptedKeyJSON := EncryptedKeyJSON{
		Alias:    key.Alias,
		KeyID:    key.KeyID,
		Accounts: accountJSONS,
		Time:     uint64(time.Now().Unix()),
		Version:  version,
	}

	if key.Mnemonic != "" {
		cryptoJSON, err1 := EncryptKey(key.Mnemonic, auth, ks.scryptN, ks.scryptP)
		if err1 != nil {
			return err1
		}

		encryptedKeyJSON.Mnemonic = &cryptoJSON
	}

	jsonFile, err := json.MarshalIndent(encryptedKeyJSON, "", " ")
	if err != nil {
		return err
	}

	// Write into temporary file
	tmpName, err := writeTemporaryKeyFile(ks.JoinPath(key.FileName()), jsonFile)
	if err != nil {
		return err
	}
	if !ks.skipKeyFileVerification {
		// Verify that we can decrypt the file with the given password.
		_, err = getKey(tmpName, key.KeyID, auth)
		if err != nil {
			msg := "An error was encountered when saving and verifying the keystore file. \n" +
				"This indicates that the keystore is corrupted. \n" +
				"The corrupted file is stored at \n%v\n" +
				"Please file a ticket at:\n\n" +
				"https://github.com/ethereum/go-ethereum/issues." +
				"The error was : %s"
			//lint:ignore ST1005 This is a message for the user
			return fmt.Errorf(msg, tmpName, err)
		}
	}
	return os.Rename(tmpName, ks.JoinPath(key.FileName()))
}

func (ks StorePassphrase) JoinPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(ks.keysDirPath, filename)
}

// EncryptDataV3 encrypts the data given as 'data' with the password 'auth'.
func EncryptDataV3(data, auth []byte, scryptN, scryptP int) (CryptoJSON, error) {
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	derivedKey, err := scrypt.Key(auth, salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return CryptoJSON{}, err
	}
	encryptKey := derivedKey[:16]

	iv := make([]byte, aes.BlockSize) // 16
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	cipherText, err := aesCTRXOR(encryptKey, data, iv)
	if err != nil {
		return CryptoJSON{}, err
	}
	mac := crypto.Keccak256(derivedKey[16:32], cipherText)

	scryptParamsJSON := make(map[string]interface{}, 5)
	scryptParamsJSON["n"] = scryptN
	scryptParamsJSON["r"] = scryptR
	scryptParamsJSON["p"] = scryptP
	scryptParamsJSON["dklen"] = scryptDKLen
	scryptParamsJSON["salt"] = hex.EncodeToString(salt)
	cipherParamsJSON := cipherParamsJSON{
		IV: hex.EncodeToString(iv),
	}

	cryptoStruct := CryptoJSON{
		Cipher:       "aes-128-ctr",
		CipherText:   hex.EncodeToString(cipherText),
		CipherParams: cipherParamsJSON,
		KDF:          keyHeaderKDF,
		KDFParams:    scryptParamsJSON,
		MAC:          hex.EncodeToString(mac),
	}
	return cryptoStruct, nil
}

// EncryptAccount encrypts a key using the specified scrypt parameters into a json
// blob that can be decrypted later on.
func EncryptAccount(account *Account, auth string, scryptN, scryptP int) (EncryptedAccountJSON, error) {

	cryptoStruct, err := EncryptDataV3([]byte(account.PrivateKey), []byte(auth), scryptN, scryptP)
	if err != nil {
		return EncryptedAccountJSON{}, err
	}

	return EncryptedAccountJSON{
		Address: account.Address,
		Crypto:  cryptoStruct,
	}, nil
}

func EncryptKey(data, auth string, scryptN, scryptP int) (CryptoJSON, error) {

	cryptoStruct, err := EncryptDataV3([]byte(data), []byte(auth), scryptN, scryptP)
	if err != nil {
		return CryptoJSON{}, err
	}

	return cryptoStruct, nil

}

// DecryptKey decrypts a key from a json blob, returning the private key itself.
func DecryptKey(keyJson []byte, auth string) (*Key, error) {
	// Parse the json into a simple map to fetch the key version
	m := make(map[string]interface{})
	if err := json.Unmarshal(keyJson, &m); err != nil {
		return nil, err
	}
	// Depending on the version try to parse one way or another

	k := new(EncryptedKeyJSON)
	if err := json.Unmarshal(keyJson, k); err != nil {
		return nil, err
	}

	key, err := decryptKeyV3(k, auth)

	// Handle any decryption errors and return the key
	if err != nil {
		return nil, err
	}

	return key, nil
}

func decryptDataV3(cryptoJson CryptoJSON, auth string) ([]byte, error) {
	if cryptoJson.Cipher != "aes-128-ctr" {
		return nil, fmt.Errorf("cipher not supported: %v", cryptoJson.Cipher)
	}
	mac, err := hex.DecodeString(cryptoJson.MAC)
	if err != nil {
		return nil, err
	}

	iv, err := hex.DecodeString(cryptoJson.CipherParams.IV)
	if err != nil {
		return nil, err
	}

	cipherText, err := hex.DecodeString(cryptoJson.CipherText)
	if err != nil {
		return nil, err
	}

	derivedKey, err := getKDFKey(cryptoJson, auth)
	if err != nil {
		return nil, err
	}

	calculatedMAC := crypto.Keccak256(derivedKey[16:32], cipherText)
	if !bytes.Equal(calculatedMAC, mac) {
		return nil, utils.ErrDecrypt
	}

	plainText, err := aesCTRXOR(derivedKey[:16], cipherText, iv)
	if err != nil {
		return nil, err
	}
	return plainText, err
}

func decryptKeyV3(keyProtected *EncryptedKeyJSON, auth string) (*Key, error) {
	if keyProtected.Version != version {
		return nil, fmt.Errorf("version not supported: %v", keyProtected.Version)
	}

	var mPlainText []byte

	if keyProtected.Mnemonic != nil {
		text, err := decryptDataV3(*keyProtected.Mnemonic, auth)
		if err != nil {
			return nil, err
		}

		mPlainText = text
	} else {
		mPlainText = nil
	}

	var accounts = make(Accounts)
	for key, acc := range keyProtected.Accounts {
		aPlainText, err := decryptDataV3(acc.Crypto, auth)
		if err != nil {
			return nil, err
		}

		account := &Account{
			Address:    acc.Address,
			PrivateKey: string(aPlainText),
		}

		accounts[key] = account
	}

	// 将EncryptedKeyJSON 转为 KeyID
	return &Key{
		Alias:    keyProtected.Alias,
		KeyID:    keyProtected.KeyID,
		Mnemonic: string(mPlainText),
		Time:     keyProtected.Time,
		Accounts: accounts,
	}, nil
}

func getKDFKey(cryptoJSON CryptoJSON, auth string) ([]byte, error) {
	authArray := []byte(auth)
	salt, err := hex.DecodeString(cryptoJSON.KDFParams["salt"].(string))
	if err != nil {
		return nil, err
	}
	dkLen := ensureInt(cryptoJSON.KDFParams["dklen"])

	if cryptoJSON.KDF == keyHeaderKDF {
		n := ensureInt(cryptoJSON.KDFParams["n"])
		r := ensureInt(cryptoJSON.KDFParams["r"])
		p := ensureInt(cryptoJSON.KDFParams["p"])
		return scrypt.Key(authArray, salt, n, r, p, dkLen)
	} else if cryptoJSON.KDF == "pbkdf2" {
		c := ensureInt(cryptoJSON.KDFParams["c"])
		prf := cryptoJSON.KDFParams["prf"].(string)
		if prf != "hmac-sha256" {
			return nil, fmt.Errorf("unsupported PBKDF2 PRF: %s", prf)
		}
		key := pbkdf2.Key(authArray, salt, c, dkLen, sha256.New)
		return key, nil
	}

	return nil, fmt.Errorf("unsupported KDF: %s", cryptoJSON.KDF)
}

// TODO: can we do without this when unmarshalling dynamic JSON?
// why do integers in KDF params end up as float64 and not int after
// unmarshal?
func ensureInt(x interface{}) int {
	res, ok := x.(int)
	if !ok {
		res = int(x.(float64))
	}
	return res
}

func writeTemporaryKeyFile(file string, content []byte) (string, error) {
	// Create the keystore directory with appropriate permissions
	// in case it is not present yet.
	const dirPerm = 0700
	if err := os.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return "", err
	}
	// Atomic write: create a temporary hidden file first
	// then move it into place. TempFile assigns mode 0600.
	f, err := os.CreateTemp(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return "", err
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", err
	}
	f.Close()
	return f.Name(), nil
}

func aesCTRXOR(key, inText, iv []byte) ([]byte, error) {
	// AES-128 is selected due to size of encryptKey.
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	outText := make([]byte, len(inText))
	stream.XORKeyStream(outText, inText)
	return outText, err
}
