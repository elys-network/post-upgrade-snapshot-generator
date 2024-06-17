package utils

import (
	"log"
	"os/exec"
	"strings"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func QueryNodeId(cmdPath, home string) string {
	// Command and arguments
	args := []string{"tendermint", "show-node-id", "--home", home}

	// Execute the command
	output, err := exec.Command(cmdPath, args...).CombinedOutput()
	if err != nil {
		log.Fatalf(types.ColorRed+"Failed to query node id: %v", err)
	}

	// trim the output
	outputStr := strings.TrimSpace(string(output))

	return outputStr
}
