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

// BIP86Generator implements Generator interface for BIP86 (Taproot addresses)
type BIP86Generator struct {
	coinType     domainCoin.CoinType
	coinTypeCode domainCoin.CoinTypeCode
	conf         *chaincfg.Params
}

// NewBIP86Generator returns BIP86Generator
func NewBIP86Generator(coinTypeCode domainCoin.CoinTypeCode, conf *chaincfg.Params) *BIP86Generator {
	return &BIP86Generator{
		coinType:     domainCoin.GetCoinType(coinTypeCode, conf),
		coinTypeCode: coinTypeCode,
		conf:         conf,
	}
}

// KeyType returns the key type this generator supports
func (*BIP86Generator) KeyType() domainKey.KeyType {
	return domainKey.KeyTypeBIP86
}

// CreateKey creates keys based on BIP86 standard
// TODO: Implement BIP86 key generation for Taproot addresses
func (*BIP86Generator) CreateKey(
	_ []byte,
	_ domainAccount.AccountType,
	_, _ uint32,
) ([]domainKey.WalletKey, error) {
	return nil, errors.New("BIP86 (Taproot) key generation not yet implemented")
}

// SupportsAddressType checks if this generator supports the given address type
func (*BIP86Generator) SupportsAddressType(addrType address.AddrType) bool {
	return addrType == address.AddrTypeTaproot
}

// GetDerivationPath returns the BIP86 derivation path
func (g *BIP86Generator) GetDerivationPath(accountType domainAccount.AccountType, index uint32) string {
	return fmt.Sprintf("m/86'/%d'/%d'/0/%d",
		g.coinType.Uint32(),
		accountType.Uint32(),
		index)
}
