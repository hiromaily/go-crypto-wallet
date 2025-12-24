package key

import (
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainKey "github.com/hiromaily/go-crypto-wallet/pkg/domain/key"
)

// Generator is key generator interface
type Generator interface {
	CreateKey(seed []byte, actType domainAccount.AccountType, idxFrom, count uint32) ([]WalletKey, error)
}

// WalletKey keys
// - [BTC] P2PKHAddr is not used anywhere, P2SHSegWitAddr should be used.
// - [BCH] P2SHSegWitAddr is invalid. P2PKHAddr should be used.
//
// Deprecated: Use domain/key.WalletKey
type WalletKey = domainKey.WalletKey
