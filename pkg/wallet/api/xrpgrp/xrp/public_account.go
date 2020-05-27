package xrp

import (
	"context"

	"github.com/pkg/errors"
)

// https://xrpl.org/account-methods.html
// error: https://xrpl.org/error-formatting.html#universal-errors

// AccountChannels is request data for account_channels method
type AccountChannels struct {
	ID                 int    `json:"id"`
	Command            string `json:"command"`
	Account            string `json:"account"`
	DestinationAccount string `json:"destination_account"`
	LedgerIndex        string `json:"ledger_index"`
}

// ResponseAccountChannels is response data for account_channels method
type ResponseAccountChannels struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result struct {
		Account  string `json:"account"`
		Channels []struct {
			Account            string `json:"account"`
			Amount             string `json:"amount"`
			Balance            string `json:"balance"`
			ChannelID          string `json:"channel_id"`
			DestinationAccount string `json:"destination_account"`
			DestinationTag     int    `json:"destination_tag"`
			Expiration         int    `json:"expiration"`
			PublicKey          string `json:"public_key"`
			PublicKeyHex       string `json:"public_key_hex"`
			SettleDelay        int    `json:"settle_delay"`
		} `json:"channels"`
		LedgerHash  string `json:"ledger_hash"`
		LedgerIndex int    `json:"ledger_index"`
		Validated   bool   `json:"validated"`
	} `json:"result"`
}

// AccountChannels calls account_channels method
func (r *Ripple) AccountChannels(sender, receiver string) (*ResponseAccountChannels, error) {
	req := AccountChannels{
		ID:                 1,
		Command:            "account_channels",
		Account:            sender,
		DestinationAccount: receiver,
		LedgerIndex:        "validated",
	}
	var res ResponseAccountChannels
	if err := r.wsPublic.Call(context.Background(), &req, &res); err != nil {
		return nil, errors.Wrap(err, "fail to call wsClient.Call()")
	}
	return &res, nil
}
