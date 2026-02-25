package bch

import (
	"context"
	"errors"
	"fmt"
	"strings"

	appdto "github.com/hiromaily/go-crypto-wallet/internal/application/dto"
	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	usecaseshared "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/shared"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type importAddressUseCase struct {
	bchClient    apibtc.AddressImportClient
	addrRepo     repowatch.AddressRepositorier
	addrFileRepo file.AddressFileRepositorier
	coinTypeCode domainCoin.CoinTypeCode
	addrType     domainAddress.AddrType
}

// NewImportBCHAddressUseCase creates a new BCH-specific ImportAddressUseCase
func NewImportBCHAddressUseCase(
	bchClient apibtc.AddressImportClient,
	addrRepo repowatch.AddressRepositorier,
	addrFileRepo file.AddressFileRepositorier,
	coinTypeCode domainCoin.CoinTypeCode,
	addrType domainAddress.AddrType,
) watchusecase.ImportAddressUseCase {
	return &importAddressUseCase{
		bchClient:    bchClient,
		addrRepo:     addrRepo,
		addrFileRepo: addrFileRepo,
		coinTypeCode: coinTypeCode,
		addrType:     addrType,
	}
}

// Execute imports BCH addresses from a file with optional rescan
func (u *importAddressUseCase) Execute(ctx context.Context, input watchusecase.ImportAddressInput) error {
	// Read addresses from file
	pubKeys, err := u.addrFileRepo.ImportAddress(input.FileName)
	if err != nil {
		return fmt.Errorf("failed to import addresses from file: %w", err)
	}

	pubKeyData := make([]*domainAddress.Address, 0, len(pubKeys))
	importStats := struct {
		total   int
		success int
		skipped int
	}{
		total: len(pubKeys),
	}

	// Helper function to handle recoverable import errors (e.g., address already exists)
	handleRecoverableError := func(err error, targetAddr string, accountType domainAccount.AccountType) bool {
		if !usecaseshared.IsRecoverableImportError(err) {
			return false
		}

		logger.Warn(
			"address already exists in wallet, skipping",
			"address", targetAddr,
			"account_type", accountType.String())
		importStats.skipped++
		return true
	}

	for _, key := range pubKeys {
		// Parse CSV line
		inner := strings.Split(key, ",")

		// Convert address format
		addrFmt, err := address.ConvertLine(u.bchClient.CoinTypeCode(), inner)
		if err != nil {
			return fmt.Errorf("failed to convert address format: %w", err)
		}

		// Select target address based on account type (BCH-specific logic)
		targetAddr, err := u.selectTargetAddress(addrFmt)
		if err != nil {
			return err
		}

		// Import single address
		imported, err := u.importSingleAddress(targetAddr, addrFmt, input.Rescan, handleRecoverableError)
		if err != nil {
			return err
		}

		// Skip if address was already imported
		if !imported {
			continue
		}

		importStats.success++

		// Add to batch for database insertion
		pubKeyData = append(pubKeyData, &domainAddress.Address{
			CoinTypeCode:  u.coinTypeCode,
			AccountType:   addrFmt.AccountType,
			WalletAddress: targetAddr,
			IsAllocated:   false,
		})

		// Verify address was imported correctly
		u.verifyImportedAddress(targetAddr)
	}

	// Log import statistics
	logger.Info("BCH address import completed",
		"total", importStats.total,
		"imported", importStats.success,
		"skipped", importStats.skipped)

	// Insert all addresses into database
	if len(pubKeyData) > 0 {
		if err := u.addrRepo.InsertBulk(ctx, pubKeyData); err != nil {
			return fmt.Errorf("failed to insert addresses into database: %w", err)
		}
	}

	return nil
}

// importSingleAddress imports a single BCH address into Bitcoin Cash Core
// Returns (imported=true, nil) if successfully imported
// Returns (imported=false, nil) if address already exists (recoverable error)
// Returns (imported=false, error) for critical errors
func (u *importAddressUseCase) importSingleAddress(
	targetAddr string,
	addrFmt *appdto.AddressFormat,
	rescan bool,
	handleRecoverableError func(error, string, domainAccount.AccountType) bool,
) (bool, error) {
	// BCH: Use ImportMulti with pubkey to make UTXOs solvable
	// For multisig, use redeem script
	if addrFmt.RedeemScript != "" {
		return u.importWithRedeemScript(targetAddr, addrFmt, rescan, handleRecoverableError)
	}
	return u.importWithPubkey(targetAddr, addrFmt, rescan, handleRecoverableError)
}

// importWithRedeemScript imports P2SH multisig address with redeem script using ImportMulti
func (u *importAddressUseCase) importWithRedeemScript(
	targetAddr string,
	addrFmt *appdto.AddressFormat,
	rescan bool,
	handleRecoverableError func(error, string, domainAccount.AccountType) bool,
) (bool, error) {
	logger.Debug("importing BCH multisig address with redeem script",
		"address", targetAddr,
		"account", addrFmt.AccountType.String())

	requests := []dtobtc.ImportMultiRequest{
		{
			ScriptPubKey: map[string]string{"address": targetAddr},
			Timestamp:    "now", // Skip rescanning for faster import
			RedeemScript: addrFmt.RedeemScript,
			WatchOnly:    true,
			Label:        addrFmt.AccountType.String(),
		},
	}

	responses, err := u.bchClient.ImportMulti(requests, &dtobtc.ImportMultiOptions{Rescan: rescan})
	if err != nil {
		return false, fmt.Errorf("failed to call ImportMulti for BCH address %s: %w", targetAddr, err)
	}

	if len(responses) == 0 || !responses[0].Success {
		errMsg := "unknown error"
		if len(responses) > 0 && responses[0].Error != nil {
			errMsg = responses[0].Error.Message
		}
		// Check if error is recoverable using helper
		if handleRecoverableError(fmt.Errorf("%s", errMsg), targetAddr, addrFmt.AccountType) {
			return false, nil // Address already exists, skip
		}
		return false, fmt.Errorf("failed to import BCH multisig address %s: %s", targetAddr, errMsg)
	}

	return true, nil
}

// importWithPubkey imports BCH P2PKH address with pubkey using ImportMulti
// BCH requires pubkey to make UTXOs solvable for spending
func (u *importAddressUseCase) importWithPubkey(
	targetAddr string,
	addrFmt *appdto.AddressFormat,
	rescan bool,
	handleRecoverableError func(error, string, domainAccount.AccountType) bool,
) (bool, error) {
	if addrFmt.FullPublicKey == "" {
		return false, fmt.Errorf("BCH requires full public key for address %s", targetAddr)
	}

	logger.Debug("importing BCH address with pubkey using ImportMulti",
		"pubkey", addrFmt.FullPublicKey,
		"address", targetAddr,
		"account", addrFmt.AccountType.String())

	// Set timestamp based on rescan parameter:
	// - rescan=true: Use 0 to rescan from beginning
	// - rescan=false: Use "now" to only watch new transactions
	timestamp := "now"
	if rescan {
		timestamp = "0"
	}

	requests := []dtobtc.ImportMultiRequest{
		{
			ScriptPubKey: map[string]string{"address": targetAddr},
			Timestamp:    timestamp,
			PubKeys:      []string{addrFmt.FullPublicKey},
			WatchOnly:    true,
			Label:        addrFmt.AccountType.String(),
		},
	}

	responses, err := u.bchClient.ImportMulti(requests, &dtobtc.ImportMultiOptions{Rescan: rescan})
	if err != nil {
		return false, fmt.Errorf("failed to call ImportMulti for BCH address %s: %w", targetAddr, err)
	}

	if len(responses) == 0 || !responses[0].Success {
		errMsg := "unknown error"
		if len(responses) > 0 && responses[0].Error != nil {
			errMsg = responses[0].Error.Message
		}
		// Check if error is recoverable using helper
		if handleRecoverableError(fmt.Errorf("%s", errMsg), targetAddr, addrFmt.AccountType) {
			return false, nil // Address already exists, skip
		}
		return false, fmt.Errorf("failed to import BCH address %s: %s", targetAddr, errMsg)
	}

	return true, nil
}

// selectTargetAddress determines which address format to use for BCH
// BCH-specific logic:
// - Client accounts: Always use P2PKH (CashAddr format)
// - Single-sig accounts (Pattern 1): Use P2PKH for all account types
// - Multisig accounts (Pattern 2, 3): Use multisig address or P2SH
func (*importAddressUseCase) selectTargetAddress(addrFmt *appdto.AddressFormat) (string, error) {
	// For client accounts, always use P2PKH
	if addrFmt.AccountType == domainAccount.AccountTypeClient {
		return addrFmt.P2PKHAddress, nil
	}

	// For non-client accounts (deposit, payment, etc.)
	// Check if this is a single-sig account (no multisig address or redeem script)
	if addrFmt.MultisigAddress == "" && addrFmt.RedeemScript == "" {
		// BCH Pattern 1 (Single-sig): use P2PKH address
		return addrFmt.P2PKHAddress, nil
	}

	// For multisig accounts (Pattern 2, 3), use multisig address
	if addrFmt.MultisigAddress != "" {
		return addrFmt.MultisigAddress, nil
	}
	// Fallback to P2SH if MultisigAddress is empty
	if addrFmt.P2SHSegwitAddress != "" {
		return addrFmt.P2SHSegwitAddress, nil
	}

	return "", errors.New("no valid address found for BCH account")
}

// verifyImportedAddress confirms the BCH address was imported correctly as watch-only
func (u *importAddressUseCase) verifyImportedAddress(addr string) {
	addrInfo, err := u.bchClient.GetAddressInfo(addr)
	if err != nil {
		logger.Error(
			"failed to verify imported BCH address",
			"address", addr,
			"error", err)
		return
	}

	labelName := ""
	if len(addrInfo.Labels) != 0 {
		labelName = addrInfo.Labels[0]
	}
	logger.Debug("BCH address verified",
		"account", labelName,
		"address", addr)

	// Warn if not watch-only (should always be watch-only for watch wallets)
	if !addrInfo.IsWatchOnly {
		logger.Warn("BCH address should be watch-only",
			"address", addr)
	}
}
