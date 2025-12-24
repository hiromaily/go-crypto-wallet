package shared

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen"
	sharedkeygensrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/keygen/shared"
)

type exportAddressUseCase struct {
	addressExporter *sharedkeygensrv.AddressExport
}

// NewExportAddressUseCase creates a new ExportAddressUseCase
func NewExportAddressUseCase(addressExporter *sharedkeygensrv.AddressExport) keygen.ExportAddressUseCase {
	return &exportAddressUseCase{
		addressExporter: addressExporter,
	}
}

func (u *exportAddressUseCase) Export(
	ctx context.Context,
	input keygen.ExportAddressInput,
) (keygen.ExportAddressOutput, error) {
	fileName, err := u.addressExporter.ExportAddress(input.AccountType)
	if err != nil {
		return keygen.ExportAddressOutput{}, fmt.Errorf("failed to export address: %w", err)
	}

	return keygen.ExportAddressOutput{
		FileName: fileName,
	}, nil
}
