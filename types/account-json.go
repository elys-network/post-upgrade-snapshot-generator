package types

import (
	"encoding/json"
	"fmt"
)

func (a *Account) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set the Type field from the raw data
	typeStr, ok := raw["@type"].(string)
	if !ok {
		return fmt.Errorf("type field is missing or invalid")
	}
	a.Type = typeStr

	switch a.Type {
	case "/cosmos.vesting.v1beta1.PeriodicVestingAccount":
		var va VestingAccount
		if err := json.Unmarshal(data, &va); err != nil {
			return err
		}
		a.VestingAccount = &va
	case "/cosmos.vesting.v1beta1.ContinuousVestingAccount":
		var va VestingAccount
		if err := json.Unmarshal(data, &va); err != nil {
			return err
		}
		a.VestingAccount = &va
	case "/cosmos.auth.v1beta1.BaseAccount":
		var ba BaseAccount
		if err := json.Unmarshal(data, &ba); err != nil {
			return err
		}
		a.BaseAccount = &ba
	case "/cosmos.auth.v1beta1.ModuleAccount":
		var ma ModuleAccount
		if err := json.Unmarshal(data, &ma); err != nil {
			return err
		}
		a.ModuleAccount = &ma
	case "/ibc.applications.interchain_accounts.v1.InterchainAccount":
		var ica InterchainAccount
		if err := json.Unmarshal(data, &ica); err != nil {
			return err
		}
		a.InterchainAccount = &ica
	default:
		return fmt.Errorf("unknown account type found during unmarshalling: %s", a.Type)
	}
	return nil
}

func (a Account) MarshalJSON() ([]byte, error) {
	// Helper function to add the "@type" field
	addTypeField := func(aux interface{}) ([]byte, error) {
		type wrapper struct {
			Type string `json:"@type"`
			Data interface{}
		}
		return json.Marshal(wrapper{
			Type: a.Type,
			Data: aux,
		})
	}

	// Delegate marshalling based on the Type field
	switch a.Type {
	case "/cosmos.vesting.v1beta1.PeriodicVestingAccount",
		"/cosmos.vesting.v1beta1.ContinuousVestingAccount":
		if a.VestingAccount == nil {
			return nil, fmt.Errorf("missing VestingAccount data")
		}
		return addTypeField(*a.VestingAccount)

	case "/cosmos.auth.v1beta1.BaseAccount":
		if a.BaseAccount == nil {
			return nil, fmt.Errorf("missing BaseAccount data")
		}
		return addTypeField(*a.BaseAccount)

	case "/cosmos.auth.v1beta1.ModuleAccount":
		if a.ModuleAccount == nil {
			return nil, fmt.Errorf("missing ModuleAccount data")
		}
		return addTypeField(*a.ModuleAccount)

	case "/ibc.applications.interchain_accounts.v1.InterchainAccount":
		if a.InterchainAccount == nil {
			return nil, fmt.Errorf("missing InterchainAccount data")
		}
		return addTypeField(*a.InterchainAccount)

	default:
		return nil, fmt.Errorf("unknown account type found during Marshalling: %s", a.Type)
	}
}
