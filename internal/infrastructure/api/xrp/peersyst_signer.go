package xrp

import (
	"context"
	"errors"
	"fmt"

	binarycodec "github.com/Peersyst/xrpl-go/binary-codec"
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
		return combineMultiSigSignatures(txInput, secret, *existingSignedBlob)
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

// combineMultiSigSignatures combines an existing multi-sig transaction blob with a new signature.
//
// This function implements the multi-signature accumulation workflow:
//  1. Decode the existing signed blob to extract the Signers array
//  2. Sign the original transaction with the new wallet to get a new signature
//  3. Decode the newly signed blob to extract the new Signer entry
//  4. Combine both Signers arrays (existing + new)
//  5. Create the final transaction with the combined Signers array
//  6. Re-encode and calculate the transaction hash
//
// Parameters:
//   - txInput: Original unsigned transaction data
//   - secret: Seed/secret for the new signer
//   - existingSignedBlob: Hex-encoded signed transaction blob with existing signatures
//
// Returns:
//   - string: Transaction hash (64-character hex string)
//   - string: Combined signed transaction blob (hex-encoded, ready for submission)
//   - error: Returns error if decoding, signing, or combination fails
func combineMultiSigSignatures(
	txInput *dtoxrp.TxInput,
	secret string,
	existingSignedBlob string,
) (string, string, error) {
	// Step 1: Decode existing signed blob to extract Signers array
	existingTx, err := binarycodec.Decode(existingSignedBlob)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode existing signed blob: %w", err)
	}

	// Extract existing Signers array
	existingSigners, ok := existingTx["Signers"].([]any)
	if !ok {
		return "", "", errors.New("existing signed blob does not contain valid Signers array")
	}

	// Step 2: Sign the original transaction with the new wallet
	w, err := wallet.FromSeed(secret, "")
	if err != nil {
		return "", "", fmt.Errorf("failed to derive wallet from seed: %w", err)
	}

	// Convert txInput to transaction map for signing
	tx := convertToPeersystTransaction(txInput, true)

	// Sign with the new wallet to get a new signature
	// Note: Multisign returns (blob, hash, error) - we need the hash for the final result
	newSignedBlob, txHash, err := w.Multisign(tx)
	if err != nil {
		return "", "", fmt.Errorf("failed to multisign transaction: %w", err)
	}

	// Step 3: Decode the newly signed blob to extract the new Signer entry
	newTx, err := binarycodec.Decode(newSignedBlob)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode newly signed blob: %w", err)
	}

	newSigners, ok := newTx["Signers"].([]any)
	if !ok || len(newSigners) == 0 {
		return "", "", errors.New("newly signed blob does not contain valid Signers array")
	}

	// Step 4: Combine Signers arrays (existing + new)
	// The Signers array should contain all signatures
	existingSigners = append(existingSigners, newSigners[0])

	// Step 5: Create the final transaction with combined Signers array
	// Use the original transaction structure but with combined Signers
	finalTx := convertToPeersystTransaction(txInput, true)
	finalTx["Signers"] = existingSigners

	// Step 6: Re-encode the combined transaction
	combinedBlob, err := binarycodec.Encode(finalTx)
	if err != nil {
		return "", "", fmt.Errorf("failed to encode combined transaction: %w", err)
	}

	// Return the transaction hash from the new signature (same for all signers of the same tx)
	// and the combined blob with all signatures
	return txHash, combinedBlob, nil
}
