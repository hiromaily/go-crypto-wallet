package xrp

import (
	"context"
	"database/sql"
	"fmt"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ripple"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// XRPKeygen keygen wallet object
type XRPKeygen struct {
	XRP                     ripple.Rippler
	dbConn                  *sql.DB
	wtype                   domainWallet.WalletType
	generateSeedUseCase     keygenusecase.GenerateSeedUseCase
	generateHDWalletUseCase keygenusecase.GenerateHDWalletUseCase
	generateKeyUseCase      keygenusecase.GenerateKeyUseCase
	exportAddressUseCase    keygenusecase.ExportAddressUseCase
	signTxUseCase           keygenusecase.SignTransactionUseCase
}

// NewXRPKeygen returns XRPKeygen object
func NewXRPKeygen(
	xrp ripple.Rippler,
	dbConn *sql.DB,
	walletType domainWallet.WalletType,
	generateSeedUseCase keygenusecase.GenerateSeedUseCase,
	generateHDWalletUseCase keygenusecase.GenerateHDWalletUseCase,
	generateKeyUseCase keygenusecase.GenerateKeyUseCase,
	exportAddressUseCase keygenusecase.ExportAddressUseCase,
	signTxUseCase keygenusecase.SignTransactionUseCase,
) *XRPKeygen {
	return &XRPKeygen{
		XRP:                     xrp,
		dbConn:                  dbConn,
		wtype:                   walletType,
		generateSeedUseCase:     generateSeedUseCase,
		generateHDWalletUseCase: generateHDWalletUseCase,
		generateKeyUseCase:      generateKeyUseCase,
		exportAddressUseCase:    exportAddressUseCase,
		signTxUseCase:           signTxUseCase,
	}
}

// GenerateSeed generates seed
func (k *XRPKeygen) GenerateSeed() ([]byte, error) {
	// k.logger.Info("no functionality for GenerateSeed() in XRP")
	output, err := k.generateSeedUseCase.Generate(context.Background())
	if err != nil {
		return nil, err
	}
	return output.Seed, nil
}

// StoreSeed stores seed
func (k *XRPKeygen) StoreSeed(strSeed string) ([]byte, error) {
	// k.logger.Info("no functionality for StoreSeed() in XRP")
	output, err := k.generateSeedUseCase.Store(context.Background(), keygenusecase.StoreSeedInput{
		Seed: strSeed,
	})
	if err != nil {
		return nil, err
	}
	return output.Seed, nil
}

// GenerateAccountKey generates account keys
func (k *XRPKeygen) GenerateAccountKey(
	accountType domainAccount.AccountType, seed []byte, count uint32, isKeyPair bool,
) ([]domainKey.WalletKey, error) {
	// First, generate HD wallet keys
	output, err := k.generateHDWalletUseCase.Generate(context.Background(), keygenusecase.GenerateHDWalletInput{
		AccountType: accountType,
		Seed:        seed,
		Count:       count,
	})
	if err != nil {
		return nil, fmt.Errorf("fail to call HDWalleter.Generate(): %w", err)
	}

	// Then, generate XRP-specific keys
	// Note: We pass nil for keys since use case handles retrieving them
	err = k.generateKeyUseCase.Generate(context.Background(), keygenusecase.GenerateKeyInput{
		AccountType: accountType,
		IsKeyPair:   isKeyPair,
		WalletKeys:  nil, // Will be retrieved from database
	})
	if err != nil {
		return nil, err
	}

	// Return nil since keys are stored in database
	// The output only has count, not the keys themselves
	_ = output
	return nil, nil
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
	output, err := k.exportAddressUseCase.Export(context.Background(), keygenusecase.ExportAddressInput{
		AccountType: accountType,
	})
	if err != nil {
		return "", err
	}
	return output.FileName, nil
}

// SignTx signs on transaction
func (k *XRPKeygen) SignTx(filePath string) (string, bool, string, error) {
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
func (k *XRPKeygen) Done() {
	_ = k.dbConn.Close() // Best effort cleanup

	_ = k.XRP.Close() // Best effort cleanup
}
