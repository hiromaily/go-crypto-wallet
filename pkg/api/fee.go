package api

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

// EstimateSmartFeeResult estimatesmartfeeをcallしたresponseの型
type EstimateSmartFeeResult struct {
	FeeRate float64  `json:"feerate"`
	Errors  []string `json:"errors"`
	Blocks  int64    `json:"blocks"`
}

//Making Sense of Bitcoin Transaction Fees
//https://bitzuma.com/posts/making-sense-of-bitcoin-transaction-fees/

// EstimateSmartFee bitcoin coreの`estimatesmartfee`APIをcallする
// 戻り値はBTC/kB(float64)
func (b *Bitcoin) EstimateSmartFee() (float64, error) {
	input, err := json.Marshal(uint64(b.confirmationBlock)) //ここは固定(6)でいいはず
	if err != nil {
		return 0, errors.Errorf("json.Marchal(): error: %v", err)
	}
	rawResult, err := b.client.RawRequest("estimatesmartfee", []json.RawMessage{input})
	if err != nil {
		return 0, errors.Errorf("json.RawRequest(estimatesmartfee): error: %v", err)
	}

	estimateResult := EstimateSmartFeeResult{}
	err = json.Unmarshal([]byte(rawResult), &estimateResult)
	if err != nil {
		return 0, errors.Errorf("json.Unmarshal(): error: %v", err)
	}
	if len(estimateResult.Errors) != 0 {
		return 0, errors.Errorf("json.RawRequest(estimatesmartfee): error: %v", estimateResult.Errors[0])
	}

	//log.Printf("[Debug]Estimatesmartfee: %v: %f\n", estimateResult, estimateResult.FeeRate)
	//1.116e-05
	//0.000011 per 1kb

	return estimateResult.FeeRate, nil
}

// GetTransactionFee トランザクションサイズからfeeを算出する
func (b *Bitcoin) GetTransactionFee(tx *wire.MsgTx) (btcutil.Amount, error) {
	feePerKB, err := b.EstimateSmartFee()
	if err != nil {
		return 0, errors.Errorf("EstimateSmartFee(): error: %v", err)
	}
	fee := fmt.Sprintf("%f", feePerKB*float64(tx.SerializeSize())/1000)

	//To Amount
	feeAsBit, err := b.CastStrBitToAmount(fee)
	if err != nil {
		return 0, errors.Errorf("CastStrToSatoshi(%s): error: %v", fee, err)
	}

	return feeAsBit, nil
}

// ValidateAdjustmentFee 起動時に渡されたfeeの適用範囲をValidateする
func (b *Bitcoin) ValidateAdjustmentFee(fee float64) bool {
	//Rangeの範囲内であればOK
	if fee >= b.FeeRangeMin() && fee <= b.FeeRangeMax() {
		return true
	}
	return false
}

// CalculateNewFee 手数料を調整する
func (b *Bitcoin) CalculateNewFee(fee btcutil.Amount, adjustmentFee float64) (btcutil.Amount, error) {
	newFee, err := b.FloatBitToAmount(fee.ToBTC() * adjustmentFee)
	if err != nil {
		return 0, errors.Errorf("FloatBitToAmount() error: %v", err)
	}
	return newFee, nil
}
