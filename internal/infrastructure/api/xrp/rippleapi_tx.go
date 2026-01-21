package xrp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	dtoRipple "github.com/hiromaily/go-crypto-wallet/internal/application/dto/ripple"
	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/protogen"
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

// PrepareTransaction calls PrepareTransaction API
func (r *Ripple) PrepareTransaction(
	ctx context.Context, senderAccount, receiverAccount string, amount float64, instructions *dtoRipple.Instructions,
) (*dtoRipple.TxInput, string, error) {
	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	req := protogen.RequestPrepareTransaction_builder{
		TxType:          protogen.EnumTransactionType_TX_PAYMENT,
		SenderAccount:   senderAccount,
		Amount:          amount,
		ReceiverAccount: receiverAccount,
		Instructions:    infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction(): %w", err)
	}
	logger.Debug("response",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput TxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(txJSON): %w", err)
	}

	// Convert infrastructure type to DTO
	return ToDTOTxInput(&txInput), unquotedJSON, nil
}

// SetRegularKeyTxInput is the transaction input for SetRegularKey
type SetRegularKeyTxInput struct {
	TransactionType    string `json:"TransactionType"`
	Account            string `json:"Account"`
	RegularKey         string `json:"RegularKey,omitempty"` // Omit to remove regular key
	Fee                string `json:"Fee"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence"`
	Sequence           uint64 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey,omitempty"`
	TxnSignature       string `json:"TxnSignature,omitempty"`
	Hash               string `json:"hash,omitempty"`
}

// PrepareSetRegularKeyTransaction prepares a SetRegularKey transaction
// - regularKey: the address to authorize as regular key, or empty to remove
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/setregularkey
func (r *Ripple) PrepareSetRegularKeyTransaction(
	ctx context.Context, senderAccount, regularKey string, instructions *dtoRipple.Instructions,
) (*dtoRipple.SetRegularKeyTxInput, string, error) {
	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	req := protogen.RequestPrepareTransaction_builder{
		TxType:        protogen.EnumTransactionType_TX_SET_REGULAR_KEY,
		SenderAccount: senderAccount,
		RegularKey:    regularKey,
		Instructions:  infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction() for SetRegularKey: %w", err)
	}
	logger.Debug("response SetRegularKey",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput SetRegularKeyTxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(SetRegularKeyTxJSON): %w", err)
	}

	// Convert to DTO type
	return &dtoRipple.SetRegularKeyTxInput{
		TransactionType:    txInput.TransactionType,
		Account:            txInput.Account,
		RegularKey:         txInput.RegularKey,
		Fee:                txInput.Fee,
		Flags:              txInput.Flags,
		LastLedgerSequence: txInput.LastLedgerSequence,
		Sequence:           txInput.Sequence,
		SigningPubKey:      txInput.SigningPubKey,
		TxnSignature:       txInput.TxnSignature,
		Hash:               txInput.Hash,
	}, unquotedJSON, nil
}

// AccountSetTxInput is the transaction input for AccountSet
type AccountSetTxInput struct {
	TransactionType    string `json:"TransactionType"`
	Account            string `json:"Account"`
	SetFlag            uint32 `json:"SetFlag,omitempty"`
	ClearFlag          uint32 `json:"ClearFlag,omitempty"`
	Fee                string `json:"Fee"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence"`
	Sequence           uint64 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey,omitempty"`
	TxnSignature       string `json:"TxnSignature,omitempty"`
	Hash               string `json:"hash,omitempty"`
}

// AccountSet flag constants
const (
	// AsfDisableMaster disables the master key from signing transactions
	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/accountset#accountset-flags
	AsfDisableMaster uint32 = 4
)

// SignerListEntry represents a single signer in the SignerList
type SignerListEntry struct {
	SignerEntry struct {
		Account      string `json:"Account"`
		SignerWeight uint32 `json:"SignerWeight"`
	} `json:"SignerEntry"`
}

// SignerListSetTxInput is the transaction input for SignerListSet
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/signerlistset
type SignerListSetTxInput struct {
	TransactionType    string            `json:"TransactionType"`
	Account            string            `json:"Account"`
	SignerQuorum       uint32            `json:"SignerQuorum"`
	SignerEntries      []SignerListEntry `json:"SignerEntries,omitempty"`
	Fee                string            `json:"Fee"`
	Flags              uint64            `json:"Flags"`
	LastLedgerSequence uint64            `json:"LastLedgerSequence"`
	Sequence           uint64            `json:"Sequence"`
	SigningPubKey      string            `json:"SigningPubKey,omitempty"`
	TxnSignature       string            `json:"TxnSignature,omitempty"`
	Hash               string            `json:"hash,omitempty"`
}

// PrepareAccountSetTransaction prepares an AccountSet transaction
// - setFlag: flag to set (e.g., AsfDisableMaster = 4)
// - clearFlag: flag to clear
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/accountset
func (r *Ripple) PrepareAccountSetTransaction(
	ctx context.Context, senderAccount string, setFlag, clearFlag uint32, instructions *dtoRipple.Instructions,
) (*dtoRipple.AccountSetTxInput, string, error) {
	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	req := protogen.RequestPrepareTransaction_builder{
		TxType:        protogen.EnumTransactionType_TX_ACCOUNT_SET,
		SenderAccount: senderAccount,
		SetFlag:       setFlag,
		ClearFlag:     clearFlag,
		Instructions:  infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction() for AccountSet: %w", err)
	}
	logger.Debug("response AccountSet",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput AccountSetTxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(AccountSetTxJSON): %w", err)
	}

	// Convert to DTO type
	return &dtoRipple.AccountSetTxInput{
		TransactionType:    txInput.TransactionType,
		Account:            txInput.Account,
		SetFlag:            txInput.SetFlag,
		ClearFlag:          txInput.ClearFlag,
		Fee:                txInput.Fee,
		Flags:              txInput.Flags,
		LastLedgerSequence: txInput.LastLedgerSequence,
		Sequence:           txInput.Sequence,
		SigningPubKey:      txInput.SigningPubKey,
		TxnSignature:       txInput.TxnSignature,
		Hash:               txInput.Hash,
	}, unquotedJSON, nil
}

