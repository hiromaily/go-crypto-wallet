package btc

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/sign"
	btcsignsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/sign/btc"
)

type exportFullPubkeyUseCase struct {
	fullPubkeyExporter *btcsignsrv.FullPubkeyExport
}

// NewExportFullPubkeyUseCase creates a new ExportFullPubkeyUseCase
func NewExportFullPubkeyUseCase(fullPubkeyExporter *btcsignsrv.FullPubkeyExport) sign.ExportFullPubkeyUseCase {
	return &exportFullPubkeyUseCase{
		fullPubkeyExporter: fullPubkeyExporter,
	}
}

func (u *exportFullPubkeyUseCase) Export(ctx context.Context) (sign.ExportFullPubkeyOutput, error) {
	fileName, err := u.fullPubkeyExporter.ExportFullPubkey()
	if err != nil {
		return sign.ExportFullPubkeyOutput{}, fmt.Errorf("failed to export full pubkey: %w", err)
	}

	return sign.ExportFullPubkeyOutput{
		FileName: fileName,
	}, nil
}
