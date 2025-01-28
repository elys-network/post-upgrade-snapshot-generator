package flags

const (
	// global
	FlagSkipProposal        = "skip-proposal"
	FlagSkipBinary          = "skip-binary"
	FlagSkipUnbondValidator = "skip-unbond-validator"
	FlagChainId             = "chain-id"
	FlagKeyringBackend      = "keyring-backend"
	FlagGenesisFilePath     = "genesis-file-path"
	FlagBroadcastMode       = "broadcast-mode"
	FlagDbEngine            = "db-engine"

	// timeout
	FlagTimeOutToWaitForNode = "timeout-wait-for-node"
	FlagTimeOutNextBlock     = "timeout-next-block"

	// node 1
	FlagHome                    = "home"
	FlagMoniker                 = "moniker"
	FlagValidatorKeyName        = "validator-key-name"
	FlagValidatorSelfDelegation = "validator-self-delegation"
	FlagValidatorMnemonic       = "validator-mnemonic"
	FlagRpc                     = "rpc"
	FlagP2p                     = "p2p"
	FlagPprof                   = "pprof"
	FlagApi                     = "api"

	// node 2
	FlagHome2                    = "home-2"
	FlagMoniker2                 = "moniker-2"
	FlagValidatorKeyName2        = "validator-key-name-2"
	FlagValidatorSelfDelegation2 = "validator-self-delegation-2"
	FlagValidatorMnemonic2       = "validator-mnemonic-2"
	FlagRpc2                     = "rpc-2"
	FlagP2p2                     = "p2p-2"
	FlagPprof2                   = "pprof-2"
	FlagApi2                     = "api-2"

	FlagValidatorBalances = "validator-balances"

	// price feeder flags
	FlagPriceFeederEnable     = "pricefeeder.enable"
	FlagPriceFeederConfigPath = "pricefeeder.config_path"
	FlagPriceFeederLogLevel   = "pricefeeder.log_level"
)
