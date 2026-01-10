package btc

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"

	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/btc"
	infraKey "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
)

type generateDescriptorUseCase struct {
	descriptorService  *btc.DescriptorService
	chainConfig        *chaincfg.Params
	authFullPubKeyRepo persistence.AuthFullPubkeyRepositorier
	accountKeyRepo     persistence.BTCAccountKeyRepositorier
	multisigConfig     *domainAccount.MultisigConfig
}

// NewGenerateDescriptorUseCase creates a descriptor generation use case.
func NewGenerateDescriptorUseCase(
	descriptorService *btc.DescriptorService,
	chainConfig *chaincfg.Params,
	authFullPubKeyRepo persistence.AuthFullPubkeyRepositorier,
	accountKeyRepo persistence.BTCAccountKeyRepositorier,
	multisigConfig *domainAccount.MultisigConfig,
) keygenusecase.GenerateDescriptorUseCase {
	return &generateDescriptorUseCase{
		descriptorService:  descriptorService,
		chainConfig:        chainConfig,
		authFullPubKeyRepo: authFullPubKeyRepo,
		accountKeyRepo:     accountKeyRepo,
		multisigConfig:     multisigConfig,
	}
}

func (u *generateDescriptorUseCase) Generate(
	ctx context.Context,
	input keygenusecase.GenerateDescriptorInput,
) (keygenusecase.GenerateDescriptorOutput, error) {
	_ = ctx

	isMultisig := u.multisigConfig.IsMultisigAccount(input.AccountType)

	var (
		descriptor string
		err        error
	)

	if isMultisig {
		descriptor, err = u.generateMultisigDescriptor(input)
	} else {
		descriptor, err = u.generateSingleSigDescriptor(input)
	}

	if err != nil {
		return keygenusecase.GenerateDescriptorOutput{},
			fmt.Errorf("failed to generate descriptor for %s: %w", input.AccountType.String(), err)
	}

	return keygenusecase.GenerateDescriptorOutput{
		Descriptor:  descriptor,
		AccountType: input.AccountType,
		AddressType: input.AddressType,
		IsMultisig:  isMultisig,
	}, nil
}

func (u *generateDescriptorUseCase) generateSingleSigDescriptor(
	input keygenusecase.GenerateDescriptorInput,
) (string, error) {
	accountKey, err := u.accountKeyRepo.GetOneMaxID(input.AccountType)
	if err != nil {
		return "", fmt.Errorf("get account key: %w", err)
	}
	if accountKey == nil {
		return "", fmt.Errorf("no account key found for %s", input.AccountType.String())
	}

	xpub, err := hdkeychain.NewKeyFromString(accountKey.FullPublicKey)
	if err != nil {
		return "", fmt.Errorf("parse extended public key: %w", err)
	}

	fp, err := infraKey.FingerprintFromExtendedKey(accountKey.FullPublicKey)
	if err != nil {
		return "", fmt.Errorf("calculate fingerprint: %w", err)
	}

	derivationPath, err := derivationPathForAddress(input.AddressType, false, input.AccountType, u.chainConfig)
	if err != nil {
		return "", err
	}

	switch input.AddressType {
	case domainAddress.AddrTypeTaproot:
		return u.descriptorService.GenerateTaprootDescriptor(fp.String(), derivationPath, xpub, input.IsChange)
	case domainAddress.AddrTypeBech32:
		return u.descriptorService.GenerateBech32Descriptor(fp.String(), derivationPath, xpub, input.IsChange)
	case domainAddress.AddrTypeP2shSegwit:
		return u.descriptorService.GenerateP2SHSegWitDescriptor(fp.String(), derivationPath, xpub, input.IsChange)
	case domainAddress.AddrTypeLegacy:
		return u.descriptorService.GenerateP2PKHDescriptor(fp.String(), derivationPath, xpub, input.IsChange)
	case domainAddress.AddrTypeBCHCashAddr, domainAddress.AddrTypeETH:
		return "", fmt.Errorf("unsupported address type for Bitcoin descriptors: %s", input.AddressType)
	default:
		return "", fmt.Errorf("unsupported address type for single-sig: %s", input.AddressType)
	}
}

func (u *generateDescriptorUseCase) generateMultisigDescriptor(
	input keygenusecase.GenerateDescriptorInput,
) (string, error) {
	multiConfig := u.multisigConfig.MultiAccounts()[input.AccountType]
	if len(multiConfig) == 0 {
		return "", fmt.Errorf("multisig account not configured for %s", input.AccountType.String())
	}

	requiredSigs, authTypes, err := selectRequiredSigConfig(multiConfig, input.RequiredSigs)
	if err != nil {
		return "", err
	}

	signers, err := u.buildMultisigSigners(authTypes, input.AddressType, input.AccountType)
	if err != nil {
		return "", err
	}

	if input.AddressType == domainAddress.AddrTypeTaproot {
		return u.descriptorService.GenerateTaprootScriptPathDescriptor(signers, input.IsChange)
	}

	return u.descriptorService.GenerateMultisigDescriptor(requiredSigs, signers, input.IsChange)
}

