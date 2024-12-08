package utils

import (
	"encoding/json"
	"log"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func UpdateGenesis(validatorBalance, cmdPath, homePath, genesisFilePath string) {
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

	newValidatorBalance, ok := math.NewIntFromString(validatorBalance)
	if !ok {
		panic(types.ColorRed + "invalid number")
	}

	// update supply
	genesis.AppState.Bank.Supply = genesis.AppState.Bank.Supply.Sub(coinsToRemove...)

	// add node 1 supply
	genesis.AppState.Bank.Supply = genesis.AppState.Bank.Supply.
		Add(sdk.NewCoin("uelys", newValidatorBalance)).
		Add(sdk.NewCoin("ibc/F082B65C88E4B6D5EF1DB243CDA1D331D002759E938A0F5CD3FFDC5D53B3E349", newValidatorBalance)).
		Add(sdk.NewCoin("ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9", newValidatorBalance))

	// add node 2 supply
	genesis.AppState.Bank.Supply = genesis.AppState.Bank.Supply.
		Add(sdk.NewCoin("uelys", newValidatorBalance)).
		Add(sdk.NewCoin("ibc/F082B65C88E4B6D5EF1DB243CDA1D331D002759E938A0F5CD3FFDC5D53B3E349", newValidatorBalance)).
		Add(sdk.NewCoin("ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9", newValidatorBalance))

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

	// write genesis file
	outputFilePath := homePath + "/config/genesis.json"
	if err := WriteGenesisFile(outputFilePath, genesis); err != nil {
		log.Fatalf(types.ColorRed+"Error writing genesis file: %v", err)
	}
}
