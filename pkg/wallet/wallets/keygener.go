package wallets

import (
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainKey "github.com/hiromaily/go-crypto-wallet/pkg/domain/key"
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
