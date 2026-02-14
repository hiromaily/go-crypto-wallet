package xrp

import (
	"context"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
)

// TransactionSigner defines an interface for offline XRP transaction signing.
//
// This interface follows the Interface Segregation Principle by providing only
// signing operations, with NO network dependencies. This enables offline signing
// on air-gapped systems (keygen and sign wallets).
//
// # Responsibilities
//   - Sign XRP transactions (legacy gRPC-based and native Go implementation)
//   - Support both single-signature and multi-signature transactions
//   - Generate signed transaction blob and transaction hash
//
// # Constraints
//   - ZERO network calls (offline-capable interface)
//   - No account query methods (those belong to AccountInfoProvider)
//   - No transaction submission methods (those belong to TransactionSubmitter)
//   - Deterministic signing (same input + secret = same signature)
//
// # Implementation Notes
//   - SignTransaction: Legacy gRPC-based signing (existing implementation)
//   - SignTransactionNative: New native Go signing via Peersyst/xrpl-go (task 3.1)
//   - Used by SignTransactionUseCase on offline keygen/sign wallets
//
// # Migration Path
//   - Current code uses SignTransaction (gRPC)
//   - Task 5.1 will migrate to SignTransactionNative (native Go)
//   - Both methods supported during transition
//
// # Security
//   - Never logs secret parameter
//   - Operates entirely offline (no network access)
//   - Supports air-gapped signing workflows
//
// # Related Interfaces
//   - AccountInfoProvider: For account queries (online only)
//   - TransactionSubmitter: For submitting signed transactions (online only)
type TransactionSigner interface {
	// SignTransaction signs a transaction using gRPC-based signing (legacy method).
	//
	// Deprecated: Use SignTransactionNative for offline native Go signing.
	// This method is kept for backward compatibility with existing code.
	// Task 5.1 will migrate existing use cases to SignTransactionNative.
	//
	// Parameters:
	//   - ctx: Context for cancellation control
	//   - txJSON: Unsigned transaction data
	//   - secret: XRP seed/secret
	//
	// Returns:
	//   - string: Transaction hash
	//   - string: Signed transaction blob
	//   - error: Signing error
	SignTransaction(ctx context.Context, txJSON *dtoxrp.TxInput, secret string) (string, string, error)

	// SignTransactionNative signs an XRP transaction using native Go implementation.
	//
	// This method uses the Peersyst/xrpl-go library to sign transactions entirely
	// offline without any gRPC or network dependencies. It supports both single-signature
	// and multi-signature transactions.
	//
	// Parameters:
	//   - ctx: Context for cancellation control (not used for network, purely offline)
	//   - txInput: Unsigned transaction data with all required fields populated
	//   - secret: XRP seed/secret in rXXX family seed or hex format
	//   - isMultiSig: Explicit indicator - true for multi-signature, false for single-signature
	//
	// Returns:
	//   - string: Transaction hash (64-character hex string, unique identifier)
	//   - string: Signed transaction blob (hex-encoded, ready for submission)
	//   - error: Returns error if:
	//     * txInput has missing or invalid required fields
	//     * secret is invalid or malformed
	//     * transaction type not supported
	//     * signing operation fails
	//
	// Preconditions:
	//   - txInput must have all required fields:
	//     * Account (sender address)
	//     * Destination (receiver address for Payment transactions)
	//     * Amount (transaction amount)
	//     * Fee (transaction fee)
	//     * Sequence (account sequence number)
	//     * LastLedgerSequence (transaction expiration)
	//   - secret must be valid XRP seed format:
	//     * rXXX family seed format, OR
	//     * ed25519/secp256k1 hex seed format
	//   - isMultiSig must correctly indicate the signature type
	//     * true: uses wallet.Multisign() for multi-signature transactions
	//     * false: uses wallet.Sign() for single-signature transactions
	//
	// Postconditions:
	//   - Returns 64-character hex transaction hash
	//   - Returns hex-encoded signed transaction blob
	//   - Signed blob includes signature in TxnSignature field (single-sig)
	//     or Signers array (multi-sig)
	//   - NO network calls made (offline guarantee)
	//   - Deterministic: same txInput + secret produces same signature
	//
	// Multi-Signature Behavior:
	//   - Detects multi-sig if txInput.SigningPubKey is empty or Signers array present
	//   - Uses Peersyst wallet.Multisign() for multi-sig transactions
	//   - Uses Peersyst wallet.Sign() for single-sig transactions
	//   - Multi-sig fee: (N+1) × base fee where N = number of signatures
	//
	// Example Usage (Single Signature):
	//   txInput := &dtoxrp.TxInput{
	//       Account:            "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH",
	//       Destination:        "rLHzPsX6oXkzU9J1kSwq7X4X8J1Xa8AkDC",
	//       Amount:             "1000000", // 1 XRP in drops
	//       Fee:                "12",      // Base fee
	//       Sequence:           42,
	//       LastLedgerSequence: 1000000,
	//   }
	//   txID, signedBlob, err := signer.SignTransactionNative(ctx, txInput, "sEdTM1uX8pu2do5XvTnutH6HsouMaM2", false)
	//   if err != nil {
	//       return fmt.Errorf("failed to sign transaction: %w", err)
	//   }
	//
	// Example Usage (Multi-Signature):
	//   txInput := &dtoxrp.TxInput{
	//       Account:            "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH",
	//       Destination:        "rLHzPsX6oXkzU9J1kSwq7X4X8J1Xa8AkDC",
	//       Amount:             "1000000",
	//       Fee:                "36",  // (2+1) × 12 for 2-of-3 multisig
	//       Sequence:           42,
	//       LastLedgerSequence: 1000000,
	//   }
	//   txID, signedBlob, err := signer.SignTransactionNative(ctx, txInput, secret, true)
	//
	// Security Notes:
	//   - secret parameter is NEVER logged
	//   - Operates entirely offline (air-gapped signing)
	//   - No network access during signing operation
	SignTransactionNative(
		ctx context.Context, txInput *dtoxrp.TxInput, secret string, isMultiSig bool,
	) (string, string, error)
}
