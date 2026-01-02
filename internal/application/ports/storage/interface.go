package storage

import (
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

// TransactionFileRepositorier is file storager for tx info
type TransactionFileRepositorier interface {
	CreateFilePath(actionType domainTx.ActionType, txType domainTx.TxType, txID int64, signedCount int) string
	GetFileNameType(filePath string) (*FileName, error)
	ValidateFilePath(
		filePath string,
		expectedTxType domainTx.TxType,
	) (domainTx.ActionType, domainTx.TxType, int64, int, error)
	ReadFile(path string) (string, error)
	ReadFileSlice(path string) ([]string, error)
	WriteFile(path, hexTx string) (string, error)
	WriteFileSlice(path string, data []string) (string, error)

	// PSBT-specific methods (BIP174)
	ReadPSBTFile(path string) (string, error)
	WritePSBTFile(path, psbtBase64 string) (string, error)
}

// FileName is object for items in fine name
type FileName struct {
	ActionType  domainTx.ActionType
	TxType      domainTx.TxType
	TxID        int64
	SignedCount int
}

// AddressFileRepositorier is address file storage interface
type AddressFileRepositorier interface {
	CreateFilePath(accountType domainAccount.AccountType) string
	ValidateFilePath(fileName string, accountType domainAccount.AccountType) error
	ImportAddress(fileName string) ([]string, error)
}
