package sign

import (
	"errors"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

func runSignature(wallet wallets.Keygener, filePath string) error {
	fmt.Println("sign on unsigned transaction (account would be found from file name)")

	// validator
	if filePath == "" {
		return errors.New("file path option [-file] is required")
	}

	// sign on unsigned transactions, action(deposit/payment) could be found from file name
	hexTx, isSigned, generatedFileName, err := wallet.SignTx(filePath)
	if err != nil {
		return fmt.Errorf("fail to call SignTx() %w", err)
	}

	// TODO: output should be json if json option is true
	fmt.Printf("[hex]: %s\n[isCompleted]: %t\n[fileName]: %s\n", hexTx, isSigned, generatedFileName)

	return nil
}
