// Package btc provides Bitcoin-specific use cases for Keygen wallet operations.
//
// # MuSig2 Address Creation
//
// The CreateMuSig2AddressUseCase generates MuSig2 Taproot addresses by aggregating
// multiple public keys from different signers (Keygen, Sign1, Sign2, etc.) into a
// single aggregated public key, which is then used to create a P2TR address.
//
// # Architecture Overview
//
// MuSig2 address creation follows these steps:
//  1. Validate that the account type is configured as a multisig account
//  2. Retrieve all signer public keys from auth_fullpubkey table
//  3. Retrieve account keys from account_key table with status AddrStatusPrivKeyImported
//  4. For each account key:
//     - Parse the private key (WIF format) and public key
//     - Combine all public keys (auth keys + account key)
//     - Create MuSig2 context with Taproot tweaking
//     - Generate aggregated public key
//     - Create Taproot address (P2TR) from aggregated key
//     - Update database with Taproot address and status AddrStatusMultisigAddressGenerated
//
// # Key Generation Workflow
//
//	┌─────────────────────────────────────────────────────────┐
//	│ Input: AccountType (e.g., deposit, payment)            │
//	└───────────────────┬─────────────────────────────────────┘
//	                    │
//	                    ▼
//	┌───────────────────────────────────────────────────────┐
//	│ 1. Get Auth Public Keys from auth_fullpubkey         │
//	│    - Keygen's full public key                        │
//	│    - Sign1's full public key                         │
//	│    - Sign2's full public key                         │
//	└───────────────────┬───────────────────────────────────┘
//	                    │
//	                    ▼
//	┌───────────────────────────────────────────────────────┐
//	│ 2. Get Account Keys from account_key                 │
//	│    - Filter: addr_status = AddrStatusPrivKeyImported │
//	│    - Contains: WIF, full public key                  │
//	└───────────────────┬───────────────────────────────────┘
//	                    │
//	                    ▼
//	┌───────────────────────────────────────────────────────┐
//	│ 3. For Each Account Key:                             │
//	│    a. Parse private key (WIF format)                 │
//	│    b. Parse account public key                       │
//	│    c. Combine: [auth1, auth2, auth3, account]        │
//	└───────────────────┬───────────────────────────────────┘
//	                    │
//	                    ▼
//	┌───────────────────────────────────────────────────────┐
//	│ 4. Create MuSig2 Context with Taproot Tweak          │
//	│    - Input: account private key, all public keys     │
//	│    - Applies BIP86 Taproot tweak                     │
//	│    - Returns: musig2.Context                         │
//	└───────────────────┬───────────────────────────────────┘
//	                    │
//	                    ▼
//	┌───────────────────────────────────────────────────────┐
//	│ 5. Get Aggregated Public Key                         │
//	│    - Single public key representing all signers      │
//	│    - Used for Taproot address generation             │
//	└───────────────────┬───────────────────────────────────┘
//	                    │
//	                    ▼
//	┌───────────────────────────────────────────────────────┐
//	│ 6. Create Taproot Address (P2TR)                     │
//	│    - Format: bc1p... (mainnet) / tb1p... (testnet)  │
//	│    - Indistinguishable from single-sig Taproot       │
//	└───────────────────┬───────────────────────────────────┘
//	                    │
//	                    ▼
//	┌───────────────────────────────────────────────────────┐
//	│ 7. Update Database                                    │
//	│    - Store Taproot address                           │
//	│    - Update status: AddrStatusMultisigAddressGenerated│
//	└───────────────────────────────────────────────────────┘
//
// # Security Considerations
//
//   - Private keys are only used locally for context creation
//   - Private keys never leave the Keygen wallet
//   - Aggregated public keys are indistinguishable from single-signature keys
//   - Taproot addresses provide privacy (multisig looks like single-sig)
//
// # Example
//
// For detailed usage examples, see Example_createMuSig2Address.
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

// ExampleCreateMuSig2Address demonstrates how to use CreateMuSig2AddressUseCase
// to generate MuSig2 Taproot addresses for a multisig account.
//
// This example shows the basic workflow:
//  1. Initialize dependencies (MuSig2Service, repositories, configuration)
//  2. Create the use case with NewCreateMuSig2AddressUseCase
//  3. Execute the use case with appropriate input (account type)
//  4. Handle errors appropriately
//
// Note: This is a conceptual example. In real implementation, you would need:
//   - Properly configured Bitcoin network parameters (mainnet/testnet)
//   - Database connections for repositories
//   - Auth full public keys already stored in auth_fullpubkey table
//   - Account keys with AddrStatusPrivKeyImported status in account_key table
func ExampleCreateMuSig2Address() {
	ctx := context.Background()

	// Initialize dependencies (simplified for example)
	var (
		musig2Service      *btc.MuSig2Service              // Provides MuSig2 cryptographic operations
		chainConfig        = &chaincfg.TestNet3Params      // Bitcoin network configuration
		authFullPubKeyRepo cold.AuthFullPubkeyRepositorier // Repository for auth public keys
		accountKeyRepo     cold.AccountKeyRepositorier     // Repository for account keys
		multisigAccount    account.MultisigAccounter       // Multisig account configuration
	)

	// Create the use case
	useCase := NewCreateMuSig2AddressUseCase(
		musig2Service,
		chainConfig,
		authFullPubKeyRepo,
		accountKeyRepo,
		multisigAccount,
	)

	// Prepare input (specify which account type to generate addresses for)
	input := keygenusecase.CreateMuSig2AddressInput{
		AccountType: "deposit", // or "payment", "client", etc. (see account.AccountType constants)
	}

	// Execute the use case
	err := useCase.Create(ctx, input)
	if err != nil {
		logger.Error("failed to create MuSig2 addresses", "error", err)
		return
	}

	logger.Info("successfully created MuSig2 Taproot addresses")

	// After successful execution:
	// - All eligible account keys now have Taproot addresses generated
	// - Database is updated with:
	//   - taproot_address field populated (tb1p... format for testnet)
	//   - multisig_address field populated (same as taproot_address)
	//   - addr_status updated to AddrStatusMultisigAddressGenerated
	//
	// These addresses can now be:
	// 1. Exported to Watch wallet for monitoring
	// 2. Used to receive Bitcoin payments
	// 3. Used as inputs in MuSig2 signing workflows
}
