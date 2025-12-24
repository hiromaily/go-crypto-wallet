package xrp

import (
	"context"
	"errors"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen"
	domainKey "github.com/hiromaily/go-crypto-wallet/pkg/domain/key"
	xrpkeygensrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/keygen/xrp"
)

type generateKeyUseCase struct {
	keyGenerator *xrpkeygensrv.XRPKeyGenerate
}

// NewGenerateKeyUseCase creates a new GenerateKeyUseCase
func NewGenerateKeyUseCase(keyGenerator *xrpkeygensrv.XRPKeyGenerate) keygen.GenerateKeyUseCase {
	return &generateKeyUseCase{
		keyGenerator: keyGenerator,
	}
}

func (u *generateKeyUseCase) Generate(ctx context.Context, input keygen.GenerateKeyInput) error {
	// Convert interface{} to []domainKey.WalletKey
	walletKeys, ok := input.WalletKeys.([]domainKey.WalletKey)
	if !ok {
		return errors.New("invalid wallet keys type")
	}

	if err := u.keyGenerator.Generate(input.AccountType, input.IsKeyPair, walletKeys); err != nil {
		return fmt.Errorf("failed to generate XRP keys: %w", err)
	}
	return nil
}
