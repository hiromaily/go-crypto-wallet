package xrpwallet

import (
	"database/sql"
	"fmt"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainKey "github.com/hiromaily/go-crypto-wallet/pkg/domain/key"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ripple"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/xrp/keygensrv"
)

// XRPKeygen keygen wallet object
type XRPKeygen struct {
	XRP    ripple.Rippler
	dbConn *sql.DB
	wtype  domainWallet.WalletType
	service.Seeder
	service.HDWalleter
	keygensrv.XRPKeyGenerator
	service.AddressExporter
	service.Signer
}

// NewXRPKeygen returns XRPKeygen object
func NewXRPKeygen(
	xrp ripple.Rippler,
	dbConn *sql.DB,
	walletType domainWallet.WalletType,
	seeder service.Seeder,
	hdWallter service.HDWalleter,
	keyGenerator keygensrv.XRPKeyGenerator,
	addressExporter service.AddressExporter,
	signer service.Signer,
) *XRPKeygen {
	return &XRPKeygen{
		XRP:             xrp,
		dbConn:          dbConn,
		wtype:           walletType,
		Seeder:          seeder,
		HDWalleter:      hdWallter,
		XRPKeyGenerator: keyGenerator,
		AddressExporter: addressExporter,
		Signer:          signer,
	}
}

// GenerateSeed generates seed
func (k *XRPKeygen) GenerateSeed() ([]byte, error) {
	// k.logger.Info("no functionality for GenerateSeed() in XRP")
	return k.Seeder.Generate()
}

// StoreSeed stores seed
func (k *XRPKeygen) StoreSeed(strSeed string) ([]byte, error) {
	// k.logger.Info("no functionality for StoreSeed() in XRP")
	return k.Seeder.Store(strSeed)
}

// GenerateAccountKey generates account keys
func (k *XRPKeygen) GenerateAccountKey(
	accountType domainAccount.AccountType, seed []byte, count uint32, isKeyPair bool,
) ([]domainKey.WalletKey, error) {
	keys, err := k.HDWalleter.Generate(accountType, seed, count)
	if err != nil {
		return nil, fmt.Errorf("fail to call HDWalleter.Generate(): %w", err)
	}
	err = k.XRPKeyGenerator.Generate(accountType, isKeyPair, keys)
	return keys, err
}

// ImportPrivKey imports privKey
func (*XRPKeygen) ImportPrivKey(_ domainAccount.AccountType) error {
	logger.Info("no functionality for ImportPrivKey() in XRP")
	return nil
}

// ImportFullPubKey imports full-pubkey
func (*XRPKeygen) ImportFullPubKey(_ string) error {
	logger.Info("no functionality for ImportFullPubKey() in XRP")
	return nil
}

// CreateMultisigAddress creates multi sig address returns Multisiger interface
func (*XRPKeygen) CreateMultisigAddress(_ domainAccount.AccountType) error {
	logger.Info("no functionality for CreateMultisigAddress() in XRP")
	return nil
}

// ExportAddress exports address
func (k *XRPKeygen) ExportAddress(accountType domainAccount.AccountType) (string, error) {
	return k.AddressExporter.ExportAddress(accountType)
}

// SignTx signs on transaction
func (k *XRPKeygen) SignTx(filePath string) (string, bool, string, error) {
	return k.Signer.SignTx(filePath)
}

// Done should be called before exit
func (k *XRPKeygen) Done() {
	k.dbConn.Close()

	k.XRP.Close()
}
