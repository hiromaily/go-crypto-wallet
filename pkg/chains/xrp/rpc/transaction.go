package rpc

import (
	"context"
	"fmt"
)

// https://xrpl.org/transaction-methods.html

// SubmitRequest is the request payload for the submit command.
type SubmitRequest struct {
	ID      int    `json:"id"`
	Command string `json:"command"`
	TxBlob  string `json:"tx_blob"`
}

// SubmitTxJSON is the tx_json field in the submit response.
type SubmitTxJSON struct {
	Account            string `json:"Account"`
	Amount             string `json:"Amount"`
	Destination        string `json:"Destination"`
	Fee                string `json:"Fee"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence"`
	Sequence           uint64 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey"`
	TransactionType    string `json:"TransactionType"`
	TxnSignature       string `json:"TxnSignature"`
	Hash               string `json:"hash"`
}

// ResponseSubmit is the wire-format response of the submit command.
type ResponseSubmit struct {
	Result struct {
		Accepted             bool         `json:"accepted"`
		EngineResult         string       `json:"engine_result"`
		EngineResultCode     int          `json:"engine_result_code"`
		EngineResultMessage  string       `json:"engine_result_message"`
		TxBlob               string       `json:"tx_blob"`
		TxJSON               SubmitTxJSON `json:"tx_json"`
		ValidatedLedgerIndex uint64       `json:"validated_ledger_index"`
	} `json:"result"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Error  string `json:"error,omitempty"`
}

// Submit calls the submit WebSocket command and returns the raw wire response.
func Submit(ctx context.Context, caller WSCaller, txBlob string) (*ResponseSubmit, error) {
	req := &SubmitRequest{
		ID:      3,
		Command: "submit",
		TxBlob:  txBlob,
	}
	var res ResponseSubmit
	if err := caller.Call(ctx, req, &res); err != nil {
		return nil, fmt.Errorf("fail to call wsClient.Call(submit): %w", err)
	}
	return &res, nil
}

// TxRequest is the request payload for the tx command.
type TxRequest struct {
	ID          int    `json:"id"`
	Command     string `json:"command"`
	Transaction string `json:"transaction"`
	MinLedger   uint64 `json:"min_ledger,omitempty"`
	MaxLedger   uint64 `json:"max_ledger,omitempty"`
}

// TxMeta is the meta field in the tx response.
type TxMeta struct {
	TransactionResult string `json:"TransactionResult"`
	TransactionIndex  int    `json:"TransactionIndex"`
	DeliveredAmount   string `json:"delivered_amount,omitempty"`
}

// ResponseTx is the wire-format response of the tx command.
type ResponseTx struct {
	Result struct {
		Account            string `json:"Account"`
		Amount             string `json:"Amount"`
		Destination        string `json:"Destination"`
		Fee                string `json:"Fee"`
		Flags              uint64 `json:"Flags"`
		LastLedgerSequence uint64 `json:"LastLedgerSequence"`
		Sequence           uint64 `json:"Sequence"`
		TransactionType    string `json:"TransactionType"`
		Hash               string `json:"hash"`
		LedgerIndex        uint64 `json:"ledger_index"`
		Validated          bool   `json:"validated"`
		Meta               TxMeta `json:"meta"`
	} `json:"result"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Error  string `json:"error,omitempty"`
}

// GetTx calls the tx WebSocket command and returns the raw wire response.
func GetTx(ctx context.Context, caller WSCaller, txHash string, minLedger uint64) (*ResponseTx, error) {
	req := &TxRequest{
		ID:          4,
		Command:     "tx",
		Transaction: txHash,
		MinLedger:   minLedger,
	}
	var res ResponseTx
	if err := caller.Call(ctx, req, &res); err != nil {
		return nil, fmt.Errorf("fail to call wsClient.Call(tx): %w", err)
	}
	return &res, nil
}

// LedgerCurrentRequest is the request payload for the ledger_current command.
type LedgerCurrentRequest struct {
	ID      int    `json:"id"`
	Command string `json:"command"`
}

// ResponseLedgerCurrent is the wire-format response of the ledger_current command.
type ResponseLedgerCurrent struct {
	Result struct {
		LedgerCurrentIndex uint64 `json:"ledger_current_index"`
	} `json:"result"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Error  string `json:"error,omitempty"`
}

// LedgerCurrent calls the ledger_current WebSocket command and returns the current open ledger index.
func LedgerCurrent(ctx context.Context, caller WSCaller) (*ResponseLedgerCurrent, error) {
	req := &LedgerCurrentRequest{
		ID:      5,
		Command: "ledger_current",
	}
	var res ResponseLedgerCurrent
	if err := caller.Call(ctx, req, &res); err != nil {
		return nil, fmt.Errorf("fail to call wsClient.Call(ledger_current): %w", err)
	}
	return &res, nil
}

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
func LedgerAccept(ctx context.Context, caller WSCaller) (*ResponseLedgerAccept, error) {
	req := &LedgerAcceptRequest{
		ID:      6,
		Command: "ledger_accept",
	}
	var res ResponseLedgerAccept
	if err := caller.Call(ctx, req, &res); err != nil {
		return nil, fmt.Errorf("fail to call wsAdmin.Call(ledger_accept): %w", err)
	}
	return &res, nil
}
