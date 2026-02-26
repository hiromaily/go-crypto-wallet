package btc

import "github.com/btcsuite/btcd/btcutil/hdkeychain"

// MultisigSigner represents a signer in a multisig descriptor.
type MultisigSigner struct {
	Fingerprint    string
	DerivationPath string
	ExtendedKey    *hdkeychain.ExtendedKey
}
