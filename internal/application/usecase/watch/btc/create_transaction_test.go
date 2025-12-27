package btc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
)

// TestNewCreateTransactionUseCase tests the constructor
func TestNewCreateTransactionUseCase(t *testing.T) {
	t.Run("creates use case successfully with nil dependencies", func(t *testing.T) {
		useCase := btc.NewCreateTransactionUseCase(
			nil, // btcClient
			nil, // dbConn
			nil, // addrRepo
			nil, // txRepo
			nil, // txInputRepo
			nil, // txOutputRepo
			nil, // payReqRepo
			nil, // txFileRepo
			domainAccount.AccountTypeDeposit,
			domainAccount.AccountTypePayment,
			domainWallet.WalletTypeWatchOnly,
		)

		assert.NotNil(t, useCase, "use case should not be nil")
	})

	t.Run("returns correct interface type", func(t *testing.T) {
		useCase := btc.NewCreateTransactionUseCase(
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			domainAccount.AccountTypeDeposit,
			domainAccount.AccountTypePayment,
			domainWallet.WalletTypeWatchOnly,
		)

		// Verify it implements the interface
		assert.Implements(t, (*watchusecase.CreateTransactionUseCase)(nil), useCase)
	})
}

// Note: Full integration tests for CreateTransactionUseCase would require:
// 1. Mock Bitcoin client for PSBT creation
// 2. Test database setup for transaction storage
// 3. Mock file repository for PSBT file writing
// 4. Proper cleanup after tests
//
// These would be better suited for integration tests rather than unit tests.
// To enable proper unit testing, all dependencies should be mocked.
//
// Current test coverage focuses on:
// - Constructor validation
// - Interface compliance
//
// PSBT-specific functionality testing:
// The generatePSBTFile method now creates PSBT files instead of CSV files.
// Integration tests should verify:
// - PSBT creation for deposit transactions (client → deposit)
// - PSBT creation for payment transactions (payment → user addresses)
// - PSBT creation for transfer transactions (account → account)
// - PSBT includes all required metadata (amounts, scriptPubKeys, redeem scripts)
// - PSBT files follow naming convention: {actionType}_{txID}_{txType}_{signedCount}_{timestamp}.psbt
// - PSBT files are base64-encoded and BIP174 compliant
// - Generated PSBT files can be read by Keygen wallet (offline signing)
//
// See docs/TESTING_STRATEGY.md for comprehensive testing approach.
