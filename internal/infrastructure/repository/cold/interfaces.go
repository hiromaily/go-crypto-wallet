package cold

import (
	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlc"
)

// Type aliases for backward compatibility.
// These interfaces have been moved to pkg/application/ports/persistence.

// SeedRepositorier is SeedRepository interface
type SeedRepositorier = persistence.SeedRepositorier

// BTCAccountKeyRepositorier is BtcAccountKeyRepository interface for BTC/BCH
type BTCAccountKeyRepositorier = persistence.BTCAccountKeyRepositorier

// ETHAccountKeyRepositorier is EthAccountKeyRepository interface for ETH
type ETHAccountKeyRepositorier = persistence.ETHAccountKeyRepositorier

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
