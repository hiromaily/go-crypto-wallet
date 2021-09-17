package btcwallet

import (
	"database/sql"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
)

// BTCKeygen is keygen wallet object
type BTCKeygen struct {
	BTC      btcgrp.Bitcoiner
	dbConn   *sql.DB
	addrType address.AddrType
	wtype    wallet.WalletType
	service.Seeder
	service.HDWalleter
	service.PrivKeyer
	service.FullPubKeyImporter
	service.Multisiger
	service.AddressExporter
	service.Signer
}

// NewBTCKeygen returns Keygen object
func NewBTCKeygen(
	btc btcgrp.Bitcoiner,
	dbConn *sql.DB,
	addrType address.AddrType,
	seeder service.Seeder,
	hdWallter service.HDWalleter,
	privKeyer service.PrivKeyer,
	pubkeyImporter service.FullPubKeyImporter,
	multisiger service.Multisiger,
	addressExporter service.AddressExporter,
	signer service.Signer,
	wtype wallet.WalletType) *BTCKeygen {
	return &BTCKeygen{
		BTC:                btc,
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
func (k *BTCKeygen) GenerateSeed() ([]byte, error) {
	return k.Seeder.Generate()
}

// StoreSeed stores seed
func (k *BTCKeygen) StoreSeed(strSeed string) ([]byte, error) {
	return k.Seeder.Store(strSeed)
}

// GenerateAccountKey generates account keys
func (k *BTCKeygen) GenerateAccountKey(accountType account.AccountType, seed []byte, count uint32, _ bool) ([]key.WalletKey, error) {
	return k.HDWalleter.Generate(accountType, seed, count)
}

// ImportPrivKey imports privKey
func (k *BTCKeygen) ImportPrivKey(accountType account.AccountType) error {
	return k.PrivKeyer.Import(accountType)
}

// ImportFullPubKey imports full-pubkey
func (k *BTCKeygen) ImportFullPubKey(fileName string) error {
	return k.FullPubKeyImporter.ImportFullPubKey(fileName)
}

// CreateMultisigAddress creates multi sig address returns Multisiger interface
func (k *BTCKeygen) CreateMultisigAddress(accountType account.AccountType) error {
	return k.Multisiger.AddMultisigAddress(accountType, k.addrType)
}

// ExportAddress exports address
func (k *BTCKeygen) ExportAddress(accountType account.AccountType) (string, error) {
	return k.AddressExporter.ExportAddress(accountType)
}

// SignTx signs on transaction
func (k *BTCKeygen) SignTx(filePath string) (string, bool, string, error) {
	return k.Signer.SignTx(filePath)
}

// Done should be called before exit
func (k *BTCKeygen) Done() {
	k.dbConn.Close()
	k.BTC.Close()
}
