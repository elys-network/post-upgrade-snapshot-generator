package types

import (
	"encoding/json"
	"time"

	evidencetypes "cosmossdk.io/x/evidence/types"
	feegranttypes "cosmossdk.io/x/feegrant"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	feeibctypes "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibctypes "github.com/cosmos/ibc-go/v8/modules/core/types"
	accountedpooltypes "github.com/elys-network/elys/x/accountedpool/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	burnertypes "github.com/elys-network/elys/x/burner/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	epochstypes "github.com/elys-network/elys/x/epochs/types"
	leveragelptypes "github.com/elys-network/elys/x/leveragelp/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	parametertypes "github.com/elys-network/elys/x/parameter/types"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
	tiertypes "github.com/elys-network/elys/x/tier/types"
	tokenomicstypes "github.com/elys-network/elys/x/tokenomics/types"
	tradeshieldtypes "github.com/elys-network/elys/x/tradeshield/types"
	transferhooktypes "github.com/elys-network/elys/x/transferhook/types"

	cometbfttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	estakingtypes "github.com/elys-network/elys/x/estaking/types"
	mastercheftypes "github.com/elys-network/elys/x/masterchef/types"
)

type Genesis struct {
	genutiltypes.AppGenesis

	GenesisTime   string      `json:"genesis_time"`
	InitialHeight json.Number `json:"initial_height"`
	AppHash       string      `json:"app_hash"`
	AppState      AppState    `json:"app_state"`
	Consensus     *Consensus  `json:"consensus"`
}

type Consensus struct {
	// genutiltypes.ConsensusGenesis

	// Validators []GenesisValidator `json:"validators"`
	Params *ConsensusParams `json:"params"`
}

type GenesisValidator struct {
	cometbfttypes.GenesisValidator

	Address string `json:"address"`
	PubKey  string `json:"pub_key"`
	Power   string `json:"power"`
	Name    string `json:"name"`
}

type ConsensusParams struct {
	cometbfttypes.ConsensusParams

	Block    BlockParams    `json:"block"`
	Evidence EvidenceParams `json:"evidence"`
	Version  VersionParams  `json:"version"`
	ABCI     ABCIParams     `json:"abci"`
}

type BlockParams struct {
	cometbfttypes.BlockParams

	MaxBytes string `json:"max_bytes"`
	MaxGas   string `json:"max_gas"`
}

type EvidenceParams struct {
	cometbfttypes.EvidenceParams

	MaxAgeNumBlocks string `json:"max_age_num_blocks"`
	MaxAgeDuration  string `json:"max_age_duration"`
	MaxBytes        string `json:"max_bytes,omitempty"`
}

type VersionParams struct {
	cometbfttypes.VersionParams

	App string `json:"app"`
}

type ABCIParams struct {
	cometbfttypes.ABCIParams

	VoteExtensionsEnableHeight string `json:"vote_extensions_enable_height"`
}

