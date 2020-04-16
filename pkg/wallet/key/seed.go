package key

import (
	"encoding/base64"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/pkg/errors"
)

// GenerateSeed generate seed as []byte
func GenerateSeed() ([]byte, error) {
	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		return nil, err
	}
	return seed, nil
}

// SeedToString encode by base64 to string
func SeedToString(seed []byte) string {
	base64seed := base64.StdEncoding.EncodeToString(seed)
	return base64seed
}

// SeedToByte decode string to []byte
func SeedToByte(seed string) ([]byte, error) {
	unbase64, err := base64.StdEncoding.DecodeString(seed)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call base64.StdEncoding.DecodeString()")
	}
	return unbase64, nil
}
