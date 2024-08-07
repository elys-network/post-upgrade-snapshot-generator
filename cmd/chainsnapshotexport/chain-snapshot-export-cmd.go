package chainsnapshotexport

import (
	"log"

	"github.com/spf13/cobra"

	elysdcmd "github.com/elys-network/elys/cmd/elysd/cmd"
	"github.com/elys-network/post-upgrade-snapshot-generator/flags"
	"github.com/elys-network/post-upgrade-snapshot-generator/types"
	"github.com/elys-network/post-upgrade-snapshot-generator/utils"
)

func ChainSnapshotExportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "chain-snapshot-export [snapshot_url] [old_binary_url] [flags]",
		Short: "Export the chain snapshot",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// set prefix
			elysdcmd.InitSDKConfig()

			// get args
			snapshotUrl := args[0]
			if snapshotUrl == "" {
				log.Fatalf(types.ColorRed + "snapshot url is required")
			}

			oldBinaryUrl := args[1]
			if oldBinaryUrl == "" {
				log.Fatalf(types.ColorRed + "old binary url is required")
			}

			// get flags values
			skipBinary, _ := cmd.Flags().GetBool(flags.FlagSkipBinary)
			if skipBinary {
				log.Printf(types.ColorYellow + "skipping binary download")
			}

			chainId, _ := cmd.Flags().GetString(flags.FlagChainId)
			if chainId == "" {
				log.Fatalf(types.ColorRed + "chain id is required")
			}

			genesisFilePath, _ := cmd.Flags().GetString(flags.FlagGenesisFilePath)
			if genesisFilePath == "" {
				log.Fatalf(types.ColorRed + "genesis file path is required")
			}

			homePath, _ := cmd.Flags().GetString(flags.FlagHome)
			if homePath == "" {
				log.Fatalf(types.ColorRed + "home path is required")
			}

			moniker, _ := cmd.Flags().GetString(flags.FlagMoniker)
			if moniker == "" {
				log.Fatalf(types.ColorRed + "moniker is required")
			}

			// download and run old binary
			oldBinaryPath, oldVersion, err := utils.DownloadAndRunVersion(oldBinaryUrl, skipBinary)
			if err != nil {
				log.Fatalf(types.ColorRed+"Error downloading and running old binary: %v", err)
			}

			// print old binary path and version
			log.Printf(types.ColorGreen+"Old binary path: %v and version: %v", oldBinaryPath, oldVersion)

			// remove home path
			utils.RemoveHome(homePath)

			// init chain
			utils.InitNode(oldBinaryPath, moniker, chainId, homePath)

			// update config files
			utils.UpdateConfig(homePath, "pebbledb")

			// retrieve the snapshot
			utils.RetrieveSnapshot(snapshotUrl, homePath)

			// export genesis file
			utils.Export(oldBinaryPath, homePath, genesisFilePath)
		},
	}
}
