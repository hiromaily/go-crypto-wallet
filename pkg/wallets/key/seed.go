package key

import (
	"encoding/base64"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/pkg/errors"
)

// GenerateSeed seedを生成する []byte
func GenerateSeed() ([]byte, error) {
	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		return nil, err
	}

	return seed, nil
}

// SeedToString stringにエンコードする
func SeedToString(seed []byte) string {
	base64seed := base64.StdEncoding.EncodeToString(seed)
	//logger.Debug("generated seed(string):", base64seed)

	return base64seed
}

// SeedToByte byte型にデコードする
func SeedToByte(seed string) ([]byte, error) {
	unbase64, err := base64.StdEncoding.DecodeString(seed)
	if err != nil {
		return nil, errors.Errorf("[Error] base64.StdEncoding.DecodeString(): error: %v", err)
	}
	return unbase64, nil
}
