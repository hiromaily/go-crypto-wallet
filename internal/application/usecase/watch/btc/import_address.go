package btc

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
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// importAddrBTCClient defines the minimal interface needed for BTC address import.
// This follows the Interface Segregation Principle - depend only on methods actually used.
type importAddrBTCClient interface {
	apibtc.ChainConfigProvider   // CoinTypeCode
	apibtc.LegacyAddressImporter // ImportAddressWithLabel, ImportMulti
	apibtc.AddressOperator       // GetAddressInfo
}

// ImportAddressUseCase handles BTC address imports with rescan support
type ImportAddressUseCase interface {
	Execute(ctx context.Context, input watchusecase.ImportAddressInput) error
}

type importAddressUseCase struct {
	btcClient    importAddrBTCClient
	addrRepo     repowatch.AddressRepositorier
	addrFileRepo file.AddressFileRepositorier
	coinTypeCode domainCoin.CoinTypeCode
	addrType     domainAddress.AddrType
}

// NewImportAddressUseCase creates a new BTC-specific ImportAddressUseCase
func NewImportAddressUseCase(
	btcClient apibtc.Bitcoiner,
	addrRepo repowatch.AddressRepositorier,
	addrFileRepo file.AddressFileRepositorier,
	coinTypeCode domainCoin.CoinTypeCode,
	addrType domainAddress.AddrType,
) ImportAddressUseCase {
	return &importAddressUseCase{
		btcClient:    btcClient,
		addrRepo:     addrRepo,
		addrFileRepo: addrFileRepo,
		coinTypeCode: coinTypeCode,
		addrType:     addrType,
	}
}

