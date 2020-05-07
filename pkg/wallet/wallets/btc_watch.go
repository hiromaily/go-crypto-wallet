package wallets

import (
	"database/sql"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	wtype "github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/btcgrp"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service/watchsrv"
)

// Watcher is for watch only wallet service interface
type Watcher interface {
	ImportAddress(fileName string, isRescan bool) error
	CreateDepositTx(adjustmentFee float64) (string, string, error)
	CreatePaymentTx(adjustmentFee float64) (string, string, error)
	CreateTransferTx(sender, receiver account.AccountType, floatAmount, adjustmentFee float64) (string, string, error)
	SendTx(filePath string) (string, error)
	UpdateTxStatus() error
	CreatePaymentRequest() error

	Done()
	//GetBTC() btcgrp.Bitcoiner
}

// BTCWatch watch only wallet object
type BTCWatch struct {
	BTC      btcgrp.Bitcoiner
	dbConn   *sql.DB
	logger   *zap.Logger
	tracer   opentracing.Tracer
	addrType address.AddrType
	wtype    wtype.WalletType
	watchsrv.AddressImporter
	watchsrv.TxCreator
	watchsrv.TxSender
	watchsrv.TxMonitorer
	watchsrv.PaymentRequestCreator
}

// NewBTCWatch returns Watch object
func NewBTCWatch(
	btc btcgrp.Bitcoiner,
	dbConn *sql.DB,
	logger *zap.Logger,
	tracer opentracing.Tracer,
	addrType address.AddrType,
	addrImporter watchsrv.AddressImporter,
	txCreator watchsrv.TxCreator,
	txSender watchsrv.TxSender,
	txMonitorer watchsrv.TxMonitorer,
	paymentRequestCreator watchsrv.PaymentRequestCreator,
	wtype wtype.WalletType) *BTCWatch {

	return &BTCWatch{
		BTC:                   btc,
		logger:                logger,
		dbConn:                dbConn,
		tracer:                tracer,
		addrType:              addrType,
		wtype:                 wtype,
		AddressImporter:       addrImporter,
		TxCreator:             txCreator,
		TxSender:              txSender,
		TxMonitorer:           txMonitorer,
		PaymentRequestCreator: paymentRequestCreator,
	}
}

// ImportAddress imports address
func (w *BTCWatch) ImportAddress(fileName string, isRescan bool) error {
	return w.AddressImporter.ImportAddress(fileName, isRescan)
}

// CreateDepositTx creates deposit unsigned transaction
func (w *BTCWatch) CreateDepositTx(adjustmentFee float64) (string, string, error) {
	return w.TxCreator.CreateDepositTx(adjustmentFee)
}

// CreatePaymentTx creates payment unsigned transaction
func (w *BTCWatch) CreatePaymentTx(adjustmentFee float64) (string, string, error) {
	return w.TxCreator.CreatePaymentTx(adjustmentFee)
}

// CreateTransferTx creates transfer unsigned transaction
func (w *BTCWatch) CreateTransferTx(sender, receiver account.AccountType, floatAmount, adjustmentFee float64) (string, string, error) {
	return w.TxCreator.CreateTransferTx(sender, receiver, floatAmount, adjustmentFee)
}

// UpdateTxStatus updates transaction status
func (w *BTCWatch) UpdateTxStatus() error {
	return w.TxMonitorer.UpdateTxStatus()
}

// SendTx sends signed transaction
func (w *BTCWatch) SendTx(filePath string) (string, error) {
	return w.TxSender.SendTx(filePath)
}

// CreatePaymentRequest creates payment_request dummy data for development
func (w *BTCWatch) CreatePaymentRequest() error {
	return w.PaymentRequestCreator.CreatePaymentRequest()
}

// Done should be called before exit
func (w *BTCWatch) Done() {
	w.dbConn.Close()
	w.BTC.Close()
}

// GetBTC gets btc
//func (w *BTCWatch) GetBTC() btcgrp.Bitcoiner {
//	return w.BTC
//}
