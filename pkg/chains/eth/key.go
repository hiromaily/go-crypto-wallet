package eth

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// ToECDSA converts a hex-encoded private key string to an ECDSA private key.
func ToECDSA(privKey string) (*ecdsa.PrivateKey, error) {
	bytePrivKey, err := hexutil.Decode(privKey)
	if err != nil {
		return nil, fmt.Errorf("fail to call hexutil.Decode(): %w", err)
	}
	return crypto.ToECDSA(bytePrivKey)
}
