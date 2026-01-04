package btc

import "github.com/btcsuite/btcd/btcutil/hdkeychain"

// GenerateTaprootDescriptor generates a Taproot (tr) descriptor following BIP380.
//
// Format: tr([fingerprint/86'/0'/0']xpub.../0/*)
func (d *DescriptorService) GenerateTaprootDescriptor(
	fingerprint string,
	derivationPath string,
	xpub *hdkeychain.ExtendedKey,
	isChange bool,
) (string, error) {
	return d.formatDescriptor(descriptorFormatTR, fingerprint, derivationPath, xpub, isChange)
}
