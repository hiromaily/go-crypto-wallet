package imports

import (
	"errors"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

func runAddress(wallet wallets.Watcher, filePath string, isRescan bool) error {
	fmt.Println("-file: " + filePath)

	// validator
	if filePath == "" {
		return errors.New("file path option [-file] is required")
	}

	// import public address
	err := wallet.ImportAddress(filePath, isRescan)
	if err != nil {
		return fmt.Errorf("fail to call ImportAddress() %w", err)
	}
	fmt.Println("Done!")

	return nil
}
