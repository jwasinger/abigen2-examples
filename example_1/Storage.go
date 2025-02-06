// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// StorageMetaData contains all meta data concerning the Storage contract.
var StorageMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"retrieve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "Storage",
	Bin: "0x6080604052348015600e575f80fd5b506101438061001c5f395ff3fe608060405234801561000f575f80fd5b5060043610610034575f3560e01c80632e64cec1146100385780636057361d14610056575b5f80fd5b610040610072565b60405161004d919061009b565b60405180910390f35b610070600480360381019061006b91906100e2565b61007a565b005b5f8054905090565b805f8190555050565b5f819050919050565b61009581610083565b82525050565b5f6020820190506100ae5f83018461008c565b92915050565b5f80fd5b6100c181610083565b81146100cb575f80fd5b50565b5f813590506100dc816100b8565b92915050565b5f602082840312156100f7576100f66100b4565b5b5f610104848285016100ce565b9150509291505056fea2646970667358221220bb1add3845e8d6556c69512a3badf1160140f85779530246e037c720d649074d64736f6c634300081a0033",
}

// Storage is an auto generated Go binding around an Ethereum contract.
type Storage struct {
	abi abi.ABI
}

// NewStorage creates a new instance of Storage.
func NewStorage() *Storage {
	parsed, err := StorageMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Storage{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Storage) Instance(backend bind.ContractBackend, addr common.Address) bind.BoundContract {
	return bind.NewBoundContract(backend, addr, c.abi)
}

// Retrieve is a free data retrieval call binding the contract method 0x2e64cec1.
//
// Solidity: function retrieve() view returns(uint256)
func (storage *Storage) PackRetrieve() []byte {
	enc, err := storage.abi.Pack("retrieve")
	if err != nil {
		panic(err)
	}
	return enc
}

func (storage *Storage) UnpackRetrieve(data []byte) (*big.Int, error) {
	out, err := storage.abi.Unpack("retrieve", data)

	if err != nil {
		return new(big.Int), err
	}

	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)

	return out0, err

}

// Store is a free data retrieval call binding the contract method 0x6057361d.
//
// Solidity: function store(uint256 number) returns()
func (storage *Storage) PackStore(Number *big.Int) []byte {
	enc, err := storage.abi.Pack("store", Number)
	if err != nil {
		panic(err)
	}
	return enc
}
