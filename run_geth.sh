#! /usr/bin/env bash

rm -rf ./datadir

geth --datadir ./datadir init example_data/genesis.json
geth --datadir ./datadir --dev console
