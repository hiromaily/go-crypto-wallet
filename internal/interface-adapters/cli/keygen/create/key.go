package create

import (
	"fmt"

	wallets "github.com/hiromaily/go-crypto-wallet/internal/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/wallet/btcwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/key"
)

func runKey(wallet wallets.Keygener) error {
	fmt.Println("create one key for debug use")

	// create one key for debug use - BTC only
	// This is debug code that uses utility function directly
	// Not migrated to use case layer as it's a simple utility operation
	if v, ok := wallet.(*btcwallet.BTCKeygen); ok {
		wif, pubAddress, err := key.GenerateWIF(v.BTC.GetChainConf())
		if err != nil {
			return fmt.Errorf("fail to generate WIF key: %w", err)
		}
		fmt.Printf("[WIF] %s - [Pub Address] %s\n", wif.String(), pubAddress)
	} else {
		fmt.Println("for only BTC")
	}
	return nil
}
