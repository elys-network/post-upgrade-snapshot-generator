package utils

import (
	"log"
	"os/exec"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func BackupGenesisInitFile(homePath string) {
	// Copy genesis.json to genesis_init.json
	args := []string{
		homePath + "/config/genesis.json",
		homePath + "/config/genesis_init.json",
	}

	if err := exec.Command("cp", args...).Run(); err != nil {
		log.Fatalf(types.ColorRed+"Failed to copy genesis.json to genesis_init.json: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(types.ColorYellow + "Genesis file copied to genesis_init.json")
}
