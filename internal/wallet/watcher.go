package wallets

import (
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// Watcher is for watch only wallet service interface
type Watcher interface {
	ImportAddress(fileName string, isRescan bool) error
	CreateDepositTx(adjustmentFee float64) (string, string, error)
	CreatePaymentTx(adjustmentFee float64) (string, string, error)
	CreateTransferTx(
		sender, receiver domainAccount.AccountType, floatAmount, adjustmentFee float64,
	) (string, string, error)
	SendTx(filePath string) (string, error)
	UpdateTxStatus() error
	MonitorBalance(confirmationNum uint64) error
	CreatePaymentRequest() error
	Done()
	CoinTypeCode() domainCoin.CoinTypeCode
}
