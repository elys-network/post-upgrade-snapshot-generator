package utils

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

// generate upgrade version from the current version (v999999.999999.999999 => v999999)
func generateUpgradeVersion(currentVersion string) string {
	parts := strings.Split(currentVersion, ".")
	if len(parts) != 3 {
		panic(fmt.Sprintf("Invalid version format: %s. Expected format: vX.Y.Z", currentVersion))
	}
	majorVersion := strings.TrimPrefix(parts[0], "v")
	return fmt.Sprintf("v%s", majorVersion)
}

func SubmitUpgradeProposal(cmdPath, name, newVersion, upgradeHeight, homePath, keyringBackend, chainId, node, broadcastMode string) string {
	newVersion = generateUpgradeVersion(newVersion)

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
	log.Printf(types.ColorYellow+"args: %v", args)
	output, err := exec.Command(cmdPath, args...).CombinedOutput()
	if err != nil {
		log.Fatalf(types.ColorRed+"Command execution failed: %v", err)
	}

	log.Printf(types.ColorYellow+"Tx Response: %v", string(output))

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
