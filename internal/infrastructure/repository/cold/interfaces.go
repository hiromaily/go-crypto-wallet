package cold

import (
	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
)

// Type aliases for backward compatibility.
// These interfaces have been moved to pkg/application/ports/persistence.

// SeedRepositorier is SeedRepository interface
type SeedRepositorier = persistence.SeedRepositorier

// AccountKeyRepositorier is AccountKeyRepository interface
type AccountKeyRepositorier = persistence.AccountKeyRepositorier

// XRPAccountKeyRepositorier is XRPAccountKeyRepository interface
type XRPAccountKeyRepositorier = persistence.XRPAccountKeyRepositorier

// AuthFullPubkeyRepositorier is AuthFullPubkeyRepository interface
type AuthFullPubkeyRepositorier = persistence.AuthFullPubkeyRepositorier

// AuthAccountKeyRepositorier is AuthAccountKeyRepository interface
type AuthAccountKeyRepositorier = persistence.AuthAccountKeyRepositorier

// GetRedeemScriptByAddress returns redeem script by address
func GetRedeemScriptByAddress(accountKeys []*models.AccountKey, addr string) string {
	for _, val := range accountKeys {
		if val.MultisigAddress == addr {
			return val.RedeemScript
		}
	}
	return ""
}
