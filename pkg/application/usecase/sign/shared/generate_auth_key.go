package shared

import (
	"context"
	"fmt"

	signusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/sign"
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/pkg/domain/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/keygen/shared"
)

type generateAuthKeyUseCase struct {
	repo         shared.HDWalletRepo
	keygen       key.Generator
	coinTypeCode domainCoin.CoinTypeCode
}

// NewGenerateAuthKeyUseCase creates a new GenerateAuthKeyUseCase for sign wallet
func NewGenerateAuthKeyUseCase(
	repo shared.HDWalletRepo,
	keygen key.Generator,
	coinTypeCode domainCoin.CoinTypeCode,
) signusecase.GenerateAuthKeyUseCase {
	return &generateAuthKeyUseCase{
		repo:         repo,
		keygen:       keygen,
		coinTypeCode: coinTypeCode,
	}
}

func (u *generateAuthKeyUseCase) Generate(
	ctx context.Context,
	input signusecase.GenerateAuthKeyInput,
) (signusecase.GenerateAuthKeyOutput, error) {
	accountType := input.AuthType.AccountType()
	logger.Debug("generate HDWallet", "account_type", accountType.String())

	// Get latest index
	idxFrom, err := u.repo.GetMaxIndex(accountType)
	if err != nil {
		logger.Info(err.Error())
		return signusecase.GenerateAuthKeyOutput{
			GeneratedCount: 0,
		}, nil
	}
	logger.Debug("max_index",
		"account_type", accountType.String(),
		"current_index", idxFrom,
	)

	// Generate HD wallet keys
	walletKeys, err := u.generateHDKey(accountType, input.Seed, uint32(idxFrom), input.Count)
	if err != nil {
		return signusecase.GenerateAuthKeyOutput{}, fmt.Errorf("fail to generate HD key: %w", err)
	}

	// Insert key information to auth_account_key_table
	err = u.repo.Insert(walletKeys, idxFrom, u.coinTypeCode, accountType)
	if err != nil {
		return signusecase.GenerateAuthKeyOutput{}, fmt.Errorf("fail to call repo.Insert(): %w", err)
	}

	return signusecase.GenerateAuthKeyOutput{
		GeneratedCount: len(walletKeys),
	}, nil
}

// generateHDKey generates HD wallet keys
func (u *generateAuthKeyUseCase) generateHDKey(
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
