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

	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	xrpkg "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp"
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
// without requiring any API calls to external services.
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
	items := make([]*domainXRP.XRPAccountKey, 0, len(walletKeys))
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
		xrpKey, err := domainXRP.NewXRPAccountKey(
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
// The WIF field is expected to contain either:
// - A hex-encoded 32-byte private key
// - Raw private key bytes that can be parsed by btcec
//
// SECURITY NOTE: This function will fail explicitly if no valid private key
// can be extracted. We do NOT fall back to public key data as that would be
// cryptographically insecure.
func (*offlineGenerateKeyUseCase) extractPrivateKeyBytes(wk domainKey.WalletKey) ([]byte, error) {
	if wk.WIF == "" {
		return nil, errors.New("WIF field is empty: no private key available")
	}

	// Try to decode as hex-encoded private key
	privKeyBytes, err := hex.DecodeString(wk.WIF)
	if err == nil && len(privKeyBytes) == 32 {
		return privKeyBytes, nil
	}

	// If hex decode failed or wrong length, the WIF field format is not supported
	return nil, fmt.Errorf(
		"invalid private key format in WIF field: expected 32-byte hex-encoded key, got %d bytes",
		len(privKeyBytes),
	)
}

// convertAlgorithmToKeyType converts xrpkg.KeyAlgorithm to domainXRP.XRPKeyType.
func (*offlineGenerateKeyUseCase) convertAlgorithmToKeyType(alg xrpkg.KeyAlgorithm) domainXRP.XRPKeyType {
	switch alg {
	case xrpkg.AlgorithmEd25519:
		return domainXRP.XRPKeyTypeEd25519
	case xrpkg.AlgorithmSecp256k1:
		return domainXRP.XRPKeyTypeSecp256k1
	default:
		return domainXRP.XRPKeyTypeSecp256k1
	}
}
