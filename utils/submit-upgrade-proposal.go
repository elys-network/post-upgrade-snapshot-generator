package utils

import (
	"log"
	"os/exec"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func SubmitUpgradeProposal(cmdPath, name, newVersion, upgradeHeight, homePath, keyringBackend, chainId, node, broadcastMode string) string {
	// Command and arguments
	args := []string{
		"software-upgrade-tx",
		newVersion,
		upgradeHeight,
		"10000000uelys",
		newVersion,
		newVersion,
		"false",
		"--from", name,
		"--keyring-backend", keyringBackend,
		"--chain-id", chainId,
		"--node", node,
		"--broadcast-mode", broadcastMode,
		"--fees", "1000000uelys",
		"--gas", "1000000",
		"--home", homePath,
		"--output", "json",
		"--yes",
	}

	// Execute the command
	output, err := exec.Command(cmdPath, args...).CombinedOutput()
	if err != nil {
		log.Fatalf(types.ColorRed+"Command execution failed: %v", err)
	}

	// Parse output to find the transaction hash
	txHash, err := ParseTxHash(output)
	if err != nil {
		log.Fatalf(types.ColorRed+"Failed to parse transaction hash: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(types.ColorYellow+"Submitted upgrade proposal: %s, upgrade block height: %s", newVersion, upgradeHeight)

	// Return the transaction hash
	return txHash
}
