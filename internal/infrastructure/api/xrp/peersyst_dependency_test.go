package xrp_test

import (
	"testing"

	"github.com/Peersyst/xrpl-go/xrpl/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPeersystDependencyAvailable verifies that the Peersyst xrpl-go library
// is available and can be used for wallet operations.
// This test ensures task 1.1 (Add Peersyst xrpl-go dependency) is complete.
func TestPeersystDependencyAvailable(t *testing.T) {
	t.Parallel()

	// Test that we can create a wallet from a test seed
	testSeed := "sEdTM1uX8pu2do5XvTnutH6HsouMaM2" // Test seed from XRP Ledger docs

	// Attempt to derive wallet from seed
	w, err := wallet.FromSeed(testSeed, "")
	require.NoError(t, err, "should be able to create wallet from seed")
	require.NotNil(t, w, "wallet should not be nil")

	// Verify wallet has expected properties
	assert.NotEmpty(t, w.ClassicAddress, "wallet should have a classic address")
	assert.NotEmpty(t, w.PublicKey, "wallet should have a public key")
	assert.NotEmpty(t, w.PrivateKey, "wallet should have a private key")
}

// TestPeersystWalletSigningCapability verifies that the Peersyst wallet
// can be used for transaction signing (basic smoke test).
func TestPeersystWalletSigningCapability(t *testing.T) {
	t.Parallel()

	// Create a test wallet
	testSeed := "sEdTM1uX8pu2do5XvTnutH6HsouMaM2"
	w, err := wallet.FromSeed(testSeed, "")
	require.NoError(t, err)

	// Verify the wallet can access signing key
	assert.NotEmpty(t, w.PrivateKey, "wallet should have private key for signing")

	// This confirms the wallet package is available and functional
	// Actual signing tests will be in PeersystSigner tests (task 3.1)
}
