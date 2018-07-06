package token

import (
	"github.com/Akagi201/eosgo/types"
)

func NewIssue(to types.AccountName, quantity types.Asset, memo string) *types.Action {
	return &types.Action{
		Account: types.AN("eosio.token"),
		Name:    types.ActN("issue"),
		Authorization: []types.PermissionLevel{
			{Actor: types.AN("eosio"), Permission: types.PN("active")},
		},
		ActionData: types.NewActionData(Issue{
			To:       to,
			Quantity: quantity,
			Memo:     memo,
		}),
	}
}

// Issue represents the `issue` struct on the `eosio.token` contract.
type Issue struct {
	To       types.AccountName `json:"to"`
	Quantity types.Asset       `json:"quantity"`
	Memo     string            `json:"memo"`
}
