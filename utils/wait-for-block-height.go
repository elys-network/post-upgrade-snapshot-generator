package utils

import (
	"log"
	"strconv"
	"time"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func WaitForBlockHeight(cmdPath, node, height string) {
	targetBlockHeight, err := strconv.Atoi(height)
	if err != nil {
		log.Fatalf(types.ColorRed+"Error converting target block height to integer: %v", err)
	}

	// Now, wait for the block height
	for {
		var blockHeightStr string
		blockHeightStr, err = QueryBlockHeight(cmdPath, node)
		if err == nil {
			newBlockHeight, err := strconv.Atoi(blockHeightStr)
			if err == nil && newBlockHeight >= targetBlockHeight {
				break
			}
		}
		log.Println(types.ColorYellow+"Waiting for block height", height, "...")
		time.Sleep(5 * time.Second) // Wait 5 seconds before retrying
	}

	log.Printf(types.ColorYellow+"Block height %d reached", targetBlockHeight)
}
