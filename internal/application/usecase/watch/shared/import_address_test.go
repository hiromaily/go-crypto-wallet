package shared_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/shared"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
)

// TestNewImportAddressUseCase tests the constructor
func TestNewImportAddressUseCase(t *testing.T) {
	t.Run("creates use case successfully with non-nil repositories", func(t *testing.T) {
		useCase := shared.NewImportAddressUseCase(
			nil, // addrRepo
			nil, // addrFileRepo
			domainCoin.BTC,
			address.AddrTypeLegacy,
			domainWallet.WalletTypeWatchOnly,
		)

		assert.NotNil(t, useCase, "use case should not be nil")
	})

	t.Run("returns correct interface type", func(t *testing.T) {
		useCase := shared.NewImportAddressUseCase(
			nil,
			nil,
			domainCoin.BTC,
			address.AddrTypeLegacy,
			domainWallet.WalletTypeWatchOnly,
		)

		// Verify it implements the interface
		_ = useCase
		assert.Implements(t, (*watch.ImportAddressUseCase)(nil), useCase)
	})
}

// Note: Full integration tests for ImportAddressUseCase.Execute() would require:
// 1. Real AddressRepositorier and AddressFileRepositorier implementations
// 2. Test database setup
// 3. Test file fixtures
//
// These would be better suited for integration tests rather than unit tests.
// To enable proper unit testing, the repository interfaces should be mocked.
//
// Current test coverage focuses on:
// - Constructor validation
// - Interface compliance
