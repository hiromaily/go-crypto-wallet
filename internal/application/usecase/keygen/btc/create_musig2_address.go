package btc

import (
	"context"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/btc"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/config/account"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type createMuSig2AddressUseCase struct {
	musig2Service      *btc.MuSig2Service
	chainConfig        *chaincfg.Params
	authFullPubKeyRepo cold.AuthFullPubkeyRepositorier
	accountKeyRepo     cold.AccountKeyRepositorier
	multisigAccount    account.MultisigAccounter
}

// NewCreateMuSig2AddressUseCase creates a new CreateMuSig2AddressUseCase
func NewCreateMuSig2AddressUseCase(
	musig2Service *btc.MuSig2Service,
	chainConfig *chaincfg.Params,
	authFullPubKeyRepo cold.AuthFullPubkeyRepositorier,
	accountKeyRepo cold.AccountKeyRepositorier,
	multisigAccount account.MultisigAccounter,
) keygenusecase.CreateMuSig2AddressUseCase {
	return &createMuSig2AddressUseCase{
		musig2Service:      musig2Service,
		chainConfig:        chainConfig,
		authFullPubKeyRepo: authFullPubKeyRepo,
		accountKeyRepo:     accountKeyRepo,
		multisigAccount:    multisigAccount,
	}
}

func (u *createMuSig2AddressUseCase) Create(
	ctx context.Context,
	input keygenusecase.CreateMuSig2AddressInput,
) error {
	logger.Debug("create musig2 taproot address",
		"account_type", input.AccountType.String(),
	)

	// Validate accountType
	if !u.multisigAccount.IsMultisigAccount(input.AccountType) {
		return fmt.Errorf("account '%s' is not a multisig account", input.AccountType.String())
	}

	// Get all signer public keys from auth_fullpubkey table
	var authFullPubKeys []*btcec.PublicKey
	for _, authTypes := range u.multisigAccount.MultiAccounts()[input.AccountType] {
		for _, authType := range authTypes {
			// Get record from auth_fullpubkey table
			fullPubKeyItem, err := u.authFullPubKeyRepo.GetOne(authType)
			if err != nil {
				return fmt.Errorf("fail to call authFullPubKeyRepo.GetOne() %s: %w", authType.String(), err)
			}

			// Parse public key
			pubKey, err := btcec.ParsePubKey([]byte(fullPubKeyItem.FullPublicKey))
			if err != nil {
				return fmt.Errorf("fail to parse public key for %s: %w", authType.String(), err)
			}
			authFullPubKeys = append(authFullPubKeys, pubKey)
		}
		logger.Debug("collected auth public keys", "count", len(authFullPubKeys))
	}

	// Get target addresses from account_key table, addr_status=AddrStatusPrivKeyImported
	accountKeyItems, err := u.accountKeyRepo.GetAllAddrStatus(input.AccountType, address.AddrStatusPrivKeyImported)
	if err != nil {
		return fmt.Errorf("fail to call accountKeyRepo.GetAllAddrStatus(%s): %w", input.AccountType.String(), err)
	}

	// Process each account key to create MuSig2 Taproot address
	for _, item := range accountKeyItems {
		// Parse account private key (WIF format)
		wif, err := btcutil.DecodeWIF(item.WalletImportFormat)
		if err != nil {
			logger.Error(
				"fail to decode WIF",
				"account_key_id", item.ID,
				"error", err,
			)
			continue
		}
		accountPrivKey := wif.PrivKey

		// Parse account public key
		accountPubKey, err := btcec.ParsePubKey([]byte(item.FullPublicKey))
		if err != nil {
			logger.Error(
				"fail to parse account public key",
				"full_public_key", item.FullPublicKey,
				"error", err,
			)
			continue
		}

		// Combine all public keys (auth keys + account key)
		allPubKeys := make([]*btcec.PublicKey, len(authFullPubKeys)+1)
		copy(allPubKeys, authFullPubKeys)
		allPubKeys[len(authFullPubKeys)] = accountPubKey

		// Create MuSig2 context with Taproot tweaking
		// Use the account's private key for context creation
		musig2Ctx, err := u.musig2Service.CreateContextWithTaproot(accountPrivKey, allPubKeys, true)
		if err != nil {
			logger.Error(
				"fail to create MuSig2 context",
				"account_key_id", item.ID,
				"error", err,
			)
			continue
		}

		// Get aggregated public key
		aggregatedPubKey, err := u.musig2Service.GetCombinedKey(musig2Ctx)
		if err != nil {
			logger.Error(
				"fail to get combined key",
				"account_key_id", item.ID,
				"error", err,
			)
			continue
		}

		// Create Taproot address from aggregated key
		taprootKey := txscript.ComputeTaprootKeyNoScript(aggregatedPubKey)
		witnessProg := taprootKey.SerializeCompressed()[1:] // Remove parity byte

		taprootAddr, err := btcutil.NewAddressTaproot(witnessProg, u.chainConfig)
		if err != nil {
			logger.Error(
				"fail to create Taproot address",
				"account_key_id", item.ID,
				"error", err,
			)
			continue
		}

		// Update taproot address and addr_status in database
		item.TaprootAddress.String = taprootAddr.EncodeAddress()
		item.TaprootAddress.Valid = true
		item.MultisigAddress = taprootAddr.EncodeAddress() // Store in multisig_address for compatibility
		item.AddrStatus = address.AddrStatusMultisigAddressGenerated.Int8()

		_, err = u.accountKeyRepo.UpdateMultisigAddr(input.AccountType, item)
		if err != nil {
			return fmt.Errorf("fail to call accountKeyRepo.UpdateMultisigAddr(%s): %w", input.AccountType.String(), err)
		}

		logger.Debug("MuSig2 Taproot address created",
			"account_key_id", item.ID,
			"taproot_address", taprootAddr.EncodeAddress(),
		)
	}

	return nil
}
