package bch

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"

	bchutil "github.com/hiromaily/go-crypto-wallet/pkg/chains/bch"
	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// GetAddressInfo can be used as an alternative to `getaccount`, `validateaddress`.
// This override handles BCH nodes that return a legacy "label" string field instead
// of the "labels" array used by BTC nodes.
func (b *BitcoinCash) GetAddressInfo(addr string) (*btcrpc.GetAddressInfoResult, error) {
	result, err := b.GetPkgRPC().GetAddressInfo(addr)
	if err != nil {
		return nil, fmt.Errorf("fail to call pkgrpc.GetAddressInfo() in bch: %w", err)
	}
	// BCH nodes may return a legacy single "label" field instead of a "labels" array.
	if len(result.Labels) == 0 && result.Label != "" {
		result.Labels = btcrpc.FlexibleLabels{result.Label}
	}
	return result, nil
}

// GetAddressesByLabel overrides Bitcoin's GetAddressesByLabel to use BCH address decoding.
// This override is necessary because the parent Bitcoin.GetAddressesByLabel calls
// Bitcoin.DecodeAddress internally, which doesn't understand BCH CashAddr format.
func (b *BitcoinCash) GetAddressesByLabel(labelName string) ([]btcutil.Address, error) {
	logger.Debug("BCH GetAddressesByLabel called", "label", labelName)

	labels, err := b.GetPkgRPC().GetAddressesByLabel(labelName)
	if err != nil {
		logger.Debug("getaddressesbylabel RPC failed", "label", labelName, "error", err)
		return nil, fmt.Errorf(
			"fail to call pkgrpc.GetAddressesByLabel() for label %s in bch: %w",
			labelName, err)
	}

	if len(labels) == 0 {
		logger.Debug("no addresses found for label", "label", labelName)
		return nil, nil
	}

	logger.Debug("found addresses for label",
		"label", labelName,
		"raw_count", len(labels))

	// Decode each address key using BCH CashAddr format.
	resAddrs := make([]btcutil.Address, 0, len(labels))
	for key := range labels {
		address, err := b.DecodeAddress(key)
		if err != nil {
			logger.Error(
				"failed to decode address",
				"address", key,
				"label", labelName,
				"error", err)
			continue
		}
		resAddrs = append(resAddrs, address)
	}

	logger.Debug("successfully decoded addresses",
		"label", labelName,
		"decoded_count", len(resAddrs))

	return resAddrs, nil
}

// DecodeAddress overrides Bitcoin's DecodeAddress to handle BCH CashAddr format
func (b *BitcoinCash) DecodeAddress(addr string) (btcutil.Address, error) {
	// Use BCH-specific decoder that handles CashAddr format
	logger.Debug("BCH DecodeAddress called", "address", addr)
	address, err := bchutil.DecodeAddress(addr, b.GetChainConf())
	if err != nil {
		return nil, fmt.Errorf("fail to call bchutil.DecodeAddress() in bch: %w", err)
	}
	return address, nil
}
