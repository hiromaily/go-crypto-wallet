package address

//----------------------------------------------------
// AddressType
//----------------------------------------------------

//AddressType address type for bitcoin
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

// AddressTypeValue value
var AddressTypeValue = map[AddressType]uint8{
	AddressTypeLegacy:     0,
	AddressTypeP2shSegwit: 1,
	AddressTypeBech32:     2,
}

//----------------------------------------------------
// AddressStatus
//----------------------------------------------------

// FIXME: name should be just Status
// AddressStatus address generation progress for records in database
type AddressStatus string

// address_status
const (
	AddressStatusHDKeyGenerated       AddressStatus = "generated"              // hd_walletによってkeyが生成された
	AddressStatusPrivKeyImported      AddressStatus = "importprivkey"          // importprivkeyが実行された
	AddressStatusPubkeyExported       AddressStatus = "pubkey_exported"        // pubkeyがexportされた(receipt/payment)
	AddressStatusMultiAddressImported AddressStatus = "multi_address_imported" // multiaddがimportされた(receipt/payment)
	AddressStatusAddressExported      AddressStatus = "address_exported"       // addressがexportされた
)

func (a AddressStatus) String() string {
	return string(a)
}

// AddressStatusValue value
var AddressStatusValue = map[AddressStatus]uint8{
	AddressStatusHDKeyGenerated:       0,
	AddressStatusPrivKeyImported:      1,
	AddressStatusPubkeyExported:       2,
	AddressStatusMultiAddressImported: 3,
	AddressStatusAddressExported:      4,
}

// ValidateAddressStatus validates AddressStatus
func ValidateAddressStatus(val string) bool {
	if _, ok := AddressStatusValue[AddressStatus(val)]; ok {
		return true
	}
	return false
}
