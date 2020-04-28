package key

import (
	"github.com/hiromaily/go-bitcoin/pkg/account"
)

// Generator is key generator interface
type Generator interface {
	CreateKey(seed []byte, actType account.AccountType, idxFrom, count uint32) ([]WalletKey, error)
}

// WalletKey keys
// - P2PKHAddr is not used anywhere, P2SHSegWitAddr should be used.
type WalletKey struct {
	WIF            string
	P2PKHAddr      string
	P2SHSegWitAddr string
	FullPubKey     string
	RedeemScript   string
}
