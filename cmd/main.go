package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/elys-network/post-upgrade-snapshot-generator/cmd/chaininit"
	"github.com/elys-network/post-upgrade-snapshot-generator/cmd/chainsnapshotexport"
	"github.com/elys-network/post-upgrade-snapshot-generator/cmd/createsecondvalidator"
	"github.com/elys-network/post-upgrade-snapshot-generator/cmd/createsnapshot"
	"github.com/elys-network/post-upgrade-snapshot-generator/cmd/deletesnapshot"
	"github.com/elys-network/post-upgrade-snapshot-generator/cmd/preparevalidatordata"
	"github.com/elys-network/post-upgrade-snapshot-generator/cmd/submitnewproposal"
	"github.com/elys-network/post-upgrade-snapshot-generator/cmd/upgradetonewbinary"
	"github.com/elys-network/post-upgrade-snapshot-generator/cmd/uploadsnapshot"
	"github.com/elys-network/post-upgrade-snapshot-generator/cmd/version"
	"github.com/elys-network/post-upgrade-snapshot-generator/flags"
	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "upgrade-assure [flags]",
		Short: "Upgrade Assure is a tool for running a chain from a snapshot and to test out the upgrade process.",
		Long:  `A tool for running a chain from a snapshot.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(1)
		},
	}

	// get HOME environment variable
	homeEnv, _ := os.LookupEnv("HOME")

	// global flags
	rootCmd.PersistentFlags().Bool(flags.FlagSkipProposal, false, "skip proposal")
	rootCmd.PersistentFlags().Bool(flags.FlagSkipBinary, false, "skip binary download")
	rootCmd.PersistentFlags().Bool(flags.FlagSkipUnbondValidator, false, "skip unbond validator")
	rootCmd.PersistentFlags().String(flags.FlagChainId, "elysicstestnet-1", "chain id")
	rootCmd.PersistentFlags().String(flags.FlagKeyringBackend, "test", "keyring backend")
	rootCmd.PersistentFlags().String(flags.FlagGenesisFilePath, "/tmp/genesis.json", "genesis file path")
	rootCmd.PersistentFlags().String(flags.FlagBroadcastMode, "sync", "broadcast mode")
	rootCmd.PersistentFlags().String(flags.FlagDbEngine, "pebbledb", "database engine to use")

	rootCmd.PersistentFlags().Int(flags.FlagTimeOutToWaitForNode, 600, "set the maximum timeout in (seconds) to wait for the node starting")
	rootCmd.PersistentFlags().Int(flags.FlagTimeOutNextBlock, 5, "set the maximum timeout in (minutes) to wait for the next block")
	// node 1 flags
	rootCmd.PersistentFlags().String(flags.FlagHome, homeEnv+"/.elys", "home directory")
	rootCmd.PersistentFlags().String(flags.FlagMoniker, "alice", "moniker")
	rootCmd.PersistentFlags().String(flags.FlagValidatorKeyName, "validator", "validator key name")
	rootCmd.PersistentFlags().String(flags.FlagValidatorSelfDelegation, "50000000000000", "validator self delegation")
	rootCmd.PersistentFlags().String(flags.FlagValidatorMnemonic, "shrug census ancient uniform sausage own oil boss tool captain ride year conduct welcome siren protect mutual zero funny universe candy gown rack sister", "validator mnemonic")
	rootCmd.PersistentFlags().String(flags.FlagRpc, "tcp://0.0.0.0:26657", "rpc")
	rootCmd.PersistentFlags().String(flags.FlagP2p, "tcp://0.0.0.0:26656", "p2p")
	rootCmd.PersistentFlags().String(flags.FlagPprof, "localhost:6060", "pprof")
	rootCmd.PersistentFlags().String(flags.FlagApi, "tcp://localhost:1317", "api")

	// node 2 flags
	rootCmd.PersistentFlags().String(flags.FlagHome2, homeEnv+"/.elys2", "home directory 2")
	rootCmd.PersistentFlags().String(flags.FlagMoniker2, "bob", "moniker 2")
	rootCmd.PersistentFlags().String(flags.FlagValidatorKeyName2, "validator-2", "validator key name 2")
	rootCmd.PersistentFlags().String(flags.FlagValidatorSelfDelegation2, "1000000", "validator self delegation 2")
	rootCmd.PersistentFlags().String(flags.FlagValidatorMnemonic2, "august viable pet tone normal below almost blush portion example trick circle pumpkin citizen conduct outdoor universe wolf ankle asthma deliver correct pool juice", "validator mnemonic 2")
	rootCmd.PersistentFlags().String(flags.FlagRpc2, "tcp://0.0.0.0:26667", "rpc")
	rootCmd.PersistentFlags().String(flags.FlagP2p2, "tcp://0.0.0.0:26666", "p2p")
	rootCmd.PersistentFlags().String(flags.FlagPprof2, "localhost:6061", "pprof")
	rootCmd.PersistentFlags().String(flags.FlagApi2, "tcp://localhost:1318", "api")

	rootCmd.PersistentFlags().StringSlice(flags.FlagValidatorBalances, []string{
		"200000000000000uelys",
		"100000000000000ibc/F082B65C88E4B6D5EF1DB243CDA1D331D002759E938A0F5CD3FFDC5D53B3E349",
		"100000000000000ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9",
		"10000000000000000wbtc-satoshi",
		"100000000000000000000000000weth-wei",
	}, "list of balance+denomination pairs for validator accounts (e.g., '100000000uelys,200000000uatom')")

	// price feeder flags
	//rootCmd.PersistentFlags().Bool(flags.FlagPriceFeederEnable, true, "Enable price feeder")
	//rootCmd.PersistentFlags().String(flags.FlagPriceFeederConfigPath, "cmd/elysd/price-feeder.example.toml", "Price feeder config path")
	//rootCmd.PersistentFlags().String(flags.FlagPriceFeederLogLevel, "error", "Price feeder log level")

	rootCmd.AddCommand(version.VersionCmd())
	rootCmd.AddCommand(chainsnapshotexport.ChainSnapshotExportCmd())
	rootCmd.AddCommand(chaininit.ChainInitCmd())
	rootCmd.AddCommand(createsecondvalidator.CreateSecondValidatorCmd())
	rootCmd.AddCommand(preparevalidatordata.PrepareValidatorDataCmd())
	rootCmd.AddCommand(submitnewproposal.SubmitNewProposalCmd())
	rootCmd.AddCommand(upgradetonewbinary.UpgradeToNewBinaryCmd())
	rootCmd.AddCommand(uploadsnapshot.UploadSnapshotCmd())
	rootCmd.AddCommand(deletesnapshot.DeleteSnapshotCmd())
	rootCmd.AddCommand(createsnapshot.CreateSnapshotCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf(types.ColorRed+"Error executing command: %v", err)
	}
}
