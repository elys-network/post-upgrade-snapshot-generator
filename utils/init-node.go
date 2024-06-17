package utils

import (
	"log"
	"os/exec"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func InitNode(cmdPath, moniker, chainId, homePath string) {
	// Command and arguments
	args := []string{"init", moniker, "--chain-id", chainId, "--home", homePath}

	// Execute the command
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf(types.ColorRed+"Command execution failed: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(types.ColorYellow+"init node with moniker %s, chain id %s and home path: %s successfully", moniker, chainId, homePath)
}
