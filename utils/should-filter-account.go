package utils

import (
	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func shouldFilterAccount(account types.Account, filterAddresses map[string]struct{}) bool {
	if account.BaseAccount != nil {
		if _, exists := filterAddresses[account.BaseAccount.Address]; exists {
			return true
		}
	}
	if account.ModuleAccount != nil {
		if _, exists := filterAddresses[account.ModuleAccount.BaseAccount.Address]; exists {
			return true
		}
	}
	return false
}
