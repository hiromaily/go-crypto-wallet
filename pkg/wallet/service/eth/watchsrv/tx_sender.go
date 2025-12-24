package watchsrv

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// TxSend type
type TxSend struct {
	eth          ethereum.Ethereumer
	dbConn       *sql.DB
	addrRepo     watch.AddressRepositorier // not used
	txRepo       watch.TxRepositorier      // not used
	txDetailRepo watch.EthDetailTxRepositorier
	txFileRepo   file.TransactionFileRepositorier
	wtype        domainWallet.WalletType
}

// NewTxSend returns TxSend object
func NewTxSend(
	eth ethereum.Ethereumer,
	dbConn *sql.DB,
	addrRepo watch.AddressRepositorier,
	txRepo watch.TxRepositorier,
	txDetailRepo watch.EthDetailTxRepositorier,
	txFileRepo file.TransactionFileRepositorier,
	wtype domainWallet.WalletType,
) *TxSend {
	return &TxSend{
		eth:          eth,
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
		sentTx, err = t.eth.SendSignedRawTransaction(context.TODO(), signedTx)
		if err != nil {
			logger.Warn("fail to call eth.SendSignedRawTransaction()",
				"error", err,
			)
			continue
		}
		if sentTx == "" {
			logger.Warn("no sentTx by calling eth.SendSignedRawTransaction()",
				"error", err,
			)
			continue
		}

		// update eth_detail_tx
		var affectedNum int64
		affectedNum, err = t.txDetailRepo.UpdateAfterTxSent(uuid, domainTx.TxTypeSent, signedTx, sentTx)
		if err != nil {
			// TODO: even if error occurred, tx is already sent. so db should be corrected manually
			logger.Warn(
				"fail to call repo.Tx().UpdateAfterTxSent() but tx is already sent. "+
					"So database should be updated manually",
				"tx_id", txID,
				"tx_type", domainTx.TxTypeSent.String(),
				"tx_type_value", domainTx.TxTypeSent.Int8(),
				"signed_hex_tx", signedTx,
				"sent_hash_tx", sentTx,
			)
			continue
		}
		if affectedNum == 0 {
			logger.Info("no records to update tx_table",
				"tx_id", txID,
				"tx_type", domainTx.TxTypeSent.String(),
				"tx_type_value", domainTx.TxTypeSent.Int8(),
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
