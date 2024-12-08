package utils

import (
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func GetValidatorConsensusAddress(cmdPath string) string {
	// retrieve cons address:
	cmd := exec.Command(cmdPath, "cometbft", "show-validator")
	validatorPubkey, err := cmd.Output()
	if err != nil {
		log.Fatalf(types.ColorRed+"Error getting validator pubkey: %v", err)
	}

	cmd = exec.Command(cmdPath, "debug", "pubkey", strings.TrimSpace(string(validatorPubkey)))
	pubkeyOutput, err := cmd.Output()
	if err != nil {
		log.Fatalf(types.ColorRed+"Error getting validator address: %v", err)
	}

	// Parse the address from pubkey output
	re := regexp.MustCompile(`Address:\s+(\S+)`)
	matches := re.FindStringSubmatch(string(pubkeyOutput))
	if len(matches) < 2 {
		log.Fatalf(types.ColorRed + "Could not find validator address in output")
	}
	validatorAddress := matches[1]

	// Get consensus address
	cmd = exec.Command(cmdPath, "debug", "addr", validatorAddress)
	addrOutput, err := cmd.Output()
	if err != nil {
		log.Fatalf(types.ColorRed+"Error getting consensus address: %v", err)
	}

	re = regexp.MustCompile(`Bech32 Con:\s+(\S+)`)
	matches = re.FindStringSubmatch(string(addrOutput))
	if len(matches) < 2 {
		log.Fatalf(types.ColorRed + "Could not find consensus address in output")
	}
	return matches[1]
}
