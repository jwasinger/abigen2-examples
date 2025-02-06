package main

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	bind2 "github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"os"
)

var key = "{\"address\":\"85755d82a3adc23598b887a6c33a2508b4b71e5a\",\"crypto\":{\"cipher\":\"aes-128-ctr\",\"ciphertext\":\"9631b0f23356c128bb8116207fda102951ba802da2a4e6f4cc1fc4f80c5e424a\",\"cipherparams\":{\"iv\":\"d3c8ceabd915abd50e85deafe15979aa\"},\"kdf\":\"scrypt\",\"kdfparams\":{\"dklen\":32,\"n\":262144,\"p\":1,\"r\":8,\"salt\":\"217cae8c2d765f3d6bfdbab3a8f145364fd2adf10183e9fbaea4b3ffb00404d9\"},\"mac\":\"4ccf3a9847f41af5ae15d0c67a3929e3375476f8b1982ff24041c810d7aa76ab\"},\"id\":\"ee4c195b-c2ea-4951-9f33-28502fa85734\",\"version\":3}"
var passphrase = "asdfasdfasdf"

func main() {
	// create a logger to output errors
	log.SetDefault(log.NewLogger(log.NewTerminalHandler(os.Stdout, true)))
	key, err := crypto.GenerateKey()
	if err != nil {
		log.Crit("failed to generate key", "error", err)
	}

	// Since we are using a simulated backend, we will get the chain ID
	// from the same place that the simulated backend gets it.
	chainID := params.AllDevChainProtocolChanges.ChainID

	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		log.Crit("failed to create transaction auth", "error", err)
	}

	sim := simulated.NewBackend(map[common.Address]types.Account{
		auth.From: {Balance: big.NewInt(9e18)},
	})

	// set up params to deploy an instance of the Storage contract
	deployParams := bind2.DeploymentParams{
		Contracts: []*bind2.MetaData{&StorageMetaData},
	}

	// use the default deployer: it simply creates, signs and submits the deployment transactions
	deployer := bind2.DefaultDeployer(auth, sim.Client())

	// create and submit the contract deployment
	deployRes, err := bind2.LinkAndDeploy(&deployParams, deployer)
	if err != nil {
		log.Crit("error submitting contract", "error", err)
	}

	address, tx := deployRes.Addresses[StorageMetaData.ID], deployRes.Txs[StorageMetaData.ID]

	sim.Commit()

	// wait for the pending contract to be deployed on-chain
	if _, err := bind2.WaitDeployed(context.Background(), sim.Client(), tx.Hash()); err != nil {
		log.Crit("failed waiting for contract deployment", "error", err)
	}
	log.Info("contract deployed", "address", address)

	// create a BoundContract instance to interact with the pending contract
	storageABI, _ := StorageMetaData.ParseABI()
	contract := Storage{*storageABI}
	instance := contract.Instance(sim.Client(), address)

	// perform an eth_call on the pending contract
	val, err := bind2.Call(instance, &bind2.CallOpts{Pending: true}, contract.PackRetrieve(), contract.UnpackRetrieve)
	if err != nil {
		log.Crit("call returned error", "error", err)
	}
	log.Info("call to method retrieve returned result", "value", val)

}
