package utils

import (
	"log"
	"os/exec"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func RestoreGenesisInitFile(homePath string) {
	// Copy genesis_init.json to genesis.json
	args := []string{
		homePath + "/config/genesis_init.json",
		homePath + "/config/genesis.json",
	}

	if err := exec.Command("cp", args...).Run(); err != nil {
		log.Fatalf(types.ColorRed+"Failed to copy genesis_init.json to genesis.json: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(types.ColorYellow + "Genesis file copied to genesis.json")
}
