package bch

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"

	btcpkg "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

const (
	// minFeeRateSatsPerByte is the minimum fee rate for BCH transactions (1 sat/byte)
	minFeeRateSatsPerByte = 1
	// minTxFeeSats is the absolute minimum fee for any BCH transaction (1000 satoshis)
	minTxFeeSats = 1000
)

// GetFee overrides BTC GetFee to properly handle BCH fee calculation.
// BCH differences from BTC:
//   - No SegWit, so all bytes count equally (no witness discount)
//   - Relay fee requirements may differ from BTC
//   - Uses sat/byte instead of sat/vByte
func (b *BitcoinCash) GetFee(tx *wire.MsgTx, adjustmentFee float64) (btcutil.Amount, error) {
	logger.Debug("BCH GetFee called", "tx_size", tx.SerializeSize(), "adjustment", adjustmentFee)

	// BCH regtest and testnets often have unreliable fee estimation
	// Use a more conservative approach: always start with a safe minimum

	// Calculate minimum fee: 1 sat/byte using integer arithmetic
	txSize := tx.SerializeSize()
	fee := btcutil.Amount(txSize * minFeeRateSatsPerByte)

	logger.Debug("BCH calculated minimum fee",
		"tx_size_bytes", txSize,
		"fee_bch", fee.ToBTC(),
		"fee_satoshis", fee)

	// Try to get network relay fee and use it if it's higher
	relayFee, err := b.getMinRelayFeeBCH()
	if err != nil {
		logger.Debug("could not get BCH relay fee, using calculated minimum", "error", err)
	} else if relayFee > fee {
		logger.Debug("using higher relay fee", "relay_fee", relayFee.ToBTC())
		fee = relayFee
	}

	// Apply adjustment fee multiplier if provided
	if b.validateAdjustmentFee(adjustmentFee) {
		var newFee btcutil.Amount
		newFee, err = btcpkg.CalculateNewFee(fee, adjustmentFee)
		if err != nil {
			logger.Warn("fail to call btcpkg.CalculateNewFee() but continue", "error", err)
		} else {
			logger.Debug("applied BCH fee adjustment",
				"original_fee", fee.ToBTC(),
				"adjustment", adjustmentFee,
				"new_fee", newFee.ToBTC())
			fee = newFee
		}
	}

	// BCH minimum: ensure at least 1000 satoshis (0.00001 BCH)
	// This is a safety net to prevent "min relay fee not met" errors
	minFee := btcutil.Amount(minTxFeeSats)
	if fee < minFee {
		logger.Debug("increasing fee to BCH minimum",
			"calculated_fee", fee.ToBTC(),
			"min_fee", minFee.ToBTC())
		fee = minFee
	}

	logger.Debug("BCH final fee calculation",
		"fee_bch", fee.ToBTC(),
		"fee_satoshis", fee,
		"tx_size_bytes", tx.SerializeSize(),
		"fee_per_byte", float64(fee)/float64(tx.SerializeSize()))

	return fee, nil
}

// getMinRelayFeeBCH gets the minimum relay fee from BCH network.
// This is BCH-specific and may return different values than BTC.
func (b *BitcoinCash) getMinRelayFeeBCH() (btcutil.Amount, error) {
	res, err := b.GetNetworkInfo()
	if err != nil {
		return 0, fmt.Errorf("fail to call bch.GetNetworkInfo(): %w", err)
	}
	if res.RelayFee == 0 {
		// BCH regtest may return 0 relay fee, use conservative default
		logger.Debug("BCH relay fee is 0 from getnetworkinfo, using default minimum")
		return btcutil.Amount(minTxFeeSats), nil
	}
	fee, err := btcpkg.FloatToAmount(res.RelayFee)
	if err != nil {
		return 0, err
	}

	// Ensure minimum 1000 satoshis (0.00001 BCH)
	minFee := btcutil.Amount(minTxFeeSats)
	if fee < minFee {
		return minFee, nil
	}
	return fee, nil
}

// validateAdjustmentFee validates adjustment fee parameter.
// Uses inherited BTC implementation.
func (b *BitcoinCash) validateAdjustmentFee(fee float64) bool {
	return b.Bitcoin.FeeRangeMin() <= fee && fee <= b.Bitcoin.FeeRangeMax()
}
