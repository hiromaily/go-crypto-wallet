package shared

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	portsWallet "github.com/hiromaily/go-crypto-wallet/internal/application/ports/wallet"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type generateHDWalletUseCase struct {
	repo         persistence.HDWalletRepo
	keygen       portsWallet.Generator
	coinTypeCode domainCoin.CoinTypeCode
}

// NewGenerateHDWalletUseCase creates a new GenerateHDWalletUseCase
func NewGenerateHDWalletUseCase(
	repo persistence.HDWalletRepo,
	keygen portsWallet.Generator,
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

	// Generate HD wallet keys with account xpriv
	walletKeys, accountXpriv, err := u.generateHDKeyWithAccountXpriv(
		input.AccountType, input.Seed, uint32(idxFrom), input.Count,
	)
	if err != nil {
		return keygenusecase.GenerateHDWalletOutput{}, fmt.Errorf("fail to generate HD key: %w", err)
	}

	// Insert key information to account_key_table / auth_account_key_table with account xpriv
	err = u.repo.Insert(walletKeys, accountXpriv, idxFrom, u.coinTypeCode, input.AccountType, u.keygen.KeyType())
	if err != nil {
		return keygenusecase.GenerateHDWalletOutput{}, fmt.Errorf("fail to call repo.Insert(): %w", err)
	}

	return keygenusecase.GenerateHDWalletOutput{
		GeneratedCount: len(walletKeys),
	}, nil
}

// generateHDKeyWithAccountXpriv generates HD wallet keys and returns account-level extended private key.
// The account xpriv can be used for BIP32 derivation at different address indices.
func (u *generateHDWalletUseCase) generateHDKeyWithAccountXpriv(
	accountType domainAccount.AccountType,
	seed []byte,
	idxFrom,
	count uint32,
) (keys []domainKey.WalletKey, accountXpriv string, err error) {
	// Type assert to access CreateKeyWithAccountXpriv
	// This method is specific to HDKey implementation
	type accountXprivGenerator interface {
		CreateKeyWithAccountXpriv(
			seed []byte,
			accountType domainAccount.AccountType,
			idxFrom, count uint32,
		) ([]domainKey.WalletKey, string, error)
	}

	gen, ok := u.keygen.(accountXprivGenerator)
	if !ok {
		// Fallback: generate keys without xpriv (backward compatibility)
		logger.Warn(
			"generator does not implement CreateKeyWithAccountXpriv, " +
				"falling back to CreateKey (no xpriv will be stored)",
		)
		keys, err = u.keygen.CreateKey(seed, accountType, idxFrom, count)
		if err != nil {
			return nil, "", fmt.Errorf("fail to call keygen.CreateKey(): %w", err)
		}
		return keys, "", nil
	}

	logger.Debug("using CreateKeyWithAccountXpriv to generate keys with account xpriv")

	// Generate keys with account xpriv
	keys, accountXpriv, err = gen.CreateKeyWithAccountXpriv(seed, accountType, idxFrom, count)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call keygen.CreateKeyWithAccountXpriv(): %w", err)
	}
	return keys, accountXpriv, nil
}
