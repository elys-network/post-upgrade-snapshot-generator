#!/bin/bash

# download from https://tools.highstakes.ch/snapshots/elys
url="https://tools.highstakes.ch/snapshots/elys"
html_content=$(curl -s "$url")

SNAPSHOT_DOWNLOAD_URL=$(echo "$html_content" | grep 'class="a_custom"' | sed -n 's/.*href="\([^"]*\)".*/\1/p')
# if ubuntu:
# SNAPSHOT_DOWNLOAD_URL=$(echo "$html_content" | grep -oP '(?<=<a class="a_custom" href=")[^"]*' | tail -n 1)

SNAPSHOT_FILE_PATH=$HOME/snapshot-testnet.tar.gz

curl -L $SNAPSHOT_DOWNLOAD_URL -o $SNAPSHOT_FILE_PATH

# should be current testnet version
OLD_BINARY_PATH=$GOBIN/old/elysd

# should be PR branch or main branch
# make sure to tag them as v999999.999999.999999
NEW_BINARY_PATH=$GOBIN/elysd

post-upgrade-snapshot-generator chain-snapshot-export $SNAPSHOT_FILE_PATH $OLD_BINARY_PATH --timeout-next-block 100000 --timeout-wait-for-node 100000
post-upgrade-snapshot-generator chain-init $OLD_BINARY_PATH --timeout-next-block 100000 --timeout-wait-for-node 100000
post-upgrade-snapshot-generator create-second-validator $OLD_BINARY_PATH --timeout-next-block 100000 --timeout-wait-for-node 100000

post-upgrade-snapshot-generator prepare-validator-data $OLD_BINARY_PATH --timeout-next-block 100000 --timeout-wait-for-node 100000
post-upgrade-snapshot-generator submit-new-proposal $OLD_BINARY_PATH $NEW_BINARY_PATH --timeout-next-block 100000 --timeout-wait-for-node 100000
post-upgrade-snapshot-generator upgrade-to-new-binary $NEW_BINARY_PATH --timeout-next-block 100000 --timeout-wait-for-node 100000
