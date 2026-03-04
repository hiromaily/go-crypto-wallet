package xrp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	xrpsigner "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/signer"
	xrpkg "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp"
	xrpclient "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/client"
	"github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/protogen"
	xrprpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// - Send XRP https://xrpl.org/send-xrp.html
// - Payment System Basics https://xrpl.org/payment-system-basics.html

// unquoteJSON attempts to unquote a JSON string. If unquoting fails,
// it returns the original string. This handles cases where the gRPC
// response may or may not include extra quotes around the JSON.
func unquoteJSON(s string) string {
	unquoted, err := strconv.Unquote(s)
	if err != nil {
		// If unquoting fails, assume the string is not quoted
		return s
	}
	return unquoted
}

// TxInput is transaction input json type
type TxInput struct {
	TransactionType    string `json:"TransactionType"`
	Account            string `json:"Account"`
	Amount             string `json:"Amount"`
	Destination        string `json:"Destination"`
	Fee                string `json:"Fee"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence"`
	Sequence           uint64 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey,omitempty"`
	TxnSignature       string `json:"TxnSignature,omitempty"`
	Hash               string `json:"hash,omitempty"`
}

// SentTx is result transaction json type after sending
type SentTx struct {
	ResultCode          string  `json:"resultCode"`
	ResultMessage       string  `json:"resultMessage"`
	EngineResult        string  `json:"engine_result"`
	EngineResultCode    int     `json:"engine_result_code"`
	EngineResultMessage string  `json:"engine_result_message"`
	TxBlob              string  `json:"tx_blob"`
	TxJSON              TxInput `json:"tx_json"`
}

// TxInfo is result transaction json type after sending
type TxInfo struct {
	Type          string          `json:"type"`
	Address       string          `json:"address"`
	Sequence      int             `json:"sequence"`
	ID            string          `json:"id"`
	Specification TxSpecification `json:"specification"`
	Outcome       TxOutcome       `json:"outcome"`
}

// TxSpecification is part of TxInfo
type TxSpecification struct {
	Source      TxSpecSource      `json:"source"`
	Destination TxSpecDestination `json:"destination"`
}

// TxSpecSource is part of TxInfo
type TxSpecSource struct {
	Address   string   `json:"address"`
	MaxAmount TxAmount `json:"maxAmount"`
}

// TxAmount is part of TxInfo
type TxAmount struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

// TxTotalPrice is part of TxInfo
type TxTotalPrice struct {
	Currency     string `json:"currency"`
	Counterparty string `json:"counterparty"`
	Value        string `json:"value"`
}

// TxSpecDestination is part of TxInfo
type TxSpecDestination struct {
	Address string `json:"address"`
}

// TxOutcome is part of TxInfo
type TxOutcome struct {
	Result           string                         `json:"result"`
	Timestamp        time.Time                      `json:"timestamp"`
	Fee              string                         `json:"fee"`
	BalanceChanges   map[string][]TxAmount          `json:"balanceChanges"`
	OrderbookChanges map[string][]TxOrderbookChange `json:"orderbookChanges"`
	LedgerVersion    int                            `json:"ledgerVersion"`
	IndexInLedger    int                            `json:"indexInLedger"`
	DeliveredAmount  TxAmount                       `json:"deliveredAmount"`
}

// TxOrderbookChange is part of TxInfo
type TxOrderbookChange struct {
	Direction         string       `json:"direction"`
	Quantity          TxAmount     `json:"quantity"`
	TotalPrice        TxTotalPrice `json:"totalPrice"`
	MakerExchangeRate string       `json:"makerExchangeRate"`
	Sequence          int          `json:"sequence"`
	Status            string       `json:"status"`
}

// PrepareTransaction builds an unsigned Payment transaction using WebSocket account_info.
func (w *WSClient) PrepareTransaction(
	ctx context.Context, senderAccount, receiverAccount string, amount float64, instructions *dtoxrp.Instructions,
) (*dtoxrp.TxInput, string, error) {
	// Get account info for sequence number and current ledger index
	accInfo, err := xrprpc.AccountInfo(ctx, w.public, senderAccount)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction(): %w", err)
	}
	if accInfo.Error != "" {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction(): %s", accInfo.Error)
	}

	sequence := uint64(accInfo.Result.AccountData.Sequence)
	lastLedgerSequence := uint64(accInfo.Result.LedgerCurrentIndex) + xrpkg.MaxLedgerVersionOffset
	fee := "12" // minimum fee in drops

	if instructions != nil {
		if instructions.Fee != "" {
			fee = instructions.Fee
		}
		if instructions.Sequence != 0 {
			sequence = instructions.Sequence
		}
		if instructions.MaxLedgerVersion != 0 {
			lastLedgerSequence = instructions.MaxLedgerVersion
		} else if instructions.MaxLedgerVersionOffset != 0 {
			lastLedgerSequence = uint64(accInfo.Result.LedgerCurrentIndex) + instructions.MaxLedgerVersionOffset
		}
	}

	// Convert XRP amount to drops (XRPL requires Amount in drops for XRP payments)
	amountDrops := strconv.FormatInt(int64(amount*float64(dropsPerXRP)), 10)

	txInput := &TxInput{
		TransactionType:    "Payment",
		Account:            senderAccount,
		Amount:             amountDrops,
		Destination:        receiverAccount,
		Fee:                fee,
		Flags:              0,
		LastLedgerSequence: lastLedgerSequence,
		Sequence:           sequence,
	}

	jsonBytes, err := json.Marshal(txInput)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call json.Marshal(txInput): %w", err)
	}

	logger.Debug("PrepareTransaction", "TxJSON", string(jsonBytes))

	// Convert infrastructure type to DTO
	return ToDTOTxInput(txInput), string(jsonBytes), nil
}

