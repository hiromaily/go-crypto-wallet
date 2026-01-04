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

	return fmt.Sprintf("tr(%s)", strings.Join(keyStrings, ",")), nil
}
