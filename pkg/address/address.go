package address

//AddressType address種別
type AddressType string

// address type
const (
	AddressTypeLegacy     AddressType = "legacy"
	AddressTypeP2shSegwit AddressType = "p2sh-segwit"
	AddressTypeBech32     AddressType = "bech32"
)

func (a AddressType) String() string {
	return string(a)
}

//AddressTypeValue address_typeの値
var AddressTypeValue = map[AddressType]uint8{
	AddressTypeLegacy:     0,
	AddressTypeP2shSegwit: 1,
	AddressTypeBech32:     2,
}
