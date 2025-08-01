package preparevalidatordata

import (
	"log"

	"github.com/spf13/cobra"

	elysdcmd "github.com/elys-network/elys/v6/cmd/elysd/cmd"
	"github.com/elys-network/post-upgrade-snapshot-generator/flags"
	"github.com/elys-network/post-upgrade-snapshot-generator/types"
	"github.com/elys-network/post-upgrade-snapshot-generator/utils"
)

// prepare-validator-data command does the following:
// 1. restore genesis init file
// 2. copy data from node 1 to node 2
// 3. generate priv_validator_state.json file for node 2
func PrepareValidatorDataCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prepare-validator-data",
		Short: "Prepare validator data",
		Run: func(cmd *cobra.Command, args []string) {
			// set prefix
			elysdcmd.InitSDKConfig()

			// node 1
			homePath, _ := cmd.Flags().GetString(flags.FlagHome)
			if homePath == "" {
				log.Fatalf(types.ColorRed + "home path is required")
			}

			// node 2
			homePath2, _ := cmd.Flags().GetString(flags.FlagHome2)
			if homePath2 == "" {
				log.Fatalf(types.ColorRed + "home path 2 is required")
			}

			// restore genesis init file
			utils.RestoreGenesisInitFile(homePath)

			// copy data from node 1 to node 2
			utils.CopyDataFromNodeToNode(homePath, homePath2)

			// generate priv_validator_state.json file for node 2
			utils.GeneratePrivValidatorState(homePath2)
		},
	}
}