// Execute imports addresses from a file with optional rescan
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
		if !isRecoverableImportError(err) {
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
		addrFmt, err := address.ConvertLine(u.btcClient.CoinTypeCode(), inner)
		if err != nil {
			return fmt.Errorf("failed to convert address format: %w", err)
		}

		// Select target address based on account type and address type
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
	logger.Info("address import completed",
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

// importSingleAddress imports a single address into Bitcoin Core
// Returns (imported=true, nil) if successfully imported
// Returns (imported=false, nil) if address already exists (recoverable error)
// Returns (imported=false, error) for critical errors
func (u *importAddressUseCase) importSingleAddress(
	targetAddr string,
	addrFmt *appdto.AddressFormat,
	rescan bool,
	handleRecoverableError func(error, string, domainAccount.AccountType) bool,
) (bool, error) {
	// Use ImportMulti if redeem script is available (for P2SH multisig)
	// Otherwise use legacy ImportAddressWithLabel
	if addrFmt.RedeemScript != "" {
		return u.importWithRedeemScript(targetAddr, addrFmt, rescan, handleRecoverableError)
	}
	return u.importLegacyAddress(targetAddr, addrFmt, rescan, handleRecoverableError)
}

// importWithRedeemScript imports P2SH multisig address with redeem script using ImportMulti
func (u *importAddressUseCase) importWithRedeemScript(
	targetAddr string,
	addrFmt *appdto.AddressFormat,
	rescan bool,
	handleRecoverableError func(error, string, domainAccount.AccountType) bool,
) (bool, error) {
	logger.Debug("importing multisig address with redeem script",
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

	responses, err := u.btcClient.ImportMulti(requests, &dtobtc.ImportMultiOptions{Rescan: rescan})
	if err != nil {
		return false, fmt.Errorf("failed to call ImportMulti for address %s: %w", targetAddr, err)
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
		return false, fmt.Errorf("failed to import multisig address %s: %s", targetAddr, errMsg)
	}

	return true, nil
}

// importLegacyAddress imports regular address using legacy ImportAddressWithLabel
func (u *importAddressUseCase) importLegacyAddress(
	targetAddr string,
	addrFmt *appdto.AddressFormat,
	rescan bool,
	handleRecoverableError func(error, string, domainAccount.AccountType) bool,
) (bool, error) {
	// For BCH P2PKH single-sig addresses, use ImportMulti with pubkey to make UTXOs solvable
	if u.btcClient.CoinTypeCode() == domainCoin.BCH && addrFmt.FullPublicKey != "" {
		logger.Debug("importing pubkey with address for BCH using ImportMulti",
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

		responses, err := u.btcClient.ImportMulti(requests, &dtobtc.ImportMultiOptions{Rescan: rescan})
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

	// For BTC or other cases, use legacy ImportAddressWithLabel
	err := u.btcClient.ImportAddressWithLabel(targetAddr, addrFmt.AccountType.String(), rescan)
	if err != nil {
		// Check if error is recoverable using helper
		if handleRecoverableError(err, targetAddr, addrFmt.AccountType) {
			return false, nil // Address already exists, skip
		}

		// All other errors are critical - fail the import
		return false, fmt.Errorf("failed to import address %s (account: %s): %w",
			targetAddr, addrFmt.AccountType.String(), err)
	}

	return true, nil
}

// selectTargetAddress determines which address format to use based on account type and address type
func (u *importAddressUseCase) selectTargetAddress(addrFmt *appdto.AddressFormat) (string, error) {
	// For client accounts, use specific address format
	if addrFmt.AccountType == domainAccount.AccountTypeClient {
		switch u.btcClient.CoinTypeCode() {
		case domainCoin.BTC:
			switch u.addrType {
			case domainAddress.AddrTypeBech32:
				return addrFmt.Bech32Address, nil
			case domainAddress.AddrTypeTaproot:
				if addrFmt.TaprootAddress == "" {
					return "", errors.New("taproot address is empty in the imported data")
				}
				return addrFmt.TaprootAddress, nil
			case domainAddress.AddrTypeLegacy:
				return addrFmt.P2PKHAddress, nil
			case domainAddress.AddrTypeP2shSegwit,
				domainAddress.AddrTypeBCHCashAddr, domainAddress.AddrTypeETH:
				return addrFmt.P2SHSegwitAddress, nil
			default:
				return addrFmt.P2SHSegwitAddress, nil
			}
		case domainCoin.BCH:
			return addrFmt.P2PKHAddress, nil
		case domainCoin.LTC, domainCoin.ETH, domainCoin.XRP, domainCoin.ERC20, domainCoin.HYT:
			return "", fmt.Errorf("unsupported coin type: %s", u.btcClient.CoinTypeCode().String())
		default:
			return "", fmt.Errorf("unknown coin type: %s", u.btcClient.CoinTypeCode().String())
		}
	}

	// For non-client accounts (deposit, payment, etc.), determine address type based on multisig status
	// BCH: For single-sig accounts (Pattern 1), use P2PKH even for payment/deposit accounts
	if u.btcClient.CoinTypeCode() == domainCoin.BCH && addrFmt.MultisigAddress == "" && addrFmt.RedeemScript == "" {
		// BCH single-sig: use P2PKH address
		return addrFmt.P2PKHAddress, nil
	}

	// For multisig accounts, use multisig address
	// Fallback to P2SHSegwit if MultisigAddress is empty (traditional multisig)
	if addrFmt.MultisigAddress != "" {
		return addrFmt.MultisigAddress, nil
	}
	// For traditional multisig (non-MuSig2), use P2SH-Segwit address
	if addrFmt.P2SHSegwitAddress != "" {
		return addrFmt.P2SHSegwitAddress, nil
	}
	return "", errors.New("no valid address found for non-client account")
}

// verifyImportedAddress confirms the address was imported correctly as watch-only
func (u *importAddressUseCase) verifyImportedAddress(addr string) {
	addrInfo, err := u.btcClient.GetAddressInfo(addr)
	if err != nil {
		logger.Error(
			"failed to verify imported address",
			"address", addr,
			"error", err)
		return
	}

	labelName := ""
	if len(addrInfo.Labels) != 0 {
		labelName = addrInfo.Labels[0]
	}
	logger.Debug("address verified",
		"account", labelName,
		"address", addr)

	// Warn if not watch-only (should always be watch-only for watch wallets)
	if !addrInfo.IsWatchOnly {
		logger.Warn("address should be watch-only",
			"address", addr)
	}
}

// isRecoverableImportError checks if an import error is recoverable (address already exists)
// Returns true for errors that can be safely ignored (address duplicates)
// Returns false for critical errors that should fail the import
func isRecoverableImportError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())

	// Check for common "address already exists" error patterns
	// Note: "already" covers all variants like "already have", "already imported",
	// "already in wallet", and "label already exists"
	recoverablePatterns := []string{
		"already",   // Catches all "already..." variants, e.g., "address already exists", "label already exists"
		"duplicate", // For "duplicate" errors
		"exists",    // For generic "exists" errors
	}

	for _, pattern := range recoverablePatterns {
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}

	return false
}
