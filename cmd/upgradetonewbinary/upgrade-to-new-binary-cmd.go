package upgradetonewbinary

import (
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/elys-network/post-upgrade-snapshot-generator/flags"
	"github.com/elys-network/post-upgrade-snapshot-generator/types"
	"github.com/elys-network/post-upgrade-snapshot-generator/utils"
)

func UpgradeToNewBinaryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade-to-new-binary [new_binary_url] [flags]",
		Short: "Upgrade to the new binary",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// get args
			newBinaryUrl := args[0]
			if newBinaryUrl == "" {
				log.Fatalf(types.ColorRed + "new binary url is required")
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

			keyringBackend, _ := cmd.Flags().GetString(flags.FlagKeyringBackend)
			if keyringBackend == "" {
				log.Fatalf(types.ColorRed + "keyring backend is required")
			}

			broadcastMode, _ := cmd.Flags().GetString(flags.FlagBroadcastMode)
			if broadcastMode == "" {
				log.Fatalf(types.ColorRed + "broadcast mode is required")
			}

			// node 1
			homePath, _ := cmd.Flags().GetString(flags.FlagHome)
			if homePath == "" {
				log.Fatalf(types.ColorRed + "home path is required")
			}

			moniker, _ := cmd.Flags().GetString(flags.FlagMoniker)
			if moniker == "" {
				log.Fatalf(types.ColorRed + "moniker is required")
			}

			rpc, _ := cmd.Flags().GetString(flags.FlagRpc)
			if rpc == "" {
				log.Fatalf(types.ColorRed + "rpc is required")
			}

			p2p, _ := cmd.Flags().GetString(flags.FlagP2p)
			if p2p == "" {
				log.Fatalf(types.ColorRed + "p2p is required")
			}

			pprof, _ := cmd.Flags().GetString(flags.FlagPprof)
			if pprof == "" {
				log.Fatalf(types.ColorRed + "pprof is required")
			}

			api, _ := cmd.Flags().GetString(flags.FlagApi)
			if api == "" {
				log.Fatalf(types.ColorRed + "api is required")
			}

			// node 2
			homePath2, _ := cmd.Flags().GetString(flags.FlagHome2)
			if homePath2 == "" {
				log.Fatalf(types.ColorRed + "home path 2 is required")
			}

			moniker2, _ := cmd.Flags().GetString(flags.FlagMoniker2)
			if moniker2 == "" {
				log.Fatalf(types.ColorRed + "moniker 2 is required")
			}

			validatorKeyName2, _ := cmd.Flags().GetString(flags.FlagValidatorKeyName2)
			if validatorKeyName2 == "" {
				log.Fatalf(types.ColorRed + "validator key name 2 is required")
			}

			validatorSelfDelegation2, _ := cmd.Flags().GetString(flags.FlagValidatorSelfDelegation2)
			if validatorSelfDelegation2 == "" {
				log.Fatalf(types.ColorRed + "validator self delegation 2 is required")
			}

			rpc2, _ := cmd.Flags().GetString(flags.FlagRpc2)
			if rpc2 == "" {
				log.Fatalf(types.ColorRed + "rpc 2 is required")
			}

			p2p2, _ := cmd.Flags().GetString(flags.FlagP2p2)
			if p2p2 == "" {
				log.Fatalf(types.ColorRed + "p2p 2 is required")
			}

			pprof2, _ := cmd.Flags().GetString(flags.FlagPprof2)
			if pprof2 == "" {
				log.Fatalf(types.ColorRed + "pprof 2 is required")
			}

			api2, _ := cmd.Flags().GetString(flags.FlagApi2)
			if api2 == "" {
				log.Fatalf(types.ColorRed + "api 2 is required")
			}

			timeOutWaitForNode, err := cmd.Flags().GetInt(flags.FlagTimeOutToWaitForNode)

			if err != nil {
				log.Fatalf(types.ColorRed + err.Error())
			}

			if timeOutWaitForNode == 0 {
				log.Fatalf(types.ColorRed + "time out to wait for service is required")
			}

			timeOutToWaitForNextBlock, err := cmd.Flags().GetInt(flags.FlagTimeOutNextBlock)

			if err != nil {
				log.Fatalf(types.ColorRed + err.Error())
			}

			if timeOutToWaitForNextBlock == 0 {
				log.Fatalf(types.ColorRed + "time out next block is required")
			}

			timeOutForNextBlock := time.Duration(timeOutToWaitForNextBlock) * time.Minute

			skipUnbondValidator, _ := cmd.Flags().GetBool(flags.FlagSkipUnbondValidator)
			if skipUnbondValidator {
				log.Printf(types.ColorYellow + "skipping unbond validator")
			}

			// download and run new binary
			newBinaryPath, newVersion, err := utils.DownloadAndRunVersion(newBinaryUrl, skipBinary)
			if err != nil {
				log.Fatalf(types.ColorRed+"Error downloading and running new binary: %v", err)
			}

			// print new binary path and version
			log.Printf(types.ColorGreen+"New binary path: %v and version: %v", newBinaryPath, newVersion)

			// wait 5 seconds
			time.Sleep(5 * time.Second)

			// start new binary
			newBinaryCmd := utils.Start(newBinaryPath, homePath, rpc, p2p, pprof, api, moniker, "\033[32m", "\033[31m")
			newBinaryCmd2 := utils.Start(newBinaryPath, homePath2, rpc2, p2p2, pprof2, api2, moniker2, "\033[32m", "\033[31m")

			// wait for node to start
			utils.WaitForServiceToStart(rpc, moniker, timeOutWaitForNode)
			utils.WaitForServiceToStart(rpc2, moniker2, timeOutWaitForNode)

			// wait for next block
			utils.WaitForNextBlock(newBinaryPath, rpc, moniker, timeOutForNextBlock)
			utils.WaitForNextBlock(newBinaryPath, rpc2, moniker2, timeOutForNextBlock)

			// check if the upgrade was successful
			utils.QueryUpgradeApplied(newBinaryPath, rpc, newVersion)
			utils.QueryUpgradeApplied(newBinaryPath, rpc2, newVersion)

			if !skipUnbondValidator {
				operatorAddress2 := utils.QueryOperatorAddress(newBinaryPath, homePath, keyringBackend, validatorKeyName2)

				// print operator address 2
				log.Printf(types.ColorGreen+"Operator address 2: %v", operatorAddress2)

				// unbound the second validator power
				utils.UnbondValidator(newBinaryPath, validatorKeyName2, operatorAddress2, validatorSelfDelegation2, keyringBackend, chainId, rpc, broadcastMode, homePath)

				// wait for next block
				utils.WaitForNextBlock(newBinaryPath, rpc, moniker, timeOutForNextBlock)
			}

			// stop new binaries
			utils.Stop(newBinaryCmd, newBinaryCmd2)
		},
	}
}
