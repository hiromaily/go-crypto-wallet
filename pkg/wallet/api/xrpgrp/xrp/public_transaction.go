package xrp

import (
	"context"

	"github.com/pkg/errors"
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

// As one of workaround
// - create node.js server to run ripple-lib

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
func (r *Ripple) Sign(txJSON *PublicTxType, secret string, offline bool) (*ResponseSign, error) {
	req := Sign{
		ID:      2,
		Command: "sign",
		TxJSON:  *txJSON,
		Secret:  secret,
		Offline: offline,
	}
	var res ResponseSign
	if err := r.wsPublic.Call(context.Background(), &req, &res); err != nil {
		return nil, errors.Wrap(err, "fail to call wsClient.Call(sign)")
	}
	return &res, nil
}
