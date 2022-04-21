package btcwallet

import (
	"database/sql"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	wtype "github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
)

// BTCWatch watch only wallet object
type BTCWatch struct {
	BTC      btcgrp.Bitcoiner
	dbConn   *sql.DB
	logger   *zap.Logger
	tracer   opentracing.Tracer
	addrType address.AddrType
	wtype    wtype.WalletType
	service.AddressImporter
	service.TxCreator
	service.TxSender
	service.TxMonitorer
	service.PaymentRequestCreator
}

// NewBTCWatch returns Watch object
func NewBTCWatch(
	btc btcgrp.Bitcoiner,
	dbConn *sql.DB,
	logger *zap.Logger,
	tracer opentracing.Tracer,
	addrType address.AddrType,
	addrImporter service.AddressImporter,
	txCreator service.TxCreator,
	txSender service.TxSender,
	txMonitorer service.TxMonitorer,
	paymentRequestCreator service.PaymentRequestCreator,
	wtype wtype.WalletType,
) *BTCWatch {
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

// MonitorBalance monitors balance
func (w *BTCWatch) MonitorBalance(confirmationNum uint64) error {
	return w.TxMonitorer.MonitorBalance(confirmationNum)
}

// SendTx sends signed transaction
func (w *BTCWatch) SendTx(filePath string) (string, error) {
	return w.TxSender.SendTx(filePath)
}

// CreatePaymentRequest creates payment_request dummy data for development
func (w *BTCWatch) CreatePaymentRequest() error {
	amtList := []float64{
		0.00001,
		0.00002,
		0.000025,
		0.000015,
		0.00003,
	}

	return w.PaymentRequestCreator.CreatePaymentRequest(amtList)
}

// Done should be called before exit
func (w *BTCWatch) Done() {
	w.dbConn.Close()
	w.BTC.Close()
}

// CoinTypeCode returns coin.CoinTypeCode
func (w *BTCWatch) CoinTypeCode() coin.CoinTypeCode {
	return w.BTC.CoinTypeCode()
}

// GetBTC gets btc
//func (w *BTCWatch) GetBTC() btcgrp.Bitcoiner {
//	return w.BTC
//}
