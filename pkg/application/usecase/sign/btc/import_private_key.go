package btc

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/sign"
	btcsignsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/sign/btc"
)

type importPrivateKeyUseCase struct {
	privKey *btcsignsrv.PrivKey
}

// NewImportPrivateKeyUseCase creates a new ImportPrivateKeyUseCase
func NewImportPrivateKeyUseCase(privKey *btcsignsrv.PrivKey) sign.ImportPrivateKeyUseCase {
	return &importPrivateKeyUseCase{
		privKey: privKey,
	}
}

func (u *importPrivateKeyUseCase) Import(ctx context.Context, input sign.ImportPrivateKeyInput) error {
	// Note: BTC PrivKey.Import() doesn't take authType as parameter
	// The authType is already set in the PrivKey struct during construction
	if err := u.privKey.Import(); err != nil {
		return fmt.Errorf("failed to import private key: %w", err)
	}
	return nil
}
