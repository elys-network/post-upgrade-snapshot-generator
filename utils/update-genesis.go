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

func UpdateGenesis(validatorBalance, homePath, genesisFilePath string) {
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
		Add(sdk.NewCoin("ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65", newValidatorBalance)).
		Add(sdk.NewCoin("ibc/E2D2F6ADCC68AA3384B2F5DFACCA437923D137C14E86FB8A10207CF3BED0C8D4", newValidatorBalance)).
		Add(sdk.NewCoin("ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953", newValidatorBalance))

	// add node 2 supply
	genesis.AppState.Bank.Supply = genesis.AppState.Bank.Supply.
		Add(sdk.NewCoin("uelys", newValidatorBalance)).
		Add(sdk.NewCoin("ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65", newValidatorBalance)).
		Add(sdk.NewCoin("ibc/E2D2F6ADCC68AA3384B2F5DFACCA437923D137C14E86FB8A10207CF3BED0C8D4", newValidatorBalance)).
		Add(sdk.NewCoin("ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953", newValidatorBalance))

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
	validatorConsAddr := GetValidatorConsensusAddress()
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
