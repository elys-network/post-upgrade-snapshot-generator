package submitnewproposal

import (
	"log"
	"time"

	"github.com/spf13/cobra"

	elysdcmd "github.com/elys-network/elys/cmd/elysd/cmd"
	"github.com/elys-network/post-upgrade-snapshot-generator/flags"
	"github.com/elys-network/post-upgrade-snapshot-generator/types"
	"github.com/elys-network/post-upgrade-snapshot-generator/utils"
)

// submit-new-proposal command does the following:
// 1. download and run old binary
// 2. download and run new binary
// 3. start node 1 and 2
// 4. query and calculate upgrade block height
// 5. query next proposal id
// 6. submit upgrade proposal
// 7. vote on upgrade proposal
// 8. wait for upgrade block height
// 9. stop old binaries
func SubmitNewProposalCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "submit-new-proposal [old_binary_url] [new_binary_url] [flags]",
		Short: "Submit a new proposal",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// set prefix
			elysdcmd.InitSDKConfig()

			// get args
			oldBinaryUrl := args[0]
			if oldBinaryUrl == "" {
				log.Fatalf(types.ColorRed + "old binary url is required")
			}

			newBinaryUrl := args[1]
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

			validatorKeyName, _ := cmd.Flags().GetString(flags.FlagValidatorKeyName)
			if validatorKeyName == "" {
				log.Fatalf(types.ColorRed + "validator key name is required")
			}

			validatorSelfDelegation, _ := cmd.Flags().GetString(flags.FlagValidatorSelfDelegation)
			if validatorSelfDelegation == "" {
				log.Fatalf(types.ColorRed + "validator self delegation is required")
			}

			validatorMnemonic, _ := cmd.Flags().GetString(flags.FlagValidatorMnemonic)
			if validatorMnemonic == "" {
				log.Fatalf(types.ColorRed + "validator mnemonic is required")
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

			validatorMnemonic2, _ := cmd.Flags().GetString(flags.FlagValidatorMnemonic2)
			if validatorMnemonic2 == "" {
				log.Fatalf(types.ColorRed + "validator mnemonic 2 is required")
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

			validatorBalances, _ := cmd.Flags().GetStringSlice(flags.FlagValidatorBalances)
			if len(validatorBalances) == 0 {
				log.Fatalf(types.ColorRed + "validator balances are required")
			}

			timeOutWaitForNode, err := cmd.Flags().GetInt(flags.FlagTimeOutToWaitForNode)

			if err != nil {
				log.Fatalf(types.ColorRed + err.Error())
			}

			if timeOutWaitForNode == 0 {
				log.Fatalf(types.ColorRed + "time out to wait for service is required")
			}

			// download and run old binary
			oldBinaryPath, oldVersion, err := utils.DownloadAndRunVersion(oldBinaryUrl, skipBinary)
			if err != nil {
				log.Fatalf(types.ColorRed+"Error downloading and running old binary: %v", err)
			}

			// print old binary path and version
			log.Printf(types.ColorGreen+"Old binary path: %v and version: %v", oldBinaryPath, oldVersion)

			// download and run new binary
			newBinaryPath, newVersion, err := utils.DownloadAndRunVersion(newBinaryUrl, skipBinary)
			if err != nil {
				log.Fatalf(types.ColorRed+"Error downloading and running new binary: %v", err)
			}

			// print new binary path and version
			log.Printf(types.ColorGreen+"New binary path: %v and version: %v", newBinaryPath, newVersion)

			// start node 1 and 2
			oldBinaryCmd := utils.Start(oldBinaryPath, homePath, rpc, p2p, pprof, api, moniker, types.ColorGreen, types.ColorRed)
			oldBinaryCmd2 := utils.Start(oldBinaryPath, homePath2, rpc2, p2p2, pprof2, api2, moniker2, types.ColorGreen, types.ColorRed)

			// wait for rpc 1 and 2 to start
			utils.WaitForServiceToStart(rpc, moniker, timeOutWaitForNode)
			utils.WaitForServiceToStart(rpc2, moniker2, timeOutWaitForNode)

			// query and calculate upgrade block height
			upgradeBlockHeight := utils.QueryAndCalcUpgradeBlockHeight(oldBinaryPath, rpc)

			// query next proposal id
			proposalId, err := utils.QueryNextProposalId(oldBinaryPath, rpc)
			if err != nil {
				log.Printf(types.ColorYellow+"Error querying next proposal id: %v", err)
				log.Printf(types.ColorYellow + "Setting proposal id to 1")
				proposalId = "1"
			}

			// submit upgrade proposal
			txHash := utils.SubmitUpgradeProposal(oldBinaryPath, validatorKeyName, newVersion, upgradeBlockHeight, homePath, keyringBackend, chainId, rpc, broadcastMode)

			utils.WaitForTxConfirmation(oldBinaryPath, rpc, txHash, 5*time.Minute)

			// vote on upgrade proposal
			txHash = utils.VoteOnUpgradeProposal(oldBinaryPath, validatorKeyName, proposalId, homePath, keyringBackend, chainId, rpc, broadcastMode)

			utils.WaitForTxConfirmation(oldBinaryPath, rpc, txHash, 5*time.Minute)

			// wait for upgrade block height
			utils.WaitForBlockHeight(oldBinaryPath, rpc, upgradeBlockHeight)

			// wait 5 seconds
			time.Sleep(5 * time.Second)

			// stop old binaries
			utils.Stop(oldBinaryCmd, oldBinaryCmd2)
		},
	}
}
