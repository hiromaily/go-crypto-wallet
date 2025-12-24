package shared

import (
	"context"
	"fmt"

	signusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/sign"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
)

type generateAuthKeyUseCase struct {
	hdWalleter service.HDWalleter
}

// NewGenerateAuthKeyUseCase creates a new GenerateAuthKeyUseCase for sign wallet
func NewGenerateAuthKeyUseCase(hdWalleter service.HDWalleter) signusecase.GenerateAuthKeyUseCase {
	return &generateAuthKeyUseCase{
		hdWalleter: hdWalleter,
	}
}

func (u *generateAuthKeyUseCase) Generate(
	ctx context.Context, input signusecase.GenerateAuthKeyInput,
) (signusecase.GenerateAuthKeyOutput, error) {
	keys, err := u.hdWalleter.Generate(input.AuthType.AccountType(), input.Seed, input.Count)
	if err != nil {
		return signusecase.GenerateAuthKeyOutput{}, fmt.Errorf("failed to generate auth keys: %w", err)
	}

	return signusecase.GenerateAuthKeyOutput{
		GeneratedCount: len(keys),
	}, nil
}
