package websocket_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ripple/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/network/websocket"
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
		Account     string `json:"account"`
		Channels    []any  `json:"channels"`
		LedgerHash  string `json:"ledger_hash"`
		LedgerIndex int    `json:"ledger_index"`
		Validated   bool   `json:"validated"`
	} `json:"result"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

// TestCall is test for Call
func TestCall(t *testing.T) {
	websoc, err := websocket.New(context.Background(), xrp.PublicWSServerTestnet.String())
	require.NoError(t, err, "ws.New() should not return error")

	ac := AccountChannels{
		ID:                 1,
		Command:            "account_channels",
		Account:            "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH",
		DestinationAccount: "rf1BiGeXwwQoi8Z2ueFYTEXSwuJYfV2Jpn",
		LedgerIndex:        "validated",
	}
	var res ResponseAccountChannels
	err = websoc.Call(context.Background(), &ac, &res)
	require.NoError(t, err, "websoc.Call() should not return error")
	t.Log(res)

	_ = websoc.Close() // Best effort cleanup
}
