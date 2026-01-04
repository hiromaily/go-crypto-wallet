package btc

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// GenerateMuSig2Descriptor generates a MuSig2 Taproot descriptor.
//
// Format: tr([fp1/path]xpub1/change/*,[fp2/path]xpub2/change/*,...)
func (d *DescriptorService) GenerateMuSig2Descriptor(
	signers []MultisigSigner,
	isChange bool,
) (string, error) {
	if len(signers) < 2 {
		return "", errors.New("MuSig2 requires at least 2 signers")
	}

	changeIndex := 0
	if isChange {
		changeIndex = 1
	}

	keyStrings := make([]string, len(signers))
	for i, signer := range signers {
		keyStr, err := d.formatMultisigKey(signer, changeIndex)
		if err != nil {
			return "", fmt.Errorf("invalid signer %d: %w", i, err)
		}
		keyStrings[i] = keyStr
	}

	// Deterministic ordering for reproducible descriptors.
	sort.Strings(keyStrings)

	return fmt.Sprintf("tr(%s)", strings.Join(keyStrings, ",")), nil
}
