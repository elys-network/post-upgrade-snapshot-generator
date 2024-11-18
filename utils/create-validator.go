package utils

import (
	"log"
	"os"
	"os/exec"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func CreateValidator(cmdPath, name, selfDelegation, moniker, pubkey, homePath, keyringBackend, chainId, node, broadcastMode string) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "validator.json-*")
	if err != nil {
		return
	}
	tmpFilePath := tmpFile.Name()
	defer os.Remove(tmpFilePath) // Clean up

	// prepare the validator.json file
	validatorJSON := `{
		"pubkey": ` + pubkey + `,
		"amount": "` + selfDelegation + `uelys",
		"moniker": "` + moniker + `",
		"identity": "bob",
		"website": "https://example.com",
		"security": "secury@example.com",
		"details": "details",
		"commission-rate": "0.05",
		"commission-max-rate": "0.50",
		"commission-max-change-rate": "0.01",
		"min-self-delegation": "1"
	}`
	if _, err := tmpFile.Write([]byte(validatorJSON)); err != nil {
		log.Fatalf(types.ColorRed+"Failed to write to file: %v", err)
	}
	tmpFile.Close()

	// Command and arguments
	args := []string{
		"tx",
		"staking",
		"create-validator",
		tmpFilePath,
		"--from", name,
		"--keyring-backend", keyringBackend,
		"--chain-id", chainId,
		"--node", node,
		"--broadcast-mode", broadcastMode,
		"--fees", "1000000uelys",
		"--gas", "3000000",
		"--gas-adjustment", "1.5",
		"--home", homePath,
		"--yes",
	}

	// Execute the command
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf(types.ColorRed+"Failed to create validator: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(types.ColorYellow+"Validator %s created successfully", moniker)
}
