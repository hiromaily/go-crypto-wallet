package key

// WalletKey Walletのキーペア
type WalletKey struct {
	WIF            string
	Address        string
	EncodedAddress string
	FullPubKey     string
}