// PrepareSignerListSetTransaction prepares a SignerListSet transaction
// - signerQuorum: minimum total weight of signatures required (0 to remove signer list)
// - signerEntries: list of signers with their weights (must be empty if signerQuorum is 0)
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/signerlistset
func (r *Ripple) PrepareSignerListSetTransaction(
	ctx context.Context,
	senderAccount string,
	signerQuorum uint32,
	signerEntries []apixrp.SignerEntryInput,
	instructions *dtoRipple.Instructions,
) (*dtoRipple.SignerListSetTxInput, string, error) {
	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	// Convert SignerEntryInput to protogen.SignerEntry
	protoSignerEntries := make([]*protogen.SignerEntry, len(signerEntries))
	for i, entry := range signerEntries {
		protoSignerEntries[i] = protogen.SignerEntry_builder{
			Account: entry.Account,
			Weight:  entry.Weight,
		}.Build()
	}

	req := protogen.RequestPrepareTransaction_builder{
		TxType:        protogen.EnumTransactionType_TX_SINGER_LIST_SET,
		SenderAccount: senderAccount,
		SignerQuorum:  signerQuorum,
		SignerEntries: protoSignerEntries,
		Instructions:  infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction() for SignerListSet: %w", err)
	}
	logger.Debug("response SignerListSet",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput SignerListSetTxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(SignerListSetTxJSON): %w", err)
	}

	// Convert infrastructure type to DTO
	return ToDTOSignerListSetTxInput(&txInput), unquotedJSON, nil
}

// SignSignerListSetTransaction signs a SignerListSet transaction
func (r *Ripple) SignSignerListSetTransaction(
	ctx context.Context, txInput *SignerListSetTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "SignerListSet")
}

// IssuedCurrencyAmount represents a token amount with currency and issuer
// Used for TrustSet and other token-related transactions
type IssuedCurrencyAmount struct {
	Currency string `json:"currency"`
	Issuer   string `json:"issuer"`
	Value    string `json:"value"`
}

// TrustSetTxInput is the transaction input for TrustSet
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/trustset
type TrustSetTxInput struct {
	TransactionType    string               `json:"TransactionType"`
	Account            string               `json:"Account"`
	LimitAmount        IssuedCurrencyAmount `json:"LimitAmount"`
	QualityIn          uint32               `json:"QualityIn,omitempty"`
	QualityOut         uint32               `json:"QualityOut,omitempty"`
	Fee                string               `json:"Fee"`
	Flags              uint64               `json:"Flags"`
	LastLedgerSequence uint64               `json:"LastLedgerSequence"`
	Sequence           uint64               `json:"Sequence"`
	SigningPubKey      string               `json:"SigningPubKey,omitempty"`
	TxnSignature       string               `json:"TxnSignature,omitempty"`
	Hash               string               `json:"hash,omitempty"`
}

// PrepareTrustSetTransaction prepares a TrustSet transaction
// - limitAmount: the currency, issuer, and trust line limit (required)
// - qualityIn: value incoming balances at this ratio per 1,000,000,000 (optional, 0 to omit)
// - qualityOut: value outgoing balances at this ratio per 1,000,000,000 (optional, 0 to omit)
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/trustset
func (r *Ripple) PrepareTrustSetTransaction(
	ctx context.Context,
	senderAccount string,
	limitAmount *dtoRipple.IssuedCurrencyAmount,
	qualityIn, qualityOut uint32,
	instructions *dtoRipple.Instructions,
) (*dtoRipple.TrustSetTxInput, string, error) {
	// Validate required parameter
	if limitAmount == nil {
		return nil, "", errors.New("limitAmount is required for TrustSet transaction")
	}

	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	// Convert IssuedCurrencyAmount DTO to protogen
	protoLimitAmount := protogen.IssuedCurrencyAmount_builder{
		Currency: limitAmount.Currency,
		Issuer:   limitAmount.Issuer,
		Value:    limitAmount.Value,
	}.Build()

	req := protogen.RequestPrepareTransaction_builder{
		TxType:        protogen.EnumTransactionType_TX_TRUST_SET,
		SenderAccount: senderAccount,
		LimitAmount:   protoLimitAmount,
		QualityIn:     qualityIn,
		QualityOut:    qualityOut,
		Instructions:  infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction() for TrustSet: %w", err)
	}
	logger.Debug("response TrustSet",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput TrustSetTxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(TrustSetTxJSON): %w", err)
	}

	// Convert infrastructure type to DTO
	return ToDTOTrustSetTxInput(&txInput), unquotedJSON, nil
}

// SignTrustSetTransaction signs a TrustSet transaction
func (r *Ripple) SignTrustSetTransaction(
	ctx context.Context, txInput *TrustSetTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "TrustSet")
}

// ============================================================================
// Escrow transactions
// Reference: https://xrpl.org/docs/concepts/payment-types/escrow
// ============================================================================

// EscrowCreateTxInput is the transaction input for EscrowCreate
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/escrowcreate
type EscrowCreateTxInput struct {
	TransactionType    string `json:"TransactionType"`
	Account            string `json:"Account"`
	Amount             string `json:"Amount"`
	Destination        string `json:"Destination"`
	CancelAfter        uint32 `json:"CancelAfter,omitempty"`
	FinishAfter        uint32 `json:"FinishAfter,omitempty"`
	Condition          string `json:"Condition,omitempty"`
	DestinationTag     uint32 `json:"DestinationTag,omitempty"`
	Fee                string `json:"Fee"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence"`
	Sequence           uint64 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey,omitempty"`
	TxnSignature       string `json:"TxnSignature,omitempty"`
	Hash               string `json:"hash,omitempty"`
}

