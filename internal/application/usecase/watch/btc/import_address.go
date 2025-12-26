package btc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
	models "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// ImportAddressUseCase handles BTC address imports with rescan support
type ImportAddressUseCase interface {
	Execute(ctx context.Context, input watchusecase.ImportAddressInput) error
}

type importAddressUseCase struct {
	btcClient    bitcoin.Bitcoiner
	addrRepo     watch.AddressRepositorier
	addrFileRepo file.AddressFileRepositorier
	coinTypeCode domainCoin.CoinTypeCode
	addrType     address.AddrType
}

// NewImportAddressUseCase creates a new BTC-specific ImportAddressUseCase
func NewImportAddressUseCase(
	btcClient bitcoin.Bitcoiner,
	addrRepo watch.AddressRepositorier,
	addrFileRepo file.AddressFileRepositorier,
	coinTypeCode domainCoin.CoinTypeCode,
	addrType address.AddrType,
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

	pubKeyData := make([]*models.Address, 0, len(pubKeys))
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
			// Warning: address may already exist, continue with other addresses
			logger.Warn(
				"failed to import address but continuing",
				"address", targetAddr,
				"account_type", addrFmt.AccountType.String(),
				"error", err)
			continue
		}

		// Add to batch for database insertion
		pubKeyData = append(pubKeyData, &models.Address{
			Coin:          u.coinTypeCode.String(),
			Account:       addrFmt.AccountType.String(),
			WalletAddress: targetAddr,
		})

		// Verify address was imported correctly
		u.verifyImportedAddress(targetAddr)
	}

	// Insert all addresses into database
	if len(pubKeyData) > 0 {
		if err := u.addrRepo.InsertBulk(pubKeyData); err != nil {
			return fmt.Errorf("failed to insert addresses into database: %w", err)
		}
	}

	return nil
}

// selectTargetAddress determines which address format to use based on account type and address type
func (u *importAddressUseCase) selectTargetAddress(addrFmt *address.AddressFormat) (string, error) {
	// For client accounts, use specific address format
	if addrFmt.AccountType == domainAccount.AccountTypeClient {
		switch u.btcClient.CoinTypeCode() {
		case domainCoin.BTC:
			switch u.addrType {
			case address.AddrTypeBech32:
				return addrFmt.Bech32Address, nil
			case address.AddrTypeTaproot:
				// TODO: Implement Taproot address support
				return "", errors.New("taproot address type not yet supported")
			case address.AddrTypeLegacy, address.AddrTypeP2shSegwit,
				address.AddrTypeBCHCashAddr, address.AddrTypeETH:
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
	return addrFmt.MultisigAddress, nil
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

	logger.Debug("address verified",
		"account", addrInfo.GetLabelName(),
		"address", addr)

	// Warn if not watch-only (should always be watch-only for watch wallets)
	if !addrInfo.Iswatchonly {
		logger.Warn("address should be watch-only",
			"address", addr)
	}
}
