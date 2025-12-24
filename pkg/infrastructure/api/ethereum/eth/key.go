package eth

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// Note: key filename is different between Geth and Parity
// Geth
// - 0x71678cd07cfac46c2dc427f999abf46aae115925
// - UTC--2018-10-11T06-58-43.857846090Z--71678cd07cfac46c2dc427f999abf46aae115925

// Parity, filename includes just UUID
// - "0xcf9583c3c10cf895af95a2810243765c4fe7c038",
// - UTC--2018-10-11T06-59-28Z--2bd02735-84ec-593e-f2b2-73cce1b1862c

// File names for Parity keys
// https://ethereum.stackexchange.com/questions/13951/file-names-of-parity-keys

// So Parity key filename should be renamed to same format as Geth
// e.g. UTC--2018-10-12T01-53-58Z--fff7e98d-b3b7-08f4-65cd-3fe82416cebf--45783b86c2aa1ce81632ac2db26a91acc3ea6226

// ToECDSA converts privKey to ECDSA
func (*Ethereum) ToECDSA(privKey string) (*ecdsa.PrivateKey, error) {
	bytePrivKey, err := hexutil.Decode(privKey)
	if err != nil {
		return nil, fmt.Errorf("fail to call hexutil.Decode(): %w", err)
	}
	return crypto.ToECDSA(bytePrivKey)
}

// GetKeyDir returns keystore directory
func (e *Ethereum) GetKeyDir() string {
	return e.keyDir
}

// GetPrivKey returns keystore.Key object
func (e *Ethereum) GetPrivKey(hexAddr, password string) (*keystore.Key, error) {
	keyDir := e.GetKeyDir()
	logger.Debug("key_dir", "key_dir", keyDir)

	keyJSON, err := e.readPrivKey(hexAddr, keyDir)
	if err != nil {
		return nil, fmt.Errorf("fail to call e.readPrivKey(): %w", err)
	}
	if keyJSON == nil {
		// file is not found
		return nil, errors.New("private key file is not found")
	}

	key, err := keystore.DecryptKey(keyJSON, password)
	if err != nil {
		return nil, fmt.Errorf("fail to call keystore.DecryptKey(): %w", err)
	}
	return key, nil
}

// readPrivKey read private key file from directory
// Note: file is found out from local directory,
//
//	if node is working remotely, file is not found.
func (*Ethereum) readPrivKey(hexAddr, path string) ([]byte, error) {
	// search file
	// filename is like `UTC--2020-05-18T16-01-32.772616000Z--e52307deb1a7dc3985d2873b45ae23b91d57a36d`
	//                   UTC--2020-05-22T13-59-51.645058000Z--d4ec46122b3f0afc0287144adcca5d65b22b799c
	// Note: all letter of address in filename is a lowercase letter
	addr := strings.TrimPrefix(strings.ToLower(hexAddr), "0x")
	logger.Debug("readPrivKey",
		"hexAddr", hexAddr,
		"addr", addr,
		"path", path,
	)

	files, err := filepath.Glob(fmt.Sprintf("%s/*--%s", path, addr))
	if err != nil {
		return nil, fmt.Errorf("fail to call filepath.Glob(): %w", err)
	}
	if len(files) == 0 {
		// file is not found
		return nil, errors.New("private key file is not found")
	}
	if len(files) > 1 {
		return nil, fmt.Errorf("target private key files are found more than 1 by %s", addr)
	}

	return os.ReadFile(files[0])
}
