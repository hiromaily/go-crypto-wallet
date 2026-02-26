package btc

import (
	"encoding/base64"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
)

// GenerateSeed generates a new cryptographic seed as []byte.
func GenerateSeed() ([]byte, error) {
	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		return nil, err
	}
	return seed, nil
}

// SeedToString encodes a seed byte slice to a base64 string.
func SeedToString(seed []byte) string {
	return base64.StdEncoding.EncodeToString(seed)
}

// SeedToByte decodes a base64-encoded seed string to bytes.
func SeedToByte(seed string) ([]byte, error) {
	unbase64, err := base64.StdEncoding.DecodeString(seed)
	if err != nil {
		return nil, fmt.Errorf("fail to call base64.StdEncoding.DecodeString(): %w", err)
	}
	return unbase64, nil
}
