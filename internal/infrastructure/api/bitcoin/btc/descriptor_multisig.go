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
func (d *DescriptorService) GenerateMultisigDescriptor(
	requiredSigs int,
	signers []MultisigSigner,
	isChange bool,
	isP2SHWrapped bool,
) (string, error) {
	if len(signers) == 0 {
		return "", errors.New("no multisig signers provided")
	}

	if requiredSigs < 1 || requiredSigs > len(signers) {
		return "", fmt.Errorf("invalid required signatures: %d of %d", requiredSigs, len(signers))
	}

	keyStrings, err := d.formatAndSortMultisigKeys(signers, isChange)
	if err != nil {
		return "", err
	}

	multisigPart := fmt.Sprintf(
		"sortedmulti(%d,%s)",
		requiredSigs,
		strings.Join(keyStrings, ","),
	)

	var descriptor string
	if isP2SHWrapped {
		// P2SH-P2WSH: sh(wsh(sortedmulti(...)))
		descriptor = fmt.Sprintf("sh(wsh(%s))", multisigPart)
	} else {
		// P2WSH: wsh(sortedmulti(...))
		descriptor = fmt.Sprintf("wsh(%s)", multisigPart)
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
