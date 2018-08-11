package api

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/wire"
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
func (b *Bitcoin) EstimateSmartFee() (float64, error) {
	input, err := json.Marshal(uint64(b.confirmationBlock)) //ここは固定(6)でいいはず
	if err != nil {
		return 0, errors.Errorf("json.Marchal(): error: %v", err)
	}
	rawResult, err := b.client.RawRequest("estimatesmartfee", []json.RawMessage{input})
	if err != nil {
		//-3: Expected type number, got object
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

	//TODO:Amoutとして返したほうがいいかも
	return estimateResult.FeeRate, nil
}

// GetTransactionFee トランザクションサイズからfeeを算出する
//FIXME: WIP
func (b *Bitcoin) GetTransactionFee(tx *wire.MsgTx) (string, error) {
	feePerKB, err := b.EstimateSmartFee()
	if err != nil {
		return "", errors.Errorf("EstimateSmartFee(): error: %v", err)
	}
	fee := fmt.Sprintf("%f", feePerKB*float64(tx.SerializeSize())/1000)

	return fee, nil
}
