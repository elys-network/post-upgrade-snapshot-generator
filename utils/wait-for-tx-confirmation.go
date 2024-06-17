package utils

import (
	"log"
	"time"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func WaitForTxConfirmation(cmdPath, node, txHash string, timeout time.Duration) {
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			log.Fatalf(types.ColorRed + "timeout reached while waiting for tx confirmation")
		}
		success, err := CheckTxStatus(cmdPath, node, txHash)
		if err != nil {
			log.Printf(types.ColorRed+"error checking tx status, retrying in 5 seconds: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if success {
			break
		}
		log.Printf(types.ColorYellow+"waiting for tx confirmation %s", txHash)
		time.Sleep(5 * time.Second)
	}
	log.Printf(types.ColorGreen+"tx %s confirmed", txHash)
}