type AppState struct {
	Amm                Amm                            `json:"amm"`
	AssetProfile       AssetProfile                   `json:"assetprofile"`
	Auth               Auth                           `json:"auth"`
	AuthZ              authz.GenesisState             `json:"authz"`
	Bank               banktypes.GenesisState         `json:"bank"`
	Burner             burnertypes.GenesisState       `json:"burner"`
	Capability         Capability                     `json:"capability"`
	Commitment         Commitment                     `json:"commitment"`
	Ccvconsumer        interface{}                    `json:"ccvconsumer"`
	Crisis             crisistypes.GenesisState       `json:"crisis"`
	Distribution       Distribution                   `json:"distribution"`
	Epochs             Epochs                         `json:"epochs"`
	Estaking           Estaking                       `json:"estaking"`
	Evidence           EvidenceState                  `json:"evidence"`
	Feegrant           feegranttypes.GenesisState     `json:"feegrant"`
	Feeibc             Feeibc                         `json:"feeibc"`
	Genutil            Genutil                        `json:"genutil"`
	Gov                Gov                            `json:"gov"`
	Group              interface{}                    `json:"group"`
	Ibc                Ibc                            `json:"ibc"`
	Interchainaccounts interface{}                    `json:"interchainaccounts"`
	LeverageLP         LeverageLP                     `json:"leveragelp"`
	Masterchef         Masterchef                     `json:"masterchef"`
	Mint               Mint                           `json:"mint"`
	Oracle             Oracle                         `json:"oracle"`
	Parameter          Parameter                      `json:"parameter"`
	Params             interface{}                    `json:"params"`
	Perpetual          Perpetual                      `json:"perpetual"`
	PoolAccounted      PoolAccounted                  `json:"poolaccounted"`
	Slashing           Slashing                       `json:"slashing"`
	StableStake        StableStake                    `json:"stablestake"`
	Staking            Staking                        `json:"staking"`
	Tier               Tier                           `json:"tier"`
	Tokenomics         Tokenomics                     `json:"tokenomics"`
	Tradeshield        Tradeshield                    `json:"tradeshield"`
	Transfer           transfertypes.GenesisState     `json:"transfer"`
	TransferHook       transferhooktypes.GenesisState `json:"transferhook"`
	Upgrade            interface{}                    `json:"upgrade"`
	Vesting            interface{}                    `json:"vesting"`
	// Include other fields as needed
}

type Tradeshield struct {
	tradeshieldtypes.GenesisState

	Params                     TradeshieldParams `json:"params"`
	PendingSpotOrderList       []interface{}     `json:"pending_spot_order_list"`
	PendingSpotOrderCount      json.Number       `json:"pending_spot_order_count"`
	PendingPerpetualOrderList  []interface{}     `json:"pending_perpetual_order_list"`
	PendingPerpetualOrderCount json.Number       `json:"pending_perpetual_order_count"`
}

type TradeshieldParams struct {
	tradeshieldtypes.Params

	LimitProcessOrder json.Number `json:"limit_process_order"`
}

type Feeibc struct {
	feeibctypes.GenesisState

	IdentifiedFees  []interface{} `json:"identified_fees"`
	ForwardRelayers []interface{} `json:"forward_relayers"`
}

type PoolAccounted struct {
	accountedpooltypes.GenesisState

	AccountedPoolList []AccountedPool `json:"accounted_pool_list"`
}

type AccountedPool struct {
	accountedpooltypes.AccountedPool

	PoolId json.Number `json:"pool_id"`
}

type Tier struct {
	tiertypes.GenesisState

	Params        interface{}   `json:"params"`
	PortfolioList []interface{} `json:"portfolioList"`
}

type Masterchef struct {
	mastercheftypes.GenesisState

	Params                 MasterchefParams            `json:"params"`
	ExternalIncentiveIndex json.Number                 `json:"external_incentive_index"`
	PoolInfos              []MasterchefPoolInfo        `json:"pool_infos"`
	PoolRewardInfos        []MasterchefPoolRewardInfo  `json:"pool_reward_infos"`
	UserRewardInfos        []MasterchefUserRewardInfo  `json:"user_reward_infos"`
	PoolRewardsAccum       []MasterchefPoolRewardAccum `json:"pool_rewards_accum"`
}

type MasterchefParams struct {
	mastercheftypes.Params

	LpIncentives *MasterchefIncentiveInfo `json:"lp_incentives"`
}

type MasterchefIncentiveInfo struct {
	mastercheftypes.IncentiveInfo

	BlocksDistributed json.Number `json:"blocks_distributed"`
}

type MasterchefPoolRewardAccum struct {
	mastercheftypes.PoolRewardsAccum

	PoolId      json.Number `json:"pool_id"`
	BlockHeight json.Number `json:"block_height"`
	Timestamp   json.Number `json:"timestamp"`
}

type MasterchefUserRewardInfo struct {
	mastercheftypes.UserRewardInfo

	PoolId json.Number `json:"pool_id"`
}

type MasterchefPoolRewardInfo struct {
	mastercheftypes.PoolRewardInfo

	PoolId           json.Number `json:"pool_id"`
	LastUpdatedBlock json.Number `json:"last_updated_block"`
}

