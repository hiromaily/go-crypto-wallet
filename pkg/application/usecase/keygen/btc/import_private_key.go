package btc

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen"
	btckeygensrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/keygen/btc"
)

type importPrivateKeyUseCase struct {
	privKey *btckeygensrv.PrivKey
}

// NewImportPrivateKeyUseCase creates a new ImportPrivateKeyUseCase
func NewImportPrivateKeyUseCase(privKey *btckeygensrv.PrivKey) keygen.ImportPrivateKeyUseCase {
	return &importPrivateKeyUseCase{
		privKey: privKey,
	}
}

func (u *importPrivateKeyUseCase) Import(ctx context.Context, input keygen.ImportPrivateKeyInput) error {
	if err := u.privKey.Import(input.AccountType); err != nil {
		return fmt.Errorf("failed to import private key: %w", err)
	}
	return nil
}
