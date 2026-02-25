package file

import (
	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
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

	// Hex file methods (for raw transaction hex - used by BCH)
	WriteHexFile(path, hexTx string) (string, error)

	// XRP JSON transaction file methods
	ReadXRPJSONFile(path string) (*dtoxrp.XRPTransactionFile, error)
	WriteXRPJSONFile(path string, data *dtoxrp.XRPTransactionFile) (string, error)

	// ETH JSON transaction file methods
	ReadETHJSONFile(path string) (*dtoeth.ETHTransactionFile, error)
	WriteETHJSONFile(path string, data *dtoeth.ETHTransactionFile) (string, error)
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
	// WriteXpubLine writes a single xpub CSV line to a new file and returns the file name.
	// Format: coinTypeCode,accountType,44,extendedPubKey,derivationPath
	WriteXpubLine(accountType domainAccount.AccountType, coinTypeCode, xpub, derivationPath string) (string, error)
}

// DescriptorFileWriter writes descriptor data to a file path.
type DescriptorFileWriter interface {
	WriteFile(path string, data []byte) error
}