type MasterchefPoolInfo struct {
	mastercheftypes.PoolInfo

	PoolId json.Number `json:"pool_id"`
}

type Estaking struct {
	estakingtypes.GenesisState

	Params           EstakingParams `json:"params"`
	StakingSnapshots []interface{}  `json:"staking_snapshots"`
}

type EstakingParams struct {
	estakingtypes.Params

	StakeIncentives *EstakingIncentiveInfo `json:"stake_incentives"`
}

type EstakingIncentiveInfo struct {
	estakingtypes.IncentiveInfo

	BlocksDistributed json.Number `json:"blocks_distributed"`
}

type Tokenomics struct {
	tokenomicstypes.GenesisState

	AirdropList            []interface{}              `json:"airdrop_list"`
	GenesisInflation       TokenomicsGenesisInflation `json:"genesis_inflation"`
	TimeBasedInflationList []interface{}              `json:"time_based_inflation_list"`
}

type TokenomicsGenesisInflation struct {
	tokenomicstypes.GenesisInflation

	Inflation             TokenomicsInflationEntry `json:"inflation"`
	SeedVesting           json.Number              `json:"seed_vesting"`
	StrategicSalesVesting json.Number              `json:"strategic_sales_vesting"`
}

type TokenomicsInflationEntry struct {
	tokenomicstypes.InflationEntry

	LmRewards         json.Number `json:"lm_rewards"`
	IcsStakingRewards json.Number `json:"ics_staking_rewards"`
	CommunityFund     json.Number `json:"community_fund"`
	StrategicReserve  json.Number `json:"strategic_reserve"`
	TeamTokensVested  json.Number `json:"team_tokens_vested"`
}

type StableStake struct {
	stablestaketypes.GenesisState

	Params       StableStakeParams `json:"params"`
	DebtList     []interface{}     `json:"debt_list"`
	InterestList []interface{}     `json:"interest_list"`
}

type StableStakeParams struct {
	stablestaketypes.Params

	EpochLength json.Number `json:"epoch_length"`
}

type Epochs struct {
	epochstypes.GenesisState

	Epochs []interface{} `json:"epochs"`
}

type Commitment struct {
	commitmenttypes.GenesisState

	Params      CommitmentParams `json:"params"`
	Commitments []interface{}    `json:"commitments"`
	KolList     []KolList        `json:"kol_list"`
	AtomStakers []AtomStaker     `json:"atom_stakers"`
}

type AtomStaker struct {
	commitmenttypes.AtomStaker

	Address string `json:"address"`
	Amount  string `json:"amount"`
}

type KolList struct {
	commitmenttypes.KolList

	Address string `json:"address"`
	Amount  string `json:"amount"`
}

type CommitmentParams struct {
	commitmenttypes.Params

	VestingInfos            []CommitmentVestingInfo `json:"vesting_infos"`
	NumberOfCommitments     json.Number             `json:"number_of_commitments"`
	StartAirdropClaimHeight json.Number             `json:"start_airdrop_claim_height"`
	EndAirdropClaimHeight   json.Number             `json:"end_airdrop_claim_height"`
	StartKolClaimHeight     json.Number             `json:"start_kol_claim_height"`
	EndKolClaimHeight       json.Number             `json:"end_kol_claim_height"`
}

type CommitmentVestingInfo struct {
	commitmenttypes.VestingInfo

	NumBlocks      json.Number `json:"num_blocks"`
	NumMaxVestings json.Number `json:"num_max_vestings"`
}

type AssetProfile struct {
	assetprofiletypes.GenesisState

	EntryList []AssetProfileEntry `json:"entry_list"`
}

type AssetProfileEntry struct {
	assetprofiletypes.Entry

	Decimals json.Number `json:"decimals"`
}

type Amm struct {
	ammtypes.GenesisState

	Params         AmmParams     `json:"params"`
	PoolList       []interface{} `json:"pool_list"`
	SlippageTracks []interface{} `json:"slippage_tracks"`
}

type AmmParams struct {
	ammtypes.Params

	SlippageTrackDuration json.Number `json:"slippage_track_duration"`
	LpLockupDuration      json.Number `json:"lp_lockup_duration"`
}

