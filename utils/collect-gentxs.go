package utils

import (
	"log"
	"os/exec"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func CollectGentxs(cmdPath, homePath string) {
	// Command and arguments
	args := []string{"genesis", "collect-gentxs", "--home", homePath}

	// Execute the command
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf(types.ColorRed+"Command execution failed: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(types.ColorYellow+"collect gen txs with home path %s successfully", homePath)
}
