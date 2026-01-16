package btc

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/btcsuite/btcd/chaincfg"

	portsfile "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repository "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository"
	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/fullpubkey"
	infraKey "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
)

type exportFullPubkeyUseCase struct {
	authKeyRepo    repository.AuthAccountKeyRepositorier
	seedRepo       repository.SeedRepositorier
	pubkeyFileRepo portsfile.AddressFileRepositorier
	coinTypeCode   domainCoin.CoinTypeCode
	authType       domainAccount.AuthType
	wtype          domainWallet.WalletType
	chainConfig    *chaincfg.Params
}

// NewExportFullPubkeyUseCase creates a new ExportFullPubkeyUseCase for sign wallet
func NewExportFullPubkeyUseCase(
	authKeyRepo repository.AuthAccountKeyRepositorier,
	seedRepo repository.SeedRepositorier,
	pubkeyFileRepo portsfile.AddressFileRepositorier,
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

	// Determine coin index based on network
	// coin' = 0' for mainnet, 1' for testnet/regtest
	coinIndex := uint32(0) // mainnet
	if u.chainConfig.Net != 0x00000000 {
		// testnet, regtest, or other networks use 1'
		coinIndex = 1
	}

	// Get master fingerprint once (same for all BIP purposes)
	// Use any purpose for fingerprint calculation since it's derived from seed
	tempHdKey := infraKey.NewHDKey(infraKey.PurposeTypeBIP44, u.coinTypeCode, u.chainConfig)
	tempDescGenerator := infraKey.NewDescriptorGenerator(tempHdKey, u.chainConfig)
	fingerprint, err := tempDescGenerator.GetMasterFingerprintHex(seed)
	if err != nil {
		return "", fmt.Errorf("fail to get master fingerprint: %w", err)
	}

	// Export extended public keys for all 4 BIP purposes
	// Each sign wallet exports 4 xpubs (one for each BIP purpose)
	// to support all address types (Legacy, P2SH-SegWit, Native SegWit, Taproot)
	purposes := []infraKey.PurposeType{
		infraKey.PurposeTypeBIP44, // Legacy (P2PKH)
		infraKey.PurposeTypeBIP49, // P2SH-SegWit
		infraKey.PurposeTypeBIP84, // Native SegWit (Bech32)
		infraKey.PurposeTypeBIP86, // Taproot
	}

	for _, purpose := range purposes {
		// Create HDKey for this BIP purpose
		hdKey := infraKey.NewHDKey(purpose, u.coinTypeCode, u.chainConfig)

		// Derive extended public key from seed
		descGenerator := infraKey.NewDescriptorGenerator(hdKey, u.chainConfig)

		// Export purpose/coin-level extended public key (NOT account-level)
		// This allows keygen to derive account-specific keys for deposit, payment, stored
		// Path: m/purpose'/coin' (BIP48-style for multisig, stopping at coin level)
		//
		// Why not account-level (m/purpose'/coin'/account')?
		// - Sign wallets don't know which multisig account (deposit/payment/stored) will use these keys
		// - Each sign wallet should provide ONE extended key per BIP purpose
		//   that can derive ALL multisig accounts
		// - Keygen will derive account-specific keys:
		//   m/purpose'/coin'/0' for deposit, m/purpose'/coin'/1' for payment, etc.
		extendedPubKey, err := descGenerator.GetCoinLevelXPub(seed)
		if err != nil {
			return "", fmt.Errorf("fail to derive extended public key for BIP%d: %w", purpose, err)
		}

		// Build derivation path at purpose/coin level (NOT account level)
		// Format: m/purpose'/coin_type'
		// Examples:
		//   - BIP44: m/44'/coin'  (Legacy/P2PKH)
		//   - BIP49: m/49'/coin'  (P2SH-SegWit)
		//   - BIP84: m/84'/coin'  (Native SegWit/Bech32)
		//   - BIP86: m/86'/coin'  (Taproot)
		//
		// Keygen will extend this path to m/purpose'/coin'/account' for each multisig account
		derivationPath := fmt.Sprintf("m/%d'/%d'", purpose, coinIndex)

		// output: coinType, authType, purpose, extendedPubKey, fingerprint, derivationPath
		_, err = writer.WriteString(
			fullpubkey.CreateExtendedLine(
				u.coinTypeCode,
				authType,
				uint8(purpose),
				extendedPubKey,
				fingerprint,
				derivationPath,
			),
		)
		if err != nil {
			return "", fmt.Errorf("fail to call writer.WriteString(%s) for BIP%d: %w", fileName, purpose, err)
		}
	}

	if err = writer.Flush(); err != nil {
		return "", fmt.Errorf("fail to call writer.Flush(%s): %w", fileName, err)
	}
	return fileName, nil
}
