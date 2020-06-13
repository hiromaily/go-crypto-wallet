package xrpwallet

import (
	"database/sql"

	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	wtype "github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/common/watchsrv"
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
}

// NewXRPWatch returns XRPWatch object
func NewXRPWatch(
	xrp xrpgrp.Rippler,
	dbConn *sql.DB,
	logger *zap.Logger,
	addrImporter watchsrv.AddressImporter,
	txCreator xrpsrv.TxCreator,
	wtype wtype.WalletType) *XRPWatch {

	return &XRPWatch{
		XRP:             xrp,
		logger:          logger,
		dbConn:          dbConn,
		wtype:           wtype,
		AddressImporter: addrImporter,
		TxCreator:       txCreator,
	}
}

// ImportAddress imports address
func (w *XRPWatch) ImportAddress(fileName string, isRescan bool) error {
	return w.AddressImporter.ImportAddress(fileName)
}

// CreateDepositTx creates deposit unsigned transaction
func (w *XRPWatch) CreateDepositTx(adjustmentFee float64) (string, string, error) {
	return w.TxCreator.CreateDepositTx()
}

// CreatePaymentTx creates payment unsigned transaction
func (w *XRPWatch) CreatePaymentTx(adjustmentFee float64) (string, string, error) {
	w.logger.Info("not implemented yet")
	//return w.TxCreator.CreatePaymentTx()
	return "", "", nil
}

// CreateTransferTx creates transfer unsigned transaction
func (w *XRPWatch) CreateTransferTx(sender, receiver account.AccountType, floatAmount, adjustmentFee float64) (string, string, error) {
	w.logger.Info("not implemented yet")
	//return w.TxCreator.CreateTransferTx(sender, receiver, floatAmount)
	return "", "", nil
}

// UpdateTxStatus updates transaction status
func (w *XRPWatch) UpdateTxStatus() error {
	w.logger.Info("not implemented yet")
	//return w.TxMonitorer.UpdateTxStatus()
	return nil
}

// MonitorBalance monitors balance
func (w *XRPWatch) MonitorBalance(confirmationNum uint64) error {
	w.logger.Info("not implemented yet")
	//return w.TxMonitorer.MonitorBalance()
	return nil
}

// SendTx sends signed transaction
func (w *XRPWatch) SendTx(filePath string) (string, error) {
	w.logger.Info("not implemented yet")
	//return w.TxSender.SendTx(filePath)
	return "", nil
}

// CreatePaymentRequest creates payment_request dummy data for development
func (w *XRPWatch) CreatePaymentRequest() error {
	w.logger.Info("not implemented yet")
	//return w.PaymentRequestCreator.CreatePaymentRequest()
	return nil
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
