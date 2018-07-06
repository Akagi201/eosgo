package token

import (
	"github.com/Akagi201/eosgo/types"
)

func NewTransfer(from, to types.AccountName, quantity types.Asset, memo string) *types.Action {
	return &types.Action{
		Account: types.AN("eosio.token"),
		Name:    types.ActN("transfer"),
		Authorization: []types.PermissionLevel{
			{Actor: from, Permission: types.PN("active")},
		},
		ActionData: types.NewActionData(Transfer{
			From:     from,
			To:       to,
			Quantity: quantity,
			Memo:     memo,
		}),
	}
}

// Transfer represents the `transfer` struct on `eosio.token` contract.
type Transfer struct {
	From     types.AccountName `json:"from"`
	To       types.AccountName `json:"to"`
	Quantity types.Asset       `json:"quantity"`
	Memo     string            `json:"memo"`
}
