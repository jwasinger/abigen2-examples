package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	"io"
	"os"
	"path"
	"strings"
)

var key = "{\"address\":\"85755d82a3adc23598b887a6c33a2508b4b71e5a\",\"crypto\":{\"cipher\":\"aes-128-ctr\",\"ciphertext\":\"9631b0f23356c128bb8116207fda102951ba802da2a4e6f4cc1fc4f80c5e424a\",\"cipherparams\":{\"iv\":\"d3c8ceabd915abd50e85deafe15979aa\"},\"kdf\":\"scrypt\",\"kdfparams\":{\"dklen\":32,\"n\":262144,\"p\":1,\"r\":8,\"salt\":\"217cae8c2d765f3d6bfdbab3a8f145364fd2adf10183e9fbaea4b3ffb00404d9\"},\"mac\":\"4ccf3a9847f41af5ae15d0c67a3929e3375476f8b1982ff24041c810d7aa76ab\"},\"id\":\"ee4c195b-c2ea-4951-9f33-28502fa85734\",\"version\":3}"
var passphrase = "asdfasdfasdf"

func main() {
	cwd, _ := os.Getwd()
	// Create an IPC-based RPC connection to a remote node and an authorized transactor
	conn, err := ethclient.Dial(path.Join(path.Join(cwd, "datadir/geth.ipc")))
	if err != nil {
		panic(fmt.Errorf("Failed to connect to the Ethereum client: %v", err))
	}
	// Retrieve the current chain ID
	chainID, err := conn.ChainID(context.Background())
	if err != nil {
		panic(fmt.Errorf("Failed to retrieve chain ID: %v", err))
	}

	// create auth for tx signing from the key on disk
	json, err := io.ReadAll(strings.NewReader(key))
	if err != nil {
		panic(fmt.Errorf("failed to read key: %v", err))
	}
	key, err := keystore.DecryptKey(json, passphrase)
	if err != nil {
		panic(fmt.Errorf("failed to decrypt key: %v", err))
	}
	auth := bind.NewKeyedTransactor(key.PrivateKey, chainID)

	// set up params to deploy an instance of the Storage contract
	deployParams := bind.DeploymentParams{
		Contracts: []*bind.MetaData{&StorageMetaData},
	}

	// use the default deployer: it simply creates, signs and submits the deployment transactions
	deployer := bind.DefaultDeployer(auth, conn)

	// create and submit the contract deployment
	deployRes, err := bind.LinkAndDeploy(&deployParams, deployer)
	if err != nil {
		panic(fmt.Errorf("error submitting contract: %v", err))
	}

	address, tx := deployRes.Addresses[StorageMetaData.ID], deployRes.Txs[StorageMetaData.ID]

	fmt.Printf("contract pending deploy: 0x%x\n", address)
	fmt.Printf("transaction waiting to be mined: 0x%x\n", tx.Hash())

	// create a BoundContract instance to interact with the pending contract
	storageABI, _ := StorageMetaData.ParseABI()
	contract := Storage{*storageABI}
	instance := contract.Instance(conn, address)

	// perform an eth_call on the pending contract
	val, err := bind.Call(instance, &bind.CallOpts{Pending: true}, contract.PackRetrieve(), contract.UnpackRetrieve)
	if err != nil {
		panic(fmt.Errorf("call returned error: %v", err))
	}
	fmt.Printf("call to method retrieve returned result: %d\n", val)

	// wait for the pending contract to be deployed on-chain
	if _, err := bind.WaitDeployed(context.Background(), conn, tx.Hash()); err != nil {
		panic(fmt.Errorf("failed waiting for contract deployment: %v", err))
	}
	fmt.Println("contract deployed successfully")
}
