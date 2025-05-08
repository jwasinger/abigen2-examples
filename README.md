# Abigen V2 Examples

This repository contains reproduceable examples from the [Go Contract Bindings (v2)](https://geth.ethereum.org/docs/developers/dapp-developer/native-bindings-v2) page on the go-ethereum documentation website.


## Instructions

### Compile the Examples

Compile the example contract to ABI definition with deployer bytecode:
```
> solc --combined-json abi,bin Storage.sol > Storage.json 
```


Compile the ABI definition and deployer bytecode to Go contract bindings:
```
~/projects/go-ethereum/build/bin/abigen --v2 --combined-json Storage.json --pkg main > Storage.go
```

Copy `Storage.go` into `example_1`, `example_2`, `example_bc_simulate`.

Build Go executables in the examples' directories.

### Run the Examples

Start a dev-mode Geth instance: 

ensure the `geth` executable is on the path and run `run_geth.sh`

Run the examples from the top-level directory. e.g. `./example_1/example_1`.
