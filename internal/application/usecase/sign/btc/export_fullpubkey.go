package btc

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/btcsuite/btcd/chaincfg"

	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	portsStorage "github.com/hiromaily/go-crypto-wallet/internal/application/ports/storage"
	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/fullpubkey"
	infraKey "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
)

type exportFullPubkeyUseCase struct {
	authKeyRepo    persistence.AuthAccountKeyRepositorier
	seedRepo       persistence.SeedRepositorier
	pubkeyFileRepo portsStorage.AddressFileRepositorier
	coinTypeCode   domainCoin.CoinTypeCode
	authType       domainAccount.AuthType
	wtype          domainWallet.WalletType
	chainConfig    *chaincfg.Params
}

// NewExportFullPubkeyUseCase creates a new ExportFullPubkeyUseCase for sign wallet
func NewExportFullPubkeyUseCase(
	authKeyRepo persistence.AuthAccountKeyRepositorier,
	seedRepo persistence.SeedRepositorier,
	pubkeyFileRepo portsStorage.AddressFileRepositorier,
	coinTypeCode domainCoin.CoinTypeCode,
	authType domainAccount.AuthType,
	wtype domainWallet.WalletType,
	chainConfig *chaincfg.Params,
) signusecase.ExportFullPubkeyUseCase {
	return &exportFullPubkeyUseCase{
		authKeyRepo:    authKeyRepo,
		seedRepo:       seedRepo,
		pubkeyFileRepo: pubkeyFileRepo,
		coinTypeCode:   coinTypeCode,
		authType:       authType,
		wtype:          wtype,
		chainConfig:    chainConfig,
	}
}

func (u *exportFullPubkeyUseCase) Export(ctx context.Context) (signusecase.ExportFullPubkeyOutput, error) {
	// Get seed to derive extended public key
	seedData, err := u.seedRepo.GetOne(ctx)
	if err != nil {
		return signusecase.ExportFullPubkeyOutput{},
			fmt.Errorf("fail to get seed: %w", err)
	}

	// Convert seed string to bytes
	seedBytes, err := infraKey.SeedToByte(seedData.Seed)
	if err != nil {
		return signusecase.ExportFullPubkeyOutput{},
			fmt.Errorf("fail to decode seed: %w", err)
	}

	// export csv file with extended key
	fileName, err := u.exportAccountKey(u.authType, seedBytes)
	if err != nil {
		return signusecase.ExportFullPubkeyOutput{}, err
	}

	return signusecase.ExportFullPubkeyOutput{
		FileName: fileName,
	}, nil
}

// exportAccountKey export account_key_table as csv file with extended public key
func (u *exportFullPubkeyUseCase) exportAccountKey(
	authType domainAccount.AuthType,
	seed []byte,
) (string, error) {
	// create fileName
	fileName := u.pubkeyFileRepo.CreateFilePath(u.authType.AccountType())

	file, err := os.Create(fileName) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("fail to call os.Create(%s): %w", fileName, err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	writer := bufio.NewWriter(file)

	// Create HDKey for BIP49 (P2SH-SegWit)
	// Auth accounts are used for multisig P2SH-wrapped SegWit addresses
	// BIP49 path: m/49'/coin'/account'
	hdKey := infraKey.NewHDKey(infraKey.PurposeTypeBIP49, u.coinTypeCode, u.chainConfig)

	// Derive extended public key and fingerprint from seed
	descGenerator := infraKey.NewDescriptorGenerator(hdKey, u.chainConfig)

	// Export purpose/coin-level extended public key (NOT account-level)
	// This allows keygen to derive account-specific keys for deposit, payment, stored
	// Path: m/49'/coin' (BIP48-style for multisig, stopping at coin level)
	//
	// Why not account-level (m/49'/coin'/account')?
	// - Sign wallets don't know which multisig account (deposit/payment/stored) will use these keys
	// - Each sign wallet should provide ONE extended key that can derive ALL multisig accounts
	// - Keygen will derive account-specific keys: m/49'/coin'/0' for deposit, m/49'/coin'/1' for payment, etc.
	extendedPubKey, err := descGenerator.GetCoinLevelXPub(seed)
	if err != nil {
		return "", fmt.Errorf("fail to derive extended public key: %w", err)
	}

	// Get master fingerprint
	fingerprint, err := descGenerator.GetMasterFingerprintHex(seed)
	if err != nil {
		return "", fmt.Errorf("fail to get master fingerprint: %w", err)
	}

	// Build derivation path at purpose/coin level (NOT account level)
	// Format: m/purpose'/coin_type'
	// For BIP49: m/49'/coin'
	// coin' = 0' for mainnet, 1' for testnet/regtest
	//
	// Keygen will extend this path to m/49'/coin'/account' for each multisig account
	coinIndex := uint32(0) // mainnet
	if u.chainConfig.Net != 0x00000000 {
		// testnet, regtest, or other networks use 1'
		coinIndex = 1
	}
	derivationPath := fmt.Sprintf("m/49'/%d'",
		coinIndex,
	)

	// output: coinType, authType, extendedPubKey, fingerprint, derivationPath
	_, err = writer.WriteString(
		fullpubkey.CreateExtendedLine(
			u.coinTypeCode,
			authType,
			extendedPubKey,
			fingerprint,
			derivationPath,
		),
	)
	if err != nil {
		return "", fmt.Errorf("fail to call writer.WriteString(%s): %w", fileName, err)
	}
	if err = writer.Flush(); err != nil {
		return "", fmt.Errorf("fail to call writer.Flush(%s): %w", fileName, err)
	}
	return fileName, nil
}
