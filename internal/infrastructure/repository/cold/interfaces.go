package cold

import (
	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlc"
)

// Type aliases for backward compatibility.
// These interfaces have been moved to pkg/application/ports/persistence.

// SeedRepositorier is SeedRepository interface
type SeedRepositorier = persistence.SeedRepositorier

// BtcAccountKeyRepositorier is BtcAccountKeyRepository interface for BTC/BCH
type BtcAccountKeyRepositorier = persistence.BtcAccountKeyRepositorier

// EthAccountKeyRepositorier is EthAccountKeyRepository interface for ETH
type EthAccountKeyRepositorier = persistence.EthAccountKeyRepositorier

// XRPAccountKeyRepositorier is XRPAccountKeyRepository interface
type XRPAccountKeyRepositorier = persistence.XRPAccountKeyRepositorier

// AuthFullPubkeyRepositorier is AuthFullPubkeyRepository interface
type AuthFullPubkeyRepositorier = persistence.AuthFullPubkeyRepositorier

// AuthAccountKeyRepositorier is AuthAccountKeyRepository interface
type AuthAccountKeyRepositorier = persistence.AuthAccountKeyRepositorier

// GetRedeemScriptByAddress returns redeem script by address for BTC
func GetRedeemScriptByAddress(accountKeys []*sqlc.BtcAccountKey, addr string) string {
	for _, val := range accountKeys {
		if val.MultisigAddress == addr {
			return val.RedeemScript
		}
	}
	return ""
}
