package main

import (
	"context"
	bind2 "github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"io"
	"os"
	"path"
	"strings"
)

var key = "{\"address\":\"85755d82a3adc23598b887a6c33a2508b4b71e5a\",\"crypto\":{\"cipher\":\"aes-128-ctr\",\"ciphertext\":\"9631b0f23356c128bb8116207fda102951ba802da2a4e6f4cc1fc4f80c5e424a\",\"cipherparams\":{\"iv\":\"d3c8ceabd915abd50e85deafe15979aa\"},\"kdf\":\"scrypt\",\"kdfparams\":{\"dklen\":32,\"n\":262144,\"p\":1,\"r\":8,\"salt\":\"217cae8c2d765f3d6bfdbab3a8f145364fd2adf10183e9fbaea4b3ffb00404d9\"},\"mac\":\"4ccf3a9847f41af5ae15d0c67a3929e3375476f8b1982ff24041c810d7aa76ab\"},\"id\":\"ee4c195b-c2ea-4951-9f33-28502fa85734\",\"version\":3}"
var passphrase = "asdfasdfasdf"

func main() {
	// create a logger to output errors
	log.SetDefault(log.NewLogger(log.NewTerminalHandler(os.Stdout, true)))

	cwd, _ := os.Getwd()
	// Create an IPC-based RPC connection to a remote node and an authorized transactor
	conn, err := ethclient.Dial(path.Join(path.Join(cwd, "datadir/geth.ipc")))
	if err != nil {
		log.Crit("Failed to connect to the Ethereum client", "err", err)
	}
	// Retrieve the current chain ID
	chainID, err := conn.ChainID(context.Background())
	if err != nil {
		log.Crit("Failed to retrieve chain ID", "error", err)
	}

	// create auth for tx signing from the key on disk
	json, err := io.ReadAll(strings.NewReader(key))
	if err != nil {
		log.Crit("failed to read key", "error", err)
	}
	key, err := keystore.DecryptKey(json, passphrase)
	if err != nil {
		log.Crit("failed to decrypt key", "error", err)
	}
	auth := bind2.NewKeyedTransactor(key.PrivateKey, chainID)

	// set up params to deploy an instance of the Storage contract
	deployParams := bind2.DeploymentParams{
		Contracts: []*bind2.MetaData{&StorageMetaData},
	}

	// use the default deployer: it simply creates, signs and submits the deployment transactions
	deployer := bind2.DefaultDeployer(auth, conn)

	// create and submit the contract deployment
	deployRes, err := bind2.LinkAndDeploy(&deployParams, deployer)
	if err != nil {
		log.Crit("error submitting contract", "error", err)
	}

	address, tx := deployRes.Addrs[StorageMetaData.ID], deployRes.Txs[StorageMetaData.ID]

	// create a BoundContract instance to interact with the pending contract
	storageABI, _ := StorageMetaData.ParseABI()
	contract := Storage{*storageABI}
	instance := contract.Instance(conn, address)

	// perform an eth_call on the pending contract
	val, err := bind2.Call(instance, &bind2.CallOpts{Pending: true}, contract.PackRetrieve(), contract.UnpackRetrieve)
	if err != nil {
		log.Crit("call returned error", "error", err)
	}
	log.Info("call to method retrieve returned result", "value", val)
	// wait for the pending contract to be deployed on-chain
	addr, err := bind2.WaitDeployed(context.Background(), conn, tx.Hash())
	if err != nil {
		log.Crit("failed waiting for contract deployment", "error", err)
	}
	log.Info("contract deployed", "address", addr)
}
