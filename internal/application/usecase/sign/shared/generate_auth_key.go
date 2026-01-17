package shared

import (
	"context"
	"fmt"

	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	portsWallet "github.com/hiromaily/go-crypto-wallet/internal/application/ports/wallet"
	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type generateAuthKeyUseCase struct {
	repo         repocold.HDWalletRepo
	keygen       portsWallet.Generator
	coinTypeCode domainCoin.CoinTypeCode
}

// NewGenerateAuthKeyUseCase creates a new GenerateAuthKeyUseCase for sign wallet
func NewGenerateAuthKeyUseCase(
	repo repocold.HDWalletRepo,
	keygen portsWallet.Generator,
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
	// For descriptor-based multisig, generate keys for each multisig account
	// (deposit, payment, stored) instead of using auth account index (11, 12, etc.).
	// This ensures auth keys match the descriptor derivation paths:
	//   - deposit: m/purpose'/coin'/0'/0/0
	//   - payment: m/purpose'/coin'/1'/0/0
	//   - stored: m/purpose'/coin'/2'/0/0

	multisigAccounts := []domainAccount.AccountType{
		domainAccount.AccountTypeDeposit, // BIP44 account index 0
		domainAccount.AccountTypePayment, // BIP44 account index 1
		domainAccount.AccountTypeStored,  // BIP44 account index 2
	}

	totalGenerated := 0

	for _, accountType := range multisigAccounts {
		logger.Debug("generating HD wallet keys",
			"auth_type", input.AuthType.String(),
			"account_type", accountType.String())

		// Get latest index for this account
		idxFrom, err := u.repo.GetMaxIndex(ctx, accountType)
		if err != nil {
			logger.Info(err.Error())
			continue
		}
		logger.Debug("max_index",
			"account_type", accountType.String(),
			"current_index", idxFrom,
		)

		// Generate HD wallet keys for this account
		walletKeys, accountXpriv, err := u.generateHDKeyWithAccountXpriv(
			accountType, input.Seed, uint32(idxFrom), input.Count,
		)
		if err != nil {
			return signusecase.GenerateAuthKeyOutput{},
				fmt.Errorf("fail to generate HD key for %s: %w", accountType.String(), err)
		}

		// Insert key information to auth_account_key_table with account xpriv
		err = u.repo.Insert(walletKeys, accountXpriv, idxFrom, u.coinTypeCode, accountType, u.keygen.KeyType())
		if err != nil {
			return signusecase.GenerateAuthKeyOutput{},
				fmt.Errorf("fail to insert keys for %s: %w", accountType.String(), err)
		}

		totalGenerated += len(walletKeys)
	}

	return signusecase.GenerateAuthKeyOutput{
		GeneratedCount: totalGenerated,
	}, nil
}

// generateHDKeyWithAccountXpriv generates HD wallet keys and returns account-level extended private key.
// The account xpriv can be used for BIP32 derivation at different address indices.
func (u *generateAuthKeyUseCase) generateHDKeyWithAccountXpriv(
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
		keys, err = u.keygen.CreateKey(seed, accountType, idxFrom, count)
		if err != nil {
			return nil, "", fmt.Errorf("fail to call keygen.CreateKey(): %w", err)
		}
		return keys, "", nil
	}

	// Generate keys with account xpriv
	keys, accountXpriv, err = gen.CreateKeyWithAccountXpriv(seed, accountType, idxFrom, count)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call keygen.CreateKeyWithAccountXpriv(): %w", err)
	}
	return keys, accountXpriv, nil
}
