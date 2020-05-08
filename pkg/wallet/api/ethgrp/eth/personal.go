package eth

import (
	"github.com/pkg/errors"
)

// ImportRawKey Imports the given unencrypted private key (hex string) into the key store, encrypting it with the passphrase
func (e *Ethereum) ImportRawKey(hexKey, passPhrase string) (string, error) {
	var address string
	err := e.client.CallContext(e.ctx, &address, "personal_importRawKey", hexKey, passPhrase)
	if err != nil {
		return "", errors.Wrap(err, "fail to call client.CallContext(personal_importRawKey)")
	}

	return address, nil
}
