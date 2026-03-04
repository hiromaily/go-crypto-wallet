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
func (r *XRP) PrepareNFTokenMintTransaction(
	ctx context.Context,
	senderAccount string,
	nfTokenTaxon uint32,
	issuer, uri string,
	transferFee uint32,
	instructions *dtoxrp.Instructions,
) (*dtoxrp.NFTokenMintTxInput, string, error) {
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
func (r *XRP) PrepareNFTokenBurnTransaction(
	ctx context.Context,
	senderAccount, nfTokenID, owner string,
	instructions *dtoxrp.Instructions,
) (*dtoxrp.NFTokenBurnTxInput, string, error) {
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
func (r *XRP) PrepareNFTokenCreateOfferTransaction(
	ctx context.Context,
	senderAccount, nfTokenID string,
	amount float64,
	owner, destination string,
	expiration uint32,
	instructions *dtoxrp.Instructions,
) (*dtoxrp.NFTokenCreateOfferTxInput, string, error) {
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
func (r *XRP) PrepareNFTokenAcceptOfferTransaction(
	ctx context.Context,
	senderAccount string,
	nfTokenSellOffer, nfTokenBuyOffer string,
	nfTokenBrokerFee float64,
	instructions *dtoxrp.Instructions,
) (*dtoxrp.NFTokenAcceptOfferTxInput, string, error) {
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
func (r *XRP) PrepareNFTokenCancelOfferTransaction(
	ctx context.Context,
	senderAccount string,
	nfTokenOffers []string,
	instructions *dtoxrp.Instructions,
) (*dtoxrp.NFTokenCancelOfferTxInput, string, error) {
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
func (r *XRP) SignNFTokenMintTransaction(
	ctx context.Context, txInput *NFTokenMintTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "NFTokenMint")
}

// SignNFTokenBurnTransaction signs an NFTokenBurn transaction
func (r *XRP) SignNFTokenBurnTransaction(
	ctx context.Context, txInput *NFTokenBurnTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "NFTokenBurn")
}

// SignNFTokenCreateOfferTransaction signs an NFTokenCreateOffer transaction
func (r *XRP) SignNFTokenCreateOfferTransaction(
	ctx context.Context, txInput *NFTokenCreateOfferTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "NFTokenCreateOffer")
}

// SignNFTokenAcceptOfferTransaction signs an NFTokenAcceptOffer transaction
func (r *XRP) SignNFTokenAcceptOfferTransaction(
	ctx context.Context, txInput *NFTokenAcceptOfferTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "NFTokenAcceptOffer")
}

// SignNFTokenCancelOfferTransaction signs an NFTokenCancelOffer transaction
func (r *XRP) SignNFTokenCancelOfferTransaction(
	ctx context.Context, txInput *NFTokenCancelOfferTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "NFTokenCancelOffer")
}
