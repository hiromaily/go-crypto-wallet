package wallet

import (
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

// Keygener is for keygen wallet service interface
type Keygener interface {
	GenerateSeed() ([]byte, error)
	StoreSeed(strSeed string) ([]byte, error)
	GenerateAccountKey(
		accountType domainAccount.AccountType, seed []byte, count uint32, isKeyPair bool,
	) ([]domainKey.WalletKey, error)
	ImportPrivKey(accountType domainAccount.AccountType) error
	ImportFullPubKey(fileName string) error
	CreateMultisigAddress(accountType domainAccount.AccountType) error
	ExportAddress(accountType domainAccount.AccountType) (string, error)
	SignTx(filePath string) (string, bool, string, error)
	Done()
}

// Signer is for signature wallet service interface
type Signer interface {
	GenerateSeed() ([]byte, error)
	StoreSeed(strSeed string) ([]byte, error)
	GenerateAuthKey(seed []byte, count uint32) ([]domainKey.WalletKey, error)
	ImportPrivKey() error
	ExportFullPubkey() (string, error)
	SignTx(filePath string) (string, bool, string, error)
	Done()
}

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
