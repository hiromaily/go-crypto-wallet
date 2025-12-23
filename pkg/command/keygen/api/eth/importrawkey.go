package eth

import (
	"context"
	"errors"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
)

func runImportRawKey(eth ethgrp.Ethereumer, privKey, passPhrase string) error {
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