// signTransactionJSON is a generic helper that signs any transaction type.
// It marshals the input to JSON and calls the gRPC SignTransaction API.
func (r *XRP) signTransactionJSON(
	ctx context.Context, txInput any, secret, txTypeName string,
) (string, string, error) {
	strJSON, err := json.Marshal(txInput)
	if err != nil {
		return "", "", fmt.Errorf("fail to call json.Marshal(%sTxInput): %w", txTypeName, err)
	}
	req := protogen.RequestSignTransaction_builder{
		TxJSON: string(strJSON),
		Secret: secret,
	}.Build()

	res, err := r.API.txClient.SignTransaction(ctx, req)
	if err != nil {
		return "", "", fmt.Errorf("fail to call client.SignTransaction() for %s: %w", txTypeName, err)
	}

	return res.GetTxID(), res.GetTxBlob(), nil
}

// SignTransaction signs a payment transaction offline using native Go implementation.
// Offline functionality
// - https://xrpl.org/rippleapi-reference.html#offline-functionality
func (*XRP) SignTransaction(
	ctx context.Context, txInput *dtoxrp.TxInput, secret string,
) (string, string, error) {
	s := xrpsigner.NewPeersystSigner()
	return s.SignTransactionNative(ctx, txInput, secret, false, nil)
}

// SignTransactionNative signs a transaction using native Go implementation (Peersyst/xrpl-go).
// This method provides offline signing capability without gRPC dependencies.
//
// NOTE: This is a stub implementation for task 1.2 (interface segregation).
// Full implementation will be provided in task 3.1 (PeersystSigner).
//
// TODO(task-3.1): Implement native Go signing using Peersyst/xrpl-go library.
func (*XRP) SignTransactionNative(
	_ context.Context,
	_ *dtoxrp.TxInput,
	_ string,
	_ bool,
	_ *string,
) (string, string, error) {
	// Stub implementation - to be completed in task 3.1
	return "", "", errors.New("SignTransactionNative not yet implemented (task 3.1)")
}

// CombineTransaction combines signed transactions from multiple accounts for a multisignature transaction.
// - The signed transaction must subsequently be submitted.
func (r *XRP) CombineTransaction(ctx context.Context, signedTxs []string) (string, string, error) {
	req := protogen.RequestCombineTransaction_builder{
		SignedTransactions: signedTxs,
	}.Build()

	res, err := r.API.txClient.CombineTransaction(ctx, req)
	if err != nil {
		return "", "", fmt.Errorf("fail to call client.CombineTransaction(): %w", err)
	}

	return res.GetTxID(), res.GetSignedTransaction(), nil
}

// toXRPClientSentTx converts local SentTx to xrpclient.SentTx.
func toXRPClientSentTx(local *SentTx) *xrpclient.SentTx {
	return &xrpclient.SentTx{
		ResultCode:          local.ResultCode,
		ResultMessage:       local.ResultMessage,
		EngineResult:        local.EngineResult,
		EngineResultCode:    local.EngineResultCode,
		EngineResultMessage: local.EngineResultMessage,
		TxBlob:              local.TxBlob,
		TxJSON: xrpclient.TxInput{
			TransactionType:    local.TxJSON.TransactionType,
			Account:            local.TxJSON.Account,
			Amount:             local.TxJSON.Amount,
			Destination:        local.TxJSON.Destination,
			Fee:                local.TxJSON.Fee,
			Flags:              local.TxJSON.Flags,
			LastLedgerSequence: local.TxJSON.LastLedgerSequence,
			Sequence:           local.TxJSON.Sequence,
			SigningPubKey:      local.TxJSON.SigningPubKey,
			TxnSignature:       local.TxJSON.TxnSignature,
			Hash:               local.TxJSON.Hash,
		},
	}
}

