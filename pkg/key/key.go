package key

// WalletKey Walletのキーペア
type WalletKey struct {
	WIF        string
	Address    string
	P2shSegwit string
	FullPubKey string
}
