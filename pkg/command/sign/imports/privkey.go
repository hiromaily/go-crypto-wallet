package imports

import (
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

func runPrivKey(wallet wallets.Signer) error {
	fmt.Println("import generated private key for Authorization account to database")

	// import generated private key to database
	err := wallet.ImportPrivKey()
	if err != nil {
		return fmt.Errorf("fail to call ImportPrivKey() %w", err)
	}
	fmt.Println("Done!")

	return nil
}
