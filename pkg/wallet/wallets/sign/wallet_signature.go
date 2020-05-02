package sign

import (
	"database/sql"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets/coldwalletsrv"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets/coldwalletsrv/signsrv"
)

// Sign is sign wallet object
type Sign struct {
	btc         api.Bitcoiner
	dbConn      *sql.DB
	authAccount account.AuthType
	wtype       wallet.WalletType
	coldwalletsrv.Seeder
	coldwalletsrv.HDWalleter
	signsrv.PrivKeyer
	signsrv.FullPubkeyExporter
	coldwalletsrv.Signer
}

// NewSign returns Sign object
// TODO: maybe each services should be exported variable, not embedded innterface
func NewSign(
	btc api.Bitcoiner,
	dbConn *sql.DB,
	authAccount account.AuthType,
	seeder coldwalletsrv.Seeder,
	hdWallter coldwalletsrv.HDWalleter,
	privKeyer signsrv.PrivKeyer,
	fullPubkeyExporter signsrv.FullPubkeyExporter,
	signer coldwalletsrv.Signer,
	wtype wallet.WalletType) *Sign {

	return &Sign{
		btc:                btc,
		dbConn:             dbConn,
		authAccount:        authAccount,
		Seeder:             seeder,
		HDWalleter:         hdWallter,
		PrivKeyer:          privKeyer,
		FullPubkeyExporter: fullPubkeyExporter,
		Signer:             signer,
		wtype:              wtype,
	}
}

// Done should be called before exit
func (s *Sign) Done() {
	s.dbConn.Close()
	s.btc.Close()
}

// BeginTx starts transaction
//func (k *Keygen) BeginTx() (*sql.Tx, error) {
//	return k.dbConn.Begin()
//}

// GetBTC gets btc
func (s *Sign) GetBTC() api.Bitcoiner {
	return s.btc
}

// GetType gets wallet type
func (s *Sign) GetType() wallet.WalletType {
	return s.wtype
}

// GetAuthType gets auth_type
func (s *Sign) GetAuthType() account.AuthType {
	return s.authAccount
}

// Seed returns Seeder interface
func (s *Sign) Seed() coldwalletsrv.Seeder {
	return s.Seeder
}

// HDWallet returns HDWalleter interface
func (s *Sign) HDWallet() coldwalletsrv.HDWalleter {
	return s.HDWalleter
}

// PrivKey returns PrivKeyer interface
func (s *Sign) PrivKey() signsrv.PrivKeyer {
	return s.PrivKeyer
}

// FullPubkeyExport returns FullPubkeyExporter interface
func (s *Sign) FullPubkeyExport() signsrv.FullPubkeyExporter {
	return s.FullPubkeyExporter
}

// Sign returns Signer interface
func (s *Sign) Sign() coldwalletsrv.Signer {
	return s.Signer
}
