package eth

import (
	"context"
	"errors"
	"fmt"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
)

// runImportRawKey imports a raw private key using the personal_importRawKey RPC
//
// Deprecated: This command uses personal_importRawKey RPC which is NOT supported by Anvil.
// For Anvil compatibility, use the modern "import key" command instead:
//
//	keygen -coin eth import key --account <account>
//
// The modern command uses local keystore.ImportECDSA() and works with both Geth and Anvil.
func runImportRawKey(eth apieth.Ethereumer, privKey, passPhrase string) error {
	fmt.Println("import raw key")

	// validation
	if privKey == "" {
		return errors.New("key option [-key] is invalid")
	}
	if passPhrase == "" {
		return errors.New("pass option [-pass] is invalid")
	}

	addr, err := eth.ImportRawKey(context.Background(), privKey, passPhrase)
	if err != nil {
		return fmt.Errorf("fail to call eth.ImportRawKey() %w", err)
	}

	fmt.Println("new address: " + addr)

	return nil
}
