package ethwallet

import (
	"database/sql"

	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	wtype "github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
)

// ETHKeygen keygen wallet object
type ETHKeygen struct {
	ETH    ethgrp.Ethereumer
	dbConn *sql.DB
	logger *zap.Logger
	wtype  wtype.WalletType
	service.Seeder
	service.HDWalleter
	service.PrivKeyer
	service.AddressExporter
	service.Signer
}

// NewETHKeygen returns ETHKeygen object
func NewETHKeygen(
	eth ethgrp.Ethereumer,
	dbConn *sql.DB,
	logger *zap.Logger,
	wtype wtype.WalletType,
	seeder service.Seeder,
	hdWallter service.HDWalleter,
	privKeyer service.PrivKeyer,
	addressExporter service.AddressExporter,
	signer service.Signer,
) *ETHKeygen {
	return &ETHKeygen{
		ETH:             eth,
		logger:          logger,
		dbConn:          dbConn,
		wtype:           wtype,
		Seeder:          seeder,
		HDWalleter:      hdWallter,
		PrivKeyer:       privKeyer,
		AddressExporter: addressExporter,
		Signer:          signer,
	}
}

// GenerateSeed generates seed
func (k *ETHKeygen) GenerateSeed() ([]byte, error) {
	return k.Seeder.Generate()
}

// StoreSeed stores seed
func (k *ETHKeygen) StoreSeed(strSeed string) ([]byte, error) {
	return k.Seeder.Store(strSeed)
}

// GenerateAccountKey generates account keys
func (k *ETHKeygen) GenerateAccountKey(accountType account.AccountType, seed []byte, count uint32, _ bool) ([]key.WalletKey, error) {
	return k.HDWalleter.Generate(accountType, seed, count)
}

// ImportPrivKey imports privKey
func (k *ETHKeygen) ImportPrivKey(accountType account.AccountType) error {
	return k.PrivKeyer.Import(accountType)
}

// ImportFullPubKey imports full-pubkey
func (k *ETHKeygen) ImportFullPubKey(_ string) error {
	// return k.FullPubKeyImporter.ImportFullPubKey(fileName)
	k.logger.Info("no functionality for ImportFullPubKey() in ETH")
	return nil
}

// CreateMultisigAddress creates multi sig address returns Multisiger interface
func (k *ETHKeygen) CreateMultisigAddress(_ account.AccountType) error {
	k.logger.Info("no functionality for CreateMultisigAddress() in ETH")
	return nil
}

// ExportAddress exports address
func (k *ETHKeygen) ExportAddress(accountType account.AccountType) (string, error) {
	return k.AddressExporter.ExportAddress(accountType)
}

// SignTx signs on transaction
func (k *ETHKeygen) SignTx(filePath string) (string, bool, string, error) {
	return k.Signer.SignTx(filePath)
}

// Done should be called before exit
func (k *ETHKeygen) Done() {
	// nolint:errcheck
	k.dbConn.Close()
	k.ETH.Close()
}
