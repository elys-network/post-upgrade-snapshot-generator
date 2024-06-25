package utils

import (
	"log"
	"os/exec"
	"time"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func UnbondValidator(cmdPath, validatorKeyName, operatorAddress, validatorSelfDelegation, keyringBackend, chainId, rpc, broadcastMode, homePath string) {
	// Command and arguments
	args := []string{
		"tx",
		"staking",
		"unbond",
		operatorAddress,
		validatorSelfDelegation + "uelys",
		"--from", validatorKeyName,
		"--keyring-backend", keyringBackend,
		"--chain-id", chainId,
		"--node", rpc,
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
	log.Printf(types.ColorYellow+"Unbonded validator: %s, self-delegation: %s", operatorAddress, validatorSelfDelegation)

	WaitForTxConfirmation(cmdPath, rpc, txHash, 5*time.Minute)
}
