package utils

import (
	"encoding/json"
	"os/exec"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func QueryBlockHeight(cmdPath, node string) (string, error) {
	// Command and arguments
	args := []string{"status", "--node", node}

	// Execute the command
	output, err := exec.Command(cmdPath, args...).CombinedOutput()
	if err != nil {
		return "-1", err
	}

	// Unmarshal the JSON output
	var statusOutput types.StatusOutput
	if err := json.Unmarshal(output, &statusOutput); err != nil {
		return "-1", err
	}

	return statusOutput.SyncInfo.LatestBlockHeight, nil
}