type Genutil struct {
	// genutiltypes.GenesisState

	GenTxs []interface{} `json:"gen_txs"`
}

type EvidenceState struct {
	evidencetypes.GenesisState

	Evidence []interface{} `json:"evidence"`
}

type Oracle struct {
	oracletypes.GenesisState

	Params     OracleParams      `json:"params"`
	AssetInfos []OracleAssetInfo `json:"asset_infos"`
	Prices     []OraclePrice     `json:"prices"`
}

type OracleAssetInfo struct {
	oracletypes.AssetInfo

	Decimal json.Number `json:"decimal"`
}

type OraclePrice struct {
	oracletypes.Price

	Timestamp   json.Number `json:"timestamp"`
	BlockHeight json.Number `json:"block_height"`
}

type OracleParams struct {
	oracletypes.Params

	OracleScriptID json.Number `json:"oracle_script_id"`
	Multiplier     json.Number `json:"multiplier"`
	AskCount       json.Number `json:"ask_count"`
	MinCount       json.Number `json:"min_count"`
	PrepareGas     json.Number `json:"prepare_gas"`
	ExecuteGas     json.Number `json:"execute_gas"`
	ClientID       json.Number `json:"client_id"`
	BandEpoch      json.Number `json:"band_epoch"`

	VotePeriod               json.Number `json:"vote_period"`
	RewardDistributionWindow json.Number `json:"reward_distribution_window"`
	SlashWindow              json.Number `json:"slash_window"`
	HistoricStampPeriod      json.Number `json:"historic_stamp_period"`
	MedianStampPeriod        json.Number `json:"median_stamp_period"`
	MaximumPriceStamps       json.Number `json:"maximum_price_stamps"`
	MaximumMedianStamps      json.Number `json:"maximum_median_stamps"`
	PriceExpiryTime          json.Number `json:"price_expiry_time"`
	LifeTimeInBlocks         json.Number `json:"life_time_in_blocks"`
}

type Parameter struct {
	parametertypes.GenesisState

	Params ParameterParams `json:"params"`
}

type ParameterParams struct {
	parametertypes.Params

	TotalBlocksPerYear  json.Number `json:"total_blocks_per_year"`
	RewardsDataLifetime json.Number `json:"rewards_data_lifetime"`
}

type Capability struct {
	capabilitytypes.GenesisState

	Index  json.Number   `json:"index"`
	Owners []interface{} `json:"owners"`
}

type Slashing struct {
	slashingtypes.GenesisState

	Params       SlashingParams `json:"params"`
	SigningInfos []SigningInfo  `json:"signing_infos"`
	MissedBlocks []interface{}  `json:"missed_blocks"`
}

type SigningInfo struct {
	// slashingtypes.SigningInfo

	Address              string               `json:"address"`
	ValidatorSigningInfo ValidatorSigningInfo `json:"validator_signing_info"`
}

type ValidatorSigningInfo struct {
	// slashingtypes.ValidatorSigningInfo

	Address             string      `json:"address"`
	StartHeight         json.Number `json:"start_height"`
	IndexOffset         json.Number `json:"index_offset"`
	JailedUntil         time.Time   `json:"jailed_until"`
	Tombstoned          bool        `json:"tombstoned"`
	MissedBlocksCounter json.Number `json:"missed_blocks_counter"`
}

type SlashingParams struct {
	slashingtypes.Params

	SignedBlocksWindow   json.Number `json:"signed_blocks_window"`
	DowntimeJailDuration string      `json:"downtime_jail_duration"`
}

type Mint struct {
	minttypes.GenesisState

	Params MintParams `json:"params"`
}

type MintParams struct {
	minttypes.Params

	BlocksPerYear json.Number `json:"blocks_per_year"`
}

type Gov struct {
	govv1types.GenesisState

	StartingProposalId json.Number      `json:"starting_proposal_id"`
	Deposits           []interface{}    `json:"deposits"`
	Votes              []interface{}    `json:"votes"`
	Proposals          []interface{}    `json:"proposals"`
	DepositParams      GovDepositParams `json:"deposit_params"`
	VotingParams       GovVotingParams  `json:"voting_params"`
	Params             GovParams        `json:"params"`
}

