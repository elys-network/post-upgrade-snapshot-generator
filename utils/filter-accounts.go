package utils

import (
	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func FilterAccounts(accounts []types.Account, filterAddresses []string) []types.Account {
	filterMap := make(map[string]struct{})
	for _, addr := range filterAddresses {
		filterMap[addr] = struct{}{}
	}

	newAccounts := []types.Account{}
	for _, account := range accounts {
		if shouldFilterAccount(account, filterMap) {
			continue
		}
		newAccounts = append(newAccounts, account)
	}
	return newAccounts
}
