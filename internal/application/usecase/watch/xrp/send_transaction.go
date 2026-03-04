package xrp

import (
	"context"
	"fmt"
	"strings"
	"sync"

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

// signedTxEntry holds the parsed fields from a signed transaction CSV line.
// File format written by the keygen sign use case: uuid,txHash,txBlob
type signedTxEntry struct {
	UUID   string
	TxBlob string
}

// submitOneTx submits a single signed transaction and waits for confirmation.
// Returns the transaction hash on success, or empty string on any error.
func (u *sendTransactionUseCase) submitOneTx(ctx context.Context, e signedTxEntry, txID int64) string {
	select {
	case <-ctx.Done():
		logger.Warn("transaction submission cancelled", "tx_id", txID, "uuid", e.UUID, "error", ctx.Err())
		return ""
	default:
	}

	// Submit transaction to XRP network
	// https://xrpl.org/tef-codes.html, https://xrpl.org/finality-of-results.html
	sentTx, earliestLedgerVersion, err := u.xrper.SubmitTransaction(ctx, e.TxBlob)
	if err != nil {
		logger.Warn("failed to call xrp.SubmitTransaction()", "tx_id", txID, "uuid", e.UUID, "error", err)
		return ""
	}
	if !strings.Contains(sentTx.ResultCode, "tesSUCCESS") {
		logger.Warn("transaction submission failed",
			"tx_id", txID, "uuid", e.UUID,
			"result_code", sentTx.ResultCode, "result_message", sentTx.ResultMessage)
		return ""
	}

	logger.Debug("ledger version",
		"earliestLedgerVersion", earliestLedgerVersion,
		"sentTx.TxJSON.LastLedgerSequence", sentTx.TxJSON.LastLedgerSequence)

	// Wait for transaction validation
	ledgerVer, err := u.xrper.WaitValidation(ctx, sentTx.TxJSON.LastLedgerSequence)
	if err != nil {
		logger.Warn("failed to call xrp.WaitValidation()",
			"tx_id", txID, "uuid", e.UUID,
			"lastLedgerSequence", sentTx.TxJSON.LastLedgerSequence, "ledgerVer", ledgerVer, "error", err)
		return ""
	}

	// Get transaction info for verification
	txInfo, err := u.xrper.GetTransaction(ctx, sentTx.TxJSON.Hash, earliestLedgerVersion)
	if err != nil {
		logger.Warn("failed to call xrp.GetTransaction()",
			"tx_id", txID, "uuid", e.UUID,
			"hash", sentTx.TxJSON.Hash, "earliestLedgerVersion", earliestLedgerVersion, "error", err)
		return ""
	}
	logger.Debug("transaction verified on ledger",
		"tx_id", txID, "uuid", e.UUID, "hash", sentTx.TxJSON.Hash, "result", txInfo.Outcome.Result)

	// Update xrp_detail_tx table
	affectedNum, err := u.txDetailRepo.UpdateAfterTxSent(
		e.UUID, domainTx.TxTypeSent, sentTx.TxJSON.Hash, e.TxBlob, earliestLedgerVersion)
	if err != nil {
		// TODO: even if error occurred, tx is already sent. so db should be corrected manually
		logger.Warn(
			"failed to call txDetailRepo.UpdateAfterTxSent() but tx is already sent. "+
				"So database should be updated manually",
			"tx_id", txID, "uuid", e.UUID, "hash", sentTx.TxJSON.Hash,
			"tx_type", domainTx.TxTypeSent.String(), "tx_type_value", domainTx.TxTypeSent.Int8(), "error", err)
		return ""
	}
	if affectedNum == 0 {
		logger.Info("no records to update tx_table",
			"tx_id", txID, "uuid", e.UUID, "hash", sentTx.TxJSON.Hash,
			"tx_type", domainTx.TxTypeSent.String(), "tx_type_value", domainTx.TxTypeSent.Int8())
		return ""
	}

	return sentTx.TxJSON.Hash
}

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

	// Read signed transaction file (CSV format: uuid,txHash,txBlob per line)
	lines, err := u.txFileRepo.ReadFileSlice(input.FilePath)
	if err != nil {
		logger.Error("fail to call txFileRepo.ReadFileSlice()", "file_path", input.FilePath, "error", err)
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("fail to call txFileRepo.ReadFileSlice(): %w", err)
	}

	// Parse CSV lines into signed tx entries
	entries := make([]signedTxEntry, 0, len(lines))
	for _, line := range lines {
		parts := strings.SplitN(line, ",", 3)
		if len(parts) != 3 || parts[0] == "" || parts[2] == "" {
			return watchusecase.SendTransactionOutput{}, fmt.Errorf("invalid signed tx line format: %s", line)
		}
		entries = append(entries, signedTxEntry{UUID: parts[0], TxBlob: parts[2]})
	}

	if len(entries) == 0 {
		return watchusecase.SendTransactionOutput{}, fmt.Errorf(
			"no signed transactions found in file: %s", input.FilePath)
	}

	// Process each signed transaction concurrently
	var wg sync.WaitGroup
	// Channel to collect successful transaction hashes
	txHashChan := make(chan string, len(entries))

	for _, entry := range entries {
		wg.Add(1)
		go func(e signedTxEntry) {
			defer wg.Done()
			if hash := u.submitOneTx(ctx, e, txID); hash != "" {
				txHashChan <- hash
			}
		}(entry)
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
