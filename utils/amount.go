package utils

import (
	"github.com/shopspring/decimal"
	"math/big"
	"strings"
)

type OptAmount struct {
	value   string
	decimal int16
}

func (o *OptAmount) SetDecimal(decimal int16) *OptAmount {
	o.decimal = decimal
	return o
}

func (o *OptAmount) BigInt() *big.Int {
	if o == nil {
		return nil
	}

	value := new(big.Int)
	value.SetString(o.value, 10)

	return value
}

func (o *OptAmount) AmountString() string {
	value := o.value
	decimal := int(o.decimal)

	if decimal == 0 {
		return value
	}

	n := len(value) - decimal
	if n > 0 {
		value = value[:n] + "." + value[n:]
	} else {
		value = "0." + strings.Repeat("0", n*-1) + value
	}

	value = strings.TrimRight(value, "0")
	value = strings.TrimRight(value, ".")

	return value
}

func NewOptAmount(value string, decimal int16) *OptAmount {

	_, success := new(big.Int).SetString(value, 10)
	if !success {
		return nil
	}

	return &OptAmount{
		value:   value,
		decimal: decimal,
	}
}

// ParseAmount 将string转换为OptAmount
func ParseAmount(value string, decimal int16) (*OptAmount, error) {
	if decimal <= 0 {
		return nil, ErrInvalidValue
	}

	index := strings.Index(value, ".")
	if index <= 0 {
		value = value + strings.Repeat("0", int(decimal))
		return NewOptAmount(value, decimal), nil
	}

	n := int(decimal) - (len(value) - index - 1)
	if n < 0 {
		return nil, ErrInvalidValue
	}

	value = strings.ReplaceAll(value, ".", "")

	if n > 0 {
		value = value + strings.Repeat("0", n)
	}

	return NewOptAmount(value, decimal), nil
}

// Wei2ethStr 高精度转低精度
func Wei2ethStr(value *big.Int, decimals int16) string {
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)
	return result.String()
}
