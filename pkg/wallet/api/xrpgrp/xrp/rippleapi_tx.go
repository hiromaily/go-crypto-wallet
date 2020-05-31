package xrp

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	pb "github.com/hiromaily/ripple-lib-proto/pb/go/rippleapi"
)

// Send XRP
// https://xrpl.org/send-xrp.html

// TxJSON is transaction json type
type TxJSON struct {
	TransactionType    string `json:"TransactionType"`
	Account            string `json:"Account"`
	Amount             string `json:"Amount"`
	Destination        string `json:"Destination"`
	Fee                string `json:"Fee"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence"`
	Sequence           uint64 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey,omitempty"`
	TxnSignature       string `json:"TxnSignature,omitempty"`
	Hash               string `json:"hash,omitempty"`
}

// SentTxJSON is result transaction json type after sending
type SentTxJSON struct {
	ResultCode          string `json:"resultCode"`
	ResultMessage       string `json:"resultMessage"`
	EngineResult        string `json:"engine_result"`
	EngineResultCode    int    `json:"engine_result_code"`
	EngineResultMessage string `json:"engine_result_message"`
	TxBlob              string `json:"tx_blob"`
	TxJSON              TxJSON `json:"tx_json"`
}

// PrepareTransaction calls PrepareTransaction API
func (r *Ripple) PrepareTransaction(senderAccount, receiverAccount string, amount float64) (*TxJSON, error) {

	ctx := context.Background()
	req := &pb.RequestPrepareTransaction{
		TxType:          pb.TX_PAYMENT,
		SenderAccount:   senderAccount,
		Amount:          amount,
		ReceiverAccount: receiverAccount,
		Instructions:    &pb.Instructions{MaxLedgerVersionOffset: 75},
	}

	//res: *pb.ResponsePrepareTransaction
	res, err := r.API.client.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call client.PrepareTransaction()")
	}
	r.logger.Debug("response",
		zap.String("TxJSON", res.TxJSON),
		zap.Any("Instructions", res.Instructions),
	)

	var txJSON TxJSON
	unquotedJSON, _ := strconv.Unquote(res.TxJSON)
	if err = json.Unmarshal([]byte(unquotedJSON), &txJSON); err != nil {
		return nil, errors.Wrap(err, "fail to call json.Unmarshal(txJSON)")
	}

	return &txJSON, nil
}

// SignTransaction calls SignTransaction API
// TODO: Is it possible to run offline?
func (r *Ripple) SignTransaction(txJSON *TxJSON, secret string) (string, string, error) {
	ctx := context.Background()
	strJSON, err := json.Marshal(txJSON)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call json.Marshal(txJSON)")
	}
	req := &pb.RequestSignTransaction{
		TxJSON: string(strJSON),
		Secret: secret,
	}

	res, err := r.API.client.SignTransaction(ctx, req)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call client.SignTransaction()")
	}

	return res.TxID, res.TxBlob, nil
}

// SubmitTransaction calls SubmitTransaction API
// - signedTx is returned TxBlob by SignTransaction()
func (r *Ripple) SubmitTransaction(signedTx string) (*SentTxJSON, uint32, error) {
	ctx := context.Background()
	req := &pb.RequestSubmitTransaction{
		TxBlob: signedTx,
	}
	res, err := r.API.client.SubmitTransaction(ctx, req)
	if err != nil {
		return nil, 0, errors.Wrap(err, "fail to call client.SubmitTransaction()")
	}

	r.logger.Debug("response of submitTransaction",
		zap.String("res.ResultJSONString", res.ResultJSONString),
	)

	var sentTxJSON SentTxJSON
	if err = json.Unmarshal([]byte(res.ResultJSONString), &sentTxJSON); err != nil {
		return nil, 0, errors.Wrap(err, "fail to call json.Unmarshal(sentTxJSON)")
	}

	return &sentTxJSON, res.EarliestLedgerVersion, nil
}
