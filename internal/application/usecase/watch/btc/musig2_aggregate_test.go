package btc_test

import (
	"testing"

	"github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/btc"
)

func TestNewAggregateMuSig2SignaturesUseCase(t *testing.T) {
	t.Parallel()

	// Create use case with nil dependencies (constructor test only)
	useCase := btc.NewAggregateMuSig2SignaturesUseCase(nil, nil)

	if useCase == nil {
		t.Error("NewAggregateMuSig2SignaturesUseCase returned nil")
	}
}
