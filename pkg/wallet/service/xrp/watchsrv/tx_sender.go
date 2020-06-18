package watchsrv

import (
	"database/sql"
	"strings"

	"github.com/bookerzzz/grok"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
)

// TxSend type
type TxSend struct {
	xrp          xrpgrp.Rippler
	logger       *zap.Logger
	dbConn       *sql.DB
	addrRepo     watchrepo.AddressRepositorier //not used
	txRepo       watchrepo.TxRepositorier      //not used
	txDetailRepo watchrepo.XrpDetailTxRepositorier
	txFileRepo   tx.FileRepositorier
	wtype        wallet.WalletType
}

// NewTxSend returns TxSend object
func NewTxSend(
	xrp xrpgrp.Rippler,
	logger *zap.Logger,
	dbConn *sql.DB,
	addrRepo watchrepo.AddressRepositorier,
	txRepo watchrepo.TxRepositorier,
	txDetailRepo watchrepo.XrpDetailTxRepositorier,
	txFileRepo tx.FileRepositorier,
	wtype wallet.WalletType) *TxSend {

	return &TxSend{
		xrp:          xrp,
		logger:       logger,
		dbConn:       dbConn,
		addrRepo:     addrRepo,
		txRepo:       txRepo,
		txDetailRepo: txDetailRepo,
		txFileRepo:   txFileRepo,
		wtype:        wtype,
	}
}

// SendTx send signed tx by keygen/sign walet
func (t *TxSend) SendTx(filePath string) (string, error) {

	// get tx_deposit_id from file name
	//payment_5_unsigned_1_1534466246366489473
	actionType, _, txID, _, err := t.txFileRepo.ValidateFilePath(filePath, tx.TxTypeSigned)
	if err != nil {
		return "", errors.Wrap(err, "fail to call txFileRepo.ValidateFilePath()")
	}

	t.logger.Debug("send_tx", zap.String("action_type", actionType.String()))

	// read hex from file
	data, err := t.txFileRepo.ReadFileSlice(filePath)
	if err != nil {
		return "", errors.Wrap(err, "fail to call txFileRepo.ReadFile()")
	}

	//sentTxes := make([]string, 0, len(data))
	for _, txHex := range data {
		//TODO: goroutine may be useful

		// uuid, signedTxID, txBlob
		tmp := strings.Split(txHex, ",")
		if len(tmp) != 3 {
			return "", errors.New("data format is invalid in file")
		}
		uuid := tmp[0]
		signedTxID := tmp[1]
		txBlob := tmp[2]

		// submit
		sentTx, earlistLedgerVersion, err := t.xrp.SubmitTransaction(txBlob)
		if err != nil {
			t.logger.Warn("fail to call xrp.SubmitTransaction()",
				zap.Int64("tx_id", txID),
				zap.String("uuid", uuid),
				zap.String("signed_tx_id", signedTxID),
				zap.Error(err),
				// tefMAX_LEDGER / Ledger sequence too high
				// The error message Ledger sequence too high occurs if you've waited too long to confirm a transaction in Ledger Live.
			)
			continue
		}
		if !strings.Contains(sentTx.ResultCode, "tesSUCCESS") {
			t.logger.Warn("fail to call SubmitTransaction",
				zap.Int64("tx_id", txID),
				zap.String("uuid", uuid),
				zap.String("signed_tx_id", signedTxID),
				zap.String("result_code", sentTx.ResultCode),
				zap.String("result_message", sentTx.ResultMessage),
				// tefMAX_LEDGER
				// Ledger sequence too high
			)
			continue
		}
		// validate transaction
		ledgerVer, err := t.xrp.WaitValidation(sentTx.TxJSON.LastLedgerSequence)
		if err != nil {
			t.logger.Warn("fail to call xrp.WaitValidation()",
				zap.Int64("tx_id", txID),
				zap.String("uuid", uuid),
				zap.String("signed_tx_id", signedTxID),
				zap.Uint64("lastLedgerSequence", sentTx.TxJSON.LastLedgerSequence),
				zap.Uint64("ledgerVer", ledgerVer),
				zap.Error(err),
				// Transaction has not been validated yet; try again later
			)
			continue
		}

		// get transaction info
		txInfo, err := t.xrp.GetTransaction(sentTx.TxJSON.Hash, earlistLedgerVersion)
		if err != nil {
			t.logger.Warn("fail to call xrp.GetTransaction()",
				zap.Int64("tx_id", txID),
				zap.String("uuid", uuid),
				zap.String("signed_tx_id", signedTxID),
				zap.String("hash", sentTx.TxJSON.Hash),
				zap.Uint64("earlistLedgerVersion", earlistLedgerVersion),
				zap.Error(err),
			)
			continue
		}
		// for debug (should be removed later)
		grok.Value(txInfo)

		// update eth_detail_tx
		affectedNum, err := t.txDetailRepo.UpdateAfterTxSent(uuid, tx.TxTypeSent, signedTxID, txBlob, sentTx.TxBlob)
		if err != nil {
			//TODO: even if error occurred, tx is already sent. so db should be corrected manually
			t.logger.Warn("fail to call txDetailRepo.UpdateAfterTxSent() but tx is already sent. So database should be updated manually",
				zap.Int64("tx_id", txID),
				zap.String("uuid", uuid),
				zap.String("signed_tx_id", signedTxID),
				zap.String("tx_type", tx.TxTypeSent.String()),
				zap.Int8("tx_type_value", tx.TxTypeSent.Int8()),
				zap.Error(err),
			)
			//"error":"models: unable to update all for xrp_detail_tx: Error 1406: Data too long for column 'signed_tx_blob' at row 1"
			continue
		}
		if affectedNum == 0 {
			t.logger.Info("no records to update tx_table",
				zap.Int64("tx_id", txID),
				zap.String("uuid", uuid),
				zap.String("signed_tx_id", signedTxID),
				zap.String("tx_type", tx.TxTypeSent.String()),
				zap.Int8("tx_type_value", tx.TxTypeSent.Int8()),
			)
			continue
		}
	}

	//TODO: update is_allocated in account_pubkey_table
	// Not fixed yet, Ripple may use same address because no utxo
	return "", nil
}
