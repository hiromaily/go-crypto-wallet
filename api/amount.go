package api

import (
	"github.com/btcsuite/btcutil"
)

//ToSatoshi bitcoinをSatoshiに変換する
func ToSatoshi(f float64) (btcutil.Amount, error) {
	// Amount
	// Satoshiに変換しないといけない
	// 1Satoshi＝0.00000001BTC
	return btcutil.NewAmount(f)
}
