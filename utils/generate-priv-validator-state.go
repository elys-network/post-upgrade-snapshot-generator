package utils

import (
	"log"
	"os/exec"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func GeneratePrivValidatorState(homePath string) {
	// generate priv_validator_state.json with the following content:
	// {
	// 	"height": "0",
	// 	"round": 0,
	// 	"step": 0
	// }

	// Command and arguments
	args := []string{
		"-c",
		"echo",
		"{\"height\": \"0\", \"round\": 0, \"step\": 0}",
		">",
		homePath + "/data/priv_validator_state.json",
	}

	// Execute the command
	if err := exec.Command("sh", args...).Run(); err != nil {
		log.Fatalf(types.ColorRed+"Failed to generate priv_validator_state.json: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(types.ColorYellow + "priv_validator_state.json generated successfully")
}
