package btc

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
)

// EstimateSmartFeeResult is response type of PRC `estimatesmartfee`
type EstimateSmartFeeResult struct {
	FeeRate float64  `json:"feerate"`
	Errors  []string `json:"errors"`
	Blocks  uint64   `json:"blocks"`
}

// Making Sense of Bitcoin Transaction Fees
// https://bitzuma.com/posts/making-sense-of-bitcoin-transaction-fees/

// EstimateSmartFee calls RPC `estimatesmartfee` and returns BTC/kB(float64)
func (b *Bitcoin) EstimateSmartFee() (float64, error) {
	input, err := json.Marshal(b.confirmationBlock)
	if err != nil {
		return 0, fmt.Errorf("fail to call json.Marchal(confirmationBlock): %w", err)
	}
	rawResult, err := b.Client.RawRequest("estimatesmartfee", []json.RawMessage{input})
	if err != nil {
		return 0, fmt.Errorf("fail to call json.RawRequest(estimatesmartfee): %w", err)
	}

	estimateResult := EstimateSmartFeeResult{}
	err = json.Unmarshal(rawResult, &estimateResult)
	if err != nil {
		return 0, errors.New("fail to all json.Unmarshal(rawResult)")
	}
	if len(estimateResult.Errors) != 0 {
		return 0, fmt.Errorf("response includes error: %s", estimateResult.Errors[0])
	}

	return estimateResult.FeeRate, nil
}

// GetTransactionFee calculate fee from transaction size
func (b *Bitcoin) GetTransactionFee(tx *wire.MsgTx) (btcutil.Amount, error) {
	feePerKB, err := b.EstimateSmartFee()
	if err != nil {
		return 0, fmt.Errorf("fail to call btc.EstimateSmartFee(): %w", err)
	}
	fee := fmt.Sprintf("%f", feePerKB*float64(tx.SerializeSize())/1000)

	// To Amount
	feeAsBit, err := b.StrToAmount(fee)
	if err != nil {
		return 0, err
	}

	return feeAsBit, nil
}

// GetFee get more preferable fee
func (b *Bitcoin) GetFee(tx *wire.MsgTx, adjustmentFee float64) (btcutil.Amount, error) {
	// get tx fee
	fee, err := b.GetTransactionFee(tx)
	if err != nil {
		return 0, err
	}
	// b.logger.Debug("called GetTransactionFee()", "fee", fee) //0.000208 BTC

	// if response doesn't meet minimum fee, it should be overridden
	relayFee, err := b.getMinRelayFee()
	if err != nil {
		b.logger.Warn("fail to call btc.getMinRelayFee() but continue", "error", err)
	} else if fee < relayFee {
		fee = relayFee
	}

	// if adjustmentFee param is given
	if b.validateAdjustmentFee(adjustmentFee) {
		var newFee btcutil.Amount
		newFee, err = b.calculateNewFee(fee, adjustmentFee)
		if err != nil {
			b.logger.Warn("fail to call btc.calculateNewFee() but continue", "error", err)
		}
		b.logger.Debug("called btc.calculateNewFee()", "adjusted newFee", newFee) // 0.000208 BTC
		fee = newFee
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

// CalculateNewFee adjust fee by adjustment fee
func (b *Bitcoin) calculateNewFee(fee btcutil.Amount, adjustmentFee float64) (btcutil.Amount, error) {
	newFee, err := b.FloatToAmount(fee.ToBTC() * adjustmentFee)
	if err != nil {
		return 0, err
	}
	return newFee, nil
}

func (b *Bitcoin) getMinRelayFee() (btcutil.Amount, error) {
	res, err := b.GetNetworkInfo()
	if err != nil {
		return 0, fmt.Errorf("fail to call btc.GetNetworkInfo(): %w", err)
	}
	if res.Relayfee == 0 {
		return 0, errors.New("RelayFee can not be retrieved by `getnetworkinfo`")
	}
	fee, err := b.FloatToAmount(res.Relayfee)
	if err != nil {
		return 0, err
	}
	return fee, nil
}
