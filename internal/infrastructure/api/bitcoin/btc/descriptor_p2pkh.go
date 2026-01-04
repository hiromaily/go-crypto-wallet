package btc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"

	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
)

const (
	descriptorFormatPKH    = "pkh([%s%s]%s/%d/*)"
	descriptorFormatSHWPKH = "sh(wpkh([%s%s]%s/%d/*))"
)

// DescriptorService generates descriptors for supported single-signature templates.
type DescriptorService struct {
	chainParams *chaincfg.Params
}

// NewDescriptorService creates a descriptor service for the given Bitcoin network parameters.
func NewDescriptorService(chainParams *chaincfg.Params) *DescriptorService {
	return &DescriptorService{
		chainParams: chainParams,
	}
}

// GenerateP2PKHDescriptor generates a P2PKH (legacy) descriptor following BIP380.
//
// Format: pkh([fingerprint/44'/0'/0']xpub.../0/*)
func (d *DescriptorService) GenerateP2PKHDescriptor(
	fingerprint string,
	derivationPath string,
	xpub *hdkeychain.ExtendedKey,
	isChange bool,
) (string, error) {
	return d.formatDescriptor(descriptorFormatPKH, fingerprint, derivationPath, xpub, isChange)
}

func (d *DescriptorService) formatDescriptor(
	template string,
	fingerprint string,
	derivationPath string,
	xpub *hdkeychain.ExtendedKey,
	isChange bool,
) (string, error) {
	if xpub == nil {
		return "", errors.New("extended public key is nil")
	}

	if xpub.IsPrivate() {
		return "", errors.New("extended key must be public")
	}

	normalizedFingerprint := strings.ToLower(strings.TrimSpace(fingerprint))
	normalizedPath := normalizeDerivationPath(derivationPath)

	if err := domainWallet.ValidateFingerprint(normalizedFingerprint); err != nil {
		return "", fmt.Errorf("invalid fingerprint: %w", err)
	}

	if err := domainWallet.ValidateDerivationPath(normalizedPath); err != nil {
		return "", fmt.Errorf("invalid derivation path: %w", err)
	}

	xpubStr := xpub.String()
	if err := domainWallet.ValidateExtendedPubKey(xpubStr); err != nil {
		return "", fmt.Errorf("invalid extended public key: %w", err)
	}

	changeIndex := 0
	if isChange {
		changeIndex = 1
	}

	return fmt.Sprintf(template, normalizedFingerprint, normalizedPath, xpubStr, changeIndex), nil
}

func normalizeDerivationPath(path string) string {
	trimmed := strings.TrimSpace(path)
	trimmed = strings.TrimPrefix(trimmed, "m")
	trimmed = strings.TrimPrefix(trimmed, "M")

	return trimmed
}
