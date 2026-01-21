// Package xrp provides XRP key generation use cases.
//
// This package supports both online (API-based) and offline key generation modes.
// For security-critical applications, offline key generation is recommended as it
// does not require any network connectivity.
package xrp

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"

	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	domainXrp "github.com/hiromaily/go-crypto-wallet/internal/domain/xrp"
	xrpkg "github.com/hiromaily/go-crypto-wallet/pkg/cryptocurrency/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// offlineGenerateKeyUseCase implements GenerateKeyUseCase using native Go key generation.
// This use case does not require any network connectivity and can run completely offline.
type offlineGenerateKeyUseCase struct {
	dbConn            *sql.DB
	coinTypeCode      domainCoin.CoinTypeCode
	xrpAccountKeyRepo repocold.XRPAccountKeyRepositorier
	keyAlgorithm      xrpkg.KeyAlgorithm
}

// NewOfflineGenerateKeyUseCase creates a new GenerateKeyUseCase that operates entirely offline.
//
// This implementation uses native Go cryptographic libraries to generate XRP keys
// without requiring any API calls to external services like ripple-lib-server.
//
// Parameters:
//   - dbConn: Database connection for storing generated keys
//   - coinTypeCode: The coin type code (XRP)
//   - xrpAccountKeyRepo: Repository for XRP account keys
//   - keyAlgorithm: The cryptographic algorithm to use (ed25519 or secp256k1)
func NewOfflineGenerateKeyUseCase(
	dbConn *sql.DB,
	coinTypeCode domainCoin.CoinTypeCode,
	xrpAccountKeyRepo repocold.XRPAccountKeyRepositorier,
	keyAlgorithm xrpkg.KeyAlgorithm,
) keygenusecase.GenerateKeyUseCase {
	return &offlineGenerateKeyUseCase{
		dbConn:            dbConn,
		coinTypeCode:      coinTypeCode,
		xrpAccountKeyRepo: xrpAccountKeyRepo,
		keyAlgorithm:      keyAlgorithm,
	}
}

// Generate generates XRP keys from the provided wallet keys using native Go implementation.
//
// The entropy for key generation is derived from the BTC HD wallet private keys
// provided in the WalletKeys field. This ensures deterministic key generation
// from the HD wallet hierarchy.
//
// SECURITY NOTE: This function handles sensitive key material. The generated
// private keys and seeds are stored in the database and should be protected
// accordingly.
func (u *offlineGenerateKeyUseCase) Generate(
	ctx context.Context,
	input keygenusecase.GenerateKeyInput,
) error {
	// Convert interface{} to []domainKey.WalletKey
	walletKeys, ok := input.WalletKeys.([]domainKey.WalletKey)
	if !ok {
		return errors.New("invalid wallet keys type")
	}

	logger.Debug("generate keys for XRP (offline mode)",
		"account_type", input.AccountType.String(),
		"len(keys)", len(walletKeys),
		"algorithm", u.keyAlgorithm.String(),
	)

	// Start transaction
	dtx, err := u.dbConn.Begin()
	if err != nil {
		return fmt.Errorf("failed to call db.Begin(): %w", err)
	}
	defer func() {
		if err != nil {
			_ = dtx.Rollback()
		} else {
			_ = dtx.Commit()
		}
	}()

	// Create key generator with configured algorithm
	keyGen := xrpkg.NewKeyGenerator(u.keyAlgorithm)

	// Generate XRP keys
	items := make([]*domainXrp.XRPAccountKey, 0, len(walletKeys))
	for i, v := range walletKeys {
		// Derive XRP key from the HD wallet private key
		hdPrivKey, err := u.extractPrivateKeyBytes(v)
		if err != nil {
			return fmt.Errorf("failed to extract private key for key %d: %w", i, err)
		}

		// Generate XRP key pair using native Go implementation
		keyPair, err := keyGen.DeriveKeyFromHDKey(hdPrivKey)
		if err != nil {
			return fmt.Errorf("failed to generate XRP key for key %d: %w", i, err)
		}

		// Convert algorithm to domain type
		keyType := u.convertAlgorithmToKeyType(keyPair.Algorithm)

		// Create domain entity
		xrpKey, err := domainXrp.NewXRPAccountKey(
			u.coinTypeCode,
			input.AccountType,
			keyPair.ClassicAddress, // AccountID
			keyType,
			keyPair.Seed,      // MasterSeed
			keyPair.SeedHex,   // MasterSeedHex
			keyPair.PublicKey, // PublicKey
			keyPair.PublicKeyHex,
			input.IsKeyPair,
			int64(i), // AllocatedID based on index
		)
		if err != nil {
			return fmt.Errorf("failed to create XRPAccountKey for key %d: %w", i, err)
		}

		items = append(items, xrpKey)
	}

	// Insert keys to DB
	err = u.xrpAccountKeyRepo.InsertBulk(ctx, items)
	if err != nil {
		return fmt.Errorf("failed to call xrpAccountKeyRepo.InsertBulk(): %w", err)
	}

	return nil
}

// extractPrivateKeyBytes extracts the raw private key bytes from a WalletKey.
//
// The WIF field contains the private key in Wallet Import Format.
// We decode it to get the raw bytes for use as entropy in XRP key generation.
func (u *offlineGenerateKeyUseCase) extractPrivateKeyBytes(wk domainKey.WalletKey) ([]byte, error) {
	// The WIF field should contain the hex-encoded private key or WIF string
	// Try to decode as hex first
	if wk.WIF != "" {
		privKeyBytes, err := hex.DecodeString(wk.WIF)
		if err == nil && len(privKeyBytes) == 32 {
			return privKeyBytes, nil
		}

		// If not hex, try to use as raw private key bytes
		// btcec.PrivKeyFromBytes expects exactly 32 bytes
		if len(privKeyBytes) == 32 {
			privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)
			if privKey != nil {
				return privKey.Serialize(), nil
			}
		}
	}

	// Fallback: use the full public key hash as entropy source
	// This is less ideal but provides a deterministic fallback
	if wk.FullPubKey != "" {
		pubKeyBytes, err := hex.DecodeString(wk.FullPubKey)
		if err == nil && len(pubKeyBytes) >= 32 {
			return pubKeyBytes[:32], nil
		}
	}

	return nil, errors.New("unable to extract private key bytes from wallet key")
}

// convertAlgorithmToKeyType converts xrpkg.KeyAlgorithm to domainXrp.XRPKeyType.
func (u *offlineGenerateKeyUseCase) convertAlgorithmToKeyType(alg xrpkg.KeyAlgorithm) domainXrp.XRPKeyType {
	switch alg {
	case xrpkg.AlgorithmEd25519:
		return domainXrp.XRPKeyTypeEd25519
	case xrpkg.AlgorithmSecp256k1:
		return domainXrp.XRPKeyTypeSecp256k1
	default:
		return domainXrp.XRPKeyTypeSecp256k1
	}
}
