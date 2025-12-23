package create

import (
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/btcwallet"
)

func runKey(wallet wallets.Keygener) error {
	fmt.Println("create one key for debug use")

	// create one key for debug use
	if v, ok := wallet.(*btcwallet.BTCKeygen); ok {
		wif, pubAddress, err := key.GenerateWIF(v.BTC.GetChainConf())
		if err != nil {
			return fmt.Errorf("fail to call key.GenerateKey() %w", err)
		}
		fmt.Printf("[WIF] %s - [Pub Address] %s\n", wif.String(), pubAddress)
	} else {
		fmt.Println("for only BTC")
	}
	return nil
}
