package utils

import (
	"encoding/json"
	"log"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func UpdateGenesis(cmdPath, homePath, genesisFilePath string, balances []string, validatorAddress string) {
	genesis, err := ReadGenesisFile(genesisFilePath)
	if err != nil {
		log.Fatalf(types.ColorRed+"Error reading genesis file: %v", err)
	}

	genesisInitFilePath := homePath + "/config/genesis.json"
	genesisInit, err := ReadGenesisFile(genesisInitFilePath)
	if err != nil {
		log.Fatalf(types.ColorRed+"Error reading initial genesis file: %v", err)
	}

	filterAccountAddresses := []string{
		"elys1ed2lkxujcqfckkhfwmyjqwuqp47ve37crctuus", // remove existing account 0
	}
	filterBalanceAddresses := []string{
		"elys1ed2lkxujcqfckkhfwmyjqwuqp47ve37crctuus", // remove existing account 0
		authtypes.NewModuleAddress("bonded_tokens_pool").String(),
		authtypes.NewModuleAddress("not_bonded_tokens_pool").String(),
		authtypes.NewModuleAddress("gov").String(),
	}

	var coinsToRemove sdk.Coins

	genesis.AppState.Auth.Accounts = FilterAccounts(genesis.AppState.Auth.Accounts, filterAccountAddresses)
	genesis.AppState.Bank.Balances, coinsToRemove = FilterBalances(genesis.AppState.Bank.Balances, filterBalanceAddresses)

	// update supply
	genesis.AppState.Bank.Supply = genesis.AppState.Bank.Supply.Sub(coinsToRemove...)

	// Parse balances and add to supply for both validators
	for _, balance := range balances {
		coin, err := sdk.ParseCoinNormalized(balance)
		if err != nil {
			log.Fatalf(types.ColorRed+"Error parsing balance %s: %v", balance, err)
		}

		// add node 1 supply
		genesis.AppState.Bank.Supply = genesis.AppState.Bank.Supply.Add(coin)
		// add node 2 supply
		genesis.AppState.Bank.Supply = genesis.AppState.Bank.Supply.Add(coin)
	}

	// Add new validator account and balance
	genesis.AppState.Auth.Accounts = append(genesis.AppState.Auth.Accounts, genesisInit.AppState.Auth.Accounts...)
	genesis.AppState.Bank.Balances = append(genesis.AppState.Bank.Balances, genesisInit.AppState.Bank.Balances...)

	// update bank params
	genesis.AppState.Bank.Params.DefaultSendEnabled = true

	// ColorReset staking data
	stakingParams := genesis.AppState.Staking.Params
	genesis.AppState.Staking = genesisInit.AppState.Staking
	genesis.AppState.Staking.Params = stakingParams

	// temporary fix for staking params
	genesis.AppState.Staking.Params.BondDenom = "uelys"

	// ColorReset slashing data
	genesis.AppState.Slashing = genesisInit.AppState.Slashing

	// Add validator signing info to genesis
	validatorConsAddr := GetValidatorConsensusAddress(cmdPath)
	validatorSigningInfo := types.ValidatorSigningInfo{
		Address:             validatorConsAddr,
		StartHeight:         json.Number("0"),
		IndexOffset:         json.Number("0"),
		JailedUntil:         time.Time{},
		Tombstoned:          false,
		MissedBlocksCounter: json.Number("0"),
	}
	genesis.AppState.Slashing.SigningInfos = []types.SigningInfo{
		{
			Address:              validatorConsAddr,
			ValidatorSigningInfo: validatorSigningInfo,
		},
	}

	// set genutil from genesisInit
	genesis.AppState.Genutil = genesisInit.AppState.Genutil

	// add localhost as allowed client
	genesis.AppState.Ibc.ClientGenesis.Params.AllowedClients = append(genesis.AppState.Ibc.ClientGenesis.Params.AllowedClients, "09-localhost")

	// update voting period
	votingPeriod := "60s"
	minDeposit := sdk.Coins{sdk.NewInt64Coin("uelys", 10000000)}
	genesis.AppState.Gov.Params.VotingPeriod = votingPeriod
	genesis.AppState.Gov.Params.MaxDepositPeriod = votingPeriod
	genesis.AppState.Gov.Params.MinDeposit = minDeposit
	// set deprecated settings
	genesis.AppState.Gov.VotingParams.VotingPeriod = votingPeriod
	genesis.AppState.Gov.DepositParams.MaxDepositPeriod = votingPeriod
	genesis.AppState.Gov.DepositParams.MinDeposit = minDeposit

	// update oracle price expiration
	genesis.AppState.Oracle.Params.PriceExpiryTime = "31536000"
	genesis.AppState.Oracle.Params.LifeTimeInBlocks = "8000000"

	// update ccvconsumer
	genesis.AppState.Ccvconsumer = genesisInit.AppState.Ccvconsumer

	// update assetprofile and oracle to add both eth and wbtc
	wbtc := &types.AssetProfileEntry{
		Entry: assetprofiletypes.Entry{
			Denom:           "wbtc-satoshi",
			BaseDenom:       "wbtc-satoshi",
			DisplaySymbol:   "wbtc",
			DisplayName:     "Wrapped Bitcoin",
			UnitDenom:       "wbtc",
			CommitEnabled:   true,
			WithdrawEnabled: true,
		},
		Decimals: json.Number("8"),
	}
	weth := &types.AssetProfileEntry{
		Entry: assetprofiletypes.Entry{
			Denom:           "weth-wei",
			BaseDenom:       "weth-wei",
			DisplaySymbol:   "weth",
			DisplayName:     "Wrapped Ethereum",
			UnitDenom:       "weth",
			CommitEnabled:   true,
			WithdrawEnabled: true,
		},
		Decimals: json.Number("18"),
	}

	// Only append if base denoms don't exist
	for _, entry := range genesis.AppState.AssetProfile.EntryList {
		if entry.Entry.BaseDenom == wbtc.Entry.BaseDenom {
			wbtc = nil
		}
		if entry.Entry.BaseDenom == weth.Entry.BaseDenom {
			weth = nil
		}
	}

	if wbtc != nil {
		genesis.AppState.AssetProfile.EntryList = append(genesis.AppState.AssetProfile.EntryList, *wbtc)
	}
	if weth != nil {
		genesis.AppState.AssetProfile.EntryList = append(genesis.AppState.AssetProfile.EntryList, *weth)
	}

	if wbtc != nil {
		genesis.AppState.Oracle.AssetInfos = append(genesis.AppState.Oracle.AssetInfos, types.OracleAssetInfo{
			AssetInfo: oracletypes.AssetInfo{
				Denom:      "wbtc-satoshi",
				Display:    "WBTC",
				ElysTicker: "WBTC",
				BandTicker: "WBTC",
			},
			Decimal: json.Number("8"),
		})
	}
	if weth != nil {
		genesis.AppState.Oracle.AssetInfos = append(genesis.AppState.Oracle.AssetInfos, types.OracleAssetInfo{
			AssetInfo: oracletypes.AssetInfo{
				Denom:      "weth-wei",
				Display:    "WETH",
				ElysTicker: "WETH",
				BandTicker: "WETH",
			},
			Decimal: json.Number("18"),
		})
	}

	if wbtc != nil {
		genesis.AppState.Oracle.Prices = append(genesis.AppState.Oracle.Prices, types.OraclePrice{
			Price: oracletypes.Price{
				Asset: "WBTC",
				Price: math.LegacyMustNewDecFromStr("100000"),
			},
			Timestamp:   json.Number("1733907595"),
			BlockHeight: json.Number("111045"),
		})
	}
	if weth != nil {
		genesis.AppState.Oracle.Prices = append(genesis.AppState.Oracle.Prices, types.OraclePrice{
			Price: oracletypes.Price{
				Asset: "WETH",
				Price: math.LegacyMustNewDecFromStr("4000"),
			},
			Timestamp:   json.Number("1733907595"),
			BlockHeight: json.Number("111045"),
		})
	}

	// update AMM params to whitelist validator address
	genesis.AppState.Amm.Params.AllowedPoolCreators = append(genesis.AppState.Amm.Params.AllowedPoolCreators, validatorAddress)

	// update commitment airdrop params
	// genesis.AppState.Commitment.Params.EnableClaim = true
	// genesis.AppState.Commitment.Params.StartAirdropClaimHeight = json.Number("111046")
	// genesis.AppState.Commitment.Params.EndAirdropClaimHeight = json.Number("111600")
	// genesis.AppState.Commitment.Params.StartKolClaimHeight = json.Number("111046")
	// genesis.AppState.Commitment.Params.EndKolClaimHeight = json.Number("111600")

	// genesis.AppState.Commitment.KolList = []types.KolList{
	// 	{
	// 		Address: "elys1wluheavcjknwnagskt0994rrtuvap9jcz0ng5q",
	// 		Amount:  "1000000000000",
	// 	},
	// }

	// genesis.AppState.Commitment.AtomStakers = []types.AtomStaker{
	// 	{
	// 		Address: "elys1wluheavcjknwnagskt0994rrtuvap9jcz0ng5q",
	// 		Amount:  "1000000000000",
	// 	},
	// }

	// write genesis file
	outputFilePath := homePath + "/config/genesis.json"
	if err := WriteGenesisFile(outputFilePath, genesis); err != nil {
		log.Fatalf(types.ColorRed+"Error writing genesis file: %v", err)
	}
}
