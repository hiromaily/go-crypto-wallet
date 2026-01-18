package generator

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	infraKey "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"

	portsWallet "github.com/hiromaily/go-crypto-wallet/internal/application/ports/wallet"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key/strategy"
)

// BIP49Generator implements Generator interface for BIP49 (P2SH-SegWit addresses)
type BIP49Generator struct {
	hdKey *infraKey.HDKey
}

// Compile-time check that generator implements Generator interface
var _ portsWallet.Generator = (*BIP49Generator)(nil)

// NewBIP49Generator returns BIP49Generator
func NewBIP49Generator(coinTypeCode domainCoin.CoinTypeCode, conf *chaincfg.Params) *BIP49Generator {
	// Create coin-specific strategy
	coinStrategy, err := strategy.CreateCoinKeyStrategy(coinTypeCode, conf)
	if err != nil {
		// Panic is acceptable in factory/initialization code for unsupported configurations
		panic(fmt.Sprintf("failed to create coin strategy for %s: %v", coinTypeCode.String(), err))
	}

	return &BIP49Generator{
		hdKey: infraKey.NewHDKey(infraKey.PurposeTypeBIP49, coinTypeCode, conf, coinStrategy),
	}
}

// KeyType returns the key type this generator supports
func (*BIP49Generator) KeyType() domainKey.KeyType {
	return domainKey.KeyTypeBIP49
}

// CreateKey creates keys based on BIP49 standard
func (g *BIP49Generator) CreateKey(
	seed []byte,
	accountType domainAccount.AccountType,
	idxFrom, count uint32,
) ([]domainKey.WalletKey, error) {
	return g.hdKey.CreateKey(seed, accountType, idxFrom, count)
}

// CreateKeyWithAccountXpriv creates HD keys and returns the account-level extended private key.
// This method delegates to the underlying HDKey implementation.
func (g *BIP49Generator) CreateKeyWithAccountXpriv(
	seed []byte,
	accountType domainAccount.AccountType,
	idxFrom, count uint32,
) ([]domainKey.WalletKey, string, error) {
	return g.hdKey.CreateKeyWithAccountXpriv(seed, accountType, idxFrom, count)
}

// SupportsAddressType checks if this generator supports the given address type
func (*BIP49Generator) SupportsAddressType(addrType domainAddress.AddrType) bool {
	return addrType == domainAddress.AddrTypeP2shSegwit
}

// GetDerivationPath returns the BIP49 derivation path
func (g *BIP49Generator) GetDerivationPath(accountType domainAccount.AccountType, index uint32) string {
	return g.hdKey.GetDerivationPath(accountType, index)
}
