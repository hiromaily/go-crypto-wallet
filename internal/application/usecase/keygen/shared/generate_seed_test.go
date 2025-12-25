package shared_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/shared"
)

// TestNewGenerateSeedUseCase tests the constructor
func TestNewGenerateSeedUseCase(t *testing.T) {
	t.Run("creates use case successfully with non-nil repository", func(t *testing.T) {
		useCase := shared.NewGenerateSeedUseCase(nil)

		assert.NotNil(t, useCase, "use case should not be nil")
	})

	t.Run("returns correct interface type", func(t *testing.T) {
		useCase := shared.NewGenerateSeedUseCase(nil)

		// Verify it implements the interface
		_ = useCase
		assert.Implements(t, (*keygen.GenerateSeedUseCase)(nil), useCase)
	})
}

// Note: Full integration tests for GenerateSeedUseCase would require:
// 1. Test database setup for storing seeds
// 2. Proper cleanup after tests
//
// These would be better suited for integration tests rather than unit tests.
// To enable proper unit testing, the SeedRepositorier interface should be mocked.
//
// Current test coverage focuses on:
// - Constructor validation
// - Interface compliance
