package btc

// AddressType represents Bitcoin address type for Bitcoin Core RPC communication
type AddressType string

// Address types
// Note: These values are used in Bitcoin Core RPC calls and must match Bitcoin Core's expectations.
// User-facing configuration uses domain/address.AddrType with "taproot" instead of "bech32m".
const (
	AddressTypeLegacy     AddressType = "legacy"      // P2PKH
	AddressTypeP2SHSegwit AddressType = "p2sh-segwit" // P2SH-wrapped SegWit
	AddressTypeBech32     AddressType = "bech32"      // Native SegWit (P2WPKH)
	AddressTypeTaproot    AddressType = "bech32m"     // Taproot (P2TR) - Bitcoin Core uses "bech32m"
)

// String returns string representation
func (a AddressType) String() string {
	return string(a)
}
