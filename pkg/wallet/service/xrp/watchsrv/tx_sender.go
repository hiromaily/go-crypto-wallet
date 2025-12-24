package watchsrv

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/bookerzzz/grok"

	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
)

// TxSend type
type TxSend struct {
	rippler      xrpgrp.Rippler
	dbConn       *sql.DB
	addrRepo     watchrepo.AddressRepositorier // not used
	txRepo       watchrepo.TxRepositorier      // not used
	txDetailRepo watchrepo.XrpDetailTxRepositorier
	txFileRepo   tx.FileRepositorier
	wtype        domainWallet.WalletType
}

// NewTxSend returns TxSend object
func NewTxSend(
	rippler xrpgrp.Rippler,
	dbConn *sql.DB,
	addrRepo watchrepo.AddressRepositorier,
	txRepo watchrepo.TxRepositorier,
	txDetailRepo watchrepo.XrpDetailTxRepositorier,
	txFileRepo tx.FileRepositorier,
	wtype domainWallet.WalletType,
) *TxSend {
	return &TxSend{
		rippler:      rippler,
		dbConn:       dbConn,
		addrRepo:     addrRepo,
		txRepo:       txRepo,
		txDetailRepo: txDetailRepo,
		txFileRepo:   txFileRepo,
		wtype:        wtype,
	}
}

// How to send multiple transactions
// - Question about the tefPAST_SEQ (https://www.xrpchat.com/topic/33003-question-about-the-tefpast_seq/)
// - atomical multiple transaction support?
//   (https://github.com/ripple/ripple-lib/issues/839https://github.com/ripple/ripple-lib/issues/839)
// - https://stackoverflow.com/questions/57521439/can-i-send-xrp-to-multiple-addresses
// - increment the account sequence number
// - AccountTxnID (https://xrpl.org/transaction-common-fields.html#accounttxnid)
// - Execute multiple transactions atomically
//   (https://www.xrpchat.com/topic/29175-execute-multiple-transactions-atomically/)
// - トランザクションキュー (https://xrpl.org/ja/transaction-queue.html)
// - 結果のファイナリティー (https://xrpl.org/ja/finality-of-results.html)
// - Escrow (https://xrpl.org/ja/escrow.html)

// SendTx send signed tx by keygen/sign walet
func (t *TxSend) SendTx(filePath string) (string, error) {
	// get tx_deposit_id from file name
	// payment_5_unsigned_1_1534466246366489473
	actionType, _, txID, _, err := t.txFileRepo.ValidateFilePath(filePath, domainTx.TxTypeSigned)
	if err != nil {
		return "", fmt.Errorf("fail to call txFileRepo.ValidateFilePath(): %w", err)
	}

	logger.Debug("send_tx", "action_type", actionType.String())

	// read hex from file
	data, err := t.txFileRepo.ReadFileSlice(filePath)
	if err != nil {
		return "", fmt.Errorf("fail to call txFileRepo.ReadFile(): %w", err)
	}

	var wg sync.WaitGroup

	for _, txHex := range data {
		wg.Add(1)
		go func(line string) {
			defer wg.Done()

			// uuid, signedTxID, txBlob
			tmp := strings.Split(line, ",")
			if len(tmp) != 3 {
				logger.Warn("data format is invalid in file")
				return
			}
			uuid := tmp[0]
			signedTxID := tmp[1]
			txBlob := tmp[2]

			// submit
			var sentTx *xrp.SentTx
			var earlistLedgerVersion uint64
			sentTx, earlistLedgerVersion, err = t.rippler.SubmitTransaction(context.TODO(), txBlob)
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
					//  -
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

			// debug
			logger.Debug("ledger version",
				"earlistLedgerVersion", earlistLedgerVersion, // 8123733
				"sentTx.TxJSON.LastLedgerSequence", sentTx.TxJSON.LastLedgerSequence, // 8123736
			)

			// validate transaction
			var ledgerVer uint64
			ledgerVer, err = t.rippler.WaitValidation(context.TODO(), sentTx.TxJSON.LastLedgerSequence)
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

			// get transaction info
			var txInfo *xrp.TxInfo
			txInfo, err = t.rippler.GetTransaction(context.TODO(), sentTx.TxJSON.Hash, earlistLedgerVersion)
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

			// update eth_detail_tx
			var affectedNum int64
			affectedNum, err = t.txDetailRepo.UpdateAfterTxSent(
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
	return "", nil
}
