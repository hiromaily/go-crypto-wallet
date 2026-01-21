package xrp

import (
	"context"
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
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("fail to call txFileRepo.ValidateFilePath(): %w", err)
	}

	logger.Debug("send_tx", "action_type", actionType.String())

	// Read hex from file
	data, err := u.txFileRepo.ReadFileSlice(input.FilePath)
	if err != nil {
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("fail to call txFileRepo.ReadFile(): %w", err)
	}

	// Process each signed transaction concurrently
	var wg sync.WaitGroup

	for _, txHex := range data {
		wg.Add(1)
		go func(line string) {
			defer wg.Done()

			// Parse transaction data: uuid, signedTxID, txBlob
			tmp := strings.Split(line, ",")
			if len(tmp) != 3 {
				logger.Warn("data format is invalid in file")
				return
			}
			uuid := tmp[0]
			signedTxID := tmp[1]
			txBlob := tmp[2]

			// Submit transaction to XRP network
			var sentTx *dtoxrp.SentTx
			var earlistLedgerVersion uint64
			sentTx, earlistLedgerVersion, err = u.xrper.SubmitTransaction(ctx, txBlob)
			if err != nil {
				logger.Warn("fail to call xrp.SubmitTransaction()",
					"tx_id", txID,
					"uuid", uuid,
					"signed_tx_id", signedTxID,
					"error", err,
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
				logger.Warn("fail to call SubmitTransaction",
					"tx_id", txID,
					"uuid", uuid,
					"signed_tx_id", signedTxID,
					"result_code", sentTx.ResultCode,
					"result_message", sentTx.ResultMessage,
				)
				return
			}
			// txBlob and sentTx.TxBlob is same

			// Debug ledger version info
			logger.Debug("ledger version",
				"earlistLedgerVersion", earlistLedgerVersion,
				"sentTx.TxJSON.LastLedgerSequence", sentTx.TxJSON.LastLedgerSequence,
			)

			// Wait for transaction validation
			var ledgerVer uint64
			ledgerVer, err = u.xrper.WaitValidation(ctx, sentTx.TxJSON.LastLedgerSequence)
			if err != nil {
				logger.Warn("fail to call xrp.WaitValidation()",
					"tx_id", txID,
					"uuid", uuid,
					"signed_tx_id", signedTxID,
					"lastLedgerSequence", sentTx.TxJSON.LastLedgerSequence,
					"ledgerVer", ledgerVer,
					"error", err,
					// Transaction has not been validated yet; try again later
				)
				return
			}

			// Get transaction info for verification
			var txInfo *dtoxrp.TxInfo
			txInfo, err = u.xrper.GetTransaction(ctx, sentTx.TxJSON.Hash, earlistLedgerVersion)
			if err != nil {
				logger.Warn("fail to call xrp.GetTransaction()",
					"tx_id", txID,
					"uuid", uuid,
					"signed_tx_id", signedTxID,
					"hash", sentTx.TxJSON.Hash,
					"earlistLedgerVersion", earlistLedgerVersion,
					"error", err,
				)
				return
			}
			// for debug (should be removed later)
			grok.Value(txInfo)

			// Update xrp_detail_tx table
			var affectedNum int64
			affectedNum, err = u.txDetailRepo.UpdateAfterTxSent(
				uuid, domainTx.TxTypeSent, signedTxID, txBlob, earlistLedgerVersion)
			if err != nil {
				// TODO: even if error occurred, tx is already sent. so db should be corrected manually
				logger.Warn(
					"fail to call txDetailRepo.UpdateAfterTxSent() but tx is already sent. "+
						"So database should be updated manually",
					"tx_id", txID,
					"uuid", uuid,
					"signed_tx_id", signedTxID,
					"tx_type", domainTx.TxTypeSent.String(),
					"tx_type_value", domainTx.TxTypeSent.Int8(),
					"error", err,
				)
				// "error":"models: unable to update all for xrp_detail_tx: Error 1406:
				// Data too long for column 'signed_tx_blob' at row 1"
				return
			}
			if affectedNum == 0 {
				logger.Info("no records to update tx_table",
					"tx_id", txID,
					"uuid", uuid,
					"signed_tx_id", signedTxID,
					"tx_type", domainTx.TxTypeSent.String(),
					"tx_type_value", domainTx.TxTypeSent.Int8(),
				)
				return
			}
		}(txHex)
	}
	wg.Wait()

	// TODO: update is_allocated in account_pubkey_table
	// Not fixed yet, Ripple may use same address because no utxo
	return watchusecase.SendTransactionOutput{
		TxID: "",
	}, nil
}
