package ethwallet

import (
	"database/sql"

	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	wtype "github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service/eth/watchsrv"
)

// ETHWatch watch only wallet object
type ETHWatch struct {
	ETH    ethgrp.Ethereumer
	dbConn *sql.DB
	logger *zap.Logger
	wtype  wtype.WalletType
	watchsrv.AddressImporter
	service.TxCreator
	service.TxSender
	service.PaymentRequestCreator
}

// NewETHWatch returns ETHWatch object
func NewETHWatch(
	eth ethgrp.Ethereumer,
	dbConn *sql.DB,
	logger *zap.Logger,
	addrImporter watchsrv.AddressImporter,
	txCreator service.TxCreator,
	txSender service.TxSender,
	paymentRequestCreator service.PaymentRequestCreator,
	wtype wtype.WalletType) *ETHWatch {

	return &ETHWatch{
		ETH:                   eth,
		logger:                logger,
		dbConn:                dbConn,
		wtype:                 wtype,
		AddressImporter:       addrImporter,
		TxCreator:             txCreator,
		TxSender:              txSender,
		PaymentRequestCreator: paymentRequestCreator,
	}
}

// ImportAddress imports address
func (w *ETHWatch) ImportAddress(fileName string, isRescan bool) error {
	return w.AddressImporter.ImportAddress(fileName)
}

// CreateDepositTx creates deposit unsigned transaction
func (w *ETHWatch) CreateDepositTx(adjustmentFee float64) (string, string, error) {
	return w.TxCreator.CreateDepositTx(adjustmentFee)
}

// CreatePaymentTx creates payment unsigned transaction
func (w *ETHWatch) CreatePaymentTx(adjustmentFee float64) (string, string, error) {
	return w.TxCreator.CreatePaymentTx(adjustmentFee)
}

// CreateTransferTx creates transfer unsigned transaction
func (w *ETHWatch) CreateTransferTx(sender, receiver account.AccountType, floatAmount, adjustmentFee float64) (string, string, error) {
	return w.TxCreator.CreateTransferTx(sender, receiver, floatAmount, adjustmentFee)
}

// UpdateTxStatus updates transaction status
func (w *ETHWatch) UpdateTxStatus() error {
	//return w.TxMonitorer.UpdateTxStatus()
	w.logger.Warn("not implemented yet in ETH")
	return nil
}

// SendTx sends signed transaction
func (w *ETHWatch) SendTx(filePath string) (string, error) {
	return w.TxSender.SendTx(filePath)
}

// CreatePaymentRequest creates payment_request dummy data for development
func (w *ETHWatch) CreatePaymentRequest() error {
	return w.PaymentRequestCreator.CreatePaymentRequest()
}

// Done should be called before exit
func (w *ETHWatch) Done() {
	w.dbConn.Close()
	w.ETH.Close()
}
