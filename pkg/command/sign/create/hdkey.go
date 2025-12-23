package create

import (
	"fmt"

	"github.com/bookerzzz/grok"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

func runHDKey(wallet wallets.Signer) error {
	fmt.Println("create key for hd wallet for Authorization account")

	// create seed
	bSeed, err := wallet.GenerateSeed()
	if err != nil {
		return fmt.Errorf("fail to call GenerateSeed() %w", err)
	}

	// create key for hd wallet for Authorization account
	keys, err := wallet.GenerateAuthKey(bSeed, 1)
	if err != nil {
		return fmt.Errorf("fail to call GenerateAuthKey() %w", err)
	}
	grok.Value(keys)

	return nil
}
