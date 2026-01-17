package btc

// Descriptor Bech32 - BTC ONLY (BCH does NOT support descriptors or Bech32)
// See descriptor_service.go for full warning.

import "github.com/btcsuite/btcd/btcutil/hdkeychain"

// GenerateBech32Descriptor generates a Bech32 (wpkh) descriptor following BIP380.
//
// Format: wpkh([fingerprint/84'/0'/0']xpub.../0/*)
func (d *DescriptorService) GenerateBech32Descriptor(
	fingerprint string,
	derivationPath string,
	xpub *hdkeychain.ExtendedKey,
	isChange bool,
) (string, error) {
	return d.formatDescriptor(descriptorFormatWPKH, fingerprint, derivationPath, xpub, isChange)
}
