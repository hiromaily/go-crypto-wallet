package strategy

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec/v2"
)

// getFullPubKey returns the full public key as a hex string.
// This is shared by BTC and BCH strategies since they use the same key format.
//
// Parameters:
//   - privKey: The private key to extract the public key from
//   - isCompressed: Whether to use compressed (true) or uncompressed (false) format
//
// Returns:
//   - Hex-encoded public key string
func getFullPubKey(privKey *btcec.PrivateKey, isCompressed bool) string {
	var bPubKey []byte
	if isCompressed {
		// Compressed public key (33 bytes)
		bPubKey = privKey.PubKey().SerializeCompressed()
	} else {
		// Uncompressed public key (65 bytes)
		bPubKey = privKey.PubKey().SerializeUncompressed()
	}
	return hex.EncodeToString(bPubKey)
}
