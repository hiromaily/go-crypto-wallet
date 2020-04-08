package wallet

//WalletType Walletタイプ
type WalletType string

// wallet_type
const (
	WalletTypeWatchOnly WalletType = "watch_only"
	WalletTypeKeyGen    WalletType = "keygen"
	WalletTypeSignature WalletType = "signature"
)

//WalletTypeValue env_typeの値
var WalletTypeValue = map[WalletType]uint8{
	WalletTypeWatchOnly: 1,
	WalletTypeKeyGen:    2,
	WalletTypeSignature: 3,
}
