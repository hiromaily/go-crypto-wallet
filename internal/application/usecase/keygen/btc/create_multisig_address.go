package btc

import (
	"context"
	"encoding/hex"
	"fmt"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	portsBitcoin "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	infraKey "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type createMultisigAddressUseCase struct {
	btc                portsBitcoin.Bitcoiner
	authFullPubKeyRepo repository.AuthFullPubkeyRepositorier
	accountKeyRepo     repository.BTCAccountKeyRepositorier
	multisigAccount    *domainAccount.MultisigConfig
}

// NewCreateMultisigAddressUseCase creates a new CreateMultisigAddressUseCase
func NewCreateMultisigAddressUseCase(
	btc portsBitcoin.Bitcoiner,
	authFullPubKeyRepo repository.AuthFullPubkeyRepositorier,
	accountKeyRepo repository.BTCAccountKeyRepositorier,
	multisigAccount *domainAccount.MultisigConfig,
) keygenusecase.CreateMultisigAddressUseCase {
	return &createMultisigAddressUseCase{
		btc:                btc,
		authFullPubKeyRepo: authFullPubKeyRepo,
		accountKeyRepo:     accountKeyRepo,
		multisigAccount:    multisigAccount,
	}
}

func (u *createMultisigAddressUseCase) Create(
	ctx context.Context,
	input keygenusecase.CreateMultisigAddressInput,
) error {
	logger.Debug("addmultisigaddress",
		"account_type", input.AccountType.String(),
	)

	// Validate accountType
	if !u.multisigAccount.IsMultisigAccount(input.AccountType) {
		logger.Info("only multisig account is allowed")
		return nil
	}

	var requiredSig int
	var authFullPubKeys []string

	// Get fullPubKey from auth_fullpubkey table
	// AccountTypeDeposit: { //2:5+1
	//	2: {AuthType1, AuthType2, AuthType3, AuthType4, AuthType5},
	// },
	for sigCount, authTypes := range u.multisigAccount.MultiAccounts()[input.AccountType] {
		requiredSig = sigCount
		for _, authType := range authTypes {
			// Get record from auth_fullpubkey table
			fullPubKeyItem, err := u.authFullPubKeyRepo.GetOne(authType)
			if err != nil {
				return fmt.Errorf("fail to call authFullPubKeyRepo.GetOne() %s: %w", authType.String(), err)
			}

			// Derive account-specific public key if extended format is used
			var pubKey string
			if fullPubKeyItem.ExtendedPubKey != "" {
				// Extended format: derive account-specific key
				// The extended key is at m/49'/coin' level
				// We need to derive to m/49'/coin'/account' level
				pubKey, err = u.deriveAccountPublicKey(fullPubKeyItem.ExtendedPubKey, input.AccountType)
				if err != nil {
					return fmt.Errorf("fail to derive account public key for %s: %w", authType.String(), err)
				}
				logger.Debug("derived account-specific public key",
					"auth_type", authType.String(),
					"account_type", input.AccountType.String(),
					"pubkey", pubKey,
				)
			} else {
				// Legacy format: use compressed public key directly
				pubKey = fullPubKeyItem.FullPublicKey
				logger.Debug("using legacy compressed public key",
					"auth_type", authType.String(),
					"pubkey", pubKey,
				)
			}

			authFullPubKeys = append(authFullPubKeys, pubKey)
		}
		logger.Debug("don't repeat again")
	}

	// Get target addresses from account_key table, addr_status=AddrStatusPrivKeyImported
	accountKeyItems, err := u.accountKeyRepo.GetAllAddrStatus(
		input.AccountType, domainAddress.AddrStatusPrivKeyImported)
	if err != nil {
		return fmt.Errorf("fail to call accountKeyRepo.GetAllAddrStatus(%s): %w", input.AccountType.String(), err)
	}

	// Call bitcoinAPI `addmultisigaddress`
	for _, item := range accountKeyItems {
		addrs := make([]string, len(authFullPubKeys)+1)
		copy(addrs, authFullPubKeys)
		addrs[len(authFullPubKeys)] = item.FullPublicKey

		var resAddr *dtobtc.MultisigAddress
		resAddr, err = u.btc.AddMultisigAddress(
			requiredSig,
			addrs,
			fmt.Sprintf("multi_%s", input.AccountType), // this is not important
			input.AddressType,
		)
		if err != nil {
			// [Error] -5: no full public key for address mkPmdpo59gpU7ZioGYwwoMTQJjh7MiqUvd
			logger.Error(
				"fail to call btc.AddMultisigAddress()",
				"signature_count", requiredSig,
				"full public key for accountType", item.FullPublicKey,
				"full public key for authType", authFullPubKeys,
				"error", err,
			)
			continue
		}

		// Update generated multisig address, redeemScript, addrStatus
		item.MultisigAddress = resAddr.Address
		item.RedeemScript = resAddr.RedeemScript
		item.AddrStatus = domainAddress.AddrStatusMultisigAddressGenerated

		_, err = u.accountKeyRepo.UpdateMultisigAddr(input.AccountType, item)
		if err != nil {
			return fmt.Errorf("fail to call accountKeyRepo.UpdateMultisigAddr(%s): %w", input.AccountType.String(), err)
		}
	}

	return nil
}

// deriveAccountPublicKey derives an account-specific compressed public key from a coin-level extended public key.
//
// This function uses the common infraKey.DeriveAccountKey helper and extracts the compressed public key.
//
// Parameters:
//   - extendedPubKey: Extended public key at m/purpose'/coin' level (xpub/tpub format)
//   - accountType: Account type (deposit=0, payment=1, stored=2, etc.)
//
// Returns:
//   - Compressed public key (33 bytes, hex-encoded)
//   - Error if derivation fails
func (*createMultisigAddressUseCase) deriveAccountPublicKey(
	extendedPubKey string,
	accountType domainAccount.AccountType,
) (string, error) {
	// Use common derivation helper
	accountKey, err := infraKey.DeriveAccountKey(extendedPubKey, accountType)
	if err != nil {
		return "", err
	}

	// Get the public key
	pubKey, err := accountKey.ECPubKey()
	if err != nil {
		return "", fmt.Errorf("failed to get public key: %w", err)
	}

	// Serialize as compressed public key (33 bytes)
	compressedPubKey := pubKey.SerializeCompressed()

	// Convert to hex string
	pubKeyHex := hex.EncodeToString(compressedPubKey)

	return pubKeyHex, nil
}
