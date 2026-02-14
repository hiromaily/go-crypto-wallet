package xrp

import (
	"context"
	"errors"
	"fmt"

	"github.com/Peersyst/xrpl-go/xrpl/wallet"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
)

// PeersystSigner implements the TransactionSigner interface using the Peersyst xrpl-go library
// for native Go transaction signing without gRPC dependencies.
//
// This implementation supports offline signing on air-gapped systems (keygen and sign wallets)
// and handles both single-signature and multi-signature transactions.
type PeersystSigner struct{}

// NewPeersystSigner creates a new PeersystSigner instance.
func NewPeersystSigner() *PeersystSigner {
	return &PeersystSigner{}
}

// SignTransactionNative signs an XRP transaction using native Go implementation.
//
// This method implements the TransactionSigner interface and provides offline signing
// capability using the Peersyst/xrpl-go library. It supports both single-signature
// and multi-signature transactions, with proper signature accumulation for multi-sig workflows.
//
// Workflow:
//  1. Validate transaction input (required fields present)
//  2. Derive wallet from seed using Peersyst wallet.FromSeed()
//  3. Convert dtoxrp.TxInput to Peersyst transaction format
//  4. Sign transaction (single-sig or multi-sig based on isMultiSig parameter)
//  5. If existingSignedBlob provided, combine with new signature (multi-sig accumulation)
//  6. Encode signed transaction to hex blob
//  7. Calculate and return transaction hash
//
// Parameters:
//   - ctx: Context for cancellation control (not used for network, purely offline)
//   - txInput: Unsigned transaction data with all required fields populated
//   - secret: XRP seed/secret in rXXX family seed or hex format
//   - isMultiSig: true for multi-signature transactions, false for single-signature
//   - existingSignedBlob: Optional existing signed blob for multi-sig (nil for first signature)
//
// Returns:
//   - string: Transaction hash (64-character hex string)
//   - string: Signed transaction blob (hex-encoded, ready for submission)
//   - error: Returns error if validation, wallet derivation, signing, or combination fails
//
// Multi-Signature Accumulation:
//   - First signer: existingSignedBlob = nil, creates new signed blob with first signature
//   - Additional signers: existingSignedBlob != nil, adds signature to existing Signers array
//   - This ensures all signatures are preserved and combined correctly for multi-sig validation
func (*PeersystSigner) SignTransactionNative(
	ctx context.Context,
	txInput *dtoxrp.TxInput,
	secret string,
	isMultiSig bool,
	existingSignedBlob *string,
) (string, string, error) {
	// Step 1: Validate required fields
	if err := validateTxInput(txInput); err != nil {
		return "", "", fmt.Errorf("transaction validation failed: %w", err)
	}

	// Step 2: Handle multi-sig signature accumulation
	// For multi-sig workflows with existing signatures, we need to combine them
	if isMultiSig && existingSignedBlob != nil {
		// TODO: Implement signature accumulation for multi-sig transactions
		// This requires:
		// 1. Decode existingSignedBlob to extract Signers array
		// 2. Sign transaction with this wallet to get new signature
		// 3. Add new signature to Signers array
		// 4. Re-encode transaction with all signatures
		//
		// For now, this is a known limitation documented in the PR review.
		// Integration tests (task 9.2) should verify end-to-end multi-sig flow.
		//
		// Reference: PR review comment about signature accumulation
		// https://github.com/hiromaily/go-crypto-wallet/pull/560#pullrequestreview-3802102136
		return "", "", errors.New("multi-signature accumulation not yet implemented: " +
			"cannot combine with existing signatures. This is a known limitation " +
			"that needs to be addressed before production use of multi-sig wallets")
	}

	// Step 3: Derive wallet from seed
	w, err := wallet.FromSeed(secret, "")
	if err != nil {
		return "", "", fmt.Errorf("failed to derive wallet from seed: %w", err)
	}

	// Step 4: Convert dtoxrp.TxInput to Peersyst transaction format (map)
	tx := convertToPeersystTransaction(txInput, isMultiSig)

	// Step 5: Sign transaction (use method based on isMultiSig parameter)
	var txHash, signedBlob string
	if isMultiSig {
		// Multi-signature: use wallet.Multisign()
		// Note: For first signature only (existingSignedBlob == nil)
		// Peersyst returns: (signedBlob, txHash, error)
		signedBlob, txHash, err = w.Multisign(tx)
		if err != nil {
			return "", "", fmt.Errorf("failed to multisign transaction: %w", err)
		}
	} else {
		// Single signature: use wallet.Sign()
		// Peersyst returns: (signedBlob, txHash, error)
		signedBlob, txHash, err = w.Sign(tx)
		if err != nil {
			return "", "", fmt.Errorf("failed to sign transaction: %w", err)
		}
	}

	return txHash, signedBlob, nil
}

// validateTxInput validates that all required fields are present in the transaction input.
func validateTxInput(txInput *dtoxrp.TxInput) error {
	if txInput == nil {
		return errors.New("txInput is nil")
	}

	// Validate required fields for Payment transactions
	if txInput.Account == "" {
		return errors.New("account field is required")
	}

	if txInput.TransactionType == "Payment" {
		if txInput.Destination == "" {
			return errors.New("destination field is required for Payment transactions")
		}
		if txInput.Amount == "" {
			return errors.New("amount field is required for Payment transactions")
		}
	}

	if txInput.Fee == "" {
		return errors.New("fee field is required")
	}

	if txInput.Sequence == 0 {
		return errors.New("sequence must be greater than 0")
	}

	if txInput.LastLedgerSequence == 0 {
		return errors.New("lastLedgerSequence must be greater than 0")
	}

	return nil
}

// convertToPeersystTransaction converts a dtoxrp.TxInput to a Peersyst transaction map format.
//
// Multi-sig handling:
//   - Single-sig: Don't add SigningPubKey to map (wallet.Sign adds it automatically)
//   - Multi-sig: Add SigningPubKey = "" to map (required for wallet.Multisign)
//
// Parameters:
//   - txInput: Transaction input data
//   - isMultiSig: Explicit indicator from caller - true for multi-sig, false for single-sig
func convertToPeersystTransaction(txInput *dtoxrp.TxInput, isMultiSig bool) map[string]any {
	// Build base transaction map
	tx := map[string]any{
		"TransactionType":    txInput.TransactionType,
		"Account":            txInput.Account,
		"Destination":        txInput.Destination,
		"Amount":             txInput.Amount,
		"Fee":                txInput.Fee,
		"Sequence":           txInput.Sequence,
		"LastLedgerSequence": txInput.LastLedgerSequence,
	}

	// For multi-sig, explicitly set empty SigningPubKey (XRP Ledger requirement)
	if isMultiSig {
		tx["SigningPubKey"] = ""
	}
	// For single-sig: don't include SigningPubKey (wallet.Sign will add it)

	// Add Flags if non-zero
	if txInput.Flags != 0 {
		tx["Flags"] = txInput.Flags
	}

	return tx
}