// EscrowFinishTxInput is the transaction input for EscrowFinish
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/escrowfinish
type EscrowFinishTxInput struct {
	TransactionType    string `json:"TransactionType"`
	Account            string `json:"Account"`
	Owner              string `json:"Owner"`
	OfferSequence      uint32 `json:"OfferSequence"`
	Condition          string `json:"Condition,omitempty"`
	Fulfillment        string `json:"Fulfillment,omitempty"`
	Fee                string `json:"Fee"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence"`
	Sequence           uint64 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey,omitempty"`
	TxnSignature       string `json:"TxnSignature,omitempty"`
	Hash               string `json:"hash,omitempty"`
}

// EscrowCancelTxInput is the transaction input for EscrowCancel
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/escrowcancel
type EscrowCancelTxInput struct {
	TransactionType    string `json:"TransactionType"`
	Account            string `json:"Account"`
	Owner              string `json:"Owner"`
	OfferSequence      uint32 `json:"OfferSequence"`
	Fee                string `json:"Fee"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence"`
	Sequence           uint64 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey,omitempty"`
	TxnSignature       string `json:"TxnSignature,omitempty"`
	Hash               string `json:"hash,omitempty"`
}

// PrepareEscrowCreateTransaction prepares an EscrowCreate transaction
// - destinationAccount: the address to receive escrowed funds
// - amount: XRP amount to escrow
// - cancelAfter: time (seconds since Ripple Epoch) when escrow expires (optional, 0 to omit)
// - finishAfter: time (seconds since Ripple Epoch) when escrow can be finished (optional, 0 to omit)
// - condition: PREIMAGE-SHA-256 crypto-condition (optional, empty to omit)
// - destinationTag: arbitrary tag for destination (optional, 0 to omit)
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/escrowcreate
func (r *Ripple) PrepareEscrowCreateTransaction(
	ctx context.Context,
	senderAccount, destinationAccount string,
	amount float64,
	cancelAfter, finishAfter uint32,
	condition string,
	destinationTag uint32,
	instructions *dtoRipple.Instructions,
) (*dtoRipple.EscrowCreateTxInput, string, error) {
	// Validate: at least one of cancelAfter, finishAfter, or condition must be set
	if cancelAfter == 0 && finishAfter == 0 && condition == "" {
		return nil, "", errors.New(
			"at least one of cancelAfter, finishAfter, or condition must be set for EscrowCreate")
	}

	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	req := protogen.RequestPrepareTransaction_builder{
		TxType:          protogen.EnumTransactionType_TX_ESCROW_CREATE,
		SenderAccount:   senderAccount,
		ReceiverAccount: destinationAccount,
		Amount:          amount,
		CancelAfter:     cancelAfter,
		FinishAfter:     finishAfter,
		Condition:       condition,
		DestinationTag:  destinationTag,
		Instructions:    infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction() for EscrowCreate: %w", err)
	}
	logger.Debug("response EscrowCreate",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput EscrowCreateTxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(EscrowCreateTxJSON): %w", err)
	}

	// Convert infrastructure type to DTO
	return ToDTOEscrowCreateTxInput(&txInput), unquotedJSON, nil
}

// PrepareEscrowFinishTransaction prepares an EscrowFinish transaction
// - owner: the account that funded the escrow
// - offerSequence: sequence number of the EscrowCreate transaction
// - condition: crypto-condition (required if escrow has one)
// - fulfillment: fulfillment matching the condition (required if escrow has condition)
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/escrowfinish
func (r *Ripple) PrepareEscrowFinishTransaction(
	ctx context.Context,
	senderAccount, owner string,
	offerSequence uint32,
	condition, fulfillment string,
	instructions *dtoRipple.Instructions,
) (*dtoRipple.EscrowFinishTxInput, string, error) {
	// Validate required parameters
	if owner == "" {
		return nil, "", errors.New("owner is required for EscrowFinish transaction")
	}
	if offerSequence == 0 {
		return nil, "", errors.New("offerSequence is required for EscrowFinish transaction")
	}

	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	req := protogen.RequestPrepareTransaction_builder{
		TxType:        protogen.EnumTransactionType_TX_ESCROW_FINISH,
		SenderAccount: senderAccount,
		Owner:         owner,
		OfferSequence: offerSequence,
		Condition:     condition,
		Fulfillment:   fulfillment,
		Instructions:  infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction() for EscrowFinish: %w", err)
	}
	logger.Debug("response EscrowFinish",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput EscrowFinishTxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(EscrowFinishTxJSON): %w", err)
	}

	// Convert infrastructure type to DTO
	return ToDTOEscrowFinishTxInput(&txInput), unquotedJSON, nil
}

// PrepareEscrowCancelTransaction prepares an EscrowCancel transaction
// - owner: the account that funded the escrow
// - offerSequence: sequence number of the EscrowCreate transaction
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/escrowcancel
func (r *Ripple) PrepareEscrowCancelTransaction(
	ctx context.Context,
	senderAccount, owner string,
	offerSequence uint32,
	instructions *dtoRipple.Instructions,
) (*dtoRipple.EscrowCancelTxInput, string, error) {
	// Validate required parameters
	if owner == "" {
		return nil, "", errors.New("owner is required for EscrowCancel transaction")
	}
	if offerSequence == 0 {
		return nil, "", errors.New("offerSequence is required for EscrowCancel transaction")
	}

	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	req := protogen.RequestPrepareTransaction_builder{
		TxType:        protogen.EnumTransactionType_TX_ESCROW_CANCEL,
		SenderAccount: senderAccount,
		Owner:         owner,
		OfferSequence: offerSequence,
		Instructions:  infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction() for EscrowCancel: %w", err)
	}
	logger.Debug("response EscrowCancel",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput EscrowCancelTxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(EscrowCancelTxJSON): %w", err)
	}

	// Convert infrastructure type to DTO
	return ToDTOEscrowCancelTxInput(&txInput), unquotedJSON, nil
}

