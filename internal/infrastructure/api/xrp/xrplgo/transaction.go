package xrplgo

import (
	"context"
	"fmt"

	xrpl "github.com/xrpscan/xrpl-go"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
)

// submitResponse represents the response from submit command.
type submitResponse struct {
	EngineResult        string `json:"engine_result"`
	EngineResultCode    int    `json:"engine_result_code"`
	EngineResultMessage string `json:"engine_result_message"`
	TxBlob              string `json:"tx_blob"`
	TxJSON              struct {
		TransactionType    string `json:"TransactionType"`
		Account            string `json:"Account"`
		Amount             string `json:"Amount"`
		Destination        string `json:"Destination"`
		Fee                string `json:"Fee"`
		Flags              uint64 `json:"Flags"`
		LastLedgerSequence uint64 `json:"LastLedgerSequence"`
		Sequence           uint64 `json:"Sequence"`
		SigningPubKey      string `json:"SigningPubKey"`
		TxnSignature       string `json:"TxnSignature"`
		Hash               string `json:"hash"`
	} `json:"tx_json"`
}

// txResponse represents the response from tx command.
type txResponse struct {
	Account            string `json:"Account"`
	Amount             string `json:"Amount,omitempty"`
	Destination        string `json:"Destination,omitempty"`
	Fee                string `json:"Fee"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence,omitempty"`
	Sequence           uint64 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey"`
	TransactionType    string `json:"TransactionType"`
	TxnSignature       string `json:"TxnSignature"`
	Hash               string `json:"hash"`
	LedgerIndex        int    `json:"ledger_index"`
	Meta               struct {
		TransactionIndex  int    `json:"TransactionIndex"`
		TransactionResult string `json:"TransactionResult"`
		DeliveredAmount   string `json:"delivered_amount,omitempty"`
	} `json:"meta"`
	Validated bool   `json:"validated"`
	Date      uint64 `json:"date"`
}

// SubmitTransaction submits a signed transaction to the XRPL network.
// Implements TransactionSubmitter interface.
func (c *Client) SubmitTransaction(ctx context.Context, signedTx string) (*dtoxrp.SentTx, uint64, error) {
	req := xrpl.BaseRequest{
		"command": "submit",
		"tx_blob": signedTx,
	}

	resp, err := c.request(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to submit transaction: %w", err)
	}

	var result submitResponse
	if err := extractResult(resp, &result); err != nil {
		return nil, 0, fmt.Errorf("failed to parse submit response: %w", err)
	}

	sentTx := &dtoxrp.SentTx{
		ResultCode:          result.EngineResult,
		ResultMessage:       result.EngineResultMessage,
		EngineResult:        result.EngineResult,
		EngineResultCode:    result.EngineResultCode,
		EngineResultMessage: result.EngineResultMessage,
		TxBlob:              result.TxBlob,
		TxJSON: dtoxrp.TxInput{
			TransactionType:    result.TxJSON.TransactionType,
			Account:            result.TxJSON.Account,
			Amount:             result.TxJSON.Amount,
			Destination:        result.TxJSON.Destination,
			Fee:                result.TxJSON.Fee,
			Flags:              result.TxJSON.Flags,
			LastLedgerSequence: result.TxJSON.LastLedgerSequence,
			Sequence:           result.TxJSON.Sequence,
			SigningPubKey:      result.TxJSON.SigningPubKey,
			TxnSignature:       result.TxJSON.TxnSignature,
			Hash:               result.TxJSON.Hash,
		},
	}

	// Return the last ledger sequence for waiting
	return sentTx, result.TxJSON.LastLedgerSequence, nil
}

// GetTransaction retrieves transaction information from the XRPL.
// Implements TransactionGetter interface.
func (c *Client) GetTransaction(ctx context.Context, txID string, targetLedgerVersion uint64) (*dtoxrp.TxInfo, error) {
	req := xrpl.BaseRequest{
		"command":     "tx",
		"transaction": txID,
		"binary":      false,
	}

	if targetLedgerVersion > 0 {
		// Prevent underflow when targetLedgerVersion < 10
		minLedger := uint64(0)
		if targetLedgerVersion > 10 {
			minLedger = targetLedgerVersion - 10
		}
		req["min_ledger"] = minLedger
		req["max_ledger"] = targetLedgerVersion + 10
	}

	resp, err := c.request(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	var result txResponse
	if err := extractResult(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse transaction response: %w", err)
	}

	txInfo := &dtoxrp.TxInfo{
		Type:     result.TransactionType,
		Address:  result.Account,
		Sequence: int(result.Sequence),
		ID:       result.Hash,
		Specification: dtoxrp.TxSpecification{
			Source: dtoxrp.TxSpecSource{
				Address: result.Account,
				MaxAmount: dtoxrp.TxAmount{
					Currency: "XRP",
					Value:    dropsToXRP(result.Amount),
				},
			},
			Destination: dtoxrp.TxSpecDestination{
				Address: result.Destination,
			},
		},
		Outcome: dtoxrp.TxOutcome{
			Result:        result.Meta.TransactionResult,
			Fee:           dropsToXRP(result.Fee),
			LedgerVersion: result.LedgerIndex,
			IndexInLedger: result.Meta.TransactionIndex,
			DeliveredAmount: dtoxrp.TxAmount{
				Currency: "XRP",
				Value:    dropsToXRP(result.Meta.DeliveredAmount),
			},
		},
	}

	// Convert ripple time to ISO format
	if result.Date > 0 {
		txInfo.Outcome.Timestamp = xrpl.RippleTimeToISOTime(int64(result.Date))
	}

	return txInfo, nil
}
