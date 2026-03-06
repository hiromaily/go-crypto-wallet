package public

import (
	"context"
	"fmt"
)

// https://xrpl.org/transaction-methods.html
// 1. Connect to a Test Net Server
// 2. Prepare Transaction
//{
//  "TransactionType": "Payment",
//  "Account": "rPT1Sjq2YGrBMTttX4GZHjKu9dyfzbpAYe",
//  "Amount": "2000000",
//  "Destination": "rUCzEr6jrEyMpjhs4wSdQdz4g8Y382NxfM"
//}
// Transaction Common Fields
// - https://xrpl.org/transaction-common-fields.html
// - [Required]
//   - Account
//   - TransactionType
//   - Fee (auto-fillable)
//   - Sequence (auto-fillable)
//   - LastLedgerSequence (strongly recommended)

// Is there alternative way of prepareTransaction()??
// - LastLedgerSequence look important, how is it retrieved?
// - ledger command https://xrpl.org/ledger.html
// - Is `ledger_index` LastLedgerSequence??

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
func (r *PublicRPC) Submit(ctx context.Context, txBlob string) (*ResponseSubmit, error) {
	req := &SubmitRequest{
		ID:      3,
		Command: "submit",
		TxBlob:  txBlob,
	}
	var res ResponseSubmit
	if err := r.caller.Call(ctx, req, &res); err != nil {
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
func (r *PublicRPC) GetTx(ctx context.Context, txHash string, minLedger uint64) (*ResponseTx, error) {
	req := &TxRequest{
		ID:          4,
		Command:     "tx",
		Transaction: txHash,
		MinLedger:   minLedger,
	}
	var res ResponseTx
	if err := r.caller.Call(ctx, req, &res); err != nil {
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
func (r *PublicRPC) LedgerCurrent(ctx context.Context) (*ResponseLedgerCurrent, error) {
	req := &LedgerCurrentRequest{
		ID:      5,
		Command: "ledger_current",
	}
	var res ResponseLedgerCurrent
	if err := r.caller.Call(ctx, req, &res); err != nil {
		return nil, fmt.Errorf("fail to call wsClient.Call(ledger_current): %w", err)
	}
	return &res, nil
}

// Sign is request data for sign method
type Sign struct {
	ID         int          `json:"id"`
	Command    string       `json:"command"`
	TxJSON     PublicTxType `json:"tx_json"`
	Secret     string       `json:"secret"`
	Offline    bool         `json:"offline"`
	FeeMultMax int          `json:"fee_mult_max"`
}

// PublicTxType is part of Sign request
type PublicTxType struct {
	TransactionType string `json:"TransactionType"`
	Account         string `json:"Account"`
	Destination     string `json:"Destination"`
	Amount          Amount `json:"Amount"`
}

// Amount is part of Sign request
type Amount struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
	Issuer   string `json:"issuer"`
}

// ResponseSign is response data for sign method
type ResponseSign struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result struct {
		TxBlob string `json:"tx_blob"`
		TxJSON struct {
			Account string `json:"Account"`
			Amount  struct {
				Currency string `json:"currency"`
				Issuer   string `json:"issuer"`
				Value    string `json:"value"`
			} `json:"Amount"`
			Destination     string `json:"Destination"`
			Fee             string `json:"Fee"`
			Flags           int64  `json:"Flags"`
			Sequence        int    `json:"Sequence"`
			SigningPubKey   string `json:"SigningPubKey"`
			TransactionType string `json:"TransactionType"`
			TxnSignature    string `json:"TxnSignature"`
			Hash            string `json:"hash"`
		} `json:"tx_json"`
	} `json:"result"`
	Error string `json:"error,omitempty"`
}

// Sign calls sign method
// TODO: add Sign() to proper interface
func (r *PublicRPC) Sign(
	ctx context.Context, txJSON *PublicTxType, secret string, offline bool,
) (*ResponseSign, error) {
	req := Sign{
		ID:      2,
		Command: "sign",
		TxJSON:  *txJSON,
		Secret:  secret,
		Offline: offline,
	}
	var res ResponseSign
	if err := r.caller.Call(ctx, &req, &res); err != nil {
		return nil, fmt.Errorf("fail to call wsClient.Call(sign): %w", err)
	}
	return &res, nil
}
