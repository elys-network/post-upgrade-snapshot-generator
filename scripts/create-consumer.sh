#!/bin/bash

set -e  # Exit on any error

# set variables
CHAINID="elysicstestnet-1"
MONIKER="governor"
DENOM="uelys"

# Verify binary exists
if ! command -v elysd &> /dev/null; then
    echo "elysd binary not found after installation"
    exit 1
fi

# remove existing home path
rm -rf $HOME/.elys

# init the new node and redirect all output to /dev/null
elysd init $MONIKER --chain-id $CHAINID > /dev/null 2>&1

# update config files and genesis
config_toml="$HOME/.elys/config/config.toml"
client_toml="$HOME/.elys/config/client.toml"
app_toml="$HOME/.elys/config/app.toml"
genesis_json="$HOME/.elys/config/genesis.json"

sed -i -E "s|proxy_app = \".*\"|proxy_app = \"tcp://127.0.0.1:26658\"|g" $config_toml
sed -i -E "s|laddr = \"tcp://127.0.0.1:26657\"|laddr = \"tcp://127.0.0.1:26657\"|g" $config_toml
sed -i -E "s|pprof_laddr = \"localhost:6060\"|pprof_laddr = \"localhost:6060\"|g" $config_toml

sed -i -E "s|minimum-gas-prices = \".*\"|minimum-gas-prices = \"0$DENOM\"|g" $app_toml
sed -i -E "s|address = \"tcp://localhost:1317\"|address = \"tcp://localhost:1317\"|g" $app_toml
sed -i -E "s|address = \"localhost:9090\"|address = \"localhost:9090\"|g" $app_toml

sed -i -E "s|chain-id = \".*\"|chain-id = \"$CHAINID\"|g" $client_toml
sed -i -E "s|keyring-backend = \"os\"|keyring-backend = \"test\"|g" $client_toml
sed -i -E "s|node = \".*\"|node = \"tcp://localhost:26657\"|g" $client_toml

sed -i -E "s|\"stake\"|\"${DENOM}\"|g" $genesis_json

# replace assetprofile entry_list by the one in entry_list.json
entry_list=$(cat scripts/entry_list.json | tr -d '\n')  # Remove newlines from the JSON
sed -i -E "s|\"entry_list\": \[\]|\"entry_list\": ${entry_list}|g" $genesis_json

# add key
elysd keys add governor > /dev/null 2>&1

# add genesis account
elysd genesis add-genesis-account governor 100000000000000$DENOM > /dev/null 2>&1

# generate genesis tx
elysd genesis gentx governor 100000000$DENOM > /dev/null 2>&1

# collect gentxs
elysd genesis collect-gentxs > /dev/null 2>&1

# validate genesis
elysd genesis validate-genesis > /dev/null 2>&1

# retrieve the old and new address
oldaddress="elys1ed2lkxujcqfckkhfwmyjqwuqp47ve37crctuus"
newaddress=$(elysd keys show governor -a)

# copy genesis.json to the consumer
# cp scripts/elys_genesis.json $genesis_json

# replace old address with the new one
# sed -i -E "s|$oldaddress|$newaddress|g" $genesis_json

# Add validator signing info to genesis
validator_pubkey=$(elysd tendermint show-validator)
validator_address=$(elysd debug pubkey $validator_pubkey | grep "Address:" | cut -d' ' -f2)
validator_cons_addr=$(elysd debug addr $validator_address | grep "Bech32 Con:" | cut -d' ' -f3)
echo '{
  "address": "'$validator_cons_addr'",
  "validator_signing_info": {
    "address": "'$validator_cons_addr'",
    "start_height": "0",
    "index_offset": "0",
    "jailed_until": "1970-01-01T00:00:00Z",
    "tombstoned": false,
    "missed_blocks_counter": "0"
  }
}' > signing_info.json
jq '.app_state.slashing.signing_infos += [input]' $genesis_json signing_info.json > temp.json && mv temp.json $genesis_json
