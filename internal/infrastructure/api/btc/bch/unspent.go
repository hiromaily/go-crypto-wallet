package bch

import (
	"fmt"
	"sort"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// ListUnspentByAccount overrides Bitcoin's ListUnspentByAccount for BCH-specific handling
// BCH doesn't support descriptor wallets, so this simplified version only uses label-based lookup
func (b *BitcoinCash) ListUnspentByAccount(
	accountType domainAccount.AccountType, confirmationNum uint64,
) ([]dtobtc.UnspentOutput, error) {
	logger.Debug("BCH ListUnspentByAccount called", "account_type", accountType)
	label := accountType.String()

	// Use BCH's GetAddressesByLabel which handles CashAddr format
	addrs, err := b.GetAddressesByLabel(label)
	if err != nil {
		return nil, fmt.Errorf("fail to call bch.GetAddressesByLabel(): %w", err)
	}

	if len(addrs) == 0 {
		logger.Debug("no addresses found for account", "account_type", accountType, "label", label)
		return nil, fmt.Errorf("no addresses found for account %s", accountType)
	}

	logger.Debug("found labeled addresses for account",
		"account_type", accountType,
		"label", label,
		"address_count", len(addrs))

	// Get unspent outputs for these addresses
	allUnspentList, err := b.Bitcoin.ListUnspent(confirmationNum)
	if err != nil {
		return nil, fmt.Errorf("fail to call btc.ListUnspent() in bch: %w", err)
	}

	// Filter unspent outputs for the addresses we found
	// For BCH, we need to normalize addresses because:
	// - getaddressesbylabel returns CashAddr format with prefix (e.g., "bchreg:pr...")
	// - EncodeAddress() may return legacy format without prefix (e.g., "pr...")
	// - listunspent also returns CashAddr format with prefix
	// Solution: Decode unspent address and compare using EncodeAddress() format
	var unspentList []dtobtc.UnspentOutput
	addrStrings := make(map[string]bool)
	for _, addr := range addrs {
		// Store the normalized address format for quick lookup
		addrStrings[addr.EncodeAddress()] = true
	}

	for _, unspent := range allUnspentList {
		// Decode the unspent address to normalize it
		unspentAddr, err := b.DecodeAddress(unspent.Address)
		if err != nil {
			logger.Error("failed to decode unspent address",
				"address", unspent.Address,
				"error", err)
			continue
		}

		// Compare using the normalized format
		if addrStrings[unspentAddr.EncodeAddress()] {
			unspentList = append(unspentList, unspent)
		}
	}

	if len(unspentList) == 0 {
		logger.Debug("no unspent outputs found for account",
			"account_type", accountType,
			"label", label,
			"address_count", len(addrs))
		return nil, fmt.Errorf("no unspent outputs found for account %s", accountType)
	}

	logger.Debug("found unspent outputs",
		"account_type", accountType,
		"unspent_count", len(unspentList),
		"confirmation_required", confirmationNum)

	// sort amount by ascending (small to big)
	sort.Slice(unspentList, func(i, j int) bool {
		return unspentList[i].Amount < unspentList[j].Amount
	})

	return unspentList, nil
}
