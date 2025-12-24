package watchsrv

import (
	"database/sql"
	"fmt"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// TxCreator is TxCreator interface
type TxCreator interface {
	CreateDepositTx() (string, string, error)
	CreatePaymentTx() (string, string, error)
	CreateTransferTx(sender, receiver domainAccount.AccountType, floatAmount float64) (string, string, error)
}

// TxCreate type
type TxCreate struct {
	eth             ethgrp.EtherTxCreator
	dbConn          *sql.DB
	addrRepo        watchrepo.AddressRepositorier
	txRepo          watchrepo.TxRepositorier
	txDetailRepo    watchrepo.EthDetailTxRepositorier
	payReqRepo      watchrepo.PaymentRequestRepositorier
	txFileRepo      tx.FileRepositorier
	depositReceiver domainAccount.AccountType
	paymentSender   domainAccount.AccountType
	wtype           domainWallet.WalletType
	coinTypeCode    coin.CoinTypeCode
}

// NewTxCreate returns TxCreate object
func NewTxCreate(
	eth ethgrp.EtherTxCreator,
	dbConn *sql.DB,
	addrRepo watchrepo.AddressRepositorier,
	txRepo watchrepo.TxRepositorier,
	txDetailRepo watchrepo.EthDetailTxRepositorier,
	payReqRepo watchrepo.PaymentRequestRepositorier,
	txFileRepo tx.FileRepositorier,
	depositReceiver domainAccount.AccountType,
	paymentSender domainAccount.AccountType,
	wtype domainWallet.WalletType,
	coinTypeCode coin.CoinTypeCode,
) *TxCreate {
	return &TxCreate{
		eth:             eth,
		dbConn:          dbConn,
		addrRepo:        addrRepo,
		txRepo:          txRepo,
		txDetailRepo:    txDetailRepo,
		payReqRepo:      payReqRepo,
		txFileRepo:      txFileRepo,
		depositReceiver: depositReceiver,
		paymentSender:   paymentSender,
		wtype:           wtype,
		coinTypeCode:    coinTypeCode,
	}
}

func (t *TxCreate) updateDB(
	targetAction domainTx.ActionType,
	txDetailItems []*models.EthDetailTX,
	paymentRequestIds []int64,
) (int64, error) {
	// start transaction
	dtx, err := t.dbConn.Begin()
	if err != nil {
		return 0, fmt.Errorf("fail to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			dtx.Rollback()
		} else {
			dtx.Commit()
		}
	}()

	// Insert eth_tx
	txID, err := t.txRepo.InsertUnsignedTx(targetAction)
	if err != nil {
		return 0, fmt.Errorf("fail to call txRepo.InsertUnsignedTx(): %w", err)
	}
	// Insert to eth_detail_tx
	for idx := range txDetailItems {
		txDetailItems[idx].TXID = txID
	}
	if err = t.txDetailRepo.InsertBulk(txDetailItems); err != nil {
		return 0, fmt.Errorf("fail to call txDetailRepo.InsertBulk(): %w", err)
	}

	if targetAction == domainTx.ActionTypePayment {
		_, err = t.payReqRepo.UpdatePaymentID(txID, paymentRequestIds)
		if err != nil {
			return 0, fmt.Errorf("fail to call repo.PayReq().UpdatePaymentID(txID, paymentRequestIds): %w", err)
		}
	}
	return txID, nil
}

// generateHexFile generate file for hex txID and encoded previous addresses
func (t *TxCreate) generateHexFile(
	actionType domainTx.ActionType, senderAccount domainAccount.AccountType, txID int64, serializedTxs []string,
) (string, error) {
	// add senderAccount to first line
	serializedTxs = append([]string{senderAccount.String()}, serializedTxs...)

	// create file
	path := t.txFileRepo.CreateFilePath(actionType, domainTx.TxTypeUnsigned, txID, 0)
	generatedFileName, err := t.txFileRepo.WriteFileSlice(path, serializedTxs)
	if err != nil {
		return "", fmt.Errorf("fail to call txFileRepo.WriteFile(): %w", err)
	}

	return generatedFileName, nil
}