// SubmitTransaction submits a signed transaction blob via WebSocket.
// - signedTx is the TxBlob returned by SignTransaction()
func (w *WSClient) SubmitTransaction(ctx context.Context, signedTx string) (*xrpclient.SentTx, uint64, error) {
	res, err := xrprpc.Submit(ctx, w.public, signedTx)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to call client.SubmitTransaction(): %w", err)
	}
	if res.Error != "" {
		return nil, 0, fmt.Errorf("fail to call client.SubmitTransaction(): %s", res.Error)
	}

	logger.Debug("response of submitTransaction",
		"engine_result", res.Result.EngineResult,
		"engine_result_code", res.Result.EngineResultCode,
		"LastLedgerSequence", res.Result.TxJSON.LastLedgerSequence,
		"validated_ledger_index", res.Result.ValidatedLedgerIndex,
	)

	sentTx := &SentTx{
		// Map engine_result to ResultCode so the use case check works
		ResultCode:          res.Result.EngineResult,
		ResultMessage:       res.Result.EngineResultMessage,
		EngineResult:        res.Result.EngineResult,
		EngineResultCode:    res.Result.EngineResultCode,
		EngineResultMessage: res.Result.EngineResultMessage,
		TxBlob:              res.Result.TxBlob,
		TxJSON: TxInput{
			TransactionType:    res.Result.TxJSON.TransactionType,
			Account:            res.Result.TxJSON.Account,
			Amount:             res.Result.TxJSON.Amount,
			Destination:        res.Result.TxJSON.Destination,
			Fee:                res.Result.TxJSON.Fee,
			Flags:              res.Result.TxJSON.Flags,
			LastLedgerSequence: res.Result.TxJSON.LastLedgerSequence,
			Sequence:           res.Result.TxJSON.Sequence,
			SigningPubKey:      res.Result.TxJSON.SigningPubKey,
			TxnSignature:       res.Result.TxJSON.TxnSignature,
			Hash:               res.Result.TxJSON.Hash,
		},
	}

	// ValidatedLedgerIndex is the earliest ledger where we can look for this tx
	return toXRPClientSentTx(sentTx), res.Result.ValidatedLedgerIndex, nil
}

// WaitValidation waits until the ledger advances past targetLedgerVersion.
// In standalone mode it calls ledger_accept via the admin WebSocket to advance the ledger.
// In non-standalone environments the ledger advances naturally and this polls ledger_current.
func (w *WSClient) WaitValidation(ctx context.Context, targetLedgerVersion uint64) (uint64, error) {
	// Advance the ledger via the admin connection (standalone mode).
	// Silently ignore errors — in production environments this command is unavailable.
	if w.admin != nil {
		if _, err := xrprpc.LedgerAccept(ctx, w.admin); err != nil {
			logger.Warn("ledger_accept failed (non-critical; may not be standalone mode)", "error", err)
		}
	}

	// Poll ledger_current until it reaches or exceeds targetLedgerVersion.
	// In standalone mode each iteration calls ledger_accept to advance the ledger.
	const maxRetries = 30
	for range maxRetries {
		res, err := xrprpc.LedgerCurrent(ctx, w.public)
		if err != nil {
			return 0, fmt.Errorf("fail to call ledger_current: %w", err)
		}
		currentLedger := res.Result.LedgerCurrentIndex
		logger.Info("WaitValidation polling",
			"currentLedger", currentLedger,
			"targetLedgerVersion", targetLedgerVersion,
		)
		if currentLedger >= targetLedgerVersion {
			return currentLedger, nil
		}
		// Advance the ledger in standalone mode and retry.
		if w.admin != nil {
			if _, err := xrprpc.LedgerAccept(ctx, w.admin); err != nil {
				logger.Warn("failed to advance ledger", "error", err)
			}
		}
		time.Sleep(1 * time.Second)
	}
	return 0, errors.New("timeout waiting for ledger validation")
}

// GetTransaction retrieves a validated transaction by hash via WebSocket.
func (w *WSClient) GetTransaction(
	ctx context.Context, txID string, targetLedgerVersion uint64,
) (*xrpclient.TxInfo, error) {
	res, err := xrprpc.GetTx(ctx, w.public, txID, targetLedgerVersion)
	if err != nil {
		return nil, fmt.Errorf("fail to call client.GetTransaction(): %w", err)
	}
	if res.Error != "" {
		return nil, fmt.Errorf("fail to get transaction info by %s: %s", txID, res.Error)
	}

	logger.Debug("response of getTransaction",
		"hash", res.Result.Hash,
		"ledger_index", res.Result.LedgerIndex,
		"validated", res.Result.Validated,
		"TransactionResult", res.Result.Meta.TransactionResult,
	)

	return &xrpclient.TxInfo{
		ID: res.Result.Hash,
		Outcome: xrpclient.TxOutcome{
			Result:        res.Result.Meta.TransactionResult,
			LedgerVersion: int(res.Result.LedgerIndex),
		},
	}, nil
}
