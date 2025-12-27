package btc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/btc"
)

// TestNewSignTransactionUseCase tests the constructor
func TestNewSignTransactionUseCase(t *testing.T) {
	t.Run("creates use case successfully with nil dependencies", func(t *testing.T) {
		useCase := btc.NewSignTransactionUseCase(
			nil, // btc
			nil, // accountKeyRepo
			nil, // txFileRepo
			nil, // multisigAccount
		)

		assert.NotNil(t, useCase, "use case should not be nil")
	})

	t.Run("returns correct interface type", func(t *testing.T) {
		useCase := btc.NewSignTransactionUseCase(
			nil,
			nil,
			nil,
			nil,
		)

		// Verify it implements the interface
		assert.Implements(t, (*keygenusecase.SignTransactionUseCase)(nil), useCase)
	})
}

// Note: Full integration tests for SignTransactionUseCase would require:
// 1. Mock Bitcoin client for PSBT signing (SignPSBTWithKey, ParsePSBT)
// 2. Mock account key repository for key retrieval (GetAllAddrStatus)
// 3. Mock file repository for PSBT file reading/writing (ReadPSBTFile, WritePSBTFile, ValidateFilePath)
// 4. Mock multisig account configuration
// 5. Test PSBT files with various states (unsigned, partially signed, fully signed)
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
// - PSBT file reading and validation
// - PSBT parsing to extract metadata
// - Signature generation for all address types (P2PKH, P2WPKH, P2TR)
// - Multisig partial signature handling (incrementing signedCount)
// - Signed PSBT file writing with correct naming convention
// - Account inference from action type (deposit → client, payment → payment, transfer → payment)
// - WIF retrieval for sender account (GetAllAddrStatus with AddrStatusAddressExported)
// - Offline signing with btcd library (no Bitcoin Core RPC dependency)
// - Automatic signature type selection (ECDSA for legacy/SegWit, Schnorr for Taproot)
//
// See docs/TESTING_STRATEGY.md for comprehensive testing approach.
