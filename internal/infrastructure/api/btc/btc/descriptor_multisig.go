// Descriptor Multisig - BTC ONLY (BCH does NOT support descriptors)
// See descriptor_service.go for full warning.
// BCH uses traditional P2SH multisig without descriptors.
package btc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"

	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
)

// MultisigSigner represents a signer in a multisignature descriptor.
type MultisigSigner struct {
	Fingerprint    string
	DerivationPath string
	ExtendedKey    *hdkeychain.ExtendedKey
}

// GenerateMultisigDescriptor generates a traditional multisig descriptor.
//
// Format:
//   - P2WSH (native SegWit): wsh(sortedmulti(M,[fp1/path]xpub1/change/*,[fp2/path]xpub2/change/*,...))
//   - P2SH-P2WSH (wrapped SegWit): sh(wsh(sortedmulti(M,[fp1/path]xpub1/change/*,[fp2/path]xpub2/change/*,...)))
//   - P2SH (Legacy): sh(multi(M,[fp1/path]xpub1/change/*,[fp2/path]xpub2/change/*,...))
func (d *DescriptorService) GenerateMultisigDescriptor(
	requiredSigs int,
	signers []MultisigSigner,
	isChange bool,
	descriptorType domainWallet.DescriptorType,
) (string, error) {
	if len(signers) == 0 {
		return "", errors.New("no multisig signers provided")
	}

	if requiredSigs < 1 || requiredSigs > len(signers) {
		return "", fmt.Errorf("invalid required signatures: %d of %d", requiredSigs, len(signers))
	}

	// For P2SH multi(), keys should NOT be sorted (preserve original order).
	// For P2WSH/P2SH-P2WSH sortedmulti(), keys MUST be sorted.
	shouldSortKeys := descriptorType == domainWallet.DescriptorTypeWSH ||
		descriptorType == domainWallet.DescriptorTypeSHWSH

	keyStrings, err := d.formatAndSortMultisigKeys(signers, isChange, shouldSortKeys)
	if err != nil {
		return "", err
	}

	var descriptor string
	switch descriptorType {
	case domainWallet.DescriptorTypeSH:
		// P2SH (Legacy): sh(multi(...))
		// Note: For BIP44 Legacy multisig, we use multi() not sortedmulti()
		// to match Bitcoin Core's standard behavior for non-witness scripts.
		// Keys are NOT sorted - they remain in the order provided.
		multisigPart := fmt.Sprintf(
			"multi(%d,%s)",
			requiredSigs,
			strings.Join(keyStrings, ","),
		)
		descriptor = fmt.Sprintf("sh(%s)", multisigPart)

	case domainWallet.DescriptorTypeSHWSH, domainWallet.DescriptorTypeWSH:
		// Both P2SH-P2WSH (BIP49) and P2WSH (BIP48) use sortedmulti
		// Keys are sorted lexicographically
		// Generate common sortedmulti part
		multisigPart := fmt.Sprintf(
			"sortedmulti(%d,%s)",
			requiredSigs,
			strings.Join(keyStrings, ","),
		)
		// Apply appropriate wrapper based on descriptor type
		if descriptorType == domainWallet.DescriptorTypeSHWSH {
			// P2SH-P2WSH (BIP49): sh(wsh(sortedmulti(...)))
			descriptor = fmt.Sprintf("sh(wsh(%s))", multisigPart)
		} else {
			// P2WSH (BIP48): wsh(sortedmulti(...))
			descriptor = fmt.Sprintf("wsh(%s)", multisigPart)
		}

	case domainWallet.DescriptorTypePKH, domainWallet.DescriptorTypeSHWPKH,
		domainWallet.DescriptorTypeWPKH, domainWallet.DescriptorTypeTR,
		domainWallet.DescriptorTypeUnknown:
		// These descriptor types are not for traditional multisig
		return "", fmt.Errorf("descriptor type %s is not supported for traditional multisig", descriptorType.String())

	default:
		return "", fmt.Errorf("unsupported multisig descriptor type: %s", descriptorType.String())
	}

	// Note: Checksum is NOT added here because the domain layer's BIP380 implementation
	// produces incorrect checksums. Instead, the watch wallet will add checksums using
	// Bitcoin Core's getdescriptorinfo RPC before importing (which guarantees correctness).
	// This allows keygen to remain offline while ensuring correct checksums.
	// TODO: Fix BIP380 checksum implementation in internal/domain/wallet/descriptor_builder.go

	return descriptor, nil
}

func (d *DescriptorService) formatMultisigKey(signer MultisigSigner, changeIndex int) (string, error) {
	if signer.ExtendedKey == nil {
		return "", errors.New("extended public key is nil")
	}

	if signer.ExtendedKey.IsPrivate() {
		return "", errors.New("extended key must be public")
	}

	if d.chainParams != nil && !signer.ExtendedKey.IsForNet(d.chainParams) {
		return "", fmt.Errorf("extended public key network mismatch: expected %s", d.chainParams.Name)
	}

	fp := strings.ToLower(strings.TrimSpace(signer.Fingerprint))
	path := normalizeDerivationPath(signer.DerivationPath)

	if err := domainWallet.ValidateFingerprint(fp); err != nil {
		return "", fmt.Errorf("invalid fingerprint: %w", err)
	}

	if err := domainWallet.ValidateDerivationPath(path); err != nil {
		return "", fmt.Errorf("invalid derivation path: %w", err)
	}

	return fmt.Sprintf("[%s%s]%s/%d/*", fp, path, signer.ExtendedKey.String(), changeIndex), nil
}
