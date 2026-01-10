package btc

import (
	"errors"
	"fmt"
	"strings"
)

// GenerateTaprootScriptPathDescriptor generates a Taproot descriptor with multiple script-path pubkeys.
//
// SECURITY NOTE: This is NOT MuSig2. Each key forms its own script path, resulting in a 1-of-N spend policy
// (any listed key can spend via script path). For MuSig2 key-path spends, aggregate keys first and use tr(<agg_key>).
//
// Format: tr([fp1/path]xpub1/change/*,[fp2/path]xpub2/change/*,...)
func (d *DescriptorService) GenerateTaprootScriptPathDescriptor(
	signers []MultisigSigner,
	isChange bool,
) (string, error) {
	if len(signers) < 2 {
		return "", errors.New("at least 2 signers are required for taproot script-path multi-key descriptors")
	}

	keyStrings, err := d.formatAndSortMultisigKeys(signers, isChange)
	if err != nil {
		return "", err
	}

	descriptor := fmt.Sprintf("tr(%s)", strings.Join(keyStrings, ","))

	// Note: Checksum is NOT added here because the domain layer's BIP380 implementation
	// produces incorrect checksums. Instead, the watch wallet will add checksums using
	// Bitcoin Core's getdescriptorinfo RPC before importing (which guarantees correctness).
	// This allows keygen to remain offline while ensuring correct checksums.
	// TODO: Fix BIP380 checksum implementation in internal/domain/wallet/descriptor_builder.go

	return descriptor, nil
}
