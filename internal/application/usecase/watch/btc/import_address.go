package btc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	appdto "github.com/hiromaily/go-crypto-wallet/internal/application/dto"
	portsBitcoin "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/bitcoin"
	portsFile "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// ImportAddressUseCase handles BTC address imports with rescan support
type ImportAddressUseCase interface {
	Execute(ctx context.Context, input watchusecase.ImportAddressInput) error
}

type importAddressUseCase struct {
	btcClient    portsBitcoin.Bitcoiner
	addrRepo     repository.AddressRepositorier
	addrFileRepo portsFile.AddressFileRepositorier
	coinTypeCode domainCoin.CoinTypeCode
	addrType     domainAddress.AddrType
}

// NewImportAddressUseCase creates a new BTC-specific ImportAddressUseCase
func NewImportAddressUseCase(
	btcClient portsBitcoin.Bitcoiner,
	addrRepo repository.AddressRepositorier,
	addrFileRepo portsFile.AddressFileRepositorier,
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

		// Import address into Bitcoin Core
		err = u.btcClient.ImportAddressWithLabel(targetAddr, addrFmt.AccountType.String(), input.Rescan)
		if err != nil {
			// Check if error is recoverable (address already exists)
			if isRecoverableImportError(err) {
				logger.Warn(
					"address already exists in wallet, skipping",
					"address", targetAddr,
					"account_type", addrFmt.AccountType.String())
				importStats.skipped++
				continue
			}

			// All other errors are critical - fail the import
			return fmt.Errorf("failed to import address %s (account: %s): %w",
				targetAddr, addrFmt.AccountType.String(), err)
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

	// For non-client accounts (deposit, payment, etc.), use multisig address
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
