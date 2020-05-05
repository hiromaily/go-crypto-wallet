package wallets

import (
	"database/sql"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service/coldsrv"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service/coldsrv/signsrv"
)

// Signer is for signature wallet service interface
type Signer interface {
	GenerateSeed() ([]byte, error)
	StoreSeed(strSeed string) ([]byte, error)
	GenerateAuthKey(seed []byte, count uint32) ([]key.WalletKey, error)
	ImportPrivKey() error
	ExportFullPubkey() (string, error)
	SignTx(filePath string) (string, bool, string, error)

	Done()
	GetBTC() api.Bitcoiner
}

// Sign is sign wallet object
type Sign struct {
	btc         api.Bitcoiner
	dbConn      *sql.DB
	authAccount account.AuthType
	addrType    address.AddrType
	wtype       wallet.WalletType
	coldsrv.Seeder
	coldsrv.HDWalleter
	signsrv.PrivKeyer
	signsrv.FullPubkeyExporter
	coldsrv.Signer
}

// NewSign returns Sign object
func NewSign(
	btc api.Bitcoiner,
	dbConn *sql.DB,
	authAccount account.AuthType,
	addrType address.AddrType,
	seeder coldsrv.Seeder,
	hdWallter coldsrv.HDWalleter,
	privKeyer signsrv.PrivKeyer,
	fullPubkeyExporter signsrv.FullPubkeyExporter,
	signer coldsrv.Signer,
	wtype wallet.WalletType) *Sign {

	return &Sign{
		btc:                btc,
		dbConn:             dbConn,
		authAccount:        authAccount,
		addrType:           addrType,
		wtype:              wtype,
		Seeder:             seeder,
		HDWalleter:         hdWallter,
		PrivKeyer:          privKeyer,
		FullPubkeyExporter: fullPubkeyExporter,
		Signer:             signer,
	}
}

// GenerateSeed generates seed
func (s *Sign) GenerateSeed() ([]byte, error) {
	return s.Seeder.Generate()
}

// StoreSeed stores seed
func (s *Sign) StoreSeed(strSeed string) ([]byte, error) {
	return s.Seeder.Store(strSeed)
}

// GenerateAuthKey generates account keys
func (s *Sign) GenerateAuthKey(seed []byte, count uint32) ([]key.WalletKey, error) {
	return s.HDWalleter.Generate(s.authAccount.AccountType(), seed, count)
}

// ImportPrivKey imports privKey
func (s *Sign) ImportPrivKey() error {
	return s.PrivKeyer.Import()
}

// ExportFullPubkey exports full-pubkey
func (s *Sign) ExportFullPubkey() (string, error) {
	return s.FullPubkeyExporter.ExportFullPubkey()
}

// SignTx signs on transaction
func (s *Sign) SignTx(filePath string) (string, bool, string, error) {
	return s.Signer.SignTx(filePath)
}

// Done should be called before exit
func (s *Sign) Done() {
	s.dbConn.Close()
	s.btc.Close()
}

// GetBTC gets btc
func (s *Sign) GetBTC() api.Bitcoiner {
	return s.btc
}
