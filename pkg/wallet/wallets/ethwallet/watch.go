package ethwallet

import (
	"database/sql"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
	ethsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/eth/watchsrv"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/watchsrv"
)

// ETHWatch watch only wallet object
type ETHWatch struct {
	ETH    ethereum.Ethereumer
	dbConn *sql.DB
	wtype  domainWallet.WalletType
	watchsrv.AddressImporter
	ethsrv.TxCreator
	service.TxSender
	service.TxMonitorer
	service.PaymentRequestCreator
}

// NewETHWatch returns ETHWatch object
func NewETHWatch(
	eth ethereum.Ethereumer,
	dbConn *sql.DB,
	addrImporter watchsrv.AddressImporter,
	txCreator ethsrv.TxCreator,
	txSender service.TxSender,
	txMonitorer service.TxMonitorer,
	paymentRequestCreator service.PaymentRequestCreator,
	walletType domainWallet.WalletType,
) *ETHWatch {
	return &ETHWatch{
		ETH:                   eth,
		dbConn:                dbConn,
		wtype:                 walletType,
		AddressImporter:       addrImporter,
		TxCreator:             txCreator,
		TxSender:              txSender,
		TxMonitorer:           txMonitorer,
		PaymentRequestCreator: paymentRequestCreator,
	}
}

// ImportAddress imports address
func (w *ETHWatch) ImportAddress(fileName string, _ bool) error {
	return w.AddressImporter.ImportAddress(fileName)
}

// CreateDepositTx creates deposit unsigned transaction
func (w *ETHWatch) CreateDepositTx(_ float64) (string, string, error) {
	return w.TxCreator.CreateDepositTx()
}

// CreatePaymentTx creates payment unsigned transaction
func (w *ETHWatch) CreatePaymentTx(_ float64) (string, string, error) {
	return w.TxCreator.CreatePaymentTx()
}

// CreateTransferTx creates transfer unsigned transaction
func (w *ETHWatch) CreateTransferTx(
	sender, receiver domainAccount.AccountType, floatAmount, _ float64,
) (string, string, error) {
	return w.TxCreator.CreateTransferTx(sender, receiver, floatAmount)
}

// UpdateTxStatus updates transaction status
func (w *ETHWatch) UpdateTxStatus() error {
	return w.TxMonitorer.UpdateTxStatus()
}

// MonitorBalance monitors balance
func (w *ETHWatch) MonitorBalance(confirmationNum uint64) error {
	return w.TxMonitorer.MonitorBalance(confirmationNum)
}

// SendTx sends signed transaction
func (w *ETHWatch) SendTx(filePath string) (string, error) {
	return w.TxSender.SendTx(filePath)
}

// CreatePaymentRequest creates payment_request dummy data for development
func (w *ETHWatch) CreatePaymentRequest() error {
	amtList := []float64{
		0.001,
		0.002,
		0.0025,
		0.0015,
		0.003,
	}
	return w.PaymentRequestCreator.CreatePaymentRequest(amtList)
}

// Done should be called before exit
func (w *ETHWatch) Done() {
	w.dbConn.Close()
	w.ETH.Close()
}

// CoinTypeCode returns domainCoin.CoinTypeCode
func (w *ETHWatch) CoinTypeCode() domainCoin.CoinTypeCode {
	return w.ETH.CoinTypeCode()
}
