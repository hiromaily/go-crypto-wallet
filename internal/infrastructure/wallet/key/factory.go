package key

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

// Factory creates key generators based on key type
type Factory struct{}

// NewFactory returns Factory
func NewFactory() *Factory {
	return &Factory{}
}

// CreateGenerator creates a generator for the specified key type
func (*Factory) CreateGenerator(
	keyType domainKey.KeyType,
	coinTypeCode domainCoin.CoinTypeCode,
	conf *chaincfg.Params,
) (Generator, error) {
	if err := keyType.Validate(); err != nil {
		return nil, fmt.Errorf("invalid key type: %w", err)
	}

	switch keyType {
	case domainKey.KeyTypeBIP44:
		return NewBIP44Generator(coinTypeCode, conf), nil
	case domainKey.KeyTypeBIP49:
		return NewBIP49Generator(coinTypeCode, conf), nil
	case domainKey.KeyTypeBIP84:
		return NewBIP84Generator(coinTypeCode, conf), nil
	case domainKey.KeyTypeBIP86:
		return NewBIP86Generator(coinTypeCode, conf), nil
	case domainKey.KeyTypeMuSig2:
		return nil, errors.New("MuSig2 key generation not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported key type: %s", keyType)
	}
}
