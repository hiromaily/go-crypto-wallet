package wallet

// WalletType represents the type of wallet in the cryptocurrency wallet system.
//
// The system uses three different wallet types based on their security model:
//   - Watch-only wallet: Online, holds public keys only, creates and sends transactions
//   - Keygen wallet: Offline, generates keys, provides first signature for multisig
//   - Sign wallet: Offline, provides second and subsequent signatures for multisig
type WalletType string

// Wallet type constants
const (
	// WalletTypeWatchOnly is an online wallet that holds public keys only.
	// It creates and sends transactions but cannot sign them.
	WalletTypeWatchOnly WalletType = "watch"

	// WalletTypeKeyGen is an offline wallet that generates keys and provides
	// the first signature for multisig transactions.
	WalletTypeKeyGen WalletType = "keygen"

	// WalletTypeSign is an offline wallet that provides the second and
	// subsequent signatures for multisig transactions.
	WalletTypeSign WalletType = "sign"
)

// String returns the string representation of the wallet type.
func (w WalletType) String() string {
	return string(w)
}

// WalletTypeValue provides numeric values for wallet types.
// These values are used for database storage and comparison.
var WalletTypeValue = map[WalletType]uint8{
	WalletTypeWatchOnly: 1,
	WalletTypeKeyGen:    2,
	WalletTypeSign:      3,
}

// Uint8 returns the numeric value of the wallet type.
func (w WalletType) Uint8() uint8 {
	return WalletTypeValue[w]
}

// ValidateWalletType validates that the given wallet type is valid.
func ValidateWalletType(wt WalletType) bool {
	_, ok := WalletTypeValue[wt]
	return ok
}
