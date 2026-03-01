package xrp

import (
	"context"

	xrpclient "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/client"
)

// TransactionSubmitter defines an interface for submitting and tracking XRP transactions.
//
// This interface follows the Interface Segregation Principle by providing only
// submission and validation operations needed by the SendTransactionUseCase.
//
// # Responsibilities
//   - Submit signed transactions to XRP Ledger
//   - Wait for ledger validation confirmation
//   - Retrieve transaction details and status
//
// # Constraints
//   - Online-only operations (requires network connectivity)
//   - No account query methods (those belong to AccountInfoProvider)
//   - No signing methods (those belong to TransactionSigner)
//   - Idempotent submission (submitting same blob returns same result)
//
// # Implementation Notes
//   - Implemented by xrpscan/xrpl-go client in infrastructure layer
//   - Used by SendTransactionUseCase on online watch wallet
//   - Handles XRP Ledger error codes (tefXXX, tecXXX)
//
// # Related Interfaces
//   - AccountInfoProvider: For account queries
//   - TransactionSigner: For signing operations (offline)
type TransactionSubmitter interface {
	// SubmitTransaction broadcasts a signed transaction to the XRP Ledger.
	//
	// This method submits a hex-encoded signed transaction blob to the XRP Ledger
	// and returns the transaction hash and ledger version where it was included.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - signedTx: Hex-encoded signed transaction blob (from TransactionSigner)
	//
	// Returns:
	//   - SentTx: Transaction submission result with hash and preliminary status
	//   - uint64: Ledger version where transaction was included
	//   - error: Returns error if:
	//     * signedTx is invalid or malformed hex
	//     * transaction validation fails (tefXXX, tecXXX codes)
	//     * network communication fails
	//     * ledger is unavailable
	//
	// Preconditions:
	//   - signedTx must be valid hex-encoded signed transaction blob
	//   - Transaction must have sufficient signatures (complete multi-sig if applicable)
	//   - Network connectivity required (online operation)
	//
	// Postconditions:
	//   - Transaction broadcast to XRP Ledger
	//   - Returns transaction hash (unique identifier)
	//   - Returns ledger version where included
	//   - Idempotent: submitting same blob returns same hash
	//
	// XRP Ledger Error Codes:
	//   - tefPAST_SEQ: Transaction sequence number already used
	//   - tecUNFUNDED_PAYMENT: Sender account has insufficient balance
	//   - tecNO_DST: Destination account does not exist
	//   - temBAD_FEE: Transaction fee invalid (check multi-sig fee calculation)
	//   - tefMAX_LEDGER: Transaction expired (LastLedgerSequence passed)
	//
	// Example Usage:
	//   sentTx, ledgerVersion, err := submitter.SubmitTransaction(ctx, signedBlob)
	//   if err != nil {
	//       return fmt.Errorf("failed to submit transaction: %w", err)
	//   }
	//   log.Infof("Transaction submitted: hash=%s, ledger=%d", sentTx.Hash, ledgerVersion)
	SubmitTransaction(ctx context.Context, signedTx string) (*xrpclient.SentTx, uint64, error)

	// WaitValidation waits for the XRP Ledger to reach a target ledger version.
	//
	// This method blocks until the ledger advances to at least the specified version,
	// confirming that transactions in that ledger are validated and immutable.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - targetLedgerVersion: Minimum ledger version to wait for
	//
	// Returns:
	//   - uint64: Current validated ledger version (>= targetLedgerVersion)
	//   - error: Returns error if:
	//     * context timeout or cancellation
	//     * network communication fails
	//     * ledger is unavailable
	//
	// Preconditions:
	//   - targetLedgerVersion must be > 0
	//   - Network connectivity required (online operation)
	//
	// Postconditions:
	//   - Blocks until ledger version >= targetLedgerVersion
	//   - Returns current validated ledger version
	//   - Transactions in target ledger are confirmed immutable
	//
	// Timeout Handling:
	//   - Use context with timeout to avoid indefinite blocking:
	//     ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//     defer cancel()
	//
	// Example Usage:
	//   targetLedger := sentTx.LedgerVersion + 1
	//   validatedLedger, err := submitter.WaitValidation(ctx, targetLedger)
	//   if err != nil {
	//       return fmt.Errorf("failed to wait for validation: %w", err)
	//   }
	//   log.Infof("Transaction validated at ledger %d", validatedLedger)
	WaitValidation(ctx context.Context, targetLedgerVersion uint64) (uint64, error)

	// GetTransaction retrieves detailed information about a submitted transaction.
	//
	// This method queries the XRP Ledger for transaction details including status,
	// result code, metadata, and ledger inclusion information.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - txID: Transaction hash (from SubmitTransaction or SignTransactionNative)
	//   - targetLedgerVersion: Ledger version to search (0 = search all ledgers)
	//
	// Returns:
	//   - TxInfo: Transaction details including status, result, and metadata
	//   - error: Returns error if:
	//     * txID is invalid or malformed
	//     * transaction not found in ledger
	//     * network communication fails
	//     * ledger is unavailable
	//
	// Preconditions:
	//   - txID must be valid 64-character hex transaction hash
	//   - targetLedgerVersion >= 0 (0 means search all ledgers)
	//   - Network connectivity required (online operation)
	//
	// Postconditions:
	//   - Returns transaction details if found
	//   - Returns error if transaction not found or not yet validated
	//   - No ledger state modifications
	//
	// Transaction Status Values:
	//   - tesSUCCESS: Transaction succeeded (most common success code)
	//   - tecXXX: Transaction failed but fee consumed (e.g., tecUNFUNDED_PAYMENT)
	//   - tefXXX: Transaction failed, fee not consumed (e.g., tefPAST_SEQ)
	//
	// Example Usage:
	//   txInfo, err := submitter.GetTransaction(ctx, txHash, ledgerVersion)
	//   if err != nil {
	//       return fmt.Errorf("failed to get transaction: %w", err)
	//   }
	//   if txInfo.Status != "tesSUCCESS" {
	//       return fmt.Errorf("transaction failed with status: %s", txInfo.Status)
	//   }
	//   log.Infof("Transaction confirmed: %+v", txInfo.Meta)
	GetTransaction(ctx context.Context, txID string, targetLedgerVersion uint64) (*xrpclient.TxInfo, error)
}
