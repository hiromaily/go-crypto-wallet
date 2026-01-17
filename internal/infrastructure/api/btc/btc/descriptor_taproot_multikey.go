package btc

// Descriptor Taproot Multikey - BTC ONLY (BCH does NOT support descriptors or Taproot)
// See descriptor_service.go for full warning.

import (
	"errors"
	"fmt"
	"strings"
)

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
	var descriptor string
	if len(pkScripts) == 1 {
		// Single script: tr(internal_key,pk(script_key))
		descriptor = fmt.Sprintf("tr(%s,%s)", internalKey, pkScripts[0])
	} else {
		// Multiple scripts: tr(internal_key,{pk(key1),pk(key2),...})
		descriptor = fmt.Sprintf("tr(%s,{%s})", internalKey, strings.Join(pkScripts, ","))
	}

	// Note: Checksum is NOT added here because the domain layer's BIP380 implementation
	// produces incorrect checksums. Instead, the watch wallet will add checksums using
	// Bitcoin Core's getdescriptorinfo RPC before importing (which guarantees correctness).
	// This allows keygen to remain offline while ensuring correct checksums.
	// TODO: Fix BIP380 checksum implementation in internal/domain/wallet/descriptor_builder.go

	return descriptor, nil
}
