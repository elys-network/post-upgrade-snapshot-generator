#!/bin/bash

set -e  # Exit on any error

# set variables
COSMOSID="provider-1"
ELYSID="elysicstestnet-1"

# Verify binary exists
if ! command -v rly &> /dev/null; then
    echo "rly binary not found after installation"
    exit 1
fi

# remove existing home path
rm -rf $HOME/.relayer

# create relayer
rly config init

# add config to the config.yaml file
echo "global:
    api-listen-addr: :5183
    timeout: 10s
    memo: \"\"
    light-cache-size: 20
    log-level: info
    ics20-memo-limit: 0
    max-receiver-size: 150
chains:
    atom:
        type: cosmos
        value:
            key-directory: /root/.relayer/keys/atom
            key: relayer
            chain-id: $COSMOSID
            rpc-addr: http://127.0.0.1:36657
            account-prefix: cosmos
            keyring-backend: test
            gas-adjustment: 1.5
            gas-prices: 0.01uatom
            min-gas-amount: 0
            max-gas-amount: 0
            debug: false
            timeout: 20s
            block-timeout: \"\"
            output-format: json
            sign-mode: direct
            extra-codecs: []
            coin-type: 118
            signing-algorithm: \"\"
            broadcast-mode: batch
            min-loop-duration: 0s
            extension-options: []
            feegrants: null
    elys:
        type: cosmos
        value:
            key-directory: /root/.relayer/keys/elys
            key: relayer
            chain-id: $ELYSID
            rpc-addr: http://127.0.0.1:26657
            account-prefix: elys
            keyring-backend: test
            gas-adjustment: 1.5
            gas-prices: 0.01uelys
            min-gas-amount: 0
            max-gas-amount: 0
            debug: false
            timeout: 20s
            block-timeout: \"\"
            output-format: json
            sign-mode: direct
            extra-codecs: []
            coin-type: 118
            signing-algorithm: \"\"
            broadcast-mode: batch
            min-loop-duration: 0s
            extension-options: []
            feegrants: null
paths:
    path-ab:
        src:
            chain-id: $COSMOSID
            client-id: 07-tendermint-0
            connection-id: connection-0
        dst:
            chain-id: $ELYSID
            client-id: 07-tendermint-0
            connection-id: connection-0
        src-channel-filter:
            rule: \"\"
            channel-list: []" > $HOME/.relayer/config/config.yaml

# add keys
rly keys add atom relayer > /dev/null 2>&1
rly keys add elys relayer > /dev/null 2>&1
