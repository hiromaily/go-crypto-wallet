package key

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

func ToECDSA(privKey string) (*ecdsa.PrivateKey, error) {
	bytePrivKey, err := hexutil.Decode(privKey)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call hexutil.Decode()")
	}
	return crypto.ToECDSA(bytePrivKey)
}
