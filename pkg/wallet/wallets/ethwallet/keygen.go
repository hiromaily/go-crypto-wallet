package ethwallet

import (
	"database/sql"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainKey "github.com/hiromaily/go-crypto-wallet/pkg/domain/key"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
)

// ETHKeygen keygen wallet object
type ETHKeygen struct {
	ETH    ethereum.Ethereumer
	dbConn *sql.DB
	wtype  domainWallet.WalletType
	service.Seeder
	service.HDWalleter
	service.PrivKeyer
	service.AddressExporter
	service.Signer
}

// NewETHKeygen returns ETHKeygen object
func NewETHKeygen(
	eth ethereum.Ethereumer,
	dbConn *sql.DB,
	walletType domainWallet.WalletType,
	seeder service.Seeder,
	hdWallter service.HDWalleter,
	privKeyer service.PrivKeyer,
	addressExporter service.AddressExporter,
	signer service.Signer,
) *ETHKeygen {
	return &ETHKeygen{
		ETH:             eth,
		dbConn:          dbConn,
		wtype:           walletType,
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
func (k *ETHKeygen) GenerateAccountKey(
	accountType domainAccount.AccountType, seed []byte, count uint32, _ bool,
) ([]domainKey.WalletKey, error) {
	return k.HDWalleter.Generate(accountType, seed, count)
}

// ImportPrivKey imports privKey
func (k *ETHKeygen) ImportPrivKey(accountType domainAccount.AccountType) error {
	return k.PrivKeyer.Import(accountType)
}

// ImportFullPubKey imports full-pubkey
func (*ETHKeygen) ImportFullPubKey(_ string) error {
	// return k.FullPubKeyImporter.ImportFullPubKey(fileName)
	logger.Info("no functionality for ImportFullPubKey() in ETH")
	return nil
}

// CreateMultisigAddress creates multi sig address returns Multisiger interface
func (*ETHKeygen) CreateMultisigAddress(_ domainAccount.AccountType) error {
	logger.Info("no functionality for CreateMultisigAddress() in ETH")
	return nil
}

// ExportAddress exports address
func (k *ETHKeygen) ExportAddress(accountType domainAccount.AccountType) (string, error) {
	return k.AddressExporter.ExportAddress(accountType)
}

// SignTx signs on transaction
func (k *ETHKeygen) SignTx(filePath string) (string, bool, string, error) {
	return k.Signer.SignTx(filePath)
}

// Done should be called before exit
func (k *ETHKeygen) Done() {
	k.dbConn.Close()
	k.ETH.Close()
}
