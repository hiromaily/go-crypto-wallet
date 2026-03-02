package btc

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"

	btcpkg "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// Making Sense of Bitcoin Transaction Fees
// https://bitzuma.com/posts/making-sense-of-bitcoin-transaction-fees/

// GetTransactionFee calculate fee from transaction size
func (b *Bitcoin) GetTransactionFee(tx *wire.MsgTx) (btcutil.Amount, error) {
	feePerKB, err := b.pkgrpc.EstimateSmartFee(int(b.confirmationBlock))
	if err != nil {
		return 0, fmt.Errorf("fail to call btc.EstimateSmartFee(): %w", err)
	}
	fee := fmt.Sprintf("%f", feePerKB*float64(tx.SerializeSize())/1000)

	// To Amount
	feeAsBit, err := btcpkg.StrToAmount(fee)
	if err != nil {
		return 0, err
	}

	return feeAsBit, nil
}

// GetFee get more preferable fee
func (b *Bitcoin) GetFee(tx *wire.MsgTx, adjustmentFee float64) (btcutil.Amount, error) {
	// get tx fee
	fee, err := b.GetTransactionFee(tx)

	// Treat 0 fee as an error - it's invalid for transaction broadcast
	if err != nil || fee == 0 {
		// Fee estimation failed or returned 0 (common on fresh regtest/signet networks)
		// Fall back to minimum relay fee or hardcoded fallback
		if err != nil {
			logger.Warn("fee estimation failed, using fallback", "error", err)
		} else {
			logger.Warn("fee estimation returned 0, using fallback")
		}

		// Try to get minimum relay fee from network
		relayFee, relayErr := b.getMinRelayFee()
		if relayErr != nil {
			logger.Warn("failed to get relay fee, using hardcoded fallback", "error", relayErr)
			// Use hardcoded fallback: 1 sat/vbyte = 0.00001 BTC/kB
			// This is conservative and works for regtest/signet
			fallbackFeePerKB := 0.00001 // BTC/kB
			fallbackFee := fmt.Sprintf("%f", fallbackFeePerKB*float64(tx.SerializeSize())/1000)
			fee, err = btcpkg.StrToAmount(fallbackFee)
			if err != nil {
				return 0, fmt.Errorf("failed to calculate fallback fee: %w", err)
			}
			logger.Info("using hardcoded fallback fee", "fee_btc", fee.ToBTC(), "fee_per_kb", fallbackFeePerKB)
		} else {
			// Use relay fee as base
			fee = relayFee
			logger.Info("using relay fee as fallback", "fee_btc", fee.ToBTC())
		}
	}

	// if response doesn't meet minimum fee, it should be overridden
	relayFee, err := b.getMinRelayFee()
	if err != nil {
		logger.Warn("fail to call btc.getMinRelayFee() but continue", "error", err)
	} else if fee < relayFee {
		fee = relayFee
	}

	// if adjustmentFee param is given
	if b.validateAdjustmentFee(adjustmentFee) {
		var newFee btcutil.Amount
		newFee, err = btcpkg.CalculateNewFee(fee, adjustmentFee)
		if err != nil {
			logger.Warn("fail to call btcpkg.CalculateNewFee() but continue", "error", err)
		} else {
			fee = newFee
		}
		logger.Debug("called btcpkg.CalculateNewFee()", "adjusted newFee", newFee)
	}

	return fee, nil
}

// ValidateAdjustmentFee validate adjustment fee param
func (b *Bitcoin) validateAdjustmentFee(fee float64) bool {
	if fee >= b.FeeRangeMin() && fee <= b.FeeRangeMax() {
		return true
	}
	return false
}

func (b *Bitcoin) getMinRelayFee() (btcutil.Amount, error) {
	res, err := b.pkgrpc.GetNetworkInfo()
	if err != nil {
		return 0, fmt.Errorf("fail to call btc.GetNetworkInfo(): %w", err)
	}
	if res.RelayFee == 0 {
		return 0, errors.New("RelayFee can not be retrieved by `getnetworkinfo`")
	}
	fee, err := btcpkg.FloatToAmount(res.RelayFee)
	if err != nil {
		return 0, err
	}
	return fee, nil
}
