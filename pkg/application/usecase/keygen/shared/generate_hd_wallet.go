package shared

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen"
	sharedkeygensrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/keygen/shared"
)

type generateHDWalletUseCase struct {
	hdWallet *sharedkeygensrv.HDWallet
}

// NewGenerateHDWalletUseCase creates a new GenerateHDWalletUseCase
func NewGenerateHDWalletUseCase(hdWallet *sharedkeygensrv.HDWallet) keygen.GenerateHDWalletUseCase {
	return &generateHDWalletUseCase{
		hdWallet: hdWallet,
	}
}

func (u *generateHDWalletUseCase) Generate(
	ctx context.Context,
	input keygen.GenerateHDWalletInput,
) (keygen.GenerateHDWalletOutput, error) {
	walletKeys, err := u.hdWallet.Generate(input.AccountType, input.Seed, input.Count)
	if err != nil {
		return keygen.GenerateHDWalletOutput{}, fmt.Errorf("failed to generate HD wallet keys: %w", err)
	}

	return keygen.GenerateHDWalletOutput{
		GeneratedCount: len(walletKeys),
	}, nil
}
