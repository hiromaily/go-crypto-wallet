package api

import (
	"encoding/json"
	"log"

	"github.com/pkg/errors"
)

// EstimateSmartFee input (これは不要のはず)
type EstimateSmartFee struct {
	ConfTarget   int    `json:"conf_target"`
	EstimateMode string `json:"estimate_mode"`
}

// EstimateSmartFeeResult estimatesmartfeeをcallしたresponseの型
type EstimateSmartFeeResult struct {
	FeeRate float32  `json:"feerate"`
	Errors  []string `json:"errors"`
	Blocks  int64    `json:"blocks"`
}

// EstimateSmartFee bitcoin coreの`estimatesmartfee`APIをcallする
func (b *Bitcoin) EstimateSmartFee() (float32, error) {
	//TODO:ここはオンラインでしか実行できない？？
	//TODO:送信前に手数料を取得する
	//Estimatesmartfee
	//fee, err := bit.Client.EstimateFee(1)
	//if err != nil {
	//	//estimatefee is deprecated and will be fully removed in v0.17. To use estimatefee in v0.16, restart bitcoind with -deprecatedrpc=estimatefee.
	//	log.Fatal(err)
	//}
	//log.Printf("Estimatesmartfee: %v\n", fee)

	//param := EstimateSmartFee{ConfTarget: 6}
	//b, err := json.Marshal(param)
	input, err := json.Marshal(uint64(6)) //ここは固定でいいはず
	if err != nil {
		return 0, errors.Errorf("json.Marchal(): error: %v", err)
	}
	rawResult, err := b.Client.RawRequest("estimatesmartfee", []json.RawMessage{input})
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

	log.Printf("Estimatesmartfee: %v\n", estimateResult)
	//1.116e-05
	log.Printf("%f", estimateResult.FeeRate)
	//0.000011 per 1kb

	return estimateResult.FeeRate, nil
}