// SignEscrowCreateTransaction signs an EscrowCreate transaction
func (r *Ripple) SignEscrowCreateTransaction(
	ctx context.Context, txInput *EscrowCreateTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "EscrowCreate")
}

// SignEscrowFinishTransaction signs an EscrowFinish transaction
func (r *Ripple) SignEscrowFinishTransaction(
	ctx context.Context, txInput *EscrowFinishTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "EscrowFinish")
}

// SignEscrowCancelTransaction signs an EscrowCancel transaction
func (r *Ripple) SignEscrowCancelTransaction(
	ctx context.Context, txInput *EscrowCancelTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "EscrowCancel")
}

// ============================================================================
// PaymentChannel transactions
// Reference: https://xrpl.org/docs/concepts/payment-types/payment-channels
// ============================================================================

// PaymentChannelCreateTxInput is the transaction input for PaymentChannelCreate
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/paymentchannelcreate
type PaymentChannelCreateTxInput struct {
	TransactionType    string `json:"TransactionType"`
	Account            string `json:"Account"`
	Amount             string `json:"Amount"`
	Destination        string `json:"Destination"`
	SettleDelay        uint32 `json:"SettleDelay"`
	PublicKey          string `json:"PublicKey"`
	CancelAfter        uint32 `json:"CancelAfter,omitempty"`
	DestinationTag     uint32 `json:"DestinationTag,omitempty"`
	SourceTag          uint32 `json:"SourceTag,omitempty"`
	Fee                string `json:"Fee"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence"`
	Sequence           uint64 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey,omitempty"`
	TxnSignature       string `json:"TxnSignature,omitempty"`
	Hash               string `json:"hash,omitempty"`
}

// PaymentChannelFundTxInput is the transaction input for PaymentChannelFund
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/paymentchannelfund
type PaymentChannelFundTxInput struct {
	TransactionType    string `json:"TransactionType"`
	Account            string `json:"Account"`
	Channel            string `json:"Channel"`
	Amount             string `json:"Amount"`
	Expiration         uint32 `json:"Expiration,omitempty"`
	Fee                string `json:"Fee"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence"`
	Sequence           uint64 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey,omitempty"`
	TxnSignature       string `json:"TxnSignature,omitempty"`
	Hash               string `json:"hash,omitempty"`
}

// PaymentChannelClaimTxInput is the transaction input for PaymentChannelClaim
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/paymentchannelclaim
type PaymentChannelClaimTxInput struct {
	TransactionType    string `json:"TransactionType"`
	Account            string `json:"Account"`
	Channel            string `json:"Channel"`
	Balance            string `json:"Balance,omitempty"`
	Amount             string `json:"Amount,omitempty"`
	Signature          string `json:"Signature,omitempty"`
	PublicKey          string `json:"PublicKey,omitempty"`
	Fee                string `json:"Fee"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence"`
	Sequence           uint64 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey,omitempty"`
	TxnSignature       string `json:"TxnSignature,omitempty"`
	Hash               string `json:"hash,omitempty"`
}

// PreparePaymentChannelCreateTransaction prepares a PaymentChannelCreate transaction
// - destinationAccount: the address to receive payments from the channel
// - amount: XRP amount to fund the channel
// - settleDelay: seconds source must wait to close channel if unclaimed funds remain
// - publicKey: hex-encoded public key for verifying claim signatures
// - cancelAfter: immutable expiry time (optional, 0 to omit)
// - destinationTag: tag for destination (optional, 0 to omit)
// - sourceTag: tag for source (optional, 0 to omit)
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/paymentchannelcreate
func (r *Ripple) PreparePaymentChannelCreateTransaction(
	ctx context.Context,
	senderAccount, destinationAccount string,
	amount float64,
	settleDelay uint32,
	publicKey string,
	cancelAfter, destinationTag, sourceTag uint32,
	instructions *dtoRipple.Instructions,
) (*dtoRipple.PaymentChannelCreateTxInput, string, error) {
	// Validate required parameters
	if settleDelay == 0 {
		return nil, "", errors.New("settleDelay is required for PaymentChannelCreate transaction")
	}
	if publicKey == "" {
		return nil, "", errors.New("publicKey is required for PaymentChannelCreate transaction")
	}

	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	req := protogen.RequestPrepareTransaction_builder{
		TxType:          protogen.EnumTransactionType_TX_PAYMENT_CHANNEL_CREATE,
		SenderAccount:   senderAccount,
		ReceiverAccount: destinationAccount,
		Amount:          amount,
		SettleDelay:     settleDelay,
		PublicKey:       publicKey,
		CancelAfter:     cancelAfter,
		DestinationTag:  destinationTag,
		SourceTag:       sourceTag,
		Instructions:    infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction() for PaymentChannelCreate: %w", err)
	}
	logger.Debug("response PaymentChannelCreate",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput PaymentChannelCreateTxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(PaymentChannelCreateTxJSON): %w", err)
	}

	// Convert infrastructure type to DTO
	return ToDTOPaymentChannelCreateTxInput(&txInput), unquotedJSON, nil
}

