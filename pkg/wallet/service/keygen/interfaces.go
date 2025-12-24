package keygen

import (
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainKey "github.com/hiromaily/go-crypto-wallet/pkg/domain/key"
)

// HDWalleter is HD wallet key generation service
type HDWalleter interface {
	Generate(accountType domainAccount.AccountType, seed []byte, count uint32) ([]domainKey.WalletKey, error)
}

// Seeder is Seeder service
type Seeder interface {
	Generate() ([]byte, error)
	Store(strSeed string) ([]byte, error)
}

// AddressExporter is AddressExporter service
type AddressExporter interface {
	ExportAddress(accountType domainAccount.AccountType) (string, error)
}

// PrivKeyer is PrivKeyer service
type PrivKeyer interface {
	Import(accountType domainAccount.AccountType) error
}

// FullPubKeyImporter is FullPubkeyImport service
type FullPubKeyImporter interface {
	ImportFullPubKey(fileName string) error
}

// Multisiger is Multisiger service
type Multisiger interface {
	AddMultisigAddress(accountType domainAccount.AccountType, addressType address.AddrType) error
}
