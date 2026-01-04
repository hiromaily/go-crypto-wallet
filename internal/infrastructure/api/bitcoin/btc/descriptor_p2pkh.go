package btc

import "github.com/btcsuite/btcd/btcutil/hdkeychain"

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
