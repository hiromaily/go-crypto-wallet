package mysql

import (
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
)

// GetRedeemScriptByAddress returns redeem script by address for BTC
func GetRedeemScriptByAddress(accountKeys []*sqlcgen.BtcAccountKey, addr string) string {
	for _, val := range accountKeys {
		if val.MultisigAddress == addr {
			return val.RedeemScript
		}
	}
	return ""
}
