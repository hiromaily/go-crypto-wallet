package btc_test

import (
	"testing"

	"github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/btc"
)

func TestNewGenerateMuSig2NonceUseCase(t *testing.T) {
	t.Parallel()

	// Create use case with nil dependencies (constructor test only)
	useCase := btc.NewGenerateMuSig2NonceUseCase(nil, nil)

	if useCase == nil {
		t.Error("NewGenerateMuSig2NonceUseCase returned nil")
	}
}
