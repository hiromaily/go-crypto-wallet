package ethwallet

import (
	"database/sql"

	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	wtype "github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
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
) *ETHKeygen {

	return &ETHKeygen{
		ETH:        eth,
		logger:     logger,
		dbConn:     dbConn,
		wtype:      wtype,
		Seeder:     seeder,
		HDWalleter: hdWallter,
		PrivKeyer:  privKeyer,
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
func (k *ETHKeygen) GenerateAccountKey(accountType account.AccountType, seed []byte, count uint32) ([]key.WalletKey, error) {
	return k.HDWalleter.Generate(accountType, seed, count)
}

// ImportPrivKey imports privKey
func (k *ETHKeygen) ImportPrivKey(accountType account.AccountType) error {
	return k.PrivKeyer.Import(accountType)
}

// ImportFullPubKey imports full-pubkey
func (k *ETHKeygen) ImportFullPubKey(fileName string) error {
	//return k.FullPubKeyImporter.ImportFullPubKey(fileName)
	k.logger.Info("no functionality for ImportFullPubKey() in ETH")
	return nil
}

// CreateMultisigAddress creates multi sig address returns Multisiger interface
func (k *ETHKeygen) CreateMultisigAddress(accountType account.AccountType) error {
	//return k.Multisiger.AddMultisigAddress(accountType, k.addrType)
	k.logger.Info("no functionality for CreateMultisigAddress() in ETH")
	return nil
}

// ExportAddress exports address
func (k *ETHKeygen) ExportAddress(accountType account.AccountType) (string, error) {
	//return k.AddressExporter.ExportAddress(accountType)
	k.logger.Warn("not implemented yet in ETH")
	return "", nil
}

// SignTx signs on transaction
func (k *ETHKeygen) SignTx(filePath string) (string, bool, string, error) {
	//return k.Signer.SignTx(filePath)
	k.logger.Warn("not implemented yet in ETH")
	return "", false, "", nil
}

// Done should be called before exit
func (k *ETHKeygen) Done() {
	k.dbConn.Close()
	k.ETH.Close()
}
