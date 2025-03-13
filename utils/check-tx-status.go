package utils

import (
	"fmt"
	"github.com/elys-network/post-upgrade-snapshot-generator/types"
	"log"
	"os/exec"
)

func CheckTxStatus(cmdPath, node, txHash string) (bool, error) {
	args := []string{"q", "tx", txHash, "--node", node, "--output", "json"}
	output, err := exec.Command(cmdPath, args...).CombinedOutput()
	log.Printf(types.ColorYellow+"tx %v query response: %s", txHash, string(output))
	if err != nil {
		return false, fmt.Errorf("failed to query tx status: %w", err)
	}
	return true, nil
}