type GovParams struct {
	govv1types.Params

	MaxDepositPeriod      string `json:"max_deposit_period"`
	VotingPeriod          string `json:"voting_period"`
	ExpeditedVotingPeriod string `json:"expedited_voting_period"`
}

type GovDepositParams struct {
	govv1types.DepositParams

	MaxDepositPeriod string `json:"max_deposit_period"`
}

type GovVotingParams struct {
	govv1types.VotingParams

	VotingPeriod string `json:"voting_period"`
}

type Staking struct {
	stakingtypes.GenesisState

	Params               StakingParams `json:"params"`
	LastValidatorPowers  []interface{} `json:"last_validator_powers"`
	Validators           []interface{} `json:"validators"`
	Delegations          []interface{} `json:"delegations"`
	UnbondingDelegations []interface{} `json:"unbonding_delegations"`
	ColorRedelegations   []interface{} `json:"redelegations"`
}

type StakingParams struct {
	stakingtypes.Params

	UnbondingTime     string      `json:"unbonding_time"`
	MaxValidators     json.Number `json:"max_validators"`
	MaxEntries        json.Number `json:"max_entries"`
	HistoricalEntries json.Number `json:"historical_entries"`
}

type Distribution struct {
	distributiontypes.GenesisState

	DelegatorWithdrawInfos          []interface{} `json:"delegator_withdraw_infos"`
	OutstandingRewards              []interface{} `json:"outstanding_rewards"`
	ValidatorAccumulatedCommissions []interface{} `json:"validator_accumulated_commissions"`
	ValidatorHistoricalRewards      []interface{} `json:"validator_historical_rewards"`
	ValidatorCurrentRewards         []interface{} `json:"validator_current_rewards"`
	DelegatorStartingInfos          []interface{} `json:"delegator_starting_infos"`
	ValidatorSlashEvents            []interface{} `json:"validator_slash_events"`
}

type Ibc struct {
	ibctypes.GenesisState

	ClientGenesis     ClientGenesis     `json:"client_genesis"`
	ConnectionGenesis ConnectionGenesis `json:"connection_genesis"`
	ChannelGenesis    ChannelGenesis    `json:"channel_genesis"`
}

type ClientGenesis struct {
	ibcclienttypes.GenesisState

	Clients            []interface{}         `json:"clients"`
	ClientsConsensus   []interface{}         `json:"clients_consensus"`
	ClientsMetadata    []interface{}         `json:"clients_metadata"`
	Params             ibcclienttypes.Params `json:"params"`
	NextClientSequence json.Number           `json:"next_client_sequence"`
}

type ConnectionGenesis struct {
	ibcconnectiontypes.GenesisState

	Connections            []interface{}           `json:"connections"`
	ClientConnectionPaths  []interface{}           `json:"client_connection_paths"`
	NextConnectionSequence json.Number             `json:"next_connection_sequence"`
	Params                 ConnectionGenesisParams `json:"params"`
}

type ConnectionGenesisParams struct {
	ibcconnectiontypes.Params

	MaxExpectedTimePerBlock json.Number `json:"max_expected_time_per_block"`
}

type ChannelGenesis struct {
	ibcchanneltypes.GenesisState

	Channels            []interface{}        `json:"channels"`
	Acknowledgements    []interface{}        `json:"acknowledgements"`
	Commitments         []interface{}        `json:"commitments"`
	Receipts            []interface{}        `json:"receipts"`
	SendSequences       []interface{}        `json:"send_sequences"`
	RecvSequences       []interface{}        `json:"recv_sequences"`
	AckSequences        []interface{}        `json:"ack_sequences"`
	NextChannelSequence json.Number          `json:"next_channel_sequence"`
	Params              ChannelGenesisParams `json:"params"`
}

type ChannelGenesisParams struct {
	ibcchanneltypes.Params

	UpgradeTimeout ChannelGenesisTimeout `json:"upgrade_timeout"`
}

