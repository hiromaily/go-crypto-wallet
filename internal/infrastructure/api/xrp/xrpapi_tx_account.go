package xrp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/protogen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

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
func (r *XRP) PrepareSetRegularKeyTransaction(
	ctx context.Context, senderAccount, regularKey string, instructions *dtoxrp.Instructions,
) (*dtoxrp.SetRegularKeyTxInput, string, error) {
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
	return &dtoxrp.SetRegularKeyTxInput{
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
func (r *XRP) PrepareAccountSetTransaction(
	ctx context.Context, senderAccount string, setFlag, clearFlag uint32, instructions *dtoxrp.Instructions,
) (*dtoxrp.AccountSetTxInput, string, error) {
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
	return &dtoxrp.AccountSetTxInput{
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
func (r *XRP) PrepareSignerListSetTransaction(
	ctx context.Context,
	senderAccount string,
	signerQuorum uint32,
	signerEntries []apixrp.SignerEntryInput,
	instructions *dtoxrp.Instructions,
) (*dtoxrp.SignerListSetTxInput, string, error) {
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
func (r *XRP) SignSignerListSetTransaction(
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
func (r *XRP) PrepareTrustSetTransaction(
	ctx context.Context,
	senderAccount string,
	limitAmount *dtoxrp.IssuedCurrencyAmount,
	qualityIn, qualityOut uint32,
	instructions *dtoxrp.Instructions,
) (*dtoxrp.TrustSetTxInput, string, error) {
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
func (r *XRP) SignTrustSetTransaction(
	ctx context.Context, txInput *TrustSetTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "TrustSet")
}

// SignSetRegularKeyTransaction signs a SetRegularKey transaction
func (r *XRP) SignSetRegularKeyTransaction(
	ctx context.Context, txInput *SetRegularKeyTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "SetRegularKey")
}

// SignAccountSetTransaction signs an AccountSet transaction
func (r *XRP) SignAccountSetTransaction(
	ctx context.Context, txInput *AccountSetTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "AccountSet")
}
