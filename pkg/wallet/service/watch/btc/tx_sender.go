package btc

import (
	"database/sql"
	"errors"
	"fmt"

	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// TxSend type
type TxSend struct {
	btc          bitcoin.Bitcoiner
	dbConn       *sql.DB
	addrRepo     watch.AddressRepositorier
	txRepo       watch.BTCTxRepositorier
	txOutputRepo watch.TxOutputRepositorier
	txFileRepo   file.TransactionFileRepositorier
	wtype        domainWallet.WalletType
}

// NewTxSend returns TxSend object
func NewTxSend(
	btc bitcoin.Bitcoiner,
	dbConn *sql.DB,
	addrRepo watch.AddressRepositorier,
	txRepo watch.BTCTxRepositorier,
	txOutputRepo watch.TxOutputRepositorier,
	txFileRepo file.TransactionFileRepositorier,
	wtype domainWallet.WalletType,
) *TxSend {
	return &TxSend{
		btc:          btc,
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
	actionType, _, txID, _, err := t.txFileRepo.ValidateFilePath(filePath, domainTx.TxTypeSigned)
	if err != nil {
		return "", fmt.Errorf("fail to call txFileRepo.ValidateFilePath(): %w", err)
	}

	logger.Debug("send_tx", "action_type", actionType.String())

	// read hex from file
	signedHex, err := t.txFileRepo.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("fail to call txFileRepo.ReadFile(): %w", err)
	}

	// send signed tx
	hash, err := t.btc.SendTransactionByHex(signedHex)
	if err != nil {
		// if signature is not completed
		//-26: 16: mandatory-script-verify-flag-failed (Operation not valid with the current stack size)
		return "", fmt.Errorf("fail to call btc.SendTransactionByHex(): %w", err)
	}

	if hash == nil {
		// tx is already sent
		return "", nil
	}

	// update tx_table
	affectedNum, err := t.txRepo.UpdateAfterTxSent(txID, domainTx.TxTypeSent, signedHex, hash.String())
	if err != nil {
		// TODO: even if error occurred, tx is already sent. so db should be corrected manually
		logger.Warn(
			"fail to call repo.Tx().UpdateAfterTxSent() but tx is already sent. "+
				"So database should be updated manually",
			"tx_id", txID,
			"tx_type", domainTx.TxTypeSent.String(),
			"tx_type_vakue", domainTx.TxTypeSent.Int8(),
			"signed_hex_tx", signedHex,
			"sent_hash_tx", hash.String(),
		)
		return "", fmt.Errorf("fail to call updateHexForSentTx(), but tx is sent. txID: %d: %w", txID, err)
	}
	if affectedNum == 0 {
		logger.Info("no records to update tx_table",
			"tx_id", txID,
			"tx_type", domainTx.TxTypeSent.String(),
			"tx_type_vakue", domainTx.TxTypeSent.Int8(),
			"signed_hex_tx", signedHex,
			"sent_hash_tx", hash.String(),
		)
		return "", nil
	}

	// update account_pubkey_table
	if actionType != domainTx.ActionTypePayment {
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
		return fmt.Errorf("fail to call repo.TxOutput().GetAllByTxID(): %w", err)
	}
	if len(txOutputs) == 0 {
		return errors.New("output tx could not be found in tx_deposit_output")
	}

	_, err = t.addrRepo.UpdateIsAllocated(true, txOutputs[0].OutputAddress)
	if err != nil {
		return fmt.Errorf("fail to call repo.Pubkey().UpdateIsAllocated(): %w", err)
	}

	return nil
}
