package eth

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/common"
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

// ReadPrivKey reads a keystore JSON file for the given address from path.
// Keystore filenames follow the pattern: UTC--<timestamp>--<address>
// The address in the filename is always lowercase without the 0x prefix.
// Returns (nil, error) if no file or multiple files are found.
func ReadPrivKey(hexAddr, path string) ([]byte, error) {
	if !common.IsHexAddress(hexAddr) {
		return nil, fmt.Errorf("invalid Ethereum address: %s", hexAddr)
	}
	addr := strings.TrimPrefix(strings.ToLower(hexAddr), "0x")
	files, err := filepath.Glob(fmt.Sprintf("%s/*--%s", path, addr))
	if err != nil {
		return nil, fmt.Errorf("fail to call filepath.Glob(): %w", err)
	}
	if len(files) == 0 {
		return nil, errors.New("private key file is not found")
	}
	if len(files) > 1 {
		return nil, fmt.Errorf("target private key files are found more than 1 by %s", addr)
	}
	return os.ReadFile(files[0])
}
