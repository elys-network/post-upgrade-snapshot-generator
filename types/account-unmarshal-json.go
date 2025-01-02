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
	default:
		return fmt.Errorf("unknown account type: %s", a.Type)
	}
	return nil
}
