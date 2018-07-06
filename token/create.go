package token

import (
	"github.com/Akagi201/eosgo/types"
)

func NewCreate(issuer types.AccountName, maxSupply types.Asset) *types.Action {
	return &types.Action{
		Account: types.AN("eosio.token"),
		Name:    types.ActN("create"),
		Authorization: []types.PermissionLevel{
			{Actor: types.AN("eosio.token"), Permission: types.PN("active")},
		},
		ActionData: types.NewActionData(Create{
			Issuer:        issuer,
			MaximumSupply: maxSupply,
		}),
	}
}

// Create represents the `create` struct on the `eosio.token` contract.
type Create struct {
	Issuer        types.AccountName `json:"issuer"`
	MaximumSupply types.Asset       `json:"maximum_supply"`
}
