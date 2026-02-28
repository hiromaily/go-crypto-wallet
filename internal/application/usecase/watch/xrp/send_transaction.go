package xrp

import (
	"context"
	"fmt"
	"strings"
	"sync"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// xrpSendTxClient defines the interface for XRP operations needed by sendTransactionUseCase.
// This follows the Interface Segregation Principle - depend only on what you need.
// TransactionSubmitter already includes SubmitTransaction, WaitValidation, and GetTransaction.
type xrpSendTxClient interface {
	apixrp.TransactionSubmitter
}

type sendTransactionUseCase struct {
	xrper        xrpSendTxClient
	txDetailRepo repowatch.XRPDetailTXRepositorier
	txFileRepo   file.TransactionFileRepositorier
}

// NewSendTransactionUseCase creates a new SendTransactionUseCase.
// The xrper parameter accepts any type that implements TransactionSubmitter.
func NewSendTransactionUseCase(
	xrper xrpSendTxClient,
	txDetailRepo repowatch.XRPDetailTXRepositorier,
	txFileRepo file.TransactionFileRepositorier,
) watchusecase.SendTransactionUseCase {
	return &sendTransactionUseCase{
		xrper:        xrper,
		txDetailRepo: txDetailRepo,
		txFileRepo:   txFileRepo,
	}
}

// How to send multiple transactions
// - Question about the tefPAST_SEQ (https://www.xrpchat.com/topic/33003-question-about-the-tefpast_seq/)
// - atomical multiple transaction support?
//   (https://github.com/ripple/ripple-lib/issues/839)
// - https://stackoverflow.com/questions/57521439/can-i-send-xrp-to-multiple-addresses
// - increment the account sequence number
// - AccountTxnID (https://xrpl.org/transaction-common-fields.html#accounttxnid)
// - Execute multiple transactions atomically
//   (https://www.xrpchat.com/topic/29175-execute-multiple-transactions-atomically/)
// - トランザクションキュー (https://xrpl.org/ja/transaction-queue.html)
// - 結果のファイナリティー (https://xrpl.org/ja/finality-of-results.html)
// - Escrow (https://xrpl.org/ja/escrow.html)

func (u *sendTransactionUseCase) Execute(
	ctx context.Context,
	input watchusecase.SendTransactionInput,
) (watchusecase.SendTransactionOutput, error) {
	// Validate file path and extract transaction metadata
	actionType, _, txID, _, err := u.txFileRepo.ValidateFilePath(input.FilePath, domainTx.TxTypeSigned)
	if err != nil {
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("fail to call txFileRepo.ValidateFilePath(): %w", err)
	}

	logger.Debug("send_tx", "action_type", actionType.String(), "file_path", input.FilePath)

	// Read JSON transaction file
	txFile, err := u.txFileRepo.ReadXRPJSONFile(input.FilePath)
	if err != nil {
		logger.Error("fail to call txFileRepo.ReadXRPJSONFile()", "file_path", input.FilePath, "error", err)
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("fail to call txFileRepo.ReadXRPJSONFile(): %w", err)
	}

	// Validate all transactions are complete before submission
	if err := validateTransactions(txFile.Transactions); err != nil {
		return watchusecase.SendTransactionOutput{}, err
	}

	// Process each signed transaction concurrently
	var wg sync.WaitGroup
	// Channel to collect successful transaction hashes
	txHashChan := make(chan string, len(txFile.Transactions))

	for _, tx := range txFile.Transactions {
		wg.Add(1)
		go func(txEntry dtoxrp.XRPTransactionEntry) {
			defer wg.Done()

			// Check for context cancellation
			select {
			case <-ctx.Done():
				logger.Warn("transaction submission cancelled",
					"tx_id", txID,
					"uuid", txEntry.UUID,
					"error", ctx.Err(),
				)
				return
			default:
			}

			uuid := txEntry.UUID
			txBlob := *txEntry.SignedBlob

			// Submit transaction to XRP network
			sentTx, earliestLedgerVersion, submitErr := u.xrper.SubmitTransaction(ctx, txBlob)
			if submitErr != nil {
				logger.Warn("failed to call xrp.SubmitTransaction()",
					"tx_id", txID,
					"uuid", uuid,
					"sender_account", txEntry.SenderAccount,
					"error", submitErr,
					// https://xrpl.org/tef-codes.html
					// https://xrpl.org/finality-of-results.html
					// tefMAX_LEDGER / Ledger sequence too high
					//  - The error message Ledger sequence too high occurs if you've waited too long to confirm
					//    a transaction in Ledger Live.
					// tefPAST_SEQ / This sequence number has already passed
				)
				return
			}
			if !strings.Contains(sentTx.ResultCode, "tesSUCCESS") {
				logger.Warn("transaction submission failed",
					"tx_id", txID,
					"uuid", uuid,
					"sender_account", txEntry.SenderAccount,
					"result_code", sentTx.ResultCode,
					"result_message", sentTx.ResultMessage,
				)
				return
			}

			// Debug ledger version info
			logger.Debug("ledger version",
				"earliestLedgerVersion", earliestLedgerVersion,
				"sentTx.TxJSON.LastLedgerSequence", sentTx.TxJSON.LastLedgerSequence,
			)

			// Wait for transaction validation
			ledgerVer, waitErr := u.xrper.WaitValidation(ctx, sentTx.TxJSON.LastLedgerSequence)
			if waitErr != nil {
				logger.Warn("failed to call xrp.WaitValidation()",
					"tx_id", txID,
					"uuid", uuid,
					"sender_account", txEntry.SenderAccount,
					"lastLedgerSequence", sentTx.TxJSON.LastLedgerSequence,
					"ledgerVer", ledgerVer,
					"error", waitErr,
					// Transaction has not been validated yet; try again later
				)
				return
			}

			// Get transaction info for verification
			txInfo, getErr := u.xrper.GetTransaction(ctx, sentTx.TxJSON.Hash, earliestLedgerVersion)
			if getErr != nil {
				logger.Warn("failed to call xrp.GetTransaction()",
					"tx_id", txID,
					"uuid", uuid,
					"sender_account", txEntry.SenderAccount,
					"hash", sentTx.TxJSON.Hash,
					"earliestLedgerVersion", earliestLedgerVersion,
					"error", getErr,
				)
				return
			}

			// Log transaction verification result
			logger.Debug("transaction verified on ledger",
				"tx_id", txID,
				"uuid", uuid,
				"hash", sentTx.TxJSON.Hash,
				"result", txInfo.Outcome.Result,
			)

			// Update xrp_detail_tx table
			affectedNum, updateErr := u.txDetailRepo.UpdateAfterTxSent(
				uuid, domainTx.TxTypeSent, sentTx.TxJSON.Hash, txBlob, earliestLedgerVersion)
			if updateErr != nil {
				// TODO: even if error occurred, tx is already sent. so db should be corrected manually
				logger.Warn(
					"failed to call txDetailRepo.UpdateAfterTxSent() but tx is already sent. "+
						"So database should be updated manually",
					"tx_id", txID,
					"uuid", uuid,
					"sender_account", txEntry.SenderAccount,
					"hash", sentTx.TxJSON.Hash,
					"tx_type", domainTx.TxTypeSent.String(),
					"tx_type_value", domainTx.TxTypeSent.Int8(),
					"error", updateErr,
				)
				// "error":"models: unable to update all for xrp_detail_tx: Error 1406:
				// Data too long for column 'signed_tx_blob' at row 1"
				return
			}
			if affectedNum == 0 {
				logger.Info("no records to update tx_table",
					"tx_id", txID,
					"uuid", uuid,
					"sender_account", txEntry.SenderAccount,
					"hash", sentTx.TxJSON.Hash,
					"tx_type", domainTx.TxTypeSent.String(),
					"tx_type_value", domainTx.TxTypeSent.Int8(),
				)
				return
			}

			// Send transaction hash to channel on success
			txHashChan <- sentTx.TxJSON.Hash
		}(tx)
	}
	wg.Wait()
	close(txHashChan)

	// Collect transaction hashes from channel
	var txHashes []string
	for hash := range txHashChan {
		txHashes = append(txHashes, hash)
	}

	// Return the first successful transaction hash, or empty if all failed
	var resultTxID string
	if len(txHashes) > 0 {
		resultTxID = txHashes[0]
		logger.Info("transaction(s) sent successfully",
			"tx_id", txID,
			"result_tx_hash", resultTxID,
			"total_submitted", len(txHashes),
		)
	} else {
		logger.Warn("no transactions were successfully submitted",
			"tx_id", txID,
		)
	}

	// TODO: update is_allocated in account_pubkey_table
	// Not fixed yet, Ripple may use same address because no utxo
	return watchusecase.SendTransactionOutput{
		TxID: resultTxID,
	}, nil
}

// validateTransactions validates that all transactions are complete and ready for submission.
func validateTransactions(transactions []dtoxrp.XRPTransactionEntry) error {
	for i, tx := range transactions {
		// Check completion status
		if !tx.IsComplete {
			return fmt.Errorf(
				"transaction is incomplete: transaction[%d] (uuid=%s) requires %d signatures but has %d",
				i, tx.UUID, tx.RequiredSignatures, tx.SignatureCount,
			)
		}

		// Verify signedBlob is not null
		if tx.SignedBlob == nil {
			return fmt.Errorf(
				"signedBlob is null for transaction[%d] (uuid=%s)",
				i, tx.UUID,
			)
		}

		// Verify signedBlob is not empty
		if *tx.SignedBlob == "" {
			return fmt.Errorf(
				"signedBlob is empty for transaction[%d] (uuid=%s)",
				i, tx.UUID,
			)
		}
	}
	return nil
}
