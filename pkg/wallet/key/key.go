package key

import (
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

// Generator is key generator interface
type Generator interface {
	CreateKey(seed []byte, actType domainAccount.AccountType, idxFrom, count uint32) ([]domainKey.WalletKey, error)
}
