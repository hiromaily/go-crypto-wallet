package create

import (
	"fmt"
	"os"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

func runSeed(wallet wallets.Keygener, seed string) error {
	fmt.Println("create seed")

	var (
		bSeed []byte
		err   error
	)

	if seed == "" {
		seed = os.Getenv("KEYGEN_SEED")
		if seed != "" {
			fmt.Println("seed is found from environment variable")
		}
	}

	if seed != "" {
		// store seed into database, not generate seed
		bSeed, err = wallet.StoreSeed(seed)
		if err != nil {
			return fmt.Errorf("fail to call StoreSeed() %w", err)
		}
	} else {
		// create seed
		bSeed, err = wallet.GenerateSeed()
		if err != nil {
			return fmt.Errorf("fail to call GenerateSeed() %w", err)
		}
	}
	fmt.Println("seed: " + key.SeedToString(bSeed))

	return nil
}
