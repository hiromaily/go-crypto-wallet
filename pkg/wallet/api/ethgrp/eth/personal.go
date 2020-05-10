package eth

import (
	"strings"

	"github.com/pkg/errors"
)

// ImportRawKey Imports the given unencrypted private key (hex string) into the key store, encrypting it with the passphrase
// - if address has 0x like 0x5d0a82e19564ae03ad3b834ac30b94c0ccce510e86d783d3e882efcb0e84b2af,
//    error would occur `invalid hex character 'x' in private key`
func (e *Ethereum) ImportRawKey(hexKey, passPhrase string) (string, error) {
	if strings.HasPrefix(hexKey, "0x") {
		hexKey = strings.TrimLeft(hexKey, "0x")
	}

	var address string
	err := e.rpcClient.CallContext(e.ctx, &address, "personal_importRawKey", hexKey, passPhrase)
	if err != nil {
		return "", errors.Wrap(err, "fail to call client.CallContext(personal_importRawKey)")
	}

	return address, nil
}
