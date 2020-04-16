package key

import (
	"github.com/hiromaily/go-bitcoin/pkg/account"
)

type Generator interface {
	CreateKey(seed []byte, actType account.AccountType, idxFrom, count uint32) ([]WalletKey, error)
}

// WalletKey keys
type WalletKey struct {
	WIF          string
	Address      string
	P2shSegwit   string
	FullPubKey   string
	RedeemScript string
}
