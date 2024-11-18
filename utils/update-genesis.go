package utils

import (
	"log"

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
		"elys1gpv36nyuw5a92hehea3jqaadss9smsqscr3lrp", // remove existing account 0
		// "elys173n2866wggue6znwl2vnwx9zqy7nnasjed9ydh",
	}
	filterBalanceAddresses := []string{
		"elys1gpv36nyuw5a92hehea3jqaadss9smsqscr3lrp", // remove existing account 0
		// "elys173n2866wggue6znwl2vnwx9zqy7nnasjed9ydh",
		authtypes.NewModuleAddress("distribution").String(),
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

	// distrAddr := authtypes.NewModuleAddress("distribution").String()

	// addressDenomMap := map[string][]string{
	// 	distrAddr: {"ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65", "ueden", "uedenb"},
	// }

	// genesis.AppState.Bank.Balances, coinsToRemove = FilterBalancesByDenoms(genesis.AppState.Bank.Balances, addressDenomMap)

	// // update supply
	// genesis.AppState.Bank.Supply = genesis.AppState.Bank.Supply.Sub(coinsToRemove...)

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

	// ColorReset distribution data
	genesis.AppState.Distribution = genesisInit.AppState.Distribution

	// temporary fix for distribution params
	genesis.AppState.Distribution.FeePool.CommunityPool = sdk.NewDecCoins(
		sdk.NewDecCoin("ueden", math.NewInt(595021147500)),
		sdk.NewDecCoin("uedenb", math.NewInt(1983399876344)),
	)

	log.Printf("community pool: %v", genesis.AppState.Distribution.FeePool.CommunityPool)

	// set genutil from genesisInit
	genesis.AppState.Genutil = genesisInit.AppState.Genutil

	// add localhost as allowed client
	genesis.AppState.Ibc.ClientGenesis.Params.AllowedClients = append(genesis.AppState.Ibc.ClientGenesis.Params.AllowedClients, "09-localhost")

	// reset gov as there are broken proposoals
	genesis.AppState.Gov = genesisInit.AppState.Gov

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

	// update broker address
	genesis.AppState.Parameter.Params.BrokerAddress = "elys1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrqau4f4q"

	// temporary fix for oracle
	// genesis.AppState.Oracle.Params = genesisInit.AppState.Oracle.Params
	// genesis.AppState.Oracle.PortId = genesisInit.AppState.Oracle.PortId
	// genesis.AppState.Oracle.Prices = genesisInit.AppState.Oracle.Prices
	// genesis.AppState.Oracle.PriceFeeders = genesisInit.AppState.Oracle.PriceFeeders
	// genesis.AppState.Oracle.AssetInfos = genesisInit.AppState.Oracle.AssetInfos

	// update oracle price expiration
	genesis.AppState.Oracle.Params.PriceExpiryTime = "31536000"
	genesis.AppState.Oracle.Params.LifeTimeInBlocks = "8000000"

	// update stablestake
	// genesis.AppState.StableStake = genesisInit.AppState.StableStake

	// update masterchef
	genesis.AppState.Masterchef = genesisInit.AppState.Masterchef

	outputFilePath := homePath + "/config/genesis.json"
	if err := WriteGenesisFile(outputFilePath, genesis); err != nil {
		log.Fatalf(types.ColorRed+"Error writing genesis file: %v", err)
	}
}
