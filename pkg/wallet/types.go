package wallet

//WalletType wallet type
type WalletType string

// wallet_type
const (
	WalletTypeWatchOnly WalletType = "watch"
	WalletTypeKeyGen    WalletType = "keygen"
	WalletTypeSign      WalletType = "sign"
)

// String converter
func (w WalletType) String() string {
	return string(w)
}

// WalletTypeValue value
var WalletTypeValue = map[WalletType]uint8{
	WalletTypeWatchOnly: 1,
	WalletTypeKeyGen:    2,
	WalletTypeSign:      3,
}
