package service

import "github.com/hiromaily/go-bitcoin/pkg/account"

// AddressImporter is AddressImporter interface
type AddressImporter interface {
	ImportAddress(fileName string, isRescan bool) error
}

// PaymentRequestCreator is PaymentRequestCreate interface
type PaymentRequestCreator interface {
	CreatePaymentRequest() error
}

// TxCreator is TxCreator interface
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