func (u *generateDescriptorUseCase) buildMultisigSigners(
	authTypes []domainAccount.AuthType,
	addressType domainAddress.AddrType,
	accountType domainAccount.AccountType,
) ([]btc.MultisigSigner, error) {
	derivationPath, err := derivationPathForAddress(addressType, true, accountType, u.chainConfig)
	if err != nil {
		return nil, err
	}

	signers := make([]btc.MultisigSigner, 0, len(authTypes))
	for _, authType := range authTypes {
		authKey, err := u.authFullPubKeyRepo.GetOne(authType)
		if err != nil {
			return nil, fmt.Errorf("get auth full pubkey for %s: %w", authType.String(), err)
		}
		if authKey == nil {
			return nil, fmt.Errorf("auth full pubkey not found for %s", authType.String())
		}

		// Use ExtendedPubKey (new format) if available, otherwise fall back to FullPublicKey (legacy)
		coinLevelExtendedKey := authKey.ExtendedPubKey
		if coinLevelExtendedKey == "" {
			// Legacy format: FullPublicKey contains compressed pubkey, not extended key
			return nil, fmt.Errorf(
				"extended public key not found for %s (legacy compressed pubkey format not supported for descriptors)",
				authType.String(),
			)
		}

		// The extended key is at m/49'/coin' level (from sign wallet export)
		// Derive to m/49'/coin'/account' level for this specific account
		accountExtendedKey, err := u.deriveAccountExtendedKey(coinLevelExtendedKey, accountType)
		if err != nil {
			return nil, fmt.Errorf("derive account extended key for %s: %w", authType.String(), err)
		}

		xpub, err := hdkeychain.NewKeyFromString(accountExtendedKey)
		if err != nil {
			return nil, fmt.Errorf("parse account extended key for %s: %w", authType.String(), err)
		}

		if u.chainConfig != nil && !xpub.IsForNet(u.chainConfig) {
			return nil, fmt.Errorf("account extended key network mismatch for %s", authType.String())
		}

		var fp string
		if authKey.Fingerprint != nil {
			fp = authKey.Fingerprint.String()
		} else {
			finger, err := infraKey.FingerprintFromExtendedKey(accountExtendedKey)
			if err != nil {
				return nil, fmt.Errorf("calculate fingerprint for %s: %w", authType.String(), err)
			}
			fp = finger.String()
		}

		signers = append(signers, btc.MultisigSigner{
			Fingerprint:    fp,
			DerivationPath: derivationPath,
			ExtendedKey:    xpub,
		})
	}

	return signers, nil
}

func derivationPathForAddress(
	addrType domainAddress.AddrType,
	isMultisig bool,
	accountType domainAccount.AccountType,
	chainConfig *chaincfg.Params,
) (string, error) {
	// Use BIP44AccountIndex for derivation paths (deposit=0, payment=1, stored=2)
	// NOT Uint32() which returns database storage values (deposit=1, payment=2, stored=3)
	accountIndex := accountType.BIP44AccountIndex()

	// Determine coin index based on network: 0' for mainnet, 1' for testnet/regtest
	coinIndex := "1'" // Default for testnet/regtest
	if chainConfig != nil && chainConfig.Net == chaincfg.MainNetParams.Net {
		coinIndex = "0'"
	}

	switch addrType {
	case domainAddress.AddrTypeLegacy:
		return fmt.Sprintf("/44'/%s/%d'", coinIndex, accountIndex), nil
	case domainAddress.AddrTypeP2shSegwit:
		if isMultisig {
			// For multisig with BIP49 keys, show actual derivation path
			// The xpubs are derived from m/49'/coin' (hardened coin level) to
			// m/49'/coin/account (non-hardened account level) because we're deriving
			// from extended public keys.
			// Note: Account index is non-hardened (%d, not %d') for xpub derivation.
			return fmt.Sprintf("/49'/%s/%d", coinIndex, accountIndex), nil
		}
		return fmt.Sprintf("/49'/%s/%d'", coinIndex, accountIndex), nil
	case domainAddress.AddrTypeBech32:
		if isMultisig {
			// Native SegWit multisig: show actual derivation path
			// Account index is non-hardened for xpub derivation
			return fmt.Sprintf("/84'/%s/%d", coinIndex, accountIndex), nil
		}
		return fmt.Sprintf("/84'/%s/%d'", coinIndex, accountIndex), nil
	case domainAddress.AddrTypeTaproot:
		return fmt.Sprintf("/86'/%s/%d'", coinIndex, accountIndex), nil
	case domainAddress.AddrTypeBCHCashAddr, domainAddress.AddrTypeETH:
		return "", fmt.Errorf("unsupported address type for Bitcoin: %s", addrType)
	default:
		return "", fmt.Errorf("unsupported address type: %s", addrType)
	}
}

func selectRequiredSigConfig(
	config map[int][]domainAccount.AuthType,
	requested int,
) (int, []domainAccount.AuthType, error) {
	if len(config) == 0 {
		return 0, nil, errors.New("no multisig configuration available")
	}

	if requested > 0 {
		auths, ok := config[requested]
		if !ok {
			return 0, nil, fmt.Errorf("required signatures %d not configured", requested)
		}
		return requested, auths, nil
	}

	// Default to the smallest required-sigs configuration if not specified.
	requiredKeys := make([]int, 0, len(config))
	for k := range config {
		requiredKeys = append(requiredKeys, k)
	}
	sort.Ints(requiredKeys)

	req := requiredKeys[0]
	return req, config[req], nil
}

// deriveAccountExtendedKey derives an account-specific extended public key from a coin-level extended public key.
//
// This is a wrapper around infraKey.DeriveAccountKey that returns the extended key as a string.
//
// Parameters:
//   - coinLevelExtendedKey: Extended public key at m/purpose'/coin' level (xpub/tpub format)
//   - accountType: Account type (deposit=0, payment=1, stored=2, etc.)
//
// Returns:
//   - Account-level extended public key (xpub/tpub format)
//   - Error if derivation fails
func (*generateDescriptorUseCase) deriveAccountExtendedKey(
	coinLevelExtendedKey string,
	accountType domainAccount.AccountType,
) (string, error) {
	// Use common derivation helper
	accountKey, err := infraKey.DeriveAccountKey(coinLevelExtendedKey, accountType)
	if err != nil {
		return "", err
	}

	// Convert to string (xpub format)
	return accountKey.String(), nil
}
