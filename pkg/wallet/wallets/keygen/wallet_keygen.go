package keygen

import (
	"database/sql"

	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets/coldwalletsrv"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets/coldwalletsrv/keygensrv"
)

// Keygen is keygen wallet object
// TODO: maybe each services should be exported variable, not embedded innterface
type Keygen struct {
	btc    api.Bitcoiner
	dbConn *sql.DB
	wtype  wallet.WalletType
	coldwalletsrv.Seeder
	coldwalletsrv.HDWalleter
	keygensrv.PrivKeyer
	keygensrv.FullPubKeyImporter
	keygensrv.Multisiger
	keygensrv.AddressExporter
	coldwalletsrv.Signer
	//Seeder          coldwalletsrv.Seeder
	//HdWallter       coldwalletsrv.HDWalleter
	//PrivKeyer       coldwalletsrv.PrivKeyer
	//PubkeyImporter  keygensrv.PubKeyImporter
	//Multisiger      keygensrv.Multisiger
	//AddressExporter keygensrv.AddressExporter
}

// NewKeygen returns Keygen object
func NewKeygen(
	btc api.Bitcoiner,
	dbConn *sql.DB,
	seeder coldwalletsrv.Seeder,
	hdWallter coldwalletsrv.HDWalleter,
	privKeyer keygensrv.PrivKeyer,
	pubkeyImporter keygensrv.FullPubKeyImporter,
	multisiger keygensrv.Multisiger,
	addressExporter keygensrv.AddressExporter,
	signer coldwalletsrv.Signer,
	wtype wallet.WalletType) *Keygen {

	return &Keygen{
		btc:                btc,
		dbConn:             dbConn,
		Seeder:             seeder,
		HDWalleter:         hdWallter,
		PrivKeyer:          privKeyer,
		FullPubKeyImporter: pubkeyImporter,
		Multisiger:         multisiger,
		AddressExporter:    addressExporter,
		Signer:             signer,
		wtype:              wtype,
	}
}

// Done should be called before exit
func (k *Keygen) Done() {
	k.dbConn.Close()
	k.btc.Close()
}

// BeginTx starts transaction
//func (k *Keygen) BeginTx() (*sql.Tx, error) {
//	return k.dbConn.Begin()
//}

// GetBTC gets btc
func (k *Keygen) GetBTC() api.Bitcoiner {
	return k.btc
}

// GetType gets wallet type
func (k *Keygen) GetType() wallet.WalletType {
	return k.wtype
}

// Seed returns Seeder interface
func (k *Keygen) Seed() coldwalletsrv.Seeder {
	return k.Seeder
}

// HDWallet returns HDWalleter interface
func (k *Keygen) HDWallet() coldwalletsrv.HDWalleter {
	return k.HDWalleter
}

// PrivKey returns PrivKeyer interface
func (k *Keygen) PrivKey() keygensrv.PrivKeyer {
	return k.PrivKeyer
}

// FullPubKeyImport returns PubKeyImporter interface
func (k *Keygen) FullPubKeyImport() keygensrv.FullPubKeyImporter {
	return k.FullPubKeyImporter
}

// Multisig returns Multisiger interface
func (k *Keygen) Multisig() keygensrv.Multisiger {
	return k.Multisiger
}

// AddressExport returns AddressExporter interface
func (k *Keygen) AddressExport() keygensrv.AddressExporter {
	return k.AddressExporter
}

// Sign returns Signer interface
func (k *Keygen) Sign() coldwalletsrv.Signer {
	return k.Signer
}
