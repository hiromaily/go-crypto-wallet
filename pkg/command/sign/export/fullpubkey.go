package export

import (
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

func runFullPubkey(wallet wallets.Signer) error {
	fmt.Println("export full pubkey")

	// export full pubkey as csv file
	fileName, err := wallet.ExportFullPubkey()
	if err != nil {
		return fmt.Errorf("fail to call ExportFullPubkey() %w", err)
	}
	fmt.Println("[fileName]: " + fileName)

	return nil
}
