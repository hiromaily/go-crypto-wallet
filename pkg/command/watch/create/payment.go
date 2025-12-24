package create

import (
	"fmt"

	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

func runPayment(wallet wallets.Watcher, fee float64) error {
	// Create payment transaction
	hex, fileName, err := wallet.CreatePaymentTx(fee)
	if err != nil {
		return fmt.Errorf("fail to call CreatePaymentTx() %w", err)
	}
	if (wallet.CoinTypeCode() != domainCoin.ETH && wallet.CoinTypeCode() != domainCoin.XRP) && hex == "" {
		fmt.Println("No utxo")
		return nil
	}
	// TODO: output should be json if json option is true
	fmt.Printf("[hex]: %s\n[fileName]: %s\n", hex, fileName)

	return nil
}