// PreparePaymentChannelFundTransaction prepares a PaymentChannelFund transaction
// - channel: the unique ID (Hash256) of the payment channel
// - amount: XRP to add to the channel
// - expiration: new expiration time (optional, 0 to omit)
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/paymentchannelfund
func (r *Ripple) PreparePaymentChannelFundTransaction(
	ctx context.Context,
	senderAccount, channel string,
	amount float64,
	expiration uint32,
	instructions *dtoRipple.Instructions,
) (*dtoRipple.PaymentChannelFundTxInput, string, error) {
	// Validate required parameters
	if channel == "" {
		return nil, "", errors.New("channel is required for PaymentChannelFund transaction")
	}

	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	req := protogen.RequestPrepareTransaction_builder{
		TxType:        protogen.EnumTransactionType_TX_PAYMENT_CHANNEL_FUND,
		SenderAccount: senderAccount,
		Channel:       channel,
		Amount:        amount,
		Expiration:    expiration,
		Instructions:  infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction() for PaymentChannelFund: %w", err)
	}
	logger.Debug("response PaymentChannelFund",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput PaymentChannelFundTxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(PaymentChannelFundTxJSON): %w", err)
	}

	// Convert infrastructure type to DTO
	return ToDTOPaymentChannelFundTxInput(&txInput), unquotedJSON, nil
}

// PreparePaymentChannelClaimTransaction prepares a PaymentChannelClaim transaction
// - channel: the unique ID (Hash256) of the payment channel
// - balance: total XRP in drops delivered after this claim (optional, empty to omit)
// - amount: authorized cumulative XRP amount via signature (optional, 0 to omit)
// - signature: hex-encoded signature for claim authorization (optional)
// - publicKey: hex-encoded public key for signature verification (optional)
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/paymentchannelclaim
func (r *Ripple) PreparePaymentChannelClaimTransaction(
	ctx context.Context,
	senderAccount, channel string,
	balance string,
	amount float64,
	signature, publicKey string,
	instructions *dtoRipple.Instructions,
) (*dtoRipple.PaymentChannelClaimTxInput, string, error) {
	// Validate required parameters
	if channel == "" {
		return nil, "", errors.New("channel is required for PaymentChannelClaim transaction")
	}

	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	req := protogen.RequestPrepareTransaction_builder{
		TxType:        protogen.EnumTransactionType_TX_PAYMENT_CHANNEL_CLAIM,
		SenderAccount: senderAccount,
		Channel:       channel,
		Balance:       balance,
		Amount:        amount,
		Signature:     signature,
		PublicKey:     publicKey,
		Instructions:  infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction() for PaymentChannelClaim: %w", err)
	}
	logger.Debug("response PaymentChannelClaim",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput PaymentChannelClaimTxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(PaymentChannelClaimTxJSON): %w", err)
	}

	// Convert infrastructure type to DTO
	return ToDTOPaymentChannelClaimTxInput(&txInput), unquotedJSON, nil
}

// SignPaymentChannelCreateTransaction signs a PaymentChannelCreate transaction
func (r *Ripple) SignPaymentChannelCreateTransaction(
	ctx context.Context, txInput *PaymentChannelCreateTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "PaymentChannelCreate")
}

// SignPaymentChannelFundTransaction signs a PaymentChannelFund transaction
func (r *Ripple) SignPaymentChannelFundTransaction(
	ctx context.Context, txInput *PaymentChannelFundTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "PaymentChannelFund")
}

// SignPaymentChannelClaimTransaction signs a PaymentChannelClaim transaction
func (r *Ripple) SignPaymentChannelClaimTransaction(
	ctx context.Context, txInput *PaymentChannelClaimTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "PaymentChannelClaim")
}

// ============================================================================
// NFToken transactions
// Reference: https://xrpl.org/docs/concepts/tokens/nfts
// ============================================================================

// NFTokenMintTxInput is the transaction input for NFTokenMint
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokenmint
type NFTokenMintTxInput struct {
	TransactionType    string `json:"TransactionType"`
	Account            string `json:"Account"`
	NFTokenTaxon       uint32 `json:"NFTokenTaxon"`
	Issuer             string `json:"Issuer,omitempty"`
	TransferFee        uint32 `json:"TransferFee,omitempty"`
	URI                string `json:"URI,omitempty"`
	Fee                string `json:"Fee"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence"`
	Sequence           uint64 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey,omitempty"`
	TxnSignature       string `json:"TxnSignature,omitempty"`
	Hash               string `json:"hash,omitempty"`
}

// NFTokenBurnTxInput is the transaction input for NFTokenBurn
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokenburn
type NFTokenBurnTxInput struct {
	TransactionType    string `json:"TransactionType"`
	Account            string `json:"Account"`
	NFTokenID          string `json:"NFTokenID"`
	Owner              string `json:"Owner,omitempty"`
	Fee                string `json:"Fee"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence"`
	Sequence           uint64 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey,omitempty"`
	TxnSignature       string `json:"TxnSignature,omitempty"`
	Hash               string `json:"hash,omitempty"`
}

// NFTokenCreateOfferTxInput is the transaction input for NFTokenCreateOffer
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokencreateoffer
type NFTokenCreateOfferTxInput struct {
	TransactionType    string `json:"TransactionType"`
	Account            string `json:"Account"`
	NFTokenID          string `json:"NFTokenID"`
	Amount             string `json:"Amount"`
	Owner              string `json:"Owner,omitempty"`
	Expiration         uint32 `json:"Expiration,omitempty"`
	Destination        string `json:"Destination,omitempty"`
	Fee                string `json:"Fee"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence"`
	Sequence           uint64 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey,omitempty"`
	TxnSignature       string `json:"TxnSignature,omitempty"`
	Hash               string `json:"hash,omitempty"`
}

