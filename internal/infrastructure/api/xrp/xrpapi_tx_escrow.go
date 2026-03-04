package xrp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/protogen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

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
func (r *XRP) PrepareEscrowCreateTransaction(
	ctx context.Context,
	senderAccount, destinationAccount string,
	amount float64,
	cancelAfter, finishAfter uint32,
	condition string,
	destinationTag uint32,
	instructions *dtoxrp.Instructions,
) (*dtoxrp.EscrowCreateTxInput, string, error) {
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
func (r *XRP) PrepareEscrowFinishTransaction(
	ctx context.Context,
	senderAccount, owner string,
	offerSequence uint32,
	condition, fulfillment string,
	instructions *dtoxrp.Instructions,
) (*dtoxrp.EscrowFinishTxInput, string, error) {
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
func (r *XRP) PrepareEscrowCancelTransaction(
	ctx context.Context,
	senderAccount, owner string,
	offerSequence uint32,
	instructions *dtoxrp.Instructions,
) (*dtoxrp.EscrowCancelTxInput, string, error) {
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
func (r *XRP) SignEscrowCreateTransaction(
	ctx context.Context, txInput *EscrowCreateTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "EscrowCreate")
}

// SignEscrowFinishTransaction signs an EscrowFinish transaction
func (r *XRP) SignEscrowFinishTransaction(
	ctx context.Context, txInput *EscrowFinishTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "EscrowFinish")
}

// SignEscrowCancelTransaction signs an EscrowCancel transaction
func (r *XRP) SignEscrowCancelTransaction(
	ctx context.Context, txInput *EscrowCancelTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "EscrowCancel")
}
