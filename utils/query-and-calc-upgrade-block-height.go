package utils

import (
	"log"
	"strconv"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func QueryAndCalcUpgradeBlockHeight(cmdPath, node string) string {
	// query block height
	blockHeight, err := QueryBlockHeight(cmdPath, node)
	if err != nil {
		log.Fatalf(types.ColorRed+"Failed to query block height: %v", err)
	}

	// Convert blockHeight from string to int
	blockHeightInt, err := strconv.Atoi(blockHeight)
	if err != nil {
		log.Fatalf(types.ColorRed+"Failed to convert blockHeight to integer: %v", err)
	}

	// set upgrade block height
	upgradeBlockHeight := blockHeightInt + 25

	// return upgrade block height as a string
	return strconv.Itoa(upgradeBlockHeight)
}
