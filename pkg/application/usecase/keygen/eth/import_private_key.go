package eth

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen"
	ethkeygensrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/keygen/eth"
)

type importPrivateKeyUseCase struct {
	privKey *ethkeygensrv.PrivKey
}

// NewImportPrivateKeyUseCase creates a new ImportPrivateKeyUseCase
func NewImportPrivateKeyUseCase(privKey *ethkeygensrv.PrivKey) keygen.ImportPrivateKeyUseCase {
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
