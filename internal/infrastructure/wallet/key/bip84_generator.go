package key

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
)

// BIP84Generator implements Generator interface for BIP84 (Native SegWit Bech32 addresses)
type BIP84Generator struct {
	coinType     domainCoin.CoinType
	coinTypeCode domainCoin.CoinTypeCode
	conf         *chaincfg.Params
}

// NewBIP84Generator returns BIP84Generator
func NewBIP84Generator(coinTypeCode domainCoin.CoinTypeCode, conf *chaincfg.Params) *BIP84Generator {
	return &BIP84Generator{
		coinType:     domainCoin.GetCoinType(coinTypeCode, conf),
		coinTypeCode: coinTypeCode,
		conf:         conf,
	}
}

// KeyType returns the key type this generator supports
func (*BIP84Generator) KeyType() domainKey.KeyType {
	return domainKey.KeyTypeBIP84
}

// CreateKey creates keys based on BIP84 standard
// TODO: Implement BIP84 key generation for Native SegWit (Bech32) addresses
func (*BIP84Generator) CreateKey(
	_ []byte,
	_ domainAccount.AccountType,
	_, _ uint32,
) ([]domainKey.WalletKey, error) {
	return nil, errors.New("BIP84 key generation not yet implemented")
}

// SupportsAddressType checks if this generator supports the given address type
func (*BIP84Generator) SupportsAddressType(addrType address.AddrType) bool {
	return addrType == address.AddrTypeBech32
}

// GetDerivationPath returns the BIP84 derivation path
func (g *BIP84Generator) GetDerivationPath(accountType domainAccount.AccountType, index uint32) string {
	return fmt.Sprintf("m/84'/%d'/%d'/0/%d",
		g.coinType.Uint32(),
		accountType.Uint32(),
		index)
}
