package utils

import (
	"encoding/json"
	"log"
	"os/exec"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func AddKey(cmdPath, name, mnemonic, homePath, keyringBackend string) string {
	// Prepare the command
	args := []string{"keys", "add", name, "--recover", "--home", homePath, "--keyring-backend", keyringBackend, "--output", "json"}
	cmd := exec.Command(cmdPath, args...)

	// Get the stdin pipe to send the mnemonic
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf(types.ColorRed+"Failed to create stdin pipe: %v", err)
	}

	// Write the mnemonic to the stdin pipe
	go func() {
		defer stdinPipe.Close()
		_, err := stdinPipe.Write([]byte(mnemonic + "\n"))
		if err != nil {
			log.Fatalf(types.ColorRed+"Failed to write mnemonic to stdin: %v", err)
		}
	}()

	// Run the command and wait for it to finish
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf(types.ColorRed+"Command execution failed: %v", err)
	}

	// Unmarshal the JSON output
	var keyOutput types.KeyOutput
	if err := json.Unmarshal(output, &keyOutput); err != nil {
		log.Fatalf(types.ColorRed+"Failed to unmarshal JSON output: %v", err)
	}

	// Log the address
	log.Printf(types.ColorYellow+"Added key with name %s, home path: %s, keyring backend %s and address %s successfully", name, homePath, keyringBackend, keyOutput.Address)

	return keyOutput.Address
}
