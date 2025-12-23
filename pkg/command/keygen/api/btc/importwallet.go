package btc

import (
	"errors"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
)

func runImportWallet(btc btcgrp.Bitcoiner, fileName string) error {
	fmt.Println("Imports keys from a wallet dump file")

	// validator
	if fileName == "" {
		return errors.New("filename option [-file] is required")
	}

	err := btc.ImportWallet(fileName)
	if err != nil {
		return fmt.Errorf("fail to call btc.ImportWallet() %w", err)
	}

	fmt.Println("wallet file is imported!")

	return nil
}
