package bch

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"

	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
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

	// Calculate minimum fee: 1 sat/byte
	txSize := tx.SerializeSize()
	minFeePerByte := 0.00001 // BCH/kB (1 sat/byte)
	calculatedFee := fmt.Sprintf("%f", minFeePerByte*float64(txSize)/1000)
	fee, err := b.StrToAmount(calculatedFee)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate BCH minimum fee: %w", err)
	}

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
		newFee, err = b.calculateNewFee(fee, adjustmentFee)
		if err != nil {
			logger.Warn("fail to call bch.calculateNewFee() but continue", "error", err)
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
	minFee, err := b.FloatToAmount(0.00001) // 1000 satoshis
	if err != nil {
		logger.Warn("failed to convert min fee", "error", err)
	} else if fee < minFee {
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
		return b.FloatToAmount(0.00001) // 1 sat/byte = 0.00001 BCH/kB
	}
	fee, err := b.FloatToAmount(res.RelayFee)
	if err != nil {
		return 0, err
	}

	// Ensure minimum 1000 satoshis (0.00001 BCH)
	minFee, err := b.FloatToAmount(0.00001)
	if err != nil {
		return fee, nil // Return original fee if conversion fails
	}
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

// calculateNewFee adjusts fee by adjustment fee multiplier.
// Uses inherited BTC implementation.
func (b *BitcoinCash) calculateNewFee(fee btcutil.Amount, adjustmentFee float64) (btcutil.Amount, error) {
	newFee, err := b.FloatToAmount(fee.ToBTC() * adjustmentFee)
	if err != nil {
		return 0, err
	}
	return newFee, nil
}
