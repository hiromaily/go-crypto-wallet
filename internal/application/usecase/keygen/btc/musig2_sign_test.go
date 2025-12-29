package btc_test

import (
	"testing"

	"github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/btc"
)

func TestNewMuSig2SignUseCase(t *testing.T) {
	t.Parallel()

	// Create use case with nil dependencies (constructor test only)
	useCase := btc.NewMuSig2SignUseCase(nil, nil)

	if useCase == nil {
		t.Error("NewMuSig2SignUseCase returned nil")
	}
}
