package xrpwallet

import (
	"database/sql"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	wtype "github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/xrp/keygensrv"
)

// XRPKeygen keygen wallet object
type XRPKeygen struct {
	XRP    xrpgrp.Rippler
	dbConn *sql.DB
	logger *zap.Logger
	wtype  wtype.WalletType
	service.Seeder
	service.HDWalleter
	keygensrv.XRPKeyGenerator
	service.AddressExporter
}

// NewXRPKeygen returns XRPKeygen object
func NewXRPKeygen(
	xrp xrpgrp.Rippler,
	dbConn *sql.DB,
	logger *zap.Logger,
	wtype wtype.WalletType,
	seeder service.Seeder,
	hdWallter service.HDWalleter,
	keyGenerator keygensrv.XRPKeyGenerator,
	addressExporter service.AddressExporter,
) *XRPKeygen {

	return &XRPKeygen{
		XRP:             xrp,
		logger:          logger,
		dbConn:          dbConn,
		wtype:           wtype,
		Seeder:          seeder,
		HDWalleter:      hdWallter,
		XRPKeyGenerator: keyGenerator,
		AddressExporter: addressExporter,
	}
}

// GenerateSeed generates seed
func (k *XRPKeygen) GenerateSeed() ([]byte, error) {
	//k.logger.Info("no functionality for GenerateSeed() in XRP")
	return k.Seeder.Generate()
}

// StoreSeed stores seed
func (k *XRPKeygen) StoreSeed(strSeed string) ([]byte, error) {
	//k.logger.Info("no functionality for StoreSeed() in XRP")
	return k.Seeder.Store(strSeed)
}

// GenerateAccountKey generates account keys
func (k *XRPKeygen) GenerateAccountKey(accountType account.AccountType, seed []byte, count uint32, isKeyPair bool) ([]key.WalletKey, error) {
	keys, err := k.HDWalleter.Generate(accountType, seed, count)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call HDWalleter.Generate()")
	}
	err = k.XRPKeyGenerator.Generate(accountType, isKeyPair, keys)
	return keys, err
}

// ImportPrivKey imports privKey
func (k *XRPKeygen) ImportPrivKey(accountType account.AccountType) error {
	k.logger.Info("not implemented yet")
	//return k.PrivKeyer.Import(accountType)
	return nil
}

// ImportFullPubKey imports full-pubkey
func (k *XRPKeygen) ImportFullPubKey(fileName string) error {
	//return k.FullPubKeyImporter.ImportFullPubKey(fileName)
	k.logger.Info("no functionality for ImportFullPubKey() in XRP")
	return nil
}

// CreateMultisigAddress creates multi sig address returns Multisiger interface
func (k *XRPKeygen) CreateMultisigAddress(accountType account.AccountType) error {
	//return k.Multisiger.AddMultisigAddress(accountType, k.addrType)
	k.logger.Info("no functionality for CreateMultisigAddress() in XRP")
	return nil
}

// ExportAddress exports address
func (k *XRPKeygen) ExportAddress(accountType account.AccountType) (string, error) {
	return k.AddressExporter.ExportAddress(accountType)
}

// SignTx signs on transaction
func (k *XRPKeygen) SignTx(filePath string) (string, bool, string, error) {
	k.logger.Info("not implemented yet")
	//return k.Signer.SignTx(filePath)
	return "", false, "", nil
}

// Done should be called before exit
func (k *XRPKeygen) Done() {
	k.dbConn.Close()
	k.XRP.Close()
}
