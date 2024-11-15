package utils

import (
	"log"
	"os/exec"
	"strings"
)

func Export(cmdPath, homePath, genesisFilePath string) {
	// Define modules in a slice
	modules := []string{
		"amm",
		"assetprofile",
		"auth",
		"authz",
		"bank",
		"burner",
		"capability",
		"commitment", // FIXME: optimize data prior to export
		"consensus",
		"crisis",
		// "distribution", // FIXME: optimize data prior to export as it reached 1.8GB
		"epochs",
		"estaking",
		"evidence",
		"feegrant",
		"genutil",
		// "gov", // FIXME: gov proposals broken
		"group",
		"ibc",
		"interchainaccounts",
		"leveragelp",
		// "masterchef", // FIXME: disable temporarily to attempt to fix the snap gen issue
		"perpetual",
		"oracle",
		"parameter",
		"params",
		"poolaccounted",
		"stablestake",
		"staking",
		"tier",
		"tokenomics",
		"tradeshield",
		"transfer",
		"transferhook",
		"upgrade",
		"vesting",
	}

	// Combine the modules into a comma-separated string
	modulesStr := strings.Join(modules, ",")

	// Command and arguments
	args := []string{"export", "--home", homePath, "--output-document", genesisFilePath, "--modules-to-export", modulesStr}

	// Execute the command and capture the output
	cmd := exec.Command(cmdPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Command execution failed: %v\nOutput: %s", err, out)
	}

	log.Printf("Output successfully written to %s", genesisFilePath)
}
