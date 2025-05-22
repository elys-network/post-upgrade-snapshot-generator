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
	var base interface{}

	switch a.Type {
	case "/cosmos.vesting.v1beta1.PeriodicVestingAccount",
		"/cosmos.vesting.v1beta1.ContinuousVestingAccount":
		if a.VestingAccount == nil {
			return nil, fmt.Errorf("vesting account is nil for type %s", a.Type)
		}
		base = a.VestingAccount

	case "/cosmos.auth.v1beta1.BaseAccount":
		if a.BaseAccount == nil {
			return nil, fmt.Errorf("base account is nil for type %s", a.Type)
		}
		base = a.BaseAccount

	case "/cosmos.auth.v1beta1.ModuleAccount":
		if a.ModuleAccount == nil {
			return nil, fmt.Errorf("module account is nil for type %s", a.Type)
		}
		base = a.ModuleAccount

	case "/ibc.applications.interchain_accounts.v1.InterchainAccount":
		if a.InterchainAccount == nil {
			return nil, fmt.Errorf("interchain account is nil for type %s", a.Type)
		}
		base = a.InterchainAccount

	default:
		return nil, fmt.Errorf("unknown account type found during Marshalling: %s", a.Type)
	}

	// Marshal the underlying account struct to JSON
	bz, err := json.Marshal(base)
	if err != nil {
		return nil, err
	}

	// Convert to map to inject @type
	var raw map[string]interface{}
	if err := json.Unmarshal(bz, &raw); err != nil {
		return nil, err
	}
	raw["@type"] = a.Type

	return json.Marshal(raw)
}
