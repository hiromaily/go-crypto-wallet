package address

import (
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

//----------------------------------------------------
// AddrType
//----------------------------------------------------

// AddrType address type for cryptocurrency addresses
type AddrType string

// address type
const (
	AddrTypeLegacy      AddrType = "legacy"
	AddrTypeP2shSegwit  AddrType = "p2sh-segwit"
	AddrTypeBech32      AddrType = "bech32"
	AddrTypeTaproot     AddrType = "taproot"
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
	AddrTypeTaproot:     4,
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

// AddrStatusFromInt8 converts an int8 value to AddrStatus.
// Returns an error if the value is not valid to prevent silent data corruption.
func AddrStatusFromInt8(val int8) (AddrStatus, error) {
	for addrStatus, num := range AddrStatusValue {
		if int8(num) == val {
			return addrStatus, nil
		}
	}
	return "", fmt.Errorf("invalid addr status value: %d", val)
}

// ValidateAddrStatus validates AddrStatus
func ValidateAddrStatus(val string) bool {
	if _, ok := AddrStatusValue[AddrStatus(val)]; ok {
		return true
	}
	return false
}

//----------------------------------------------------
// AddrType to KeyType conversion
//----------------------------------------------------

// ToKeyType derives the corresponding KeyType from AddrType.
// This ensures address_type and key_type are always consistent.
//
// Mapping:
//   - legacy      -> bip44
//   - p2sh-segwit -> bip49
//   - bech32      -> bip84
//   - taproot     -> bip86
//   - bch-cashaddr -> bip44 (BCH uses BIP44 derivation)
//
// Returns an error for unsupported address types.
func (a AddrType) ToKeyType() (key.KeyType, error) {
	switch a {
	case AddrTypeLegacy:
		return key.KeyTypeBIP44, nil
	case AddrTypeP2shSegwit:
		return key.KeyTypeBIP49, nil
	case AddrTypeBech32:
		return key.KeyTypeBIP84, nil
	case AddrTypeTaproot:
		return key.KeyTypeBIP86, nil
	case AddrTypeBCHCashAddr:
		// BCH uses BIP44 derivation path (same as legacy BTC)
		return key.KeyTypeBIP44, nil
	case AddrTypeETH:
		// Ethereum uses BIP44 derivation path
		return key.KeyTypeBIP44, nil
	default:
		return "", fmt.Errorf("unsupported address type for key derivation: %s", a)
	}
}
