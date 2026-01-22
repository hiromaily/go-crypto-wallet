package xrp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/bookerzzz/grok"

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
type xrpSendTxClient interface {
	apixrp.TransactionSubmitter
	apixrp.LedgerWaiter
	apixrp.TransactionGetter
}

type sendTransactionUseCase struct {
	xrper        xrpSendTxClient
	txDetailRepo repowatch.XRPDetailTXRepositorier
	txFileRepo   file.TransactionFileRepositorier
}

// NewSendTransactionUseCase creates a new SendTransactionUseCase.
// The xrper parameter accepts any type that implements xrpSendTxClient
// (TransactionSubmitter + LedgerWaiter + TransactionGetter).
// Typically, apixrp.XRPer is passed which implements all required methods.
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
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("failed to validate file path: %w", err)
	}

	logger.Debug("send_tx", "action_type", actionType.String())

	// Read and parse JSON transaction file
	jsonData, err := u.txFileRepo.ReadJSONFile(input.FilePath)
	if err != nil {
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("failed to read JSON file: %w", err)
	}

	signedFile, err := dtoxrp.SignedTransactionFileFromJSON(jsonData)
	if err != nil {
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("failed to parse signed transaction file: %w", err)
	}

	// Validate that we have signed transactions
	if len(signedFile.Transactions) == 0 {
		return watchusecase.SendTransactionOutput{}, errors.New("no signed transactions in file")
	}

	// Process each signed transaction concurrently
	var wg sync.WaitGroup

	for _, signedTx := range signedFile.Transactions {
		wg.Add(1)
		go func(tx dtoxrp.SignedTransactionEntry) {
			defer wg.Done()

			// Validate transaction has required data
			if tx.SignedTxBlob == "" {
				logger.Warn("skipping transaction with empty signed blob", "uuid", tx.UUID)
				return
			}

			// Submit transaction to XRP network
			sentTx, earliestLedgerVersion, submitErr := u.xrper.SubmitTransaction(ctx, tx.SignedTxBlob)
			if submitErr != nil {
				logger.Warn("failed to submit transaction",
					"tx_id", txID,
					"uuid", tx.UUID,
					"signed_tx_id", tx.SignedTxID,
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
					"uuid", tx.UUID,
					"signed_tx_id", tx.SignedTxID,
					"result_code", sentTx.ResultCode,
					"result_message", sentTx.ResultMessage,
				)
				return
			}

			// Debug ledger version info
			logger.Debug("ledger version",
				"earliestLedgerVersion", earliestLedgerVersion,
				"lastLedgerSequence", tx.LastLedgerSequence,
				"sentTx.TxJSON.LastLedgerSequence", sentTx.TxJSON.LastLedgerSequence,
			)

			// Wait for transaction validation
			ledgerVer, waitErr := u.xrper.WaitValidation(ctx, sentTx.TxJSON.LastLedgerSequence)
			if waitErr != nil {
				logger.Warn("failed to wait for validation",
					"tx_id", txID,
					"uuid", tx.UUID,
					"signed_tx_id", tx.SignedTxID,
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
				logger.Warn("failed to get transaction info",
					"tx_id", txID,
					"uuid", tx.UUID,
					"signed_tx_id", tx.SignedTxID,
					"hash", sentTx.TxJSON.Hash,
					"earliestLedgerVersion", earliestLedgerVersion,
					"error", getErr,
				)
				return
			}
			// for debug (should be removed later)
			grok.Value(txInfo)

			// Update xrp_detail_tx table
			affectedNum, updateErr := u.txDetailRepo.UpdateAfterTxSent(
				tx.UUID, domainTx.TxTypeSent, tx.SignedTxID, tx.SignedTxBlob, earliestLedgerVersion)
			if updateErr != nil {
				// TODO: even if error occurred, tx is already sent. so db should be corrected manually
				logger.Warn(
					"failed to update database after tx sent - manual correction required",
					"tx_id", txID,
					"uuid", tx.UUID,
					"signed_tx_id", tx.SignedTxID,
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
					"uuid", tx.UUID,
					"signed_tx_id", tx.SignedTxID,
					"tx_type", domainTx.TxTypeSent.String(),
					"tx_type_value", domainTx.TxTypeSent.Int8(),
				)
				return
			}

			logger.Info("transaction sent successfully",
				"uuid", tx.UUID,
				"signed_tx_id", tx.SignedTxID,
				"sender", tx.SenderAccount,
				"receiver", tx.ReceiverAccount,
				"amount", tx.Amount,
			)
		}(signedTx)
	}
	wg.Wait()

	// TODO: update is_allocated in account_pubkey_table
	// Not fixed yet, Ripple may use same address because no utxo
	return watchusecase.SendTransactionOutput{
		TxID: "",
	}, nil
}
