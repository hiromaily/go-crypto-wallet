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
// and multi-signature transactions.
//
// Workflow:
//  1. Validate transaction input (required fields present)
//  2. Derive wallet from seed using Peersyst wallet.FromSeed()
//  3. Convert dtoxrp.TxInput to Peersyst transaction format
//  4. Sign transaction (single-sig or multi-sig based on SigningPubKey field)
//  5. Encode signed transaction to hex blob
//  6. Calculate and return transaction hash
//
// Parameters:
//   - ctx: Context for cancellation control (not used for network, purely offline)
//   - txInput: Unsigned transaction data with all required fields populated
//   - secret: XRP seed/secret in rXXX family seed or hex format
//
// Returns:
//   - string: Transaction hash (64-character hex string)
//   - string: Signed transaction blob (hex-encoded, ready for submission)
//   - error: Returns error if validation, wallet derivation, or signing fails
func (*PeersystSigner) SignTransactionNative(
	ctx context.Context,
	txInput *dtoxrp.TxInput,
	secret string,
) (string, string, error) {
	// Step 1: Validate required fields
	if err := validateTxInput(txInput); err != nil {
		return "", "", fmt.Errorf("transaction validation failed: %w", err)
	}

	// Step 2: Derive wallet from seed
	w, err := wallet.FromSeed(secret, "")
	if err != nil {
		return "", "", fmt.Errorf("failed to derive wallet from seed: %w", err)
	}

	// Step 3: Convert dtoxrp.TxInput to Peersyst transaction format (map)
	tx, needsMultiSig := convertToPeersystTransaction(txInput)

	// Step 4: Sign transaction (use method indicated by conversion)
	var txHash, signedBlob string
	if needsMultiSig {
		// Multi-signature: transaction has SigningPubKey set to empty string
		// Peersyst returns: (signedBlob, txHash, error)
		signedBlob, txHash, err = w.Multisign(tx)
		if err != nil {
			return "", "", fmt.Errorf("failed to multisign transaction: %w", err)
		}
	} else {
		// Single signature: normal signing
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
// Returns the transaction map and a boolean indicating if multi-signature signing is needed.
//
// Multi-sig detection uses fee-based heuristic (see isMultiSigTransaction):
//   - Single-sig: Don't add SigningPubKey to map (wallet.Sign adds it automatically)
//   - Multi-sig: Add SigningPubKey = "" to map (required for wallet.Multisign)
func convertToPeersystTransaction(txInput *dtoxrp.TxInput) (map[string]any, bool) {
	// Build base transaction map
	tx := map[string]any{
		"TransactionType":    "Payment",
		"Account":            txInput.Account,
		"Destination":        txInput.Destination,
		"Amount":             txInput.Amount,
		"Fee":                txInput.Fee,
		"Sequence":           txInput.Sequence,
		"LastLedgerSequence": txInput.LastLedgerSequence,
	}

	// Detect if this is a multi-signature transaction
	// For multi-sig, we need to explicitly set SigningPubKey to empty string
	needsMultiSig := isMultiSigTransaction(txInput)

	if needsMultiSig {
		// Multi-signature: explicitly set empty SigningPubKey
		tx["SigningPubKey"] = ""
	}
	// For single-sig: don't include SigningPubKey (wallet.Sign will add it)

	// Add Flags if non-zero
	if txInput.Flags != 0 {
		tx["Flags"] = txInput.Flags
	}

	return tx, needsMultiSig
}

// isMultiSigTransaction determines if a transaction should use multi-signature signing.
// Since we can't distinguish between "not set" and "set to empty" for SigningPubKey,
// we use other transaction characteristics as indicators.
func isMultiSigTransaction(txInput *dtoxrp.TxInput) bool {
	// Heuristic: Multi-sig transactions typically have:
	// 1. Higher fees: (N+1) × base fee where base fee is typically 10-12 drops
	// 2. SigningPubKey that should be empty
	//
	// We'll use fee as the primary indicator:
	// - Single-sig: fee <= 20 drops (base fee + some buffer)
	// - Multi-sig: fee > 20 drops
	//
	// This is a practical heuristic that works for typical cases.
	// Base fee is usually 10-12 drops, so 20 is a safe threshold.

	// Parse fee to check if it's multi-sig level
	fee := parseFeeDrops(txInput.Fee)

	// If fee is greater than 20 drops, likely multi-sig
	// This catches (2+1)×12=36, (3+1)×12=48, etc.
	return fee > 20
}

// parseFeeDrops parses the fee string to an integer (drops).
// Returns 0 if parsing fails.
func parseFeeDrops(feeStr string) int {
	var fee int
	_, err := fmt.Sscanf(feeStr, "%d", &fee)
	if err != nil {
		return 0
	}
	return fee
}
