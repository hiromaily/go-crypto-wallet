package btc_test

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"

	keygenusecasebtc "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/btc"
	apibtcimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/btc/btc"
)

func TestNewCreateMuSig2AddressUseCase(t *testing.T) {
	t.Parallel()

	// Create dependencies
	musig2Service := apibtcimpl.NewMuSig2Service()
	chainConfig := &chaincfg.RegressionNetParams

	// Test constructor
	useCase := keygenusecasebtc.NewCreateMuSig2AddressUseCase(
		musig2Service,
		chainConfig,
		nil, // authFullPubKeyRepo - nil for constructor test
		nil, // accountKeyRepo - nil for constructor test
		nil, // multisigAccount - nil for constructor test
	)

	if useCase == nil {
		t.Error("expected non-nil use case")
	}
}
