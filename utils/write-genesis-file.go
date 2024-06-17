package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func WriteGenesisFile(filePath string, genesis types.Genesis) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ") // disable for now

	if err := encoder.Encode(genesis); err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	return nil
}
