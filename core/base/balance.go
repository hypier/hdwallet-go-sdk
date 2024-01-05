package base

import "hypier.fun/hdwallet/hdwallet-go-sdk/utils"

type Balance struct {
	Total  *utils.OptAmount
	Usable *utils.OptAmount
}

func EmptyBalance() *Balance {
	amount := utils.NewOptAmount("0", 0)
	return &Balance{
		Total:  amount,
		Usable: amount,
	}
}
