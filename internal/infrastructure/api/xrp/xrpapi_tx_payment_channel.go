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
func (r *XRP) PreparePaymentChannelCreateTransaction(
	ctx context.Context,
	senderAccount, destinationAccount string,
	amount float64,
	settleDelay uint32,
	publicKey string,
	cancelAfter, destinationTag, sourceTag uint32,
	instructions *dtoxrp.Instructions,
) (*dtoxrp.PaymentChannelCreateTxInput, string, error) {
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
func (r *XRP) PreparePaymentChannelFundTransaction(
	ctx context.Context,
	senderAccount, channel string,
	amount float64,
	expiration uint32,
	instructions *dtoxrp.Instructions,
) (*dtoxrp.PaymentChannelFundTxInput, string, error) {
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
func (r *XRP) PreparePaymentChannelClaimTransaction(
	ctx context.Context,
	senderAccount, channel string,
	balance string,
	amount float64,
	signature, publicKey string,
	instructions *dtoxrp.Instructions,
) (*dtoxrp.PaymentChannelClaimTxInput, string, error) {
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
func (r *XRP) SignPaymentChannelCreateTransaction(
	ctx context.Context, txInput *PaymentChannelCreateTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "PaymentChannelCreate")
}

// SignPaymentChannelFundTransaction signs a PaymentChannelFund transaction
func (r *XRP) SignPaymentChannelFundTransaction(
	ctx context.Context, txInput *PaymentChannelFundTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "PaymentChannelFund")
}

// SignPaymentChannelClaimTransaction signs a PaymentChannelClaim transaction
func (r *XRP) SignPaymentChannelClaimTransaction(
	ctx context.Context, txInput *PaymentChannelClaimTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "PaymentChannelClaim")
}
