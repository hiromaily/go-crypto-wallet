package shared_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch/shared"
	sharedwatchsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/watch/shared"
)

// TestNewImportAddressUseCase tests the constructor
func TestNewImportAddressUseCase(t *testing.T) {
	t.Run("creates use case successfully with non-nil service", func(t *testing.T) {
		mockService := &sharedwatchsrv.AddressImport{}
		useCase := shared.NewImportAddressUseCase(mockService)

		assert.NotNil(t, useCase, "use case should not be nil")
	})

	t.Run("returns correct interface type", func(t *testing.T) {
		mockService := &sharedwatchsrv.AddressImport{}
		useCase := shared.NewImportAddressUseCase(mockService)

		// Verify it implements the interface
		_ = useCase
		assert.Implements(t, (*watch.ImportAddressUseCase)(nil), useCase)
	})
}

// Note: Full integration tests for ImportAddressUseCase.Execute() would require:
// 1. Real AddressImport service with all dependencies (repository, file system)
// 2. Test database setup
// 3. Test file fixtures
//
// These would be better suited for integration tests rather than unit tests.
// To enable proper unit testing, AddressImport should be refactored to an interface
// that can be mocked.
//
// Current test coverage focuses on:
// - Constructor validation
// - Interface compliance
//
// TODO: Consider refactoring services to interfaces to enable proper unit testing
