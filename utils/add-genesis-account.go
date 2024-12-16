package utils

import (
	"log"
	"os/exec"
	"strings"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func AddGenesisAccount(cmdPath, address, homePath string, balances []string) {
	// Build the balance string
	balanceStr := strings.Join(balances, ",")

	// Command and arguments
	args := []string{
		"genesis",
		"add-genesis-account",
		address,
		balanceStr,
		"--home",
		homePath,
	}

	// Execute the command
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf(types.ColorRed+"Command execution failed: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(types.ColorYellow+"add genesis account with address %s, balance: %s and home path %s successfully", address, balanceStr, homePath)
}
