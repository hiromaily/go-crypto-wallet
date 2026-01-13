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
	domainAuth "github.com/hiromaily/go-crypto-wallet/internal/domain/auth"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/btc"
	infraKey "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type generateDescriptorUseCase struct {
	descriptorService  *btc.DescriptorService
	chainConfig        *chaincfg.Params
	authFullPubKeyRepo persistence.AuthFullPubkeyRepositorier
	accountKeyRepo     persistence.BTCAccountKeyRepositorier
	seedRepo           persistence.SeedRepositorier
	coinTypeCode       domainCoin.CoinTypeCode
	multisigConfig     *domainAccount.MultisigConfig
}

// NewGenerateDescriptorUseCase creates a descriptor generation use case.
func NewGenerateDescriptorUseCase(
	descriptorService *btc.DescriptorService,
	chainConfig *chaincfg.Params,
	authFullPubKeyRepo persistence.AuthFullPubkeyRepositorier,
	accountKeyRepo persistence.BTCAccountKeyRepositorier,
	seedRepo persistence.SeedRepositorier,
	coinTypeCode domainCoin.CoinTypeCode,
	multisigConfig *domainAccount.MultisigConfig,
) keygenusecase.GenerateDescriptorUseCase {
	return &generateDescriptorUseCase{
		descriptorService:  descriptorService,
		chainConfig:        chainConfig,
		authFullPubKeyRepo: authFullPubKeyRepo,
		accountKeyRepo:     accountKeyRepo,
		seedRepo:           seedRepo,
		coinTypeCode:       coinTypeCode,
		multisigConfig:     multisigConfig,
	}
}