// NFTokenAcceptOfferTxInput is the transaction input for NFTokenAcceptOffer
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokenacceptoffer
type NFTokenAcceptOfferTxInput struct {
	TransactionType    string `json:"TransactionType"`
	Account            string `json:"Account"`
	NFTokenSellOffer   string `json:"NFTokenSellOffer,omitempty"`
	NFTokenBuyOffer    string `json:"NFTokenBuyOffer,omitempty"`
	NFTokenBrokerFee   string `json:"NFTokenBrokerFee,omitempty"`
	Fee                string `json:"Fee"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence"`
	Sequence           uint64 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey,omitempty"`
	TxnSignature       string `json:"TxnSignature,omitempty"`
	Hash               string `json:"hash,omitempty"`
}

// NFTokenCancelOfferTxInput is the transaction input for NFTokenCancelOffer
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokencanceloffer
type NFTokenCancelOfferTxInput struct {
	TransactionType    string   `json:"TransactionType"`
	Account            string   `json:"Account"`
	NFTokenOffers      []string `json:"NFTokenOffers"`
	Fee                string   `json:"Fee"`
	Flags              uint64   `json:"Flags"`
	LastLedgerSequence uint64   `json:"LastLedgerSequence"`
	Sequence           uint64   `json:"Sequence"`
	SigningPubKey      string   `json:"SigningPubKey,omitempty"`
	TxnSignature       string   `json:"TxnSignature,omitempty"`
	Hash               string   `json:"hash,omitempty"`
}

// PrepareNFTokenMintTransaction prepares an NFTokenMint transaction
// - nfTokenTaxon: taxon identifying a collection or category (required, can be 0)
// - issuer: issuer of the NFT if different from Account (optional)
// - uri: hex-encoded URI pointing to NFT metadata (optional)
// - transferFee: fee (0-50000) charged on secondary sales (optional)
// Note: Flags are set via xrpl.js autofill based on transaction requirements
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokenmint
func (r *Ripple) PrepareNFTokenMintTransaction(
	ctx context.Context,
	senderAccount string,
	nfTokenTaxon uint32,
	issuer, uri string,
	transferFee uint32,
	instructions *dtoRipple.Instructions,
) (*dtoRipple.NFTokenMintTxInput, string, error) {
	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	req := protogen.RequestPrepareTransaction_builder{
		TxType:        protogen.EnumTransactionType_TX_NFTOKEN_MINT,
		SenderAccount: senderAccount,
		NfTokenTaxon:  nfTokenTaxon,
		Issuer:        issuer,
		Uri:           uri,
		TransferFee:   transferFee,
		Instructions:  infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction() for NFTokenMint: %w", err)
	}
	logger.Debug("response NFTokenMint",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput NFTokenMintTxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(NFTokenMintTxJSON): %w", err)
	}

	// Convert infrastructure type to DTO
	return ToDTONFTokenMintTxInput(&txInput), unquotedJSON, nil
}

// PrepareNFTokenBurnTransaction prepares an NFTokenBurn transaction
// - nfTokenID: the unique ID of the NFToken to burn (required)
// - owner: the owner of the NFT if issuer is burning (optional)
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokenburn
func (r *Ripple) PrepareNFTokenBurnTransaction(
	ctx context.Context,
	senderAccount, nfTokenID, owner string,
	instructions *dtoRipple.Instructions,
) (*dtoRipple.NFTokenBurnTxInput, string, error) {
	// Validate required parameters
	if nfTokenID == "" {
		return nil, "", errors.New("nfTokenID is required for NFTokenBurn transaction")
	}

	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	req := protogen.RequestPrepareTransaction_builder{
		TxType:        protogen.EnumTransactionType_TX_NFTOKEN_BURN,
		SenderAccount: senderAccount,
		NfTokenID:     nfTokenID,
		Owner:         owner,
		Instructions:  infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction() for NFTokenBurn: %w", err)
	}
	logger.Debug("response NFTokenBurn",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput NFTokenBurnTxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(NFTokenBurnTxJSON): %w", err)
	}

	// Convert infrastructure type to DTO
	return ToDTONFTokenBurnTxInput(&txInput), unquotedJSON, nil
}

// PrepareNFTokenCreateOfferTransaction prepares an NFTokenCreateOffer transaction
// - nfTokenID: the unique ID of the NFToken (required)
// - amount: the amount in XRP for the offer (can be 0 for free sell offers)
// - owner: owner of the NFT for buy offers (optional)
// - destination: only this account can accept the offer (optional)
// - expiration: when the offer expires (optional)
// Note: tfSellNFToken flag is determined by whether owner is set (buy offer) or not (sell offer)
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokencreateoffer
func (r *Ripple) PrepareNFTokenCreateOfferTransaction(
	ctx context.Context,
	senderAccount, nfTokenID string,
	amount float64,
	owner, destination string,
	expiration uint32,
	instructions *dtoRipple.Instructions,
) (*dtoRipple.NFTokenCreateOfferTxInput, string, error) {
	// Validate required parameters
	if nfTokenID == "" {
		return nil, "", errors.New("nfTokenID is required for NFTokenCreateOffer transaction")
	}

	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	req := protogen.RequestPrepareTransaction_builder{
		TxType:          protogen.EnumTransactionType_TX_NFTOKEN_CREATE_OFFER,
		SenderAccount:   senderAccount,
		NfTokenID:       nfTokenID,
		Amount:          amount,
		Owner:           owner,
		ReceiverAccount: destination,
		Expiration:      expiration,
		Instructions:    infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction() for NFTokenCreateOffer: %w", err)
	}
	logger.Debug("response NFTokenCreateOffer",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput NFTokenCreateOfferTxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(NFTokenCreateOfferTxJSON): %w", err)
	}

	// Convert infrastructure type to DTO
	return ToDTONFTokenCreateOfferTxInput(&txInput), unquotedJSON, nil
}

