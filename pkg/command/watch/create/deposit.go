package create

import (
	"fmt"

	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

func runDeposit(wallet wallets.Watcher, fee float64) error {
	// Detect transaction for clients from blockchain network and create deposit unsigned transaction
	// It would be run manually on the daily basis because signature is manual task
	hex, fileName, err := wallet.CreateDepositTx(fee)
	if err != nil {
		return fmt.Errorf("fail to call CreateDepositTx() %w", err)
	}
	if domainCoin.IsBTCGroup(wallet.CoinTypeCode()) && hex == "" {
		fmt.Println("No utxo")
		return nil
	}
	// TODO: output should be json if json option is true
	fmt.Printf("[hex]: %s\n[fileName]: %s\n", hex, fileName)

	return nil
}
