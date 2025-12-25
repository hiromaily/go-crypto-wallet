package address

//----------------------------------------------------
// AddrType
//----------------------------------------------------

// AddrType address type for bitcoin
type AddrType string

// address type
const (
	AddrTypeLegacy      AddrType = "legacy"
	AddrTypeP2shSegwit  AddrType = "p2sh-segwit"
	AddrTypeBech32      AddrType = "bech32"
	AddrTypeBCHCashAddr AddrType = "bch-cashaddr"
	AddrTypeETH         AddrType = "eth-address"
)

// String converter
func (a AddrType) String() string {
	return string(a)
}

// AddrTypeValue value
var AddrTypeValue = map[AddrType]uint8{
	AddrTypeLegacy:      0,
	AddrTypeP2shSegwit:  1,
	AddrTypeBech32:      2,
	AddrTypeBCHCashAddr: 3,
}

//----------------------------------------------------
// AddrStatus
//----------------------------------------------------

// AddrStatus address generation progress for records in database
type AddrStatus string

// address_status for keygen wallet
const (
	AddrStatusHDKeyGenerated           AddrStatus = "hdkey_generated"
	AddrStatusPrivKeyImported          AddrStatus = "privkey_imported"
	AddrStatusMultisigAddressGenerated AddrStatus = "multisig_address_generated"
	AddrStatusAddressExported          AddrStatus = "address_exported"
)

// String converter
func (a AddrStatus) String() string {
	return string(a)
}

// Int8 converter
func (a AddrStatus) Int8() int8 {
	return int8(AddrStatusValue[a])
}

// AddrStatusValue value
var AddrStatusValue = map[AddrStatus]uint8{
	AddrStatusHDKeyGenerated:           0,
	AddrStatusPrivKeyImported:          1,
	AddrStatusMultisigAddressGenerated: 2,
	AddrStatusAddressExported:          3,
}

// ValidateAddrStatus validates AddrStatus
func ValidateAddrStatus(val string) bool {
	if _, ok := AddrStatusValue[AddrStatus(val)]; ok {
		return true
	}
	return false
}