func (u *generateDescriptorUseCase) Generate(
	ctx context.Context,
	input keygenusecase.GenerateDescriptorInput,
) (keygenusecase.GenerateDescriptorOutput, error) {
	isMultisig := u.multisigConfig.IsMultisigAccount(input.AccountType)

	var (
		descriptor string
		err        error
	)

	if isMultisig {
		descriptor, err = u.generateMultisigDescriptor(ctx, input)
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
	logger.Debug("generating single-sig descriptor",
		"account_type", input.AccountType.String(),
		"address_type", input.AddressType.String(),
		"is_change", input.IsChange,
	)

	accountKey, err := u.accountKeyRepo.GetOneMaxID(input.AccountType)
	if err != nil {
		return "", fmt.Errorf("get account key: %w", err)
	}
	if accountKey == nil {
		return "", fmt.Errorf("no account key found for %s", input.AccountType.String())
	}

	logger.Debug("found account key",
		"account_type", input.AccountType.String(),
		"key_id", accountKey.ID,
		"key_type", accountKey.KeyType,
		"has_extended_privkey", accountKey.AccountExtendedPrivkey != nil,
	)

	// Check if AccountExtendedPrivkey is available (new format)
	if accountKey.AccountExtendedPrivkey == nil || *accountKey.AccountExtendedPrivkey == "" {
		return "", fmt.Errorf(
			"account extended private key not found for %s - "+
				"descriptors require keys generated with extended key support",
			input.AccountType.String(),
		)
	}

	// Verify that the stored key_type matches the requested address_type
	// This prevents trying to generate P2TR descriptors from BIP44 keys, etc.
	expectedKeyType, err := input.AddressType.ToKeyType()
	if err != nil {
		logger.Debug("unsupported address type for key derivation",
			"address_type", input.AddressType.String(),
			"error", err.Error(),
		)
		return "", fmt.Errorf("unsupported address type %q for key derivation: %w", input.AddressType, err)
	}

	logger.Debug("validating key type match",
		"account_type", input.AccountType.String(),
		"address_type", input.AddressType.String(),
		"stored_key_type", accountKey.KeyType,
		"expected_key_type", string(expectedKeyType),
		"match", accountKey.KeyType == string(expectedKeyType),
	)

	if accountKey.KeyType != string(expectedKeyType) {
		return "", fmt.Errorf(
			"key type mismatch: stored key_type=%s, but address_type=%s requires key_type=%s",
			accountKey.KeyType,
			input.AddressType.String(),
			expectedKeyType,
		)
	}

	// Derive account-level extended public key from extended private key
	xpriv, err := hdkeychain.NewKeyFromString(*accountKey.AccountExtendedPrivkey)
	if err != nil {
		return "", fmt.Errorf("parse account extended private key: %w", err)
	}

	// Convert to extended public key
	xpub, err := xpriv.Neuter()
	if err != nil {
		return "", fmt.Errorf("derive extended public key: %w", err)
	}

	// Verify network match
	if u.chainConfig != nil && !xpub.IsForNet(u.chainConfig) {
		return "", errors.New("extended key network mismatch")
	}

	// Calculate fingerprint from extended key
	fp, err := infraKey.FingerprintFromExtendedKey(xpub.String())
	if err != nil {
		return "", fmt.Errorf("calculate fingerprint: %w", err)
	}

	derivationPath, err := derivationPathForAddress(input.AddressType, false, input.AccountType, u.chainConfig)
	if err != nil {
		return "", err
	}

	logger.Debug("generating descriptor",
		"address_type", input.AddressType.String(),
		"derivation_path", derivationPath,
		"fingerprint", fp.String(),
	)

	var descriptor string
	fpStr := fp.String()
	switch input.AddressType {
	case domainAddress.AddrTypeTaproot:
		descriptor, err = u.descriptorService.GenerateTaprootDescriptor(
			fpStr, derivationPath, xpub, input.IsChange,
		)
	case domainAddress.AddrTypeBech32:
		descriptor, err = u.descriptorService.GenerateBech32Descriptor(
			fpStr, derivationPath, xpub, input.IsChange,
		)
	case domainAddress.AddrTypeP2shSegwit:
		descriptor, err = u.descriptorService.GenerateP2SHSegWitDescriptor(
			fpStr, derivationPath, xpub, input.IsChange,
		)
	case domainAddress.AddrTypeLegacy:
		descriptor, err = u.descriptorService.GenerateP2PKHDescriptor(
			fpStr, derivationPath, xpub, input.IsChange,
		)
	case domainAddress.AddrTypeBCHCashAddr, domainAddress.AddrTypeETH:
		return "", fmt.Errorf("unsupported address type for Bitcoin descriptors: %s", input.AddressType)
	default:
		return "", fmt.Errorf("unsupported address type for single-sig: %s", input.AddressType)
	}

	if err != nil {
		return "", fmt.Errorf("generate descriptor for %s: %w", input.AddressType, err)
	}

	logger.Debug("successfully generated descriptor",
		"account_type", input.AccountType.String(),
		"address_type", input.AddressType.String(),
	)

	return descriptor, nil
}

func (u *generateDescriptorUseCase) generateMultisigDescriptor(
	ctx context.Context,
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

	signers, err := u.buildMultisigSigners(ctx, authTypes, input.AddressType, input.AccountType)
	if err != nil {
		return "", err
	}

	if input.AddressType == domainAddress.AddrTypeTaproot {
		return u.descriptorService.GenerateTaprootScriptPathDescriptor(signers, input.IsChange)
	}

	// P2SH-SegWit (BIP49) requires sh() wrapper around wsh()
	isP2SHWrapped := input.AddressType == domainAddress.AddrTypeP2shSegwit

	return u.descriptorService.GenerateMultisigDescriptor(requiredSigs, signers, input.IsChange, isP2SHWrapped)
}

func (u *generateDescriptorUseCase) buildMultisigSigners(
	ctx context.Context,
	authTypes []domainAccount.AuthType,
	addressType domainAddress.AddrType,
	accountType domainAccount.AccountType,
) ([]btc.MultisigSigner, error) {
	derivationPath, err := derivationPathForAddress(addressType, true, accountType, u.chainConfig)
	if err != nil {
		return nil, err
	}

	// Determine BIP purpose from address type
	purpose, err := domainAuth.PurposeForAddressType(addressType.String())
	if err != nil {
		return nil, fmt.Errorf("failed to determine purpose for address type %s: %w", addressType.String(), err)
	}

	// Reserve space for keygen key + auth keys
	signers := make([]btc.MultisigSigner, 0, len(authTypes)+1)

	// CRITICAL FIX: Include keygen wallet's own key in multisig
	// Without this, descriptors only contain auth keys, resulting in (N-1)-of-(N-1) instead of N-of-N
	keygenSigner, err := u.buildKeygenSigner(ctx, accountType, addressType, derivationPath)
	if err != nil {
		return nil, fmt.Errorf("build keygen signer: %w", err)
	}
	signers = append(signers, keygenSigner)

	// Add auth user keys (sign1, sign2, etc.)
	for _, authType := range authTypes {
		// Query by purpose to get correct xpub for the address type
		authKey, err := u.authFullPubKeyRepo.GetOneByPurpose(authType, purpose)
		if err != nil {
			return nil, fmt.Errorf(
				"get auth full pubkey for %s with purpose %s: %w",
				authType.String(), purpose.String(), err,
			)
		}
		if authKey == nil {
			return nil, fmt.Errorf(
				"auth full pubkey not found for %s with purpose %s",
				authType.String(), purpose.String(),
			)
		}

		// Use ExtendedPubKey (new format) if available, otherwise fall back to FullPublicKey (legacy)
		coinLevelExtendedKey := authKey.ExtendedPubKey
		if coinLevelExtendedKey == "" {
			// Legacy format: FullPublicKey contains compressed pubkey, not extended key
			return nil, fmt.Errorf(
				"extended public key not found for %s with purpose %s "+
					"(legacy compressed pubkey format not supported for descriptors)",
				authType.String(),
				purpose.String(),
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

// buildKeygenSigner creates a MultisigSigner for the keygen wallet's own key.
// This ensures that the keygen wallet participates in multisig along with auth users.
func (u *generateDescriptorUseCase) buildKeygenSigner(
	ctx context.Context,
	accountType domainAccount.AccountType,
	addressType domainAddress.AddrType,
	derivationPath string,
) (btc.MultisigSigner, error) {
	// Get seed from repository
	seedData, err := u.seedRepo.GetOne(ctx)
	if err != nil {
		return btc.MultisigSigner{}, fmt.Errorf("get seed: %w", err)
	}
	if seedData == nil {
		return btc.MultisigSigner{}, errors.New("seed not found - run 'keygen seed' first")
	}

	// Convert seed string to bytes
	seedBytes, err := infraKey.SeedToByte(seedData.Seed)
	if err != nil {
		return btc.MultisigSigner{}, fmt.Errorf("decode seed: %w", err)
	}

	// Determine BIP purpose from address type
	purpose, err := domainAuth.PurposeForAddressType(addressType.String())
	if err != nil {
		return btc.MultisigSigner{}, fmt.Errorf("determine purpose for address type %s: %w", addressType.String(), err)
	}

	// Map BIP purpose to infraKey.PurposeType
	var purposeType infraKey.PurposeType
	switch purpose {
	case domainAuth.PurposeBIP44:
		purposeType = infraKey.PurposeTypeBIP44
	case domainAuth.PurposeBIP49:
		purposeType = infraKey.PurposeTypeBIP49
	case domainAuth.PurposeBIP84:
		purposeType = infraKey.PurposeTypeBIP84
	case domainAuth.PurposeBIP86:
		purposeType = infraKey.PurposeTypeBIP86
	default:
		return btc.MultisigSigner{}, fmt.Errorf("unsupported BIP purpose: %s", purpose.String())
	}

	// Create HD key for this purpose
	hdKey := infraKey.NewHDKey(purposeType, u.coinTypeCode, u.chainConfig)

	// Create descriptor generator
	descGenerator := infraKey.NewDescriptorGenerator(hdKey, u.chainConfig)

	// Generate account-level extended public key from seed
	accountXPub, err := descGenerator.GetAccountXPub(seedBytes, accountType)
	if err != nil {
		return btc.MultisigSigner{}, fmt.Errorf("generate account xpub: %w", err)
	}

	// Parse extended public key
	xpub, err := hdkeychain.NewKeyFromString(accountXPub)
	if err != nil {
		return btc.MultisigSigner{}, fmt.Errorf("parse keygen extended key: %w", err)
	}

	// Verify network match
	if u.chainConfig != nil && !xpub.IsForNet(u.chainConfig) {
		return btc.MultisigSigner{}, errors.New("keygen extended key network mismatch")
	}

	// Get master fingerprint
	fingerprint, err := descGenerator.GetMasterFingerprintHex(seedBytes)
	if err != nil {
		return btc.MultisigSigner{}, fmt.Errorf("calculate keygen fingerprint: %w", err)
	}

	return btc.MultisigSigner{
		Fingerprint:    fingerprint,
		DerivationPath: derivationPath,
		ExtendedKey:    xpub,
	}, nil
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
