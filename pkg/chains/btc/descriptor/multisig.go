package descriptor

// Descriptor Multisig - BTC ONLY (BCH does NOT support descriptors)
// See types.go for full warning.
// BCH uses traditional P2SH multisig without descriptors.

import (
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
)

// MultisigSigner represents a signer in a multisig descriptor.
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
	descriptorType DescriptorType,
) (string, error) {
	if len(signers) == 0 {
		return "", errors.New("no multisig signers provided")
	}

	if requiredSigs < 1 || requiredSigs > len(signers) {
		return "", fmt.Errorf("invalid required signatures: %d of %d", requiredSigs, len(signers))
	}

	// For P2SH multi(), keys should NOT be sorted (preserve original order).
	// For P2WSH/P2SH-P2WSH sortedmulti(), keys MUST be sorted.
	shouldSortKeys := descriptorType == DescriptorTypeWSH ||
		descriptorType == DescriptorTypeSHWSH

	keyStrings, err := d.formatAndSortMultisigKeys(signers, isChange, shouldSortKeys)
	if err != nil {
		return "", err
	}

	var desc string
	switch descriptorType {
	case DescriptorTypeSH:
		// P2SH (Legacy): sh(multi(...))
		// Note: For BIP44 Legacy multisig, we use multi() not sortedmulti()
		// to match Bitcoin Core's standard behavior for non-witness scripts.
		// Keys are NOT sorted - they remain in the order provided.
		multisigPart := fmt.Sprintf(
			"multi(%d,%s)",
			requiredSigs,
			strings.Join(keyStrings, ","),
		)
		desc = fmt.Sprintf("sh(%s)", multisigPart)

	case DescriptorTypeSHWSH, DescriptorTypeWSH:
		// Both P2SH-P2WSH (BIP49) and P2WSH (BIP48) use sortedmulti
		// Keys are sorted lexicographically
		// Generate common sortedmulti part
		multisigPart := fmt.Sprintf(
			"sortedmulti(%d,%s)",
			requiredSigs,
			strings.Join(keyStrings, ","),
		)
		// Apply appropriate wrapper based on descriptor type
		if descriptorType == DescriptorTypeSHWSH {
			// P2SH-P2WSH (BIP49): sh(wsh(sortedmulti(...)))
			desc = fmt.Sprintf("sh(wsh(%s))", multisigPart)
		} else {
			// P2WSH (BIP48): wsh(sortedmulti(...))
			desc = fmt.Sprintf("wsh(%s)", multisigPart)
		}

	case DescriptorTypePKH, DescriptorTypeSHWPKH,
		DescriptorTypeWPKH, DescriptorTypeTR,
		DescriptorTypeUnknown:
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

	return desc, nil
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

	if err := ValidateFingerprint(fp); err != nil {
		return "", fmt.Errorf("invalid fingerprint: %w", err)
	}

	if err := ValidateDerivationPath(path); err != nil {
		return "", fmt.Errorf("invalid derivation path: %w", err)
	}

	return fmt.Sprintf("[%s%s]%s/%d/*", fp, path, signer.ExtendedKey.String(), changeIndex), nil
}

// GenerateTaprootScriptPathDescriptor generates a Taproot descriptor with multiple script-path pubkeys.
//
// SECURITY NOTE: This is NOT MuSig2. The first key is used as the internal key (key-path spending),
// and remaining keys form script paths (script-path spending). This results in a 1-of-N spend policy
// where any listed key can spend (first key via key path, others via script path).
// For MuSig2 key-path spends, aggregate keys first and use tr(<agg_key>).
//
// Format: tr(internal_key,{pk(key1),pk(key2),...})
// - internal_key: First signer's key (can spend via key path - most efficient)
// - Script tree: Remaining signers' keys wrapped in pk() (can spend via script path)
func (d *DescriptorService) GenerateTaprootScriptPathDescriptor(
	signers []MultisigSigner,
	isChange bool,
) (string, error) {
	if len(signers) < 2 {
		return "", errors.New("at least 2 signers are required for taproot script-path multi-key descriptors")
	}

	// For Taproot script-path descriptors, keys should be sorted for deterministic output
	keyStrings, err := d.formatAndSortMultisigKeys(signers, isChange, true)
	if err != nil {
		return "", err
	}

	// First key becomes the internal key (key-path spending)
	internalKey := keyStrings[0]

	// Remaining keys become script paths (script-path spending)
	scriptKeys := keyStrings[1:]

	// Wrap each script key in pk() for Tapscript
	pkScripts := make([]string, len(scriptKeys))
	for i, key := range scriptKeys {
		pkScripts[i] = fmt.Sprintf("pk(%s)", key)
	}

	// Build descriptor based on number of script keys
	var desc string
	if len(pkScripts) == 1 {
		// Single script: tr(internal_key,pk(script_key))
		desc = fmt.Sprintf("tr(%s,%s)", internalKey, pkScripts[0])
	} else {
		// Multiple scripts: tr(internal_key,{pk(key1),pk(key2),...})
		desc = fmt.Sprintf("tr(%s,{%s})", internalKey, strings.Join(pkScripts, ","))
	}

	// Note: Checksum is NOT added here because the domain layer's BIP380 implementation
	// produces incorrect checksums. Instead, the watch wallet will add checksums using
	// Bitcoin Core's getdescriptorinfo RPC before importing (which guarantees correctness).
	// This allows keygen to remain offline while ensuring correct checksums.
	// TODO: Fix BIP380 checksum implementation in internal/domain/wallet/descriptor_builder.go

	return desc, nil
}
