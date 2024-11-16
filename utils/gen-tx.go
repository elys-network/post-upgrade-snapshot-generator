package utils

import (
	"log"
	"os/exec"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func GenTx(cmdPath, name, amount, chainId, homePath, keyringBackend string) {
	// Command and arguments
	args := []string{"genesis", "gentx", name, amount + "uelys", "--chain-id", chainId, "--home", homePath, "--keyring-backend", keyringBackend, "--gas", "1000000"}

	// Execute the command
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf(types.ColorRed+"Command execution failed: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(types.ColorYellow+"gen tx with name %s, amount: %s, chain id %s, home path %s and keyring backend %s successfully", name, amount, chainId, homePath, keyringBackend)
}
