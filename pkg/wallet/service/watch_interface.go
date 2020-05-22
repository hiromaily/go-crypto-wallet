package service

import "github.com/hiromaily/go-crypto-wallet/pkg/account"

// AddressImporter is AddressImporter interface (for now btc/bch only)
type AddressImporter interface {
	ImportAddress(fileName string, isRescan bool) error
}

// PaymentRequestCreator is PaymentRequestCreate interface
type PaymentRequestCreator interface {
	CreatePaymentRequest() error
}

// TxCreator is TxCreator interface (for now btc/bch only)
type TxCreator interface {
	CreateDepositTx(adjustmentFee float64) (string, string, error)
	CreatePaymentTx(adjustmentFee float64) (string, string, error)
	CreateTransferTx(sender, receiver account.AccountType, floatAmount, adjustmentFee float64) (string, string, error)
}

// TxMonitorer is TxMonitor interface
type TxMonitorer interface {
	UpdateTxStatus() error
}

// TxSender is TxSender interface
type TxSender interface {
	SendTx(filePath string) (string, error)
}
