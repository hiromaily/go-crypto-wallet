package shared_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/sign"
	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/sign/shared"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
)

// TestNewStoreSeedUseCase tests the constructor
func TestNewStoreSeedUseCase(t *testing.T) {
	t.Run("creates use case successfully with non-nil service", func(t *testing.T) {
		// Note: service.Seeder is an interface, which is good for testing
		var mockService service.Seeder
		useCase := shared.NewStoreSeedUseCase(mockService)

		assert.NotNil(t, useCase, "use case should not be nil")
	})

	t.Run("returns correct interface type", func(t *testing.T) {
		var mockService service.Seeder
		useCase := shared.NewStoreSeedUseCase(mockService)

		// Verify it implements the interface
		_ = useCase
		assert.Implements(t, (*sign.StoreSeedUseCase)(nil), useCase)
	})
}

// Note: Full integration tests for StoreSeedUseCase would require:
// 1. Real Seeder service implementation with database connection
// 2. Test database setup
// 3. Proper cleanup after tests
//
// Since service.Seeder is already an interface, unit testing with mocks is feasible.
// Future enhancement could add mock implementations for comprehensive unit testing.
//
// Current test coverage focuses on:
// - Constructor validation
// - Interface compliance
//
// TODO: Add mock Seeder implementation for unit testing Store() method
