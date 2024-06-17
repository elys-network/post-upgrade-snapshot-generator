package utils

import (
	"log"
	"os/exec"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func CreateValidator(cmdPath, name, selfDelegation, moniker, pubkey, homePath, keyringBackend, chainId, node, broadcastMode string) {
	// Command and arguments
	args := []string{
		"tx",
		"staking",
		"create-validator",
		"--amount", selfDelegation + "uelys",
		"--pubkey", pubkey,
		"--moniker", moniker,
		"--commission-rate", "0.05",
		"--commission-max-rate", "0.50",
		"--commission-max-change-rate", "0.01",
		"--min-self-delegation", "1",
		"--from", name,
		"--keyring-backend", keyringBackend,
		"--chain-id", chainId,
		"--node", node,
		"--broadcast-mode", broadcastMode,
		"--fees", "100000uelys",
		"--gas", "3000000",
		"--gas-adjustment", "1.5",
		"--home", homePath,
		"--yes",
	}

	// Execute the command
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf(types.ColorRed+"Failed to create validator: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(types.ColorYellow+"Validator %s created successfully", moniker)
}
