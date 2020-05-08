package wallets

import "github.com/hiromaily/go-bitcoin/pkg/account"

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
}
