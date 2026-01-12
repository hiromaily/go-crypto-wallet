package bitcoin

// AddressType represents Bitcoin address type
type AddressType string

// Address types
const (
	AddressTypeLegacy     AddressType = "legacy"      // P2PKH
	AddressTypeP2SHSegwit AddressType = "p2sh-segwit" // P2SH-wrapped SegWit
	AddressTypeBech32     AddressType = "bech32"      // Native SegWit (P2WPKH)
	AddressTypeTaproot    AddressType = "taproot"     // Taproot (P2TR) - uses bech32m encoding
)

// String returns string representation
func (a AddressType) String() string {
	return string(a)
}
