package admin

import (
	"context"
	"fmt"
)

// https://xrpl.org/transaction-methods.html

// LedgerAcceptRequest is the request payload for the ledger_accept admin command.
type LedgerAcceptRequest struct {
	ID      int    `json:"id"`
	Command string `json:"command"`
}

// ResponseLedgerAccept is the wire-format response of the ledger_accept admin command.
type ResponseLedgerAccept struct {
	Result struct {
		LedgerCurrentIndex uint64 `json:"ledger_current_index"`
	} `json:"result"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Error  string `json:"error,omitempty"`
}

// LedgerAccept calls the ledger_accept admin WebSocket command to advance the ledger in standalone mode.
// This is only available on admin connections and should be used only in standalone/testing environments.
func (r *AdminRPC) LedgerAccept(ctx context.Context) (*ResponseLedgerAccept, error) {
	req := &LedgerAcceptRequest{
		ID:      6,
		Command: "ledger_accept",
	}
	var res ResponseLedgerAccept
	if err := r.caller.Call(ctx, req, &res); err != nil {
		return nil, fmt.Errorf("fail to call wsAdmin.Call(ledger_accept): %w", err)
	}
	return &res, nil
}
