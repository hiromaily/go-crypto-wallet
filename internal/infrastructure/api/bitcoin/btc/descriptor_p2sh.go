package btc

import "github.com/btcsuite/btcd/btcutil/hdkeychain"

// GenerateP2SHSegWitDescriptor generates a P2SH-SegWit descriptor following BIP380.
//
// Format: sh(wpkh([fingerprint/49'/0'/0']xpub.../0/*))
func (d *DescriptorService) GenerateP2SHSegWitDescriptor(
	fingerprint string,
	derivationPath string,
	xpub *hdkeychain.ExtendedKey,
	isChange bool,
) (string, error) {
	return d.formatDescriptor(descriptorFormatSHWPKH, fingerprint, derivationPath, xpub, isChange)
}
