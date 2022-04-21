package watchsrv

import (
	"database/sql"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/action"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
)

// TxSend type
type TxSend struct {
	btc          btcgrp.Bitcoiner
	logger       *zap.Logger
	dbConn       *sql.DB
	addrRepo     watchrepo.AddressRepositorier
	txRepo       watchrepo.BTCTxRepositorier
	txOutputRepo watchrepo.TxOutputRepositorier
	txFileRepo   tx.FileRepositorier
	wtype        wallet.WalletType
}

// NewTxSend returns TxSend object
func NewTxSend(
	btc btcgrp.Bitcoiner,
	logger *zap.Logger,
	dbConn *sql.DB,
	addrRepo watchrepo.AddressRepositorier,
	txRepo watchrepo.BTCTxRepositorier,
	txOutputRepo watchrepo.TxOutputRepositorier,
	txFileRepo tx.FileRepositorier,
	wtype wallet.WalletType,
) *TxSend {
	return &TxSend{
		btc:          btc,
		logger:       logger,
		dbConn:       dbConn,
		addrRepo:     addrRepo,
		txRepo:       txRepo,
		txOutputRepo: txOutputRepo,
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

	t.logger.Debug("send_tx", zap.String("action_type", actionType.String()))

	// read hex from file
	signedHex, err := t.txFileRepo.ReadFile(filePath)
	if err != nil {
		return "", errors.Wrap(err, "fail to call txFileRepo.ReadFile()")
	}

	// send signed tx
	hash, err := t.btc.SendTransactionByHex(signedHex)
	if err != nil {
		// if signature is not completed
		//-26: 16: mandatory-script-verify-flag-failed (Operation not valid with the current stack size)
		return "", errors.Wrap(err, "fail to call btc.SendTransactionByHex()")
	}

	if hash == nil {
		// tx is already sent
		return "", nil
	}

	// update tx_table
	affectedNum, err := t.txRepo.UpdateAfterTxSent(txID, tx.TxTypeSent, signedHex, hash.String())
	if err != nil {
		// TODO: even if error occurred, tx is already sent. so db should be corrected manually
		t.logger.Warn("fail to call repo.Tx().UpdateAfterTxSent() but tx is already sent. So database should be updated manually",
			zap.Int64("tx_id", txID),
			zap.String("tx_type", tx.TxTypeSent.String()),
			zap.Int8("tx_type_vakue", tx.TxTypeSent.Int8()),
			zap.String("signed_hex_tx", signedHex),
			zap.String("sent_hash_tx", hash.String()),
		)
		return "", errors.Wrapf(err, "fail to call updateHexForSentTx(), but tx is sent. txID: %d", txID)
	}
	if affectedNum == 0 {
		t.logger.Info("no records to update tx_table",
			zap.Int64("tx_id", txID),
			zap.String("tx_type", tx.TxTypeSent.String()),
			zap.Int8("tx_type_vakue", tx.TxTypeSent.Int8()),
			zap.String("signed_hex_tx", signedHex),
			zap.String("sent_hash_tx", hash.String()),
		)
		return "", nil
	}

	// update account_pubkey_table
	if actionType != action.ActionTypePayment {
		// skip for that receiver address is anonymous
		err = t.updateIsAllocatedAccountPubkey(txID)
		if err != nil {
			// TODO: even if error occurred, tx is already sent. so db should be corrected manually
			return "", err
		}
	}
	return hash.String(), nil
}

func (t *TxSend) updateIsAllocatedAccountPubkey(txID int64) error {
	// get txOutputs by tx_id
	txOutputs, err := t.txOutputRepo.GetAllByTxID(txID)
	if err != nil {
		return errors.Wrap(err, "fail to call repo.TxOutput().GetAllByTxID()")
	}
	if len(txOutputs) == 0 {
		return errors.New("output tx could not be found in tx_deposit_output")
	}

	_, err = t.addrRepo.UpdateIsAllocated(true, txOutputs[0].OutputAddress)
	if err != nil {
		return errors.Wrap(err, "fail to call repo.Pubkey().UpdateIsAllocated()")
	}

	return nil
}
