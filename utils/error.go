package utils

import (
	"errors"
	"fmt"
)

var (
	ErrClientNotInitialized = NewError(100, "remote client not initialized")

	ErrInvalidValue = NewError(101, "invalid value")

	// ErrInvalidMnemonicPhrase 无效的助记词
	ErrInvalidMnemonicPhrase = NewError(102, "invalid mnemonic phrase")

	// ErrInvalidPrivateKey 无效的私钥
	ErrInvalidPrivateKey = NewError(103, "invalid private key")

	ErrDecrypt = NewError(104, "could not decrypt key with given password")

	// ErrInvalidPassword 无效的密码
	ErrInvalidPassword = NewError(105, "invalid password")

	// ErrInvalidURL 非法的URL
	ErrInvalidURL = NewError(106, "invalid url")
	//ErrAccountNoAllowance 账户没有允许的allowance
	ErrAccountNoAllowance = NewError(107, "account has no allowance")
	TransactionHashError  = NewError(108, "transaction hash error")
	RPCError              = NewError(109, "rpc error")
	FromError             = NewError(110, "from error")
	ToError               = NewError(111, "to error")
	AmountError           = NewError(112, "amount error")
	PasswordError         = NewError(113, "password error")
)

type Error struct {
	ErrCode    int
	ErrMsg     error
	IsHasStack bool
}

func NewError(errCode int, errMsg string) *Error {

	return &Error{
		ErrCode: errCode,
		ErrMsg:  errors.New(errMsg),
	}
}

func ConvertError(err error) *Error {
	var u *Error

	if !errors.As(err, &u) {
		u = NewSysError(err)
	}
	return u
}

func NewSysError(err error) *Error {
	return &Error{
		ErrCode: 500,
		ErrMsg:  err,
	}
}

func (e *Error) ErrorCode() int {
	return e.ErrCode
}

func (e *Error) Error() string {
	return fmt.Sprintf("-> Code:%d, ErrMsg:%s", e.ErrCode, e.ErrMsg.Error())
}
