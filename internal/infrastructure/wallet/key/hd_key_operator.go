package key

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	portsWallet "github.com/hiromaily/go-crypto-wallet/internal/application/ports/wallet"
)

// HDKeyOperatorImpl implements the application ports HDKeyOperator interface.
type HDKeyOperatorImpl struct {
	coinTypeCode domainCoin.CoinTypeCode
	chainConfig  *chaincfg.Params
	strategy     portsWallet.CoinKeyStrategy
}

// NewHDKeyOperator creates a new HDKeyOperatorImpl.
func NewHDKeyOperator(
	coinTypeCode domainCoin.CoinTypeCode,
	chainConfig *chaincfg.Params,
	strategy portsWallet.CoinKeyStrategy,
) *HDKeyOperatorImpl {
	return &HDKeyOperatorImpl{
		coinTypeCode: coinTypeCode,
		chainConfig:  chainConfig,
		strategy:     strategy,
	}
}

// GetMasterFingerprintHex returns the master fingerprint as a hex string.
// Uses BIP44 purpose (arbitrary; master fingerprint is the same for all purposes).
func (h *HDKeyOperatorImpl) GetMasterFingerprintHex(seed []byte) (string, error) {
	hdKey := NewHDKey(PurposeTypeBIP44, h.coinTypeCode, h.chainConfig, h.strategy)
	descGen := NewDescriptorGenerator(hdKey, h.chainConfig)
	return descGen.GetMasterFingerprintHex(seed)
}

// GetCoinLevelXPub returns the coin-level extended public key for the given BIP purpose.
func (h *HDKeyOperatorImpl) GetCoinLevelXPub(seed []byte, purpose uint8) (string, error) {
	hdKey := NewHDKey(PurposeType(purpose), h.coinTypeCode, h.chainConfig, h.strategy)
	descGen := NewDescriptorGenerator(hdKey, h.chainConfig)
	return descGen.GetCoinLevelXPub(seed)
}

// GetAccountXPub returns the account-level extended public key for the given BIP purpose and account type.
func (h *HDKeyOperatorImpl) GetAccountXPub(
	seed []byte,
	purpose uint8,
	accountType domainAccount.AccountType,
) (string, error) {
	hdKey := NewHDKey(PurposeType(purpose), h.coinTypeCode, h.chainConfig, h.strategy)
	descGen := NewDescriptorGenerator(hdKey, h.chainConfig)
	xpub, err := descGen.GetAccountXPub(seed, accountType)
	if err != nil {
		return "", fmt.Errorf("failed to get account xpub for purpose %d, account %s: %w",
			purpose, accountType.String(), err)
	}
	return xpub, nil
}
