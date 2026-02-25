package btc

import (
	"errors"
	"fmt"
)

func runDumpWallet(btc btcKeygenAPICmds, fileName string) error {
	fmt.Println("dumps all wallet keys in a human-readable format to a server-side file")

	// validator
	if fileName == "" {
		return errors.New("filename option [-file] is required")
	}

	err := btc.DumpWallet(fileName)
	if err != nil {
		return fmt.Errorf("fail to call btc.DumpWallet() %w", err)
	}

	fmt.Println("wallet file is dumped!")

	return nil
}
