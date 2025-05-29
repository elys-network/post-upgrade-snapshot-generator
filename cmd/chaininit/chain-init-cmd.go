package chaininit

import (
	"log"

	"github.com/spf13/cobra"

	elysdcmd "github.com/elys-network/elys/v5/cmd/elysd/cmd"
	"github.com/elys-network/post-upgrade-snapshot-generator/flags"
	"github.com/elys-network/post-upgrade-snapshot-generator/types"
	"github.com/elys-network/post-upgrade-snapshot-generator/utils"
)

// chain-init command does the following:
// 1. download and run old binary
// 2. init nodes
// 3. update config files to enable api and cors
// 4. query node 1 id
// 5. add peers
// 6. add validator keys to node 1
// 7. add validator keys to node 2
// 8. add genesis accounts
// 9. generate genesis tx
// 10. collect genesis txs
// 11. validate genesis
// 12. backup genesis init file
// 13. update genesis
func ChainInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain-init [old_binary_url] [flags]",
		Short: "Initialize the chain",
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

			genesisFilePath, _ := cmd.Flags().GetString(flags.FlagGenesisFilePath)
			if genesisFilePath == "" {
				log.Fatalf(types.ColorRed + "genesis file path is required")
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

			// download and run old binary
			oldBinaryPath, oldVersion, err := utils.DownloadAndRunVersion(oldBinaryUrl, skipBinary)
			if err != nil {
				log.Fatalf(types.ColorRed+"Error downloading and running old binary: %v", err)
			}

			// print old binary path and version
			log.Printf(types.ColorGreen+"Old binary path: %v and version: %v", oldBinaryPath, oldVersion)

			// remove home paths
			utils.RemoveHome(homePath)
			utils.RemoveHome(homePath2)

			// init nodes
			utils.InitNode(oldBinaryPath, moniker, chainId, homePath)
			utils.InitNode(oldBinaryPath, moniker2, chainId, homePath2)

			// update config files to enable api and cors
			utils.UpdateConfig(homePath, "goleveldb")
			utils.UpdateConfig(homePath2, "goleveldb")

			// query node 1 id
			node1Id := utils.QueryNodeId(oldBinaryPath, homePath)

			// add peers
			utils.AddPeers(homePath2, p2p, node1Id)

			// add validator keys to node 1
			validatorAddress := utils.AddKey(oldBinaryPath, validatorKeyName, validatorMnemonic, homePath, keyringBackend)
			validatorAddress2 := utils.AddKey(oldBinaryPath, validatorKeyName2, validatorMnemonic2, homePath, keyringBackend)

			// add validator keys to node 2
			_ = utils.AddKey(oldBinaryPath, validatorKeyName, validatorMnemonic, homePath2, keyringBackend)
			_ = utils.AddKey(oldBinaryPath, validatorKeyName2, validatorMnemonic2, homePath2, keyringBackend)

			// add genesis accounts
			utils.AddGenesisAccount(oldBinaryPath, validatorAddress, homePath, validatorBalances)
			utils.AddGenesisAccount(oldBinaryPath, validatorAddress2, homePath, validatorBalances)

			// generate genesis tx
			utils.GenTx(oldBinaryPath, validatorKeyName, validatorSelfDelegation, chainId, homePath, keyringBackend)

			// collect genesis txs
			utils.CollectGentxs(oldBinaryPath, homePath)

			// validate genesis
			utils.ValidateGenesis(oldBinaryPath, homePath)

			// backup genesis init file
			utils.BackupGenesisInitFile(homePath)

			// update genesis
			utils.UpdateGenesis(oldBinaryPath, homePath, genesisFilePath, validatorBalances, validatorAddress)
		},
	}

	return cmd
}
