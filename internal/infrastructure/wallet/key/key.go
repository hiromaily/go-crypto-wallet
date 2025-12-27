package key

import (
	"github.com/btcsuite/btcd/chaincfg"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
)

// Generator is key generator interface
type Generator interface {
	// KeyType returns the key type this generator supports
	KeyType() domainKey.KeyType

	// CreateKey creates keys based on the seed and account type
	CreateKey(seed []byte, actType domainAccount.AccountType, idxFrom, count uint32) ([]domainKey.WalletKey, error)

	// SupportsAddressType checks if this generator supports the given address type
	SupportsAddressType(addrType address.AddrType) bool

	// GetDerivationPath returns the derivation path for the given account and index
	GetDerivationPath(accountType domainAccount.AccountType, index uint32) string
}

// GeneratorFactory creates a Generator based on key type
type GeneratorFactory interface {
	CreateGenerator(
		keyType domainKey.KeyType,
		coinTypeCode domainCoin.CoinTypeCode,
		conf *chaincfg.Params,
	) (Generator, error)
}
