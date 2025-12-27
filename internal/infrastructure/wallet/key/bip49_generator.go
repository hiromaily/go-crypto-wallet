package key

import (
	"github.com/btcsuite/btcd/chaincfg"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
)

// BIP49Generator implements Generator interface for BIP49 (P2SH-SegWit addresses)
type BIP49Generator struct {
	hdKey *HDKey
}

// NewBIP49Generator returns BIP49Generator
func NewBIP49Generator(coinTypeCode domainCoin.CoinTypeCode, conf *chaincfg.Params) *BIP49Generator {
	return &BIP49Generator{
		hdKey: NewHDKey(PurposeTypeBIP49, coinTypeCode, conf),
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

// SupportsAddressType checks if this generator supports the given address type
func (*BIP49Generator) SupportsAddressType(addrType address.AddrType) bool {
	return addrType == address.AddrTypeP2shSegwit
}

// GetDerivationPath returns the BIP49 derivation path
func (g *BIP49Generator) GetDerivationPath(accountType domainAccount.AccountType, index uint32) string {
	return g.hdKey.GetDerivationPath(accountType, index)
}
