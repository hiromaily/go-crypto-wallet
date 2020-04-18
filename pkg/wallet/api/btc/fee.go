package btc

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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
		return 0, errors.Errorf("json.Marchal(): error: %s", err)
	}
	rawResult, err := b.client.RawRequest("estimatesmartfee", []json.RawMessage{input})
	if err != nil {
		return 0, errors.Errorf("json.RawRequest(estimatesmartfee): error: %s", err)
	}

	estimateResult := EstimateSmartFeeResult{}
	err = json.Unmarshal([]byte(rawResult), &estimateResult)
	if err != nil {
		return 0, errors.Errorf("json.Unmarshal(): error: %s", err)
	}
	if len(estimateResult.Errors) != 0 {
		return 0, errors.Errorf("json.RawRequest(estimatesmartfee): error: %s", estimateResult.Errors[0])
	}

	return estimateResult.FeeRate, nil
}

// GetTransactionFee トランザクションサイズからfeeを算出する
func (b *Bitcoin) GetTransactionFee(tx *wire.MsgTx) (btcutil.Amount, error) {
	feePerKB, err := b.EstimateSmartFee()
	if err != nil {
		return 0, errors.Errorf("BTC.EstimateSmartFee(): error: %s", err)
	}
	fee := fmt.Sprintf("%f", feePerKB*float64(tx.SerializeSize())/1000)

	//To Amount
	feeAsBit, err := b.StrToAmount(fee)
	if err != nil {
		return 0, errors.Errorf("BTC.CastStrToSatoshi(%s): error: %s", fee, err)
	}

	return feeAsBit, nil
}

// GetFee 手数料を総合的に判断し取得する
func (b *Bitcoin) GetFee(tx *wire.MsgTx, adjustmentFee float64) (btcutil.Amount, error) {
	//通常の取得
	fee, err := b.GetTransactionFee(tx)
	if err != nil {
		return 0, errors.Errorf("BTC.GetTransactionFee(): error: %s", err)
	}
	b.logger.Debug("called GetTransactionFee()", zap.Any("fee", fee)) //0.000208 BTC

	//最低に満たない場合は、上書きをする
	relayFee, err := b.getMinRelayFee()
	if err != nil {
		//only log
		b.logger.Error("fail to call btc.getMinRelayFee()", zap.Error(err))
	} else {
		if fee < relayFee {
			fee = relayFee
		}
	}

	// オプションがある場合、feeの調整
	if b.validateAdjustmentFee(adjustmentFee) {
		newFee, err := b.calculateNewFee(fee, adjustmentFee)
		if err != nil {
			//only log
			b.logger.Error("fail to call btc.calculateNewFee()", zap.Error(err))
		}
		b.logger.Debug("called btc.calculateNewFee()", zap.Any("adjusted newFee", newFee)) //0.000208 BTC
		fee = newFee
	}

	return fee, nil
}

// ValidateAdjustmentFee 起動時に渡されたfeeの適用範囲をValidateする
func (b *Bitcoin) validateAdjustmentFee(fee float64) bool {
	//Rangeの範囲内であればOK
	if fee >= b.FeeRangeMin() && fee <= b.FeeRangeMax() {
		return true
	}
	return false
}

// CalculateNewFee 手数料を調整する
func (b *Bitcoin) calculateNewFee(fee btcutil.Amount, adjustmentFee float64) (btcutil.Amount, error) {
	newFee, err := b.FloatToAmount(fee.ToBTC() * adjustmentFee)
	if err != nil {
		return 0, errors.Errorf("BTC.FloatToAmount() error: %s", err)
	}
	return newFee, nil
}

func (b *Bitcoin) getMinRelayFee() (btcutil.Amount, error) {
	res, err := b.GetNetworkInfo()
	if err != nil {
		return 0, errors.Errorf("BTC.GetNetworkInfo() error: %s", err)
	}
	if res.Relayfee == 0 {
		return 0, errors.New("GetNetworkInfo().Relayfee error: RelayFee is not retrieved")
	}
	fee, err := b.FloatToAmount(res.Relayfee)
	if err != nil {
		return 0, errors.Errorf("BTC.FloatToAmount() error: %s", err)
	}
	return fee, nil
}
