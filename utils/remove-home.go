package utils

import (
	"log"
	"os/exec"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func RemoveHome(homePath string) {
	// Command and arguments
	args := []string{"-rf", homePath}

	// Execute the command
	if err := exec.Command("rm", args...).Run(); err != nil {
		log.Fatalf(types.ColorRed+"Command execution failed: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(types.ColorYellow+"removed home path %s successfully", homePath)
}
