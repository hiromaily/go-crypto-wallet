package types

//WalletType Walletタイプ
type WalletType string

// wallet_type
const (
	WalletTypeWatchOnly WalletType = "wallet"
	WalletTypeKeyGen    WalletType = "keygen"
	WalletTypeSignature WalletType = "signature"
)

func (w WalletType) String() string {
	return string(w)
}

//WalletTypeValue env_typeの値
var WalletTypeValue = map[WalletType]uint8{
	WalletTypeWatchOnly: 1,
	WalletTypeKeyGen:    2,
	WalletTypeSignature: 3,
}
