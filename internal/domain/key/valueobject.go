package key

// WalletKey represents a complete set of keys and addresses for a wallet.
//
// This is a value object containing all forms of addresses derived from a single private key:
//   - WIF: Wallet Import Format (private key)
//   - P2PKHAddr: Pay-to-Public-Key-Hash address (legacy)
//   - P2SHSegWitAddr: Pay-to-Script-Hash SegWit address
//   - Bech32Addr: Native SegWit address
//   - TaprootAddr: Taproot address (BIP86)
//   - FullPubKey: Full public key
//   - RedeemScript: Redeem script for multisig
//
// Notes:
//   - For BTC: P2SHSegWitAddr should be used (not P2PKHAddr)
//   - For BCH: P2PKHAddr should be used (P2SHSegWitAddr is invalid)
type WalletKey struct {
	WIF            string
	P2PKHAddr      string
	P2SHSegWitAddr string
	Bech32Addr     string
	TaprootAddr    string
	FullPubKey     string
	RedeemScript   string
}
