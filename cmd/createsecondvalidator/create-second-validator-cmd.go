package createsecondvalidator

import (
	"log"
	"time"

	"github.com/spf13/cobra"

	elysdcmd "github.com/elys-network/elys/cmd/elysd/cmd"
	"github.com/elys-network/post-upgrade-snapshot-generator/flags"
	"github.com/elys-network/post-upgrade-snapshot-generator/types"
	"github.com/elys-network/post-upgrade-snapshot-generator/utils"
)

// create-second-validator command does the following:
// 1. download and run old binary
// 2. start node 1
// 3. wait for rpc to start
// 4. wait for next block
// 5. create validator node 2
// 6. wait for next block
// 7. stop old binary
func CreateSecondValidatorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-second-validator [old_binary_url] [flags]",
		Short: "Create a second validator",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// set prefix
			elysdcmd.InitSDKConfig()

			// get args
			oldBinaryUrl := args[0]
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

			timeOutWaitForNode, err := cmd.Flags().GetInt(flags.FlagTimeOutToWaitForNode)

			if err != nil {
				log.Fatalf("%sError retrieving time out wait for node parameter value: %v", types.ColorRed, err)
			}

			if timeOutWaitForNode == 0 {
				log.Fatalf(types.ColorRed + "time out to wait for service is required")
			}

			timeOutToWaitForNextBlock, err := cmd.Flags().GetInt(flags.FlagTimeOutNextBlock)

			if err != nil {
				log.Fatalf("%sError retrieving time out next block parameter value: %v", types.ColorRed, err)
			}

			if timeOutToWaitForNextBlock == 0 {
				log.Fatalf(types.ColorRed + "time out next block is required")
			}

			timeOutForNextBlock := time.Duration(timeOutToWaitForNextBlock) * time.Minute

			skipProposal, _ := cmd.Flags().GetBool(flags.FlagSkipProposal)
			if skipProposal {
				log.Printf(types.ColorYellow + "skipping proposal")
			}

			// download and run old binary
			oldBinaryPath, oldVersion, err := utils.DownloadAndRunVersion(oldBinaryUrl, skipBinary)
			if err != nil {
				log.Fatalf(types.ColorRed+"Error downloading and running old binary: %v", err)
			}

			// print old binary path and version
			log.Printf(types.ColorGreen+"Old binary path: %v and version: %v", oldBinaryPath, oldVersion)

			// prepare price feeder flags if enabled
			var startArgs []string
			priceFeederEnable, _ := cmd.Flags().GetBool(flags.FlagPriceFeederEnable)
			priceFeederConfigPath, _ := cmd.Flags().GetString(flags.FlagPriceFeederConfigPath)
			priceFeederLogLevel, _ := cmd.Flags().GetString(flags.FlagPriceFeederLogLevel)
			startArgs = []string{
				"--pricefeeder.config_path=" + priceFeederConfigPath,
				"--pricefeeder.log_level=" + priceFeederLogLevel,
			}

			if priceFeederEnable {
				startArgs = append(startArgs, "--pricefeeder.enable=true")
			}

			// start node 1
			oldBinaryCmd := utils.Start(oldBinaryPath, homePath, rpc, p2p, pprof, api, moniker, types.ColorGreen, types.ColorRed, startArgs...)

			// wait for rpc to start
			utils.WaitForServiceToStart(rpc, moniker, timeOutWaitForNode)

			// wait for next block
			utils.WaitForNextBlock(oldBinaryPath, rpc, moniker, timeOutForNextBlock)

			if skipProposal {
				// listen for signals
				utils.ListenForSignals(oldBinaryCmd)
				return
			}

			// query validator pubkey
			validatorPubkey2 := utils.QueryValidatorPubkey(oldBinaryPath, homePath2)

			// create validator node 2
			utils.CreateValidator(oldBinaryPath, validatorKeyName2, validatorSelfDelegation2, moniker2, validatorPubkey2, homePath, keyringBackend, chainId, rpc, broadcastMode)

			// wait for next block
			utils.WaitForNextBlock(oldBinaryPath, rpc, moniker, timeOutForNextBlock)

			// stop old binary
			utils.Stop(oldBinaryCmd)
		},
	}
}
