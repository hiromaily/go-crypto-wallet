package xrp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ─── KeyType ──────────────────────────────────────────────────────────────────

func TestKeyType_String(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "secp256k1", KeyTypeSECP256K1.String())
	assert.Equal(t, "ed25519", KeyTypeED25519.String())
}
