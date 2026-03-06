package xrp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// xrpSendTxClient defines the interface for XRP operations needed by sendTransactionUseCase.
// It extends TransactionSubmitter with LedgerPoller so the use case can implement
// the ledger-validation polling loop itself (Plan C).
type xrpSendTxClient interface {
	apixrp.TransactionSubmitter
	apixrp.LedgerPoller
}

type sendTransactionUseCase struct {
	xrper        xrpSendTxClient
	ledgerAdv    apixrp.LedgerAdvancer // optional; non-nil only in standalone mode
	txDetailRepo repowatch.XRPDetailTXRepositorier
	txFileRepo   file.TransactionFileRepositorier
}

// NewSendTransactionUseCase creates a new SendTransactionUseCase.
// ledgerAdv may be nil when the admin client is not configured (non-standalone environments).
func NewSendTransactionUseCase(
	xrper xrpSendTxClient,
	ledgerAdv apixrp.LedgerAdvancer,
	txDetailRepo repowatch.XRPDetailTXRepositorier,
	txFileRepo file.TransactionFileRepositorier,
) watchusecase.SendTransactionUseCase {
	return &sendTransactionUseCase{
		xrper:        xrper,
		ledgerAdv:    ledgerAdv,
		txDetailRepo: txDetailRepo,
		txFileRepo:   txFileRepo,
	}
}

// waitValidation polls ledger_current until it reaches or exceeds targetLedgerVersion.
// In standalone mode it calls ledger_accept via the admin client to advance the ledger.
func (u *sendTransactionUseCase) waitValidation(ctx context.Context, targetLedgerVersion uint64) (uint64, error) {
	// Attempt an initial ledger_accept in standalone mode.
	if u.ledgerAdv != nil {
		if _, err := u.ledgerAdv.LedgerAccept(ctx); err != nil {
			logger.Warn("ledger_accept failed (non-critical; may not be standalone mode)", "error", err)
		}
	}

	const maxRetries = 30
	for range maxRetries {
		res, err := u.xrper.LedgerCurrent(ctx)
		if err != nil {
			return 0, fmt.Errorf("fail to call ledger_current: %w", err)
		}
		currentLedger := res.LedgerCurrentIndex
		logger.Info("waitValidation polling",
			"currentLedger", currentLedger,
			"targetLedgerVersion", targetLedgerVersion,
		)
		if currentLedger >= targetLedgerVersion {
			return currentLedger, nil
		}
		if u.ledgerAdv != nil {
			if _, err := u.ledgerAdv.LedgerAccept(ctx); err != nil {
				logger.Warn("failed to advance ledger", "error", err)
			}
		}
		time.Sleep(1 * time.Second)
	}
	return 0, errors.New("timeout waiting for ledger validation")
}

// signedTxEntry holds the parsed fields from a signed transaction CSV line.
type signedTxEntry struct {
	UUID   string
	TxBlob string
}

// submitOneTx submits a single signed transaction and waits for confirmation.
func (u *sendTransactionUseCase) submitOneTx(ctx context.Context, e signedTxEntry, txID int64) string {
	select {
	case <-ctx.Done():
		logger.Warn("transaction submission cancelled", "tx_id", txID, "uuid", e.UUID, "error", ctx.Err())
		return ""
	default:
	}

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

	ledgerVer, err := u.waitValidation(ctx, sentTx.TxJSON.LastLedgerSequence)
	if err != nil {
		logger.Warn("failed to wait for ledger validation",
			"tx_id", txID, "uuid", e.UUID,
			"lastLedgerSequence", sentTx.TxJSON.LastLedgerSequence, "ledgerVer", ledgerVer, "error", err)
		return ""
	}

	txInfo, err := u.xrper.GetTransaction(ctx, sentTx.TxJSON.Hash, earliestLedgerVersion)
	if err != nil {
		logger.Warn("failed to call xrp.GetTransaction()",
			"tx_id", txID, "uuid", e.UUID,
			"hash", sentTx.TxJSON.Hash, "earliestLedgerVersion", earliestLedgerVersion, "error", err)
		return ""
	}
	logger.Debug("transaction verified on ledger",
		"tx_id", txID, "uuid", e.UUID, "hash", sentTx.TxJSON.Hash, "result", txInfo.Outcome.Result)

	affectedNum, err := u.txDetailRepo.UpdateAfterTxSent(
		e.UUID, domainTx.TxTypeSent, sentTx.TxJSON.Hash, e.TxBlob, earliestLedgerVersion)
	if err != nil {
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
	actionType, _, txID, _, err := u.txFileRepo.ValidateFilePath(input.FilePath, domainTx.TxTypeSigned)
	if err != nil {
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("fail to call txFileRepo.ValidateFilePath(): %w", err)
	}

	logger.Debug("send_tx", "action_type", actionType.String(), "file_path", input.FilePath)

	lines, err := u.txFileRepo.ReadFileSlice(input.FilePath)
	if err != nil {
		logger.Error("fail to call txFileRepo.ReadFileSlice()", "file_path", input.FilePath, "error", err)
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("fail to call txFileRepo.ReadFileSlice(): %w", err)
	}

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

	var wg sync.WaitGroup
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

	var txHashes []string
	for hash := range txHashChan {
		txHashes = append(txHashes, hash)
	}

	var resultTxID string
	if len(txHashes) > 0 {
		resultTxID = txHashes[0]
		logger.Info("transaction(s) sent successfully",
			"tx_id", txID,
			"result_tx_hash", resultTxID,
			"total_submitted", len(txHashes),
		)
	} else {
		logger.Warn("no transactions were successfully submitted", "tx_id", txID)
	}

	return watchusecase.SendTransactionOutput{
		TxID: resultTxID,
	}, nil
}
