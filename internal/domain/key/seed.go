package key

import (
	"errors"
	"time"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// Seed represents an encrypted seed for HD wallet key generation in the domain layer.
//
// SECURITY NOTE: The Seed field contains encrypted seed data and must NEVER be logged
// or exposed in error messages. Handle with extreme care as this is the root of all
// derived keys in the HD wallet hierarchy.
type Seed struct {
	ID           int8
	CoinTypeCode domainCoin.CoinTypeCode
	Seed         string // Encrypted seed data - NEVER log this field
	UpdatedAt    *time.Time
}

// NewSeed creates a new Seed entity.
//
// SECURITY: This function expects the seed parameter to already be encrypted.
// Never pass plaintext seed data to this function.
func NewSeed(
	coinTypeCode domainCoin.CoinTypeCode,
	encryptedSeed string,
) (*Seed, error) {
	if coinTypeCode == "" {
		return nil, errors.New("coin type code cannot be empty")
	}
	if encryptedSeed == "" {
		return nil, errors.New("encrypted seed cannot be empty")
	}

	now := time.Now()
	return &Seed{
		CoinTypeCode: coinTypeCode,
		Seed:         encryptedSeed,
		UpdatedAt:    &now,
	}, nil
}

// GetEncryptedSeed returns the encrypted seed.
// SECURITY: This data must remain encrypted and never be logged.
func (s *Seed) GetEncryptedSeed() string {
	return s.Seed
}
