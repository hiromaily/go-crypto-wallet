package public

import (
	"context"
	"fmt"
)

// https://xrpl.org/account-methods.html

// AccountChannelsRequest is the request payload for the account_channels command.
type AccountChannelsRequest struct {
	ID                 int    `json:"id"`
	Command            string `json:"command"`
	Account            string `json:"account"`
	DestinationAccount string `json:"destination_account"`
	LedgerIndex        string `json:"ledger_index"`
}

// ResponseAccountChannels is the wire-format response of the account_channels command.
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
			CancelAfter        int    `json:"cancel_after"`
			Expiration         int    `json:"expiration"`
			PublicKey          string `json:"public_key"`
			PublicKeyHex       string `json:"public_key_hex"`
			SettleDelay        int    `json:"settle_delay"`
		} `json:"channels"`
		LedgerHash  string `json:"ledger_hash"`
		LedgerIndex int    `json:"ledger_index"`
		Validated   bool   `json:"validated"`
	} `json:"result"`
	Error string `json:"error,omitempty"`
}

// AccountInfoRequest is the request payload for the account_info command.
type AccountInfoRequest struct {
	ID          int    `json:"id"`
	Command     string `json:"command"`
	Account     string `json:"account"`
	Strict      bool   `json:"strict"`
	LedgerIndex string `json:"ledger_index"`
	Queue       bool   `json:"queue"`
}

// AccountInfo holds the account data returned by the account_info command.
type AccountInfo struct {
	AccountData struct {
		Account           string `json:"Account"`
		Balance           string `json:"Balance"`
		Flags             int    `json:"Flags"`
		LedgerEntryType   string `json:"LedgerEntryType"`
		OwnerCount        int    `json:"OwnerCount"`
		PreviousTxnID     string `json:"PreviousTxnID"`
		PreviousTxnLgrSeq int    `json:"PreviousTxnLgrSeq"`
		Sequence          int    `json:"Sequence"`
		Index             string `json:"index"`
	} `json:"account_data"`
	LedgerCurrentIndex int `json:"ledger_current_index"`
	QueueData          struct {
		TxnCount int `json:"txn_count"`
	} `json:"queue_data"`
	Validated bool `json:"validated"`
}

// ResponseAccountInfo is the wire-format response of the account_info command.
type ResponseAccountInfo struct {
	ID     int         `json:"id"`
	Result AccountInfo `json:"result"`
	Status string      `json:"status"`
	Type   string      `json:"type"`
	Error  string      `json:"error,omitempty"`
}

// AccountChannels calls the account_channels WebSocket command and returns the raw wire response.
func (r *PublicRPC) AccountChannels(
	ctx context.Context, sender, receiver string,
) (*ResponseAccountChannels, error) {
	req := &AccountChannelsRequest{
		ID:                 1,
		Command:            "account_channels",
		Account:            sender,
		DestinationAccount: receiver,
		LedgerIndex:        "validated",
	}
	var res ResponseAccountChannels
	if err := r.caller.Call(ctx, req, &res); err != nil {
		return nil, fmt.Errorf("fail to call wsClient.Call(account_channels): %w", err)
	}
	return &res, nil
}

// AccountInfo calls the account_info WebSocket command and returns the raw wire response.
func (r *PublicRPC) AccountInfo(ctx context.Context, address string) (*ResponseAccountInfo, error) {
	req := &AccountInfoRequest{
		ID:          2,
		Command:     "account_info",
		Account:     address,
		Strict:      true,
		LedgerIndex: "current",
		Queue:       true,
	}
	var res ResponseAccountInfo
	if err := r.caller.Call(ctx, req, &res); err != nil {
		return nil, fmt.Errorf("fail to call wsClient.Call(account_info): %w", err)
	}
	return &res, nil
}
