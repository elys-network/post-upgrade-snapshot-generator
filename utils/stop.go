package utils

import (
	"log"
	"os/exec"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func Stop(cmds ...*exec.Cmd) {
	for _, cmd := range cmds {
		// Stop the process
		if cmd != nil && cmd.Process != nil {
			err := cmd.Process.Kill()
			if err != nil {
				log.Fatalf(types.ColorRed+"Failed to kill process: %v", err)
			}
			log.Println(types.ColorYellow + "Process killed successfully")
		}
	}
}
