package cold

import (
	"context"

	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

// SeedRepositorier is SeedRepository interface
type SeedRepositorier interface {
	GetOne(ctx context.Context) (*domainKey.Seed, error)
	Insert(ctx context.Context, strSeed string) error
}
