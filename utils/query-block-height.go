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

	// Try unmarshalling into StatusOutput first
	var statusOutput types.StatusOutput
	if err := json.Unmarshal(output, &statusOutput); err == nil {
		return statusOutput.SyncInfo.LatestBlockHeight, nil
	}

	// If the first unmarshal fails, try unmarshalling into LegacyStatusOutput
	var legacyStatusOutput types.LegacyStatusOutput
	if err := json.Unmarshal(output, &legacyStatusOutput); err == nil {
		return legacyStatusOutput.SyncInfo.LatestBlockHeight, nil
	}

	// If both attempts fail, return an error
	return "-1", err
}
