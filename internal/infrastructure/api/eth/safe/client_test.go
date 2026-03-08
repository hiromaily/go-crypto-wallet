package safe_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/safe"
)

// compile-time check: SafeClient must implement SafeClientDeps.
var _ apieth.SafeClientDeps = (*safe.SafeClient)(nil)

func TestNewSafeClientParams_FromField(t *testing.T) {
	from := common.HexToAddress("0x1234567890123456789012345678901234567890")
	params := safe.NewSafeClientParams{From: from}
	if params.From != from {
		t.Errorf("expected From %s, got %s", from.Hex(), params.From.Hex())
	}
}
