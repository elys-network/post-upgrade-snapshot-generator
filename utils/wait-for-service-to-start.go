package utils

import (
	"log"
	"time"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func WaitForServiceToStart(url, name string, timeoutFlag int) {
	timeout := time.Duration(timeoutFlag) * time.Second
	start := time.Now()

	// Wait for the node to be running with timout
	for !IsServiceRunning(url) {
		if time.Since(start) > timeout {
			log.Fatalf(types.ColorRed + "[" + name + "] Service did not start within the specified timeout")
		}
		log.Println(types.ColorYellow + "[" + name + "] Waiting for service to start...")
		time.Sleep(5 * time.Second)
	}
	log.Println(types.ColorYellow + "[" + name + "] Service is running.")
}
