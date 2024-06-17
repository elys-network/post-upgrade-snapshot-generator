package utils

import (
	"log"
	"os/exec"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func Sed(pattern, file string) {
	// Update config.toml for cors_allowed_origins
	var args []string

	if IsLinux() {
		args = []string{"-i", pattern, file}
	} else {
		args = []string{"-i", "", pattern, file}
	}

	// Execute the sed command
	if err := exec.Command("sed", args...).Run(); err != nil {
		log.Fatalf(types.ColorRed+"Error updating "+file+": %v\n", err)
	}
}
