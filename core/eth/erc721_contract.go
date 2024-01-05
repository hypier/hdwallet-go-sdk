package eth

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils/log"
	"math/big"
	"strings"
)

const ERC721InterfaceABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"approved","type":"address"},{"indexed":true,"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"operator","type":"address"},{"indexed":false,"internalType":"bool","name":"approved","type":"bool"}],"name":"ApprovalForAll","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":true,"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"Transfer","type":"event"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"approve","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"balance","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"getApproved","outputs":[{"internalType":"address","name":"operator","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"operator","type":"address"}],"name":"isApprovedForAll","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"ownerOf","outputs":[{"internalType":"address","name":"owner","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"safeTransferFrom","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"},{"internalType":"bytes","name":"data","type":"bytes"}],"name":"safeTransferFrom","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"operator","type":"address"},{"internalType":"bool","name":"_approved","type":"bool"}],"name":"setApprovalForAll","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes4","name":"interfaceId","type":"bytes4"}],"name":"supportsInterface","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"index","type":"uint256"}],"name":"tokenByIndex","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"index","type":"uint256"}],"name":"tokenOfOwnerByIndex","outputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"tokenURI","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"transferFrom","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

// Erc721Contract tool for contract abi
type Erc721Contract struct {
	abi             abi.ABI
	contractAddress common.Address
	backend         bind.ContractBackend
	contract        *bind.BoundContract
	opts            *bind.CallOpts
}

func NewErc721Contract(address common.Address, backend bind.ContractBackend) *Erc721Contract {
	parsed, _ := abi.JSON(strings.NewReader(ERC721InterfaceABI))
	c := bind.NewBoundContract(address, parsed, backend, backend, backend)
	return &Erc721Contract{abi: parsed, contractAddress: address, backend: backend, contract: c, opts: &bind.CallOpts{}}
}

// IERC721Enumerable
func (e *Erc721Contract) TotalSupply() (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := []interface{}{ret0}
	err := e.contract.Call(e.opts, &out, "totalSupply")
	if err != nil {
		return nil, log.WithError(err)
	}
	return *ret0, nil
}

func (e *Erc721Contract) TokenOfOwnerByIndex(owner common.Address, index *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := []interface{}{ret0}
	err := e.contract.Call(e.opts, &out, "tokenOfOwnerByIndex", owner, index)
	if err != nil {
		return nil, log.WithError(err)
	}
	return *ret0, nil
}

func (e *Erc721Contract) TokenByIndex(index *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := []interface{}{ret0}
	err := e.contract.Call(e.opts, &out, "tokenByIndex", index)
	if err != nil {
		return nil, log.WithError(err)
	}
	return *ret0, nil
}

// IERC721Metadata
func (e *Erc721Contract) Name() (string, error) {
	var (
		ret0 = new(string)
	)
	out := []interface{}{ret0}
	err := e.contract.Call(e.opts, &out, "name")
	if err != nil {
		return "", log.WithError(err)
	}
	return *ret0, nil
}

func (e *Erc721Contract) Symbol() (string, error) {
	var (
		ret0 = new(string)
	)
	out := []interface{}{ret0}
	err := e.contract.Call(e.opts, &out, "symbol")
	if err != nil {
		return "", log.WithError(err)
	}
	return *ret0, nil
}

func (e *Erc721Contract) TokenURI(tokenId *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := []interface{}{ret0}
	err := e.contract.Call(e.opts, &out, "tokenURI", tokenId)
	if err != nil {
		return "", log.WithError(err)
	}
	return *ret0, nil
}

// IERC165
func (e *Erc721Contract) SupportsInterface(interfaceId [4]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := []interface{}{ret0}
	err := e.contract.Call(e.opts, &out, "supportsInterface", interfaceId)
	if err != nil {
		return false, log.WithError(err)
	}
	return *ret0, nil
}

// IERC721
func (e *Erc721Contract) BalanceOf(owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := []interface{}{ret0}
	err := e.contract.Call(e.opts, &out, "balanceOf", owner)
	if err != nil {
		return nil, log.WithError(err)
	}
	return *ret0, nil
}

func (e *Erc721Contract) OwnerOf(tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := []interface{}{ret0}
	err := e.contract.Call(e.opts, &out, "ownerOf", tokenId)
	if err != nil {
		return *ret0, log.WithError(err)
	}
	return *ret0, nil
}

func (e *Erc721Contract) GetApproved(tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := []interface{}{ret0}
	err := e.contract.Call(e.opts, &out, "getApproved", tokenId)
	if err != nil {
		return *ret0, log.WithError(err)
	}
	return *ret0, nil
}

func (e *Erc721Contract) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := []interface{}{ret0}
	err := e.contract.Call(e.opts, &out, "isApprovedForAll", owner, operator)
	if err != nil {
		return false, log.WithError(err)
	}
	return *ret0, nil
}

func (e *Erc721Contract) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return e.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

func (e *Erc721Contract) SafeTransferFrom2(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return e.contract.Transact(opts, "safeTransferFrom", from, to, tokenId, data)
}

func (e *Erc721Contract) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return e.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

func (e *Erc721Contract) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return e.contract.Transact(opts, "approve", to, tokenId)
}

func (e *Erc721Contract) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return e.contract.Transact(opts, "setApprovalForAll", operator, approved)
}
