package wallets

import (
	"database/sql"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service/coldsrv"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service/coldsrv/keygensrv"
)

// Keygener is for keygen wallet service interface
type Keygener interface {
	GenerateSeed() ([]byte, error)
	StoreSeed(strSeed string) ([]byte, error)
	GenerateAccountKey(accountType account.AccountType, seed []byte, count uint32) ([]key.WalletKey, error)
	ImportPrivKey(accountType account.AccountType) error
	ImportFullPubKey(fileName string) error
	CreateMultisigAddress(accountType account.AccountType) error
	ExportAddress(accountType account.AccountType) (string, error)
	SignTx(filePath string) (string, bool, string, error)

	Done()
	GetBTC() api.Bitcoiner
}

// Keygen is keygen wallet object
type Keygen struct {
	btc      api.Bitcoiner
	dbConn   *sql.DB
	addrType address.AddrType
	wtype    wallet.WalletType
	coldsrv.Seeder
	coldsrv.HDWalleter
	keygensrv.PrivKeyer
	keygensrv.FullPubKeyImporter
	keygensrv.Multisiger
	keygensrv.AddressExporter
	coldsrv.Signer
}

// NewKeygen returns Keygen object
func NewKeygen(
	btc api.Bitcoiner,
	dbConn *sql.DB,
	addrType address.AddrType,
	seeder coldsrv.Seeder,
	hdWallter coldsrv.HDWalleter,
	privKeyer keygensrv.PrivKeyer,
	pubkeyImporter keygensrv.FullPubKeyImporter,
	multisiger keygensrv.Multisiger,
	addressExporter keygensrv.AddressExporter,
	signer coldsrv.Signer,
	wtype wallet.WalletType) *Keygen {

	return &Keygen{
		btc:                btc,
		dbConn:             dbConn,
		addrType:           addrType,
		wtype:              wtype,
		Seeder:             seeder,
		HDWalleter:         hdWallter,
		PrivKeyer:          privKeyer,
		FullPubKeyImporter: pubkeyImporter,
		Multisiger:         multisiger,
		AddressExporter:    addressExporter,
		Signer:             signer,
	}
}

// GenerateSeed generates seed
func (k *Keygen) GenerateSeed() ([]byte, error) {
	return k.Seeder.Generate()
}

// StoreSeed stores seed
func (k *Keygen) StoreSeed(strSeed string) ([]byte, error) {
	return k.Seeder.Store(strSeed)
}

// GenerateAccountKey generates account keys
func (k *Keygen) GenerateAccountKey(accountType account.AccountType, seed []byte, count uint32) ([]key.WalletKey, error) {
	return k.HDWalleter.Generate(accountType, seed, count)
}

// ImportPrivKey imports privKey
func (k *Keygen) ImportPrivKey(accountType account.AccountType) error {
	return k.PrivKeyer.Import(accountType)
}

// ImportFullPubKey imports full-pubkey
func (k *Keygen) ImportFullPubKey(fileName string) error {
	return k.FullPubKeyImporter.ImportFullPubKey(fileName)
}

// CreateMultisigAddress creates multi sig address returns Multisiger interface
func (k *Keygen) CreateMultisigAddress(accountType account.AccountType) error {
	return k.Multisiger.AddMultisigAddress(accountType, k.addrType)
}

// ExportAddress exports address
func (k *Keygen) ExportAddress(accountType account.AccountType) (string, error) {
	return k.AddressExporter.ExportAddress(accountType)
}

// SignTx signs on transaction
func (k *Keygen) SignTx(filePath string) (string, bool, string, error) {
	return k.Signer.SignTx(filePath)
}

// Done should be called before exit
func (k *Keygen) Done() {
	k.dbConn.Close()
	k.btc.Close()
}

// GetBTC gets btc
func (k *Keygen) GetBTC() api.Bitcoiner {
	return k.btc
}
