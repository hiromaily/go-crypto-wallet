package btc

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen"
	btckeygensrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/keygen/btc"
)

type importFullPubkeyUseCase struct {
	fullPubkeyImporter *btckeygensrv.FullPubkeyImport
}

// NewImportFullPubkeyUseCase creates a new ImportFullPubkeyUseCase
func NewImportFullPubkeyUseCase(fullPubkeyImporter *btckeygensrv.FullPubkeyImport) keygen.ImportFullPubkeyUseCase {
	return &importFullPubkeyUseCase{
		fullPubkeyImporter: fullPubkeyImporter,
	}
}

func (u *importFullPubkeyUseCase) Import(ctx context.Context, input keygen.ImportFullPubkeyInput) error {
	if err := u.fullPubkeyImporter.ImportFullPubKey(input.FileName); err != nil {
		return fmt.Errorf("failed to import full pubkey: %w", err)
	}
	return nil
}