type ChannelGenesisTimeout struct {
	ibcchanneltypes.Timeout

	Height    ChannelGenesisHeight `json:"height"`
	Timestamp json.Number          `json:"timestamp"`
}

type ChannelGenesisHeight struct {
	ibcclienttypes.Height

	RevisionNumber json.Number `json:"revision_number"`
	RevisionHeight json.Number `json:"revision_height"`
}

type LeverageLP struct {
	leveragelptypes.GenesisState

	Params       LeverageLPParams `json:"params"`
	PoolList     []interface{}    `json:"pool_list"`
	PositionList []interface{}    `json:"position_list"`
}

type LeverageLPParams struct {
	leveragelptypes.Params

	EpochLength      json.Number   `json:"epoch_length"`
	MaxOpenPositions json.Number   `json:"max_open_positions"`
	NumberPerBlock   json.Number   `json:"number_per_block"`
	EnabledPools     []json.Number `json:"enabled_pools"`
}

type Perpetual struct {
	perpetualtypes.GenesisState

	Params   PerpetualParams `json:"params"`
	PoolList []interface{}   `json:"pool_list"`
	MtpList  []interface{}   `json:"mtp_list"`
}

type PerpetualParams struct {
	perpetualtypes.Params

	MaxOpenPositions json.Number   `json:"max_open_positions"`
	MaxLimitOrder    json.Number   `json:"max_limit_order"`
	EnabledPools     []json.Number `json:"enabled_pools"`
}

type AuthParams struct {
	authtypes.Params

	MaxMemoCharacters      json.Number `json:"max_memo_characters"`
	TxSigLimit             json.Number `json:"tx_sig_limit"`
	TxSizeCostPerByte      json.Number `json:"tx_size_cost_per_byte"`
	SigVerifyCostEd25519   json.Number `json:"sig_verify_cost_ed25519"`
	SigVerifyCostSecp256K1 json.Number `json:"sig_verify_cost_secp256k1"`
}

type BaseAccount struct {
	Address       string      `json:"address"`
	PubKey        interface{} `json:"pub_key"`
	AccountNumber json.Number `json:"account_number"`
	Sequence      json.Number `json:"sequence"`
}

type ModuleAccount struct {
	BaseAccount BaseAccount `json:"base_account"`
	Name        string      `json:"name"`
	Permissions []string    `json:"permissions"`
}

type Account struct {
	*VestingAccount `json:",omitempty"`
	*BaseAccount    `json:",omitempty"`
	*ModuleAccount  `json:",omitempty"`
	Type            string `json:"@type"`
}

type Auth struct {
	authtypes.GenesisState

	Params   AuthParams `json:"params"`
	Accounts []Account  `json:"accounts"`
}

type VestingAccount struct {
	BaseVestingAccount BaseVestingAccount `json:"base_vesting_account"`
	StartTime          json.Number        `json:"start_time"`
	VestingPeriods     []VestingPeriod    `json:"vesting_periods,omitempty"`
}

type BaseVestingAccount struct {
	BaseAccount      BaseAccount `json:"base_account"`
	OriginalVesting  []sdk.Coin  `json:"original_vesting"`
	DelegatedFree    []sdk.Coin  `json:"delegated_free"`
	DelegatedVesting []sdk.Coin  `json:"delegated_vesting"`
	EndTime          json.Number `json:"end_time"`
}

type VestingPeriod struct {
	Length json.Number `json:"length"`
	Amount []sdk.Coin  `json:"amount"`
}

// KeyOutput represents the JSON structure of the output from the add key command
type KeyOutput struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Address  string `json:"address"`
	PubKey   string `json:"pubkey"`
	Mnemonic string `json:"mnemonic"`
}

// StatusOutput represents the JSON structure of the output from the status command
type StatusOutput struct {
	SyncInfo struct {
		LatestBlockHeight string `json:"latest_block_height"`
	} `json:"sync_info"`
}

// ProposalsOutput represents the JSON structure of the output from the query proposals command
type ProposalsOutput struct {
	Proposals []struct {
		Id string `json:"id"`
	} `json:"proposals"`
}
