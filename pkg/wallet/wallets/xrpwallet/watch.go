package xrpwallet

import (
	"database/sql"

	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	wtype "github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/watchsrv"
	xrpsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/xrp/watchsrv"
)

// XRPWatch watch only wallet object
type XRPWatch struct {
	XRP    xrpgrp.Rippler
	dbConn *sql.DB
	logger *zap.Logger
	wtype  wtype.WalletType
	watchsrv.AddressImporter
	xrpsrv.TxCreator
	service.TxSender
	service.TxMonitorer
	service.PaymentRequestCreator
}

// NewXRPWatch returns XRPWatch object
func NewXRPWatch(
	xrp xrpgrp.Rippler,
	dbConn *sql.DB,
	logger *zap.Logger,
	addrImporter watchsrv.AddressImporter,
	txCreator xrpsrv.TxCreator,
	txSender service.TxSender,
	txMonitorer service.TxMonitorer,
	paymentRequestCreator service.PaymentRequestCreator,
	wtype wtype.WalletType,
) *XRPWatch {
	return &XRPWatch{
		XRP:                   xrp,
		logger:                logger,
		dbConn:                dbConn,
		wtype:                 wtype,
		AddressImporter:       addrImporter,
		TxCreator:             txCreator,
		TxSender:              txSender,
		TxMonitorer:           txMonitorer,
		PaymentRequestCreator: paymentRequestCreator,
	}
}

// ImportAddress imports address
func (w *XRPWatch) ImportAddress(fileName string, _ bool) error {
	return w.AddressImporter.ImportAddress(fileName)
}

// CreateDepositTx creates deposit unsigned transaction
func (w *XRPWatch) CreateDepositTx(_ float64) (string, string, error) {
	return w.TxCreator.CreateDepositTx()
}

// CreatePaymentTx creates payment unsigned transaction
func (w *XRPWatch) CreatePaymentTx(_ float64) (string, string, error) {
	return w.TxCreator.CreatePaymentTx()
}

// CreateTransferTx creates transfer unsigned transaction
func (w *XRPWatch) CreateTransferTx(sender, receiver account.AccountType, floatAmount, _ float64) (string, string, error) {
	return w.TxCreator.CreateTransferTx(sender, receiver, floatAmount)
}

// UpdateTxStatus updates transaction status
func (w *XRPWatch) UpdateTxStatus() error {
	w.logger.Info("no functionality for XRP")
	return nil
}

// MonitorBalance monitors balance
func (w *XRPWatch) MonitorBalance(confirmationNum uint64) error {
	return w.TxMonitorer.MonitorBalance(confirmationNum)
}

// SendTx sends signed transaction
func (w *XRPWatch) SendTx(filePath string) (string, error) {
	return w.TxSender.SendTx(filePath)
}

// CreatePaymentRequest creates payment_request dummy data for development
func (w *XRPWatch) CreatePaymentRequest() error {
	amtList := []float64{
		50,
		100,
		120,
		130,
		150,
	}
	return w.PaymentRequestCreator.CreatePaymentRequest(amtList)
}

// Done should be called before exit
func (w *XRPWatch) Done() {
	w.dbConn.Close()
	w.XRP.Close()
}

// CoinTypeCode returns coin.CoinTypeCode
func (w *XRPWatch) CoinTypeCode() coin.CoinTypeCode {
	return w.XRP.CoinTypeCode()
}
