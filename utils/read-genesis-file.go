package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func ReadGenesisFile(filePath string) (types.Genesis, error) {
	var genesis types.Genesis
	file, err := os.Open(filePath)
	if err != nil {
		return genesis, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	if err := json.NewDecoder(bufio.NewReader(file)).Decode(&genesis); err != nil {
		return genesis, fmt.Errorf("error decoding JSON: %w", err)
	}

	return genesis, nil
}
