package btc

import (
	"errors"
	"fmt"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
)

func runDumpWallet(btc apibtc.PKGRPCProvider, fileName string) error {
	fmt.Println("dumps all wallet keys in a human-readable format to a server-side file")

	// validator
	if fileName == "" {
		return errors.New("filename option [-file] is required")
	}

	err := btc.GetPkgRPC().DumpWallet(fileName)
	if err != nil {
		return fmt.Errorf("fail to call btc.GetPkgRPC().DumpWallet(): %w", err)
	}

	fmt.Println("wallet file is dumped!")

	return nil
}
