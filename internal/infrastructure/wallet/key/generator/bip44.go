package generator

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"

	portsWallet "github.com/hiromaily/go-crypto-wallet/internal/application/ports/wallet"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	infraKey "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key/strategy"
)

// BIP44Generator implements Generator interface for BIP44 (Legacy P2PKH addresses)
type BIP44Generator struct {
	hdKey *infraKey.HDKey
}

// Compile-time check that generator implements Generator interface
var _ portsWallet.Generator = (*BIP44Generator)(nil)

// NewBIP44Generator returns BIP44Generator
func NewBIP44Generator(coinTypeCode domainCoin.CoinTypeCode, conf *chaincfg.Params) *BIP44Generator {
	// Create coin-specific strategy
	coinStrategy, err := strategy.CreateCoinKeyStrategy(coinTypeCode, conf)
	if err != nil {
		// Panic is acceptable in factory/initialization code for unsupported configurations
		panic(fmt.Sprintf("failed to create coin strategy for %s: %v", coinTypeCode.String(), err))
	}

	return &BIP44Generator{
		hdKey: infraKey.NewHDKey(infraKey.PurposeTypeBIP44, coinTypeCode, conf, coinStrategy),
	}
}

// KeyType returns the key type this generator supports
func (*BIP44Generator) KeyType() domainKey.KeyType {
	return domainKey.KeyTypeBIP44
}

// CreateKey creates keys based on BIP44 standard
func (g *BIP44Generator) CreateKey(
	seed []byte,
	accountType domainAccount.AccountType,
	idxFrom, count uint32,
) ([]domainKey.WalletKey, error) {
	return g.hdKey.CreateKey(seed, accountType, idxFrom, count)
}

// CreateKeyWithAccountXpriv creates HD keys and returns the account-level extended private key.
// This method delegates to the underlying HDKey implementation.
func (g *BIP44Generator) CreateKeyWithAccountXpriv(
	seed []byte,
	accountType domainAccount.AccountType,
	idxFrom, count uint32,
) ([]domainKey.WalletKey, string, error) {
	return g.hdKey.CreateKeyWithAccountXpriv(seed, accountType, idxFrom, count)
}

// SupportsAddressType checks if this generator supports the given address type
func (*BIP44Generator) SupportsAddressType(addrType domainAddress.AddrType) bool {
	return addrType == domainAddress.AddrTypeLegacy
}

// GetDerivationPath returns the BIP44 derivation path
func (g *BIP44Generator) GetDerivationPath(accountType domainAccount.AccountType, index uint32) string {
	return g.hdKey.GetDerivationPath(accountType, index)
}
