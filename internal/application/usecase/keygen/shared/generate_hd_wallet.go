package shared

import (
	"context"
	"fmt"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type generateHDWalletUseCase struct {
	repo         cold.HDWalletRepo
	keygen       key.Generator
	coinTypeCode domainCoin.CoinTypeCode
}

// NewGenerateHDWalletUseCase creates a new GenerateHDWalletUseCase
func NewGenerateHDWalletUseCase(
	repo cold.HDWalletRepo,
	keygen key.Generator,
	coinTypeCode domainCoin.CoinTypeCode,
) keygenusecase.GenerateHDWalletUseCase {
	return &generateHDWalletUseCase{
		repo:         repo,
		keygen:       keygen,
		coinTypeCode: coinTypeCode,
	}
}

func (u *generateHDWalletUseCase) Generate(
	ctx context.Context,
	input keygenusecase.GenerateHDWalletInput,
) (keygenusecase.GenerateHDWalletOutput, error) {
	logger.Debug("generate HDWallet", "account_type", input.AccountType.String())

	// Get latest index
	idxFrom, err := u.repo.GetMaxIndex(input.AccountType)
	if err != nil {
		logger.Info(err.Error())
		return keygenusecase.GenerateHDWalletOutput{
			GeneratedCount: 0,
		}, nil
	}
	logger.Debug("max_index",
		"account_type", input.AccountType.String(),
		"current_index", idxFrom,
	)

	// Generate HD wallet keys
	walletKeys, err := u.generateHDKey(input.AccountType, input.Seed, uint32(idxFrom), input.Count)
	if err != nil {
		return keygenusecase.GenerateHDWalletOutput{}, fmt.Errorf("fail to generate HD key: %w", err)
	}

	// Insert key information to account_key_table / auth_account_key_table
	err = u.repo.Insert(walletKeys, idxFrom, u.coinTypeCode, input.AccountType)
	if err != nil {
		return keygenusecase.GenerateHDWalletOutput{}, fmt.Errorf("fail to call repo.Insert(): %w", err)
	}

	return keygenusecase.GenerateHDWalletOutput{
		GeneratedCount: len(walletKeys),
	}, nil
}

// generateHDKey generates HD wallet keys
func (u *generateHDWalletUseCase) generateHDKey(
	accountType domainAccount.AccountType,
	seed []byte,
	idxFrom,
	count uint32,
) ([]domainKey.WalletKey, error) {
	// Generate key
	walletKeys, err := u.keygen.CreateKey(seed, accountType, idxFrom, count)
	if err != nil {
		return nil, fmt.Errorf("fail to call keygen.CreateKey(): %w", err)
	}
	return walletKeys, nil
}
