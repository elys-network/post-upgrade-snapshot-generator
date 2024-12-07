package utils

import (
	"log"
	"os/exec"
	"strings"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

// CompressData compresses the data directory from the given home path to the specified output path
// The compression format is determined by the output file extension (.tar.lz4, .tar.gz, or .tar)
func CompressData(homePath string, outputPath string) {
	var cmdString string

	// Check the output file type and construct the command accordingly
	if strings.HasSuffix(outputPath, ".tar.lz4") {
		cmdString = "tar -c -C " + homePath + " data | lz4 -c > " + outputPath
	} else if strings.HasSuffix(outputPath, ".tar.gz") {
		cmdString = "tar -czf " + outputPath + " -C " + homePath + " data"
	} else if strings.HasSuffix(outputPath, ".tar") {
		cmdString = "tar -cf " + outputPath + " -C " + homePath + " data"
	} else {
		log.Fatalf(types.ColorRed+"Invalid output format. Supported formats are: .tar.lz4, .tar.gz, .tar. Got: %s", outputPath)
	}

	// Print cmdString
	log.Printf(types.ColorGreen+"Compressing data directory using command: %s", cmdString)

	// Execute the command using /bin/sh
	cmd := exec.Command("/bin/sh", "-c", cmdString)
	if err := cmd.Run(); err != nil {
		log.Fatalf(types.ColorRed+"Command execution failed: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(types.ColorYellow+"Data directory compressed successfully to: %s", outputPath)
}
