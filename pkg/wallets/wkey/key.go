package wkey

// WalletKey Walletのキーペア
type WalletKey struct {
	WIF          string
	Address      string
	P2shSegwit   string
	FullPubKey   string
	RedeemScript string
}
