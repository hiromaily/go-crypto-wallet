package types

//WalletType wallet type
type WalletType string

// wallet_type
const (
	WalletTypeWatchOnly WalletType = "wallet"
	WalletTypeKeyGen    WalletType = "keygen"
	WalletTypeSignature WalletType = "signature"
)

// String converter
func (w WalletType) String() string {
	return string(w)
}

// WalletTypeValue value
var WalletTypeValue = map[WalletType]uint8{
	WalletTypeWatchOnly: 1,
	WalletTypeKeyGen:    2,
	WalletTypeSignature: 3,
}
