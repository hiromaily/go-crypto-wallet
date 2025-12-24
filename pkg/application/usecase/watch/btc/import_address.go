package btc

import (
	"context"

	watchusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch"
	btcwatchsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/watch/btc"
)

// ImportAddressUseCase handles BTC address imports with rescan support
type ImportAddressUseCase interface {
	Execute(ctx context.Context, input watchusecase.ImportAddressInput) error
}

type importAddressUseCase struct {
	addressImporter *btcwatchsrv.AddressImport
}

// NewImportAddressUseCase creates a new BTC-specific ImportAddressUseCase
func NewImportAddressUseCase(addressImporter *btcwatchsrv.AddressImport) ImportAddressUseCase {
	return &importAddressUseCase{
		addressImporter: addressImporter,
	}
}

// Execute imports addresses from a file with optional rescan
func (u *importAddressUseCase) Execute(ctx context.Context, input watchusecase.ImportAddressInput) error {
	return u.addressImporter.ImportAddress(input.FileName, input.Rescan)
}
