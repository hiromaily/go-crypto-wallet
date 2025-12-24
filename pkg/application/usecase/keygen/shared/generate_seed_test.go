package shared_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen/shared"
	sharedkeygensrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/keygen/shared"
)

// TestNewGenerateSeedUseCase tests the constructor
func TestNewGenerateSeedUseCase(t *testing.T) {
	t.Run("creates use case successfully with non-nil service", func(t *testing.T) {
		mockService := &sharedkeygensrv.Seed{}
		useCase := shared.NewGenerateSeedUseCase(mockService)

		assert.NotNil(t, useCase, "use case should not be nil")
	})

	t.Run("returns correct interface type", func(t *testing.T) {
		mockService := &sharedkeygensrv.Seed{}
		useCase := shared.NewGenerateSeedUseCase(mockService)

		// Verify it implements the interface
		_ = useCase
		assert.Implements(t, (*keygen.GenerateSeedUseCase)(nil), useCase)
	})
}

// Note: Full integration tests for GenerateSeedUseCase would require:
// 1. Real Seed service with all dependencies (seeder, repository)
// 2. Test database setup for storing seeds
// 3. Proper cleanup after tests
//
// These would be better suited for integration tests rather than unit tests.
// To enable proper unit testing, Seed service should be refactored to an interface
// that can be mocked.
//
// Current test coverage focuses on:
// - Constructor validation
// - Interface compliance
//
// TODO: Consider refactoring services to interfaces to enable proper unit testing
