package shared

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch"
	sharedwatchsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/watch/shared"
)

type importAddressUseCase struct {
	addressImporter *sharedwatchsrv.AddressImport
}

// NewImportAddressUseCase creates a new ImportAddressUseCase
func NewImportAddressUseCase(addressImporter *sharedwatchsrv.AddressImport) watch.ImportAddressUseCase {
	return &importAddressUseCase{
		addressImporter: addressImporter,
	}
}

func (u *importAddressUseCase) Execute(ctx context.Context, input watch.ImportAddressInput) error {
	if err := u.addressImporter.ImportAddress(input.FileName); err != nil {
		return fmt.Errorf("failed to import address: %w", err)
	}
	return nil
}
