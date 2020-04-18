package address

//----------------------------------------------------
// AddrType
//----------------------------------------------------

// AddrType address type for bitcoin
type AddrType string

// address type
const (
	AddrTypeLegacy     AddrType = "legacy"
	AddrTypeP2shSegwit AddrType = "p2sh-segwit"
	AddrTypeBech32     AddrType = "bech32"
)

// String
func (a AddrType) String() string {
	return string(a)
}

// AddrTypeValue value
var AddrTypeValue = map[AddrType]uint8{
	AddrTypeLegacy:     0,
	AddrTypeP2shSegwit: 1,
	AddrTypeBech32:     2,
}

//----------------------------------------------------
// AddrStatus
//----------------------------------------------------

// AddrStatus address generation progress for records in database
type AddrStatus string

// address_status
const (
	AddrStatusHDKeyGenerated       AddrStatus = "generated"              // hd_walletによってkeyが生成された
	AddrStatusPrivKeyImported      AddrStatus = "importprivkey"          // importprivkeyが実行された
	AddrStatusPubkeyExported       AddrStatus = "pubkey_exported"        // pubkeyがexportされた(receipt/payment)
	AddrStatusMultiAddressImported AddrStatus = "multi_address_imported" // multiaddがimportされた(receipt/payment)
	AddrStatusAddressExported      AddrStatus = "address_exported"       // addressがexportされた
)

// String
func (a AddrStatus) String() string {
	return string(a)
}

// AddrStatusValue value
var AddrStatusValue = map[AddrStatus]uint8{
	AddrStatusHDKeyGenerated:       0,
	AddrStatusPrivKeyImported:      1,
	AddrStatusPubkeyExported:       2,
	AddrStatusMultiAddressImported: 3,
	AddrStatusAddressExported:      4,
}

// ValidateAddrStatus validates AddrStatus
func ValidateAddrStatus(val string) bool {
	if _, ok := AddrStatusValue[AddrStatus(val)]; ok {
		return true
	}
	return false
}
