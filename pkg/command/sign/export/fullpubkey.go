package export

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/di"
)

func runFullPubkey(container di.Container) error {
	fmt.Println("export full pubkey")

	// export full pubkey as csv file
	useCase := container.NewSignExportFullPubkeyUseCase(container.AuthType())
	output, err := useCase.Export(context.Background())
	if err != nil {
		return fmt.Errorf("fail to export full pubkey: %w", err)
	}
	fmt.Println("[fileName]: " + output.FileName)

	return nil
}
