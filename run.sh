#!/bin/bash

SNAPSHOT_PATH=$HOME/data.tar.gz

OLD_BINARY_PATH=$GOBIN/old/elysd
ELYS_PATH=$GOPATH/src/github.com/elys-network/elys/
NEW_BINARY_PATH=$GOBIN/elysd

post-upgrade-snapshot-generator chain-snapshot-export $SNAPSHOT_PATH $OLD_BINARY_PATH --timeout-next-block 100000 --timeout-wait-for-node 100000
post-upgrade-snapshot-generator chain-init $OLD_BINARY_PATH --timeout-next-block 100000 --timeout-wait-for-node 100000
post-upgrade-snapshot-generator create-second-validator $OLD_BINARY_PATH --timeout-next-block 100000 --timeout-wait-for-node 100000

cd $ELYS_PATH || return
git tag -f v999999.999999.999999
make install

post-upgrade-snapshot-generator prepare-validator-data $OLD_BINARY_PATH --timeout-next-block 100000 --timeout-wait-for-node 100000
post-upgrade-snapshot-generator submit-new-proposal $OLD_BINARY_PATH $NEW_BINARY_PATH --timeout-next-block 100000 --timeout-wait-for-node 100000
post-upgrade-snapshot-generator upgrade-to-new-binary $NEW_BINARY_PATH --timeout-next-block 100000 --timeout-wait-for-node 100000
