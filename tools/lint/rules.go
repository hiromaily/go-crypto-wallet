package gorules

import "github.com/quasilyte/go-ruleguard/dsl"

// forbidCryptoNaming detects incorrect cryptocurrency naming conventions.
// Cryptocurrency names like BTC, ETH, XRP should be written in uppercase,
// not as Btc, Eth, Xrp.
//
//nolint:unused // This function is used by ruleguard, not by Go code
func forbidCryptoNaming(m dsl.Matcher) {
	// Detect variable declarations starting with Btc, Eth, or Xrp
	m.Match(`var $name $_`).
		Where(m["name"].Text.Matches(`^(Btc|Eth|Xrp)`)).
		Report("$name should use BTC, ETH, or XRP instead of Btc, Eth, or Xrp")

	// Detect short variable declarations
	m.Match(`$name := $_`).
		Where(m["name"].Text.Matches(`^(Btc|Eth|Xrp)`)).
		Report("$name should use BTC, ETH, or XRP instead of Btc, Eth, or Xrp")

	// Detect function names
	m.Match(`func $name($*_)`).
		Where(m["name"].Text.Matches(`^(Btc|Eth|Xrp)`)).
		Report("$name should use BTC, ETH, or XRP instead of Btc, Eth, or Xrp")

	// Detect type names
	m.Match(`type $name $_`).
		Where(m["name"].Text.Matches(`^(Btc|Eth|Xrp)`)).
		Report("$name should use BTC, ETH, or XRP instead of Btc, Eth, or Xrp")
}