// PrepareNFTokenAcceptOfferTransaction prepares an NFTokenAcceptOffer transaction
// - nfTokenSellOffer: ID of the sell offer to accept (optional)
// - nfTokenBuyOffer: ID of the buy offer to accept (optional)
// - nfTokenBrokerFee: broker fee in XRP for brokered mode (optional)
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokenacceptoffer
func (r *Ripple) PrepareNFTokenAcceptOfferTransaction(
	ctx context.Context,
	senderAccount string,
	nfTokenSellOffer, nfTokenBuyOffer string,
	nfTokenBrokerFee float64,
	instructions *dtoRipple.Instructions,
) (*dtoRipple.NFTokenAcceptOfferTxInput, string, error) {
	// Validate: at least one offer must be provided
	if nfTokenSellOffer == "" && nfTokenBuyOffer == "" {
		return nil, "", errors.New(
			"at least one of nfTokenSellOffer or nfTokenBuyOffer is required for NFTokenAcceptOffer")
	}

	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	req := protogen.RequestPrepareTransaction_builder{
		TxType:           protogen.EnumTransactionType_TX_NFTOKEN_ACCEPT_OFFER,
		SenderAccount:    senderAccount,
		NfTokenSellOffer: nfTokenSellOffer,
		NfTokenBuyOffer:  nfTokenBuyOffer,
		NfTokenBrokerFee: nfTokenBrokerFee,
		Instructions:     infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction() for NFTokenAcceptOffer: %w", err)
	}
	logger.Debug("response NFTokenAcceptOffer",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput NFTokenAcceptOfferTxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(NFTokenAcceptOfferTxJSON): %w", err)
	}

	// Convert infrastructure type to DTO
	return ToDTONFTokenAcceptOfferTxInput(&txInput), unquotedJSON, nil
}

// PrepareNFTokenCancelOfferTransaction prepares an NFTokenCancelOffer transaction
// - nfTokenOffers: array of NFTokenOffer IDs to cancel (required)
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokencanceloffer
func (r *Ripple) PrepareNFTokenCancelOfferTransaction(
	ctx context.Context,
	senderAccount string,
	nfTokenOffers []string,
	instructions *dtoRipple.Instructions,
) (*dtoRipple.NFTokenCancelOfferTxInput, string, error) {
	// Validate required parameters
	if len(nfTokenOffers) == 0 {
		return nil, "", errors.New("nfTokenOffers is required for NFTokenCancelOffer transaction")
	}

	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	req := protogen.RequestPrepareTransaction_builder{
		TxType:        protogen.EnumTransactionType_TX_NFTOKEN_CANCEL_OFFER,
		SenderAccount: senderAccount,
		NfTokenOffers: nfTokenOffers,
		Instructions:  infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction() for NFTokenCancelOffer: %w", err)
	}
	logger.Debug("response NFTokenCancelOffer",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput NFTokenCancelOfferTxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(NFTokenCancelOfferTxJSON): %w", err)
	}

	// Convert infrastructure type to DTO
	return ToDTONFTokenCancelOfferTxInput(&txInput), unquotedJSON, nil
}

// SignNFTokenMintTransaction signs an NFTokenMint transaction
func (r *Ripple) SignNFTokenMintTransaction(
	ctx context.Context, txInput *NFTokenMintTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "NFTokenMint")
}

// SignNFTokenBurnTransaction signs an NFTokenBurn transaction
func (r *Ripple) SignNFTokenBurnTransaction(
	ctx context.Context, txInput *NFTokenBurnTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "NFTokenBurn")
}

// SignNFTokenCreateOfferTransaction signs an NFTokenCreateOffer transaction
func (r *Ripple) SignNFTokenCreateOfferTransaction(
	ctx context.Context, txInput *NFTokenCreateOfferTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "NFTokenCreateOffer")
}

// SignNFTokenAcceptOfferTransaction signs an NFTokenAcceptOffer transaction
func (r *Ripple) SignNFTokenAcceptOfferTransaction(
	ctx context.Context, txInput *NFTokenAcceptOfferTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "NFTokenAcceptOffer")
}

// SignNFTokenCancelOfferTransaction signs an NFTokenCancelOffer transaction
func (r *Ripple) SignNFTokenCancelOfferTransaction(
	ctx context.Context, txInput *NFTokenCancelOfferTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "NFTokenCancelOffer")
}

// signTransactionJSON is a generic helper that signs any transaction type.
// It marshals the input to JSON and calls the gRPC SignTransaction API.
func (r *Ripple) signTransactionJSON(
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

// SignSetRegularKeyTransaction signs a SetRegularKey transaction
func (r *Ripple) SignSetRegularKeyTransaction(
	ctx context.Context, txInput *SetRegularKeyTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "SetRegularKey")
}

// SignAccountSetTransaction signs an AccountSet transaction
func (r *Ripple) SignAccountSetTransaction(
	ctx context.Context, txInput *AccountSetTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "AccountSet")
}

// SignTransaction calls SignTransaction API
// Offline functionality
// - https://xrpl.org/rippleapi-reference.html#offline-functionality
func (r *Ripple) SignTransaction(
	ctx context.Context, txInput *dtoRipple.TxInput, secret string,
) (string, string, error) {
	// Convert DTO to infrastructure type
	infraTxInput := ToInfraTxInput(txInput)

	strJSON, err := json.Marshal(infraTxInput)
	if err != nil {
		return "", "", fmt.Errorf("fail to call json.Marshal(txJSON): %w", err)
	}
	req := protogen.RequestSignTransaction_builder{
		TxJSON: string(strJSON),
		Secret: secret,
	}.Build()

	res, err := r.API.txClient.SignTransaction(ctx, req)
	if err != nil {
		return "", "", fmt.Errorf("fail to call client.SignTransaction(): %w", err)
	}

	return res.GetTxID(), res.GetTxBlob(), nil
}

