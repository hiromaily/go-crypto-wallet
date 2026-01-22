// Package xrp provides native Go XRP transaction signing implementation.
// This package enables offline transaction signing without requiring
// the xrpl-grpc-server to be running.
package xrp

import (
	"context"
	"errors"
	"fmt"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	xrpcrypto "github.com/hiromaily/go-crypto-wallet/pkg/cryptocurrency/xrp"
)

// NativeSigner implements XRP transaction signing using native Go code.
// This allows offline signing without requiring a gRPC connection.
type NativeSigner struct {
	algorithm xrpcrypto.KeyAlgorithm
	signer    *xrpcrypto.Signer
}

// NewNativeSigner creates a new native XRP signer with the specified algorithm.
//
// The algorithm parameter determines which cryptographic algorithm to use:
//   - xrpcrypto.AlgorithmEd25519: Recommended for new accounts
//   - xrpcrypto.AlgorithmSecp256k1: For legacy compatibility
func NewNativeSigner(algorithm xrpcrypto.KeyAlgorithm) *NativeSigner {
	return &NativeSigner{
		algorithm: algorithm,
		signer:    xrpcrypto.NewSigner(algorithm),
	}
}

// NewNativeSignerEd25519 creates a native signer using Ed25519 (recommended).
func NewNativeSignerEd25519() *NativeSigner {
	return NewNativeSigner(xrpcrypto.AlgorithmEd25519)
}

// NewNativeSignerSecp256k1 creates a native signer using secp256k1.
func NewNativeSignerSecp256k1() *NativeSigner {
	return NewNativeSigner(xrpcrypto.AlgorithmSecp256k1)
}

// SignTransaction signs an XRP transaction using native Go code.
//
// This method implements the TransactionSigner interface and can be used
// as a drop-in replacement for the gRPC-based signing.
//
// Parameters:
//   - ctx: Context for cancellation (unused in native signing but kept for interface compatibility)
//   - txInput: The transaction to sign
//   - secret: The master seed (s... format or hex)
//
// Returns:
//   - txID: The transaction hash/ID
//   - txBlob: The signed transaction blob in hex format
//   - error: Any error that occurred during signing
//
// SECURITY NOTE: The secret parameter contains sensitive key material.
// It must be handled securely and never logged.
func (s *NativeSigner) SignTransaction(
	_ context.Context, txInput *dtoxrp.TxInput, secret string,
) (string, string, error) {
	if txInput == nil {
		return "", "", errors.New("txInput cannot be nil")
	}
	if secret == "" {
		return "", "", errors.New("secret cannot be empty")
	}

	// Convert DTO to map for signing
	txFields := dtoTxInputToMap(txInput)

	// Sign the transaction
	result, err := s.signer.SignTransaction(txFields, secret)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	return result.TxID, result.TxBlob, nil
}

// SignTransactionWithDetails signs a transaction and returns detailed result.
//
// This method provides more information than SignTransaction, including
// the public key and signature separately.
func (s *NativeSigner) SignTransactionWithDetails(
	_ context.Context, txInput *dtoxrp.TxInput, secret string,
) (*SigningResult, error) {
	if txInput == nil {
		return nil, errors.New("txInput cannot be nil")
	}
	if secret == "" {
		return nil, errors.New("secret cannot be empty")
	}

	// Convert DTO to map for signing
	txFields := dtoTxInputToMap(txInput)

	// Sign the transaction
	result, err := s.signer.SignTransaction(txFields, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return &SigningResult{
		TxID:          result.TxID,
		TxBlob:        result.TxBlob,
		SigningPubKey: result.SigningPubKey,
		TxnSignature:  result.TxnSignature,
	}, nil
}

// SigningResult contains detailed signing result.
type SigningResult struct {
	// TxID is the transaction hash/ID
	TxID string
	// TxBlob is the signed transaction blob in hex format
	TxBlob string
	// SigningPubKey is the public key used for signing (hex)
	SigningPubKey string
	// TxnSignature is the signature (hex)
	TxnSignature string
}

// Algorithm returns the cryptographic algorithm used by this signer.
func (s *NativeSigner) Algorithm() xrpcrypto.KeyAlgorithm {
	return s.algorithm
}

// dtoTxInputToMap converts a DTO TxInput to a map for signing.
func dtoTxInputToMap(tx *dtoxrp.TxInput) map[string]any {
	m := make(map[string]any)

	if tx.TransactionType != "" {
		m["TransactionType"] = tx.TransactionType
	}
	if tx.Account != "" {
		m["Account"] = tx.Account
	}
	if tx.Amount != "" {
		m["Amount"] = tx.Amount
	}
	if tx.Destination != "" {
		m["Destination"] = tx.Destination
	}
	if tx.Fee != "" {
		m["Fee"] = tx.Fee
	}
	if tx.Flags != 0 {
		m["Flags"] = tx.Flags
	}
	if tx.LastLedgerSequence != 0 {
		m["LastLedgerSequence"] = tx.LastLedgerSequence
	}
	if tx.Sequence != 0 {
		m["Sequence"] = tx.Sequence
	}
	// Don't include SigningPubKey and TxnSignature - they will be added during signing
	// Don't include Hash - it's computed during signing

	return m
}

// VerifySignature verifies an XRP transaction signature.
//
// Parameters:
//   - txInput: The signed transaction to verify (must include SigningPubKey and TxnSignature)
//
// Returns true if the signature is valid.
func (s *NativeSigner) VerifySignature(txInput *dtoxrp.TxInput) (bool, error) {
	if txInput == nil {
		return false, errors.New("txInput cannot be nil")
	}
	if txInput.SigningPubKey == "" {
		return false, errors.New("missing SigningPubKey")
	}
	if txInput.TxnSignature == "" {
		return false, errors.New("missing TxnSignature")
	}

	// Reuse dtoTxInputToMap and add signature fields
	txFields := dtoTxInputToMap(txInput)
	txFields["SigningPubKey"] = txInput.SigningPubKey
	txFields["TxnSignature"] = txInput.TxnSignature

	return s.signer.VerifySignature(txFields)
}
