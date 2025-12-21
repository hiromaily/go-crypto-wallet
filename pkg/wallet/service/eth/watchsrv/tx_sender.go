package watchsrv

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
)

// TxSend type
type TxSend struct {
	eth          ethgrp.Ethereumer
	logger       logger.Logger
	dbConn       *sql.DB
	addrRepo     watchrepo.AddressRepositorier // not used
	txRepo       watchrepo.TxRepositorier      // not used
	txDetailRepo watchrepo.EthDetailTxRepositorier
	txFileRepo   tx.FileRepositorier
	wtype        wallet.WalletType
}

// NewTxSend returns TxSend object
func NewTxSend(
	eth ethgrp.Ethereumer,
	logger logger.Logger,
	dbConn *sql.DB,
	addrRepo watchrepo.AddressRepositorier,
	txRepo watchrepo.TxRepositorier,
	txDetailRepo watchrepo.EthDetailTxRepositorier,
	txFileRepo tx.FileRepositorier,
	wtype wallet.WalletType,
) *TxSend {
	return &TxSend{
		eth:          eth,
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
	// payment_5_unsigned_1_1534466246366489473
	actionType, _, txID, _, err := t.txFileRepo.ValidateFilePath(filePath, tx.TxTypeSigned)
	if err != nil {
		return "", errors.Wrap(err, "fail to call txFileRepo.ValidateFilePath()")
	}

	t.logger.Debug("send_tx", "action_type", actionType.String())

	// read hex from file
	data, err := t.txFileRepo.ReadFileSlice(filePath)
	if err != nil {
		return "", errors.Wrap(err, "fail to call txFileRepo.ReadFile()")
	}

	// sentTxes := make([]string, 0, len(data))
	for _, txHex := range data {
		// data is csv [rawTx.TxHex, signedRawTx.TxHex]
		// rawTx.TxHex is used to record status by updating database
		tmp := strings.Split(txHex, ",")
		if len(tmp) != 2 {
			return "", errors.New("data format is invalid in file")
		}
		uuid := tmp[0]
		signedTx := tmp[1]

		// sign
		var sentTx string
		sentTx, err = t.eth.SendSignedRawTransaction(signedTx)
		if err != nil {
			t.logger.Warn("fail to call eth.SendSignedRawTransaction()",
				"error", err,
			)
			continue
		}
		if sentTx == "" {
			t.logger.Warn("no sentTx by calling eth.SendSignedRawTransaction()",
				"error", err,
			)
			continue
		}

		// update eth_detail_tx
		var affectedNum int64
		affectedNum, err = t.txDetailRepo.UpdateAfterTxSent(uuid, tx.TxTypeSent, signedTx, sentTx)
		if err != nil {
			// TODO: even if error occurred, tx is already sent. so db should be corrected manually
			t.logger.Warn(
				"fail to call repo.Tx().UpdateAfterTxSent() but tx is already sent. "+
					"So database should be updated manually",
				"tx_id", txID,
				"tx_type", tx.TxTypeSent.String(),
				"tx_type_value", tx.TxTypeSent.Int8(),
				"signed_hex_tx", signedTx,
				"sent_hash_tx", sentTx,
			)
			continue
		}
		if affectedNum == 0 {
			t.logger.Info("no records to update tx_table",
				"tx_id", txID,
				"tx_type", tx.TxTypeSent.String(),
				"tx_type_value", tx.TxTypeSent.Int8(),
				"signed_hex_tx", signedTx,
				"sent_hash_tx", sentTx,
			)
			continue
		}
	}

	// TODO: update is_allocated in account_pubkey_table
	// Ethereum should use same address because no utxo
	return "", nil
}