// CombineTransaction combines signed transactions from multiple accounts for a multisignature transaction.
// - The signed transaction must subsequently be submitted.
func (r *Ripple) CombineTransaction(ctx context.Context, signedTxs []string) (string, string, error) {
	req := protogen.RequestCombineTransaction_builder{
		SignedTransactions: signedTxs,
	}.Build()

	res, err := r.API.txClient.CombineTransaction(ctx, req)
	if err != nil {
		return "", "", fmt.Errorf("fail to call client.CombineTransaction(): %w", err)
	}

	return res.GetTxID(), res.GetSignedTransaction(), nil
}

// SubmitTransaction calls SubmitTransaction API
// - signedTx is returned TxBlob by SignTransaction()
func (r *Ripple) SubmitTransaction(ctx context.Context, signedTx string) (*dtoRipple.SentTx, uint64, error) {
	req := protogen.RequestSubmitTransaction_builder{
		TxBlob: signedTx,
	}.Build()
	res, err := r.API.txClient.SubmitTransaction(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to call client.SubmitTransaction(): %w", err)
	}

	var sentTxJSON SentTx
	if err = json.Unmarshal([]byte(res.GetResultJSONString()), &sentTxJSON); err != nil {
		return nil, 0, fmt.Errorf("fail to call json.Unmarshal(sentTxJSON): %w", err)
	}

	// FIXME:
	// res.EarliestLedgerVersion may be useless because SentTxJSON includes `LastLedgerSequence` and it would be useful
	logger.Debug("response of submitTransaction",
		"res.ResultJSONString", res.GetResultJSONString(),
		"res.EarliestLedgerVersion", res.GetEarliestLedgerVersion(),
		"sentTxJSON.TxJSON.LastLedgerSequence", sentTxJSON.TxJSON.LastLedgerSequence,
	)
	// res.EarliestLedgerVersion => for when calling GetTransaction()
	// sentTxJSON.TxJSON.LastLedgerSequence => for when calling WaitValidation()

	// Convert infrastructure type to DTO
	return ToDTOSentTx(&sentTxJSON), res.GetEarliestLedgerVersion(), nil
	// return ToDTOSentTx(&sentTxJSON), sentTxJSON.TxJSON.LastLedgerSequence, nil
}

// WaitValidation calls WaitValidation API
// - handling server streaming
func (r *Ripple) WaitValidation(ctx context.Context, targetledgerVarsion uint64) (uint64, error) {
	req := &emptypb.Empty{}
	resStream, err := r.API.txClient.WaitValidation(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("fail to call client.WaitValidation(): %w", err)
	}

	defer func() {
		logger.Debug("running in defer func()")
		if closeErr := resStream.CloseSend(); closeErr != nil {
			logger.Warn("fail to call resStream.CloseSend()")
		}
	}()

	for {
		var res *protogen.ResponseWaitValidation
		res, err = resStream.Recv()
		if err == io.EOF {
			logger.Warn("server is closed in WaitValidation()")
			return 0, errors.New("server is closed")
		} else if err != nil {
			if respErr, ok := status.FromError(err); ok {
				switch respErr.Code() {
				case codes.InvalidArgument:
					logger.Warn("parameter is invalid in WaitValidation()")
				case codes.DeadlineExceeded:
					logger.Warn("timeout in WaitValidation()")
				case codes.OK, codes.Canceled, codes.Unknown, codes.NotFound, codes.AlreadyExists,
					codes.PermissionDenied, codes.ResourceExhausted, codes.FailedPrecondition,
					codes.Aborted, codes.OutOfRange, codes.Unimplemented, codes.Internal,
					codes.Unavailable, codes.DataLoss, codes.Unauthenticated:
					logger.Warn("gRPC error in WaitValidation()",
						"code", uint32(respErr.Code()),
						"message", respErr.Message(),
					)
				default:
					logger.Warn("gRPC error in WaitValidation()",
						"code", uint32(respErr.Code()),
						"message", respErr.Message(),
					)
				}
			} else {
				logger.Warn("fail to call resStream.Recv()", "error", err)
			}
			// break
			return 0, fmt.Errorf("fail to call resStream.Recv(): %w", err)
		}
		// success
		logger.Info("response in WaitValidation()", "LedgerVersion", res.GetLedgerVersion())
		if targetledgerVarsion <= res.GetLedgerVersion() {
			// done
			return res.GetLedgerVersion(), nil
		}
		// continue
	}
}

// GetTransaction calls GetTransaction API
func (r *Ripple) GetTransaction(
	ctx context.Context, txID string, targetLedgerVersion uint64,
) (*dtoRipple.TxInfo, error) {
	req := protogen.RequestGetTransaction_builder{
		TxID:             txID,
		MinLedgerVersion: targetLedgerVersion,
	}.Build()
	res, err := r.API.txClient.GetTransaction(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("fail to call client.GetTransaction(): %w", err)
	}

	if res.GetResultJSONString() == "" {
		return nil, fmt.Errorf("fail to get transaction info by %s", txID)
	}

	logger.Debug("response of getTransaction",
		"res.ResultJSONString", res.GetResultJSONString(),
	)

	var txInfo TxInfo
	if err = json.Unmarshal([]byte(res.GetResultJSONString()), &txInfo); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(txInfo): %w", err)
	}
	// TODO: check
	// txInfo.Outcome.Result : tesSUCCESS

	// Convert infrastructure type to DTO
	return ToDTOTxInfo(&txInfo), nil
}
