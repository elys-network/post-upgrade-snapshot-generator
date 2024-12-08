#!/bin/bash

set -e  # Exit on any error

# set variables
CHAINID="provider-1"
MONIKER="validator"
DENOM="uatom"

# Verify binary exists
if ! command -v gaiad &> /dev/null; then
    echo "gaiad binary not found after installation"
    exit 1
fi

# remove existing home path
rm -rf $HOME/.gaia

# init the new node and redirect all output to /dev/null
gaiad init $MONIKER --chain-id $CHAINID > /dev/null 2>&1

# update config files and genesis
config_toml="$HOME/.gaia/config/config.toml"
client_toml="$HOME/.gaia/config/client.toml"
app_toml="$HOME/.gaia/config/app.toml"
genesis_json="$HOME/.gaia/config/genesis.json"

sed -i -E "s|proxy_app = \".*\"|proxy_app = \"tcp://127.0.0.1:36658\"|g" $config_toml
sed -i -E "s|laddr = \"tcp://127.0.0.1:26657\"|laddr = \"tcp://127.0.0.1:36657\"|g" $config_toml
sed -i -E "s|laddr = \"tcp://0.0.0.0:26656\"|laddr = \"tcp://127.0.0.1:36656\"|g" $config_toml
sed -i -E "s|pprof_laddr = \"localhost:6060\"|pprof_laddr = \"localhost:36060\"|g" $config_toml

sed -i -E "s|minimum-gas-prices = \".*\"|minimum-gas-prices = \"0$DENOM\"|g" $app_toml
sed -i -E "s|address = \"tcp://localhost:1317\"|address = \"tcp://localhost:31317\"|g" $app_toml
sed -i -E "s|address = \"localhost:9090\"|address = \"localhost:39090\"|g" $app_toml

sed -i -E "s|chain-id = \".*\"|chain-id = \"$CHAINID\"|g" $client_toml
sed -i -E "s|keyring-backend = \"os\"|keyring-backend = \"test\"|g" $client_toml
sed -i -E "s|node = \".*\"|node = \"tcp://localhost:36657\"|g" $client_toml

sed -i -E "s|\"stake\"|\"${DENOM}\"|g" $genesis_json

# add key
gaiad keys add validator > /dev/null 2>&1

# add genesis account
gaiad genesis add-genesis-account validator 100000000000000$DENOM > /dev/null 2>&1

# generate genesis tx
gaiad genesis gentx validator 100000000$DENOM > /dev/null 2>&1

# collect gentxs
gaiad genesis collect-gentxs > /dev/null 2>&1

# validate genesis
gaiad genesis validate-genesis > /dev/null 2>&1
