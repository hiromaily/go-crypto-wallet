package ethwallet

import (
	"context"
	"database/sql"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ethereum"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// ETHKeygen keygen wallet object
type ETHKeygen struct {
	ETH                     ethereum.Ethereumer
	dbConn                  *sql.DB
	wtype                   domainWallet.WalletType
	generateSeedUseCase     keygenusecase.GenerateSeedUseCase
	generateHDWalletUseCase keygenusecase.GenerateHDWalletUseCase
	importPrivKeyUseCase    keygenusecase.ImportPrivateKeyUseCase
	exportAddressUseCase    keygenusecase.ExportAddressUseCase
	signTxUseCase           keygenusecase.SignTransactionUseCase
}

// NewETHKeygen returns ETHKeygen object
func NewETHKeygen(
	eth ethereum.Ethereumer,
	dbConn *sql.DB,
	walletType domainWallet.WalletType,
	generateSeedUseCase keygenusecase.GenerateSeedUseCase,
	generateHDWalletUseCase keygenusecase.GenerateHDWalletUseCase,
	importPrivKeyUseCase keygenusecase.ImportPrivateKeyUseCase,
	exportAddressUseCase keygenusecase.ExportAddressUseCase,
	signTxUseCase keygenusecase.SignTransactionUseCase,
) *ETHKeygen {
	return &ETHKeygen{
		ETH:                     eth,
		dbConn:                  dbConn,
		wtype:                   walletType,
		generateSeedUseCase:     generateSeedUseCase,
		generateHDWalletUseCase: generateHDWalletUseCase,
		importPrivKeyUseCase:    importPrivKeyUseCase,
		exportAddressUseCase:    exportAddressUseCase,
		signTxUseCase:           signTxUseCase,
	}
}

// GenerateSeed generates seed
func (k *ETHKeygen) GenerateSeed() ([]byte, error) {
	output, err := k.generateSeedUseCase.Generate(context.Background())
	if err != nil {
		return nil, err
	}
	return output.Seed, nil
}

// StoreSeed stores seed
func (k *ETHKeygen) StoreSeed(strSeed string) ([]byte, error) {
	output, err := k.generateSeedUseCase.Store(context.Background(), keygenusecase.StoreSeedInput{
		Seed: strSeed,
	})
	if err != nil {
		return nil, err
	}
	return output.Seed, nil
}

// GenerateAccountKey generates account keys
func (k *ETHKeygen) GenerateAccountKey(
	accountType domainAccount.AccountType, seed []byte, count uint32, _ bool,
) ([]domainKey.WalletKey, error) {
	_, err := k.generateHDWalletUseCase.Generate(context.Background(), keygenusecase.GenerateHDWalletInput{
		AccountType: accountType,
		Seed:        seed,
		Count:       count,
	})
	if err != nil {
		return nil, err
	}
	// Note: The original implementation returns keys but use case returns count
	// For now, return nil as the keys are stored in the database
	return nil, nil
}

// ImportPrivKey imports privKey
func (k *ETHKeygen) ImportPrivKey(accountType domainAccount.AccountType) error {
	return k.importPrivKeyUseCase.Import(context.Background(), keygenusecase.ImportPrivateKeyInput{
		AccountType: accountType,
	})
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
	output, err := k.exportAddressUseCase.Export(context.Background(), keygenusecase.ExportAddressInput{
		AccountType: accountType,
	})
	if err != nil {
		return "", err
	}
	return output.FileName, nil
}

// SignTx signs on transaction
func (k *ETHKeygen) SignTx(filePath string) (string, bool, string, error) {
	output, err := k.signTxUseCase.Sign(context.Background(), keygenusecase.SignTransactionInput{
		FilePath: filePath,
	})
	if err != nil {
		return "", false, "", err
	}

	// Determine if signing is done based on signed vs unsigned count
	isDone := output.UnsignedCount == 0

	return output.FilePath, isDone, "", nil
}

// Done should be called before exit
func (k *ETHKeygen) Done() {
	_ = k.dbConn.Close() // Best effort cleanup
	k.ETH.Close()
}
