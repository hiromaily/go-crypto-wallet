package ws_test

import (
	"context"
	"testing"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/ws"
)

// AccountChannels account_channels request
type AccountChannels struct {
	ID                 int    `json:"id"`
	Command            string `json:"command"`
	Account            string `json:"account"`
	DestinationAccount string `json:"destination_account"`
	LedgerIndex        string `json:"ledger_index"`
}

// ResponseAccountChannels account_channels response
type ResponseAccountChannels struct {
	ID     int `json:"id"`
	Result struct {
		Account     string        `json:"account"`
		Channels    []interface{} `json:"channels"`
		LedgerHash  string        `json:"ledger_hash"`
		LedgerIndex int           `json:"ledger_index"`
		Validated   bool          `json:"validated"`
	} `json:"result"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

// TestCall is test for Call
func TestCall(t *testing.T) {
	websoc := ws.New(context.Background(), xrp.PublicWSServerTestnet.String())

	ac := AccountChannels{
		ID:                 1,
		Command:            "account_channels",
		Account:            "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH",
		DestinationAccount: "rf1BiGeXwwQoi8Z2ueFYTEXSwuJYfV2Jpn",
		LedgerIndex:        "validated",
	}
	var res ResponseAccountChannels
	if err := websoc.Call(context.Background(), &ac, &res); err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
