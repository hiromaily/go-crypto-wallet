package btc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	"github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign/btc"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
)

// TestNewSignTransactionUseCase tests the constructor
func TestNewSignTransactionUseCase(t *testing.T) {
	t.Run("creates use case successfully with nil dependencies", func(t *testing.T) {
		useCase := btc.NewSignTransactionUseCase(
			nil, // btcAPI
			nil, // accountKeyRepo
			nil, // authKeyRepo
			nil, // txFileRepo
			nil, // multisigAccount
			domainWallet.WalletTypeSign,
		)

		assert.NotNil(t, useCase, "use case should not be nil")
	})

	t.Run("returns correct interface type", func(t *testing.T) {
		useCase := btc.NewSignTransactionUseCase(
			nil,
			nil,
			nil,
			nil,
			nil,
			domainWallet.WalletTypeSign,
		)

		// Verify it implements the interface
		assert.Implements(t, (*signusecase.SignTransactionUseCase)(nil), useCase)
	})
}

// Note: Full integration tests for SignTransactionUseCase would require:
// 1. Mock Bitcoin client for PSBT signing (SignPSBTWithKey)
// 2. Mock auth key repository for key retrieval (GetOne)
// 3. Mock file repository for PSBT file reading/writing (ReadPSBTFile, WritePSBTFile, ValidateFilePath)
// 4. Mock multisig account configuration
// 5. Partially signed PSBT files from Keygen wallet output
// 6. Proper cleanup after tests
//
// These would be better suited for integration tests rather than unit tests.
// To enable proper unit testing, all dependencies should be mocked.
//
// Current test coverage focuses on:
// - Constructor validation
// - Interface compliance
//
// PSBT-specific functionality testing:
// The Sign method now processes PSBT files instead of CSV files.
// Integration tests should verify:
// - PSBT file reading (partially signed from Keygen)
// - PSBT parsing with existing signatures
// - Adding second signature to PSBT (auth key from auth_account_key table)
// - Multisig completion detection (2-of-2 or 2-of-N)
// - Signature generation for all multisig address types (P2SH, P2WSH, P2TR script path)
// - Fully signed PSBT file writing with correct naming convention
// - Offline signing with btcd library (no Bitcoin Core RPC dependency)
// - Automatic signature type selection (ECDSA for P2SH/P2WSH, Schnorr for P2TR)
// - signedCount increments correctly for additional signatures needed
//
// End-to-end flow:
// 1. Watch wallet creates unsigned PSBT (0 signatures)
// 2. Keygen wallet adds first signature (1 signature)
// 3. Sign wallet adds second signature (2 signatures, fully signed)
// 4. Watch wallet finalizes and broadcasts
//
// See docs/TESTING_STRATEGY.md for comprehensive testing approach.
