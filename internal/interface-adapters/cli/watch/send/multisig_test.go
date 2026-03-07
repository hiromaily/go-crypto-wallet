package send

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseSignerEntries_ValidInput verifies parsing of a valid comma-separated signer list.
func TestParseSignerEntries_ValidInput(t *testing.T) {
	t.Parallel()

	entries, err := parseSignerEntries("rSigner1:1,rSigner2:2")
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, "rSigner1", entries[0].Account)
	assert.Equal(t, uint32(1), entries[0].Weight)
	assert.Equal(t, "rSigner2", entries[1].Account)
	assert.Equal(t, uint32(2), entries[1].Weight)
}

// TestParseSignerEntries_SingleEntry verifies that a single signer is valid.
func TestParseSignerEntries_SingleEntry(t *testing.T) {
	t.Parallel()

	entries, err := parseSignerEntries("rSolo:3")
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, "rSolo", entries[0].Account)
	assert.Equal(t, uint32(3), entries[0].Weight)
}

// TestParseSignerEntries_EmptyString verifies that an empty string returns an error.
func TestParseSignerEntries_EmptyString(t *testing.T) {
	t.Parallel()

	_, err := parseSignerEntries("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no valid signers")
}

// TestParseSignerEntries_InvalidFormat verifies that a malformed entry returns an error.
func TestParseSignerEntries_InvalidFormat(t *testing.T) {
	t.Parallel()

	_, err := parseSignerEntries("rSigner1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid signer format")
}

// TestParseSignerEntries_InvalidWeight verifies that a non-integer weight returns an error.
func TestParseSignerEntries_InvalidWeight(t *testing.T) {
	t.Parallel()

	_, err := parseSignerEntries("rSigner1:abc")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid weight")
}

// TestParseSignerEntries_WithSpaces verifies that whitespace around entries is trimmed.
func TestParseSignerEntries_WithSpaces(t *testing.T) {
	t.Parallel()

	entries, err := parseSignerEntries(" rA:1 , rB:2 ")
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, "rA", entries[0].Account)
	assert.Equal(t, "rB", entries[1].Account)
}
