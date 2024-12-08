package utils

import (
	"log"
	"os/exec"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func AddGenesisAccount(cmdPath, address, balance, homePath string) {
	// Command and arguments
	args := []string{"genesis", "add-genesis-account", address, balance + "uelys," + balance + "ibc/F082B65C88E4B6D5EF1DB243CDA1D331D002759E938A0F5CD3FFDC5D53B3E349," + balance + "ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9", "--home", homePath}

	// Execute the command
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf(types.ColorRed+"Command execution failed: %v", err) // nolint: goconst
	}

	// If execution reaches here, the command was successful
	log.Printf(types.ColorYellow+"add genesis account with address %s, balance: %s and home path %s successfully", address, balance, homePath)
}
