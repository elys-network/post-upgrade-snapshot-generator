package createsnapshot

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/elys-network/post-upgrade-snapshot-generator/flags"
	"github.com/elys-network/post-upgrade-snapshot-generator/types"
	"github.com/elys-network/post-upgrade-snapshot-generator/utils"
)

// CreateSnapshotCmd creates a compressed snapshot of the chain data
func CreateSnapshotCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-snapshot [output_path]",
		Short: "Create a compressed snapshot of the chain data",
		Long: `Create a compressed snapshot of the chain data directory.
Supported output formats are: .tar.lz4, .tar.gz, and .tar
Example: create-snapshot snapshot.tar.lz4`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			outputPath := args[0]
			if outputPath == "" {
				log.Fatalf(types.ColorRed + "output path is required")
			}

			// Get home path from flags
			homePath, _ := cmd.Flags().GetString(flags.FlagHome)
			if homePath == "" {
				log.Fatalf(types.ColorRed + "home path is required")
			}

			// Create the snapshot
			utils.CompressData(homePath, outputPath)
		},
	}
}
