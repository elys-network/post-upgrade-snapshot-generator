package utils

import (
	"log"
	"os/exec"
	"strings"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func RetrieveSnapshot(snapshotUrl, homePath string) {
	var cmdString string
	isUrl := strings.HasPrefix(snapshotUrl, "http://") || strings.HasPrefix(snapshotUrl, "https://")

	// Check the file type and construct the command accordingly
	if strings.HasSuffix(snapshotUrl, ".tar.lz4") {
		if isUrl {
			cmdString = "curl -o - -L " + snapshotUrl + " | lz4 -c -d - | tar -x -C " + homePath
		} else {
			cmdString = "lz4 -c -d " + snapshotUrl + " | tar -x -C " + homePath
		}
	} else if strings.HasSuffix(snapshotUrl, ".tar.gz") {
		if isUrl {
			cmdString = "curl -o - -L " + snapshotUrl + " | tar -xz -C " + homePath
		} else {
			cmdString = "tar -xz -f " + snapshotUrl + " -C " + homePath
		}
	} else if strings.HasSuffix(snapshotUrl, ".tar") {
		if isUrl {
			cmdString = "curl -o - -L " + snapshotUrl + " | tar -x -C " + homePath
		} else {
			cmdString = "tar -x -f " + snapshotUrl + " -C " + homePath
		}
	} else {
		log.Fatalf(types.ColorRed+"Invalid snapshot url or path: %s", snapshotUrl)
	}

	// Print cmdString
	log.Printf(types.ColorGreen+"Retrieving snapshot using command: %s", cmdString)

	// Execute the command using /bin/sh
	cmd := exec.Command("/bin/sh", "-c", cmdString)
	if err := cmd.Run(); err != nil {
		log.Fatalf(types.ColorRed+"Command execution failed: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(types.ColorYellow+"Snapshot retrieved and extracted to path: %s", homePath)
}
