package ethwallet

import (
	"database/sql"

	"go.uber.org/zap"

	wtype "github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
)

// ETHSign keygen wallet object
type ETHSign struct {
	ETH    ethgrp.Ethereumer
	dbConn *sql.DB
	logger *zap.Logger
	wtype  wtype.WalletType
}

// NewETHSign returns ETHSign object
func NewETHSign(
	eth ethgrp.Ethereumer,
	dbConn *sql.DB,
	logger *zap.Logger,
	wtype wtype.WalletType) *ETHSign {

	return &ETHSign{
		ETH:    eth,
		logger: logger,
		dbConn: dbConn,
		wtype:  wtype,
	}
}

// GenerateSeed generates seed
func (s *ETHSign) GenerateSeed() ([]byte, error) {
	s.logger.Info("no functionality for CreateMultisigAddress() in ETH")
	//return s.Seeder.Generate()
	return nil, nil
}

// StoreSeed stores seed
func (s *ETHSign) StoreSeed(strSeed string) ([]byte, error) {
	s.logger.Info("no functionality for CreateMultisigAddress() in ETH")
	//return s.Seeder.Store(strSeed)
	return nil, nil
}

// GenerateAuthKey generates account keys
func (s *ETHSign) GenerateAuthKey(seed []byte, count uint32) ([]key.WalletKey, error) {
	s.logger.Info("no functionality for CreateMultisigAddress() in ETH")
	//return s.HDWalleter.Generate(s.authAccount.AccountType(), seed, count)
	return nil, nil
}

// ImportPrivKey imports privKey
func (s *ETHSign) ImportPrivKey() error {
	s.logger.Info("no functionality for CreateMultisigAddress() in ETH")
	//return s.PrivKeyer.Import()
	return nil
}

// ExportFullPubkey exports full-pubkey
func (s *ETHSign) ExportFullPubkey() (string, error) {
	s.logger.Info("no functionality for CreateMultisigAddress() in ETH")
	//return s.FullPubkeyExporter.ExportFullPubkey()
	return "", nil
}

// SignTx signs on transaction
func (s *ETHSign) SignTx(filePath string) (string, bool, string, error) {
	s.logger.Info("no functionality for CreateMultisigAddress() in ETH")
	//return s.Signer.SignTx(filePath)
	return "", false, "", nil
}

// Done should be called before exit
func (s *ETHSign) Done() {
	s.dbConn.Close()
	s.ETH.Close()
}
