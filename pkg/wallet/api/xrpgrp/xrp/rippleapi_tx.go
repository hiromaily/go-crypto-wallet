package xrp

import (
	"context"
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// - Send XRP https://xrpl.org/send-xrp.html
// - Payment System Basics https://xrpl.org/payment-system-basics.html

// TxInput is transaction input json type
type TxInput struct {
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

// SentTx is result transaction json type after sending
type SentTx struct {
	ResultCode          string  `json:"resultCode"`
	ResultMessage       string  `json:"resultMessage"`
	EngineResult        string  `json:"engine_result"`
	EngineResultCode    int     `json:"engine_result_code"`
	EngineResultMessage string  `json:"engine_result_message"`
	TxBlob              string  `json:"tx_blob"`
	TxJSON              TxInput `json:"tx_json"`
}

// TxInfo is result transaction json type after sending
type TxInfo struct {
	Type          string          `json:"type"`
	Address       string          `json:"address"`
	Sequence      int             `json:"sequence"`
	ID            string          `json:"id"`
	Specification TxSpecification `json:"specification"`
	Outcome       TxOutcome       `json:"outcome"`
}

// TxSpecification is part of TxInfo
type TxSpecification struct {
	Source      TxSpecSource      `json:"source"`
	Destination TxSpecDestination `json:"destination"`
}

// TxSpecSource is part of TxInfo
type TxSpecSource struct {
	Address   string   `json:"address"`
	MaxAmount TxAmount `json:"maxAmount"`
}

// TxAmount is part of TxInfo
type TxAmount struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

// TxTotalPrice is part of TxInfo
type TxTotalPrice struct {
	Currency     string `json:"currency"`
	Counterparty string `json:"counterparty"`
	Value        string `json:"value"`
}

// TxSpecDestination is part of TxInfo
type TxSpecDestination struct {
	Address string `json:"address"`
}

// TxOutcome is part of TxInfo
type TxOutcome struct {
	Result           string                         `json:"result"`
	Timestamp        time.Time                      `json:"timestamp"`
	Fee              string                         `json:"fee"`
	BalanceChanges   map[string][]TxAmount          `json:"balanceChanges"`
	OrderbookChanges map[string][]TxOrderbookChange `json:"orderbookChanges"`
	LedgerVersion    int                            `json:"ledgerVersion"`
	IndexInLedger    int                            `json:"indexInLedger"`
	DeliveredAmount  TxAmount                       `json:"deliveredAmount"`
}

// TxOrderbookChange is part of TxInfo
type TxOrderbookChange struct {
	Direction         string       `json:"direction"`
	Quantity          TxAmount     `json:"quantity"`
	TotalPrice        TxTotalPrice `json:"totalPrice"`
	MakerExchangeRate string       `json:"makerExchangeRate"`
	Sequence          int          `json:"sequence"`
	Status            string       `json:"status"`
}

// PrepareTransaction calls PrepareTransaction API
func (r *Ripple) PrepareTransaction(senderAccount, receiverAccount string, amount float64, instructions *Instructions) (*TxInput, string, error) {
	ctx := context.Background()
	req := &RequestPrepareTransaction{
		TxType:          EnumTransactionType_TX_PAYMENT,
		SenderAccount:   senderAccount,
		Amount:          amount,
		ReceiverAccount: receiverAccount,
		Instructions:    instructions,
	}

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", errors.Wrap(err, "fail to call client.PrepareTransaction()")
	}
	r.logger.Debug("response",
		zap.String("TxJSON", res.TxJSON),
		zap.Any("Instructions", res.Instructions),
	)

	var txInput TxInput
	unquotedJSON, _ := strconv.Unquote(res.TxJSON)
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", errors.Wrap(err, "fail to call json.Unmarshal(txJSON)")
	}

	return &txInput, unquotedJSON, nil
}

// SignTransaction calls SignTransaction API
// Offline functionality
// - https://xrpl.org/rippleapi-reference.html#offline-functionality
func (r *Ripple) SignTransaction(txInput *TxInput, secret string) (string, string, error) {
	ctx := context.Background()
	strJSON, err := json.Marshal(txInput)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call json.Marshal(txJSON)")
	}
	req := &RequestSignTransaction{
		TxJSON: string(strJSON),
		Secret: secret,
	}

	res, err := r.API.txClient.SignTransaction(ctx, req)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call client.SignTransaction()")
	}

	return res.TxID, res.TxBlob, nil
}

// CombineTransaction combines signed transactions from multiple accounts for a multisignature transaction.
// - The signed transaction must subsequently be submitted.
func (r *Ripple) CombineTransaction(signedTxs []string) (string, string, error) {
	ctx := context.Background()
	req := &RequestCombineTransaction{
		SignedTransactions: signedTxs,
	}

	res, err := r.API.txClient.CombineTransaction(ctx, req)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call client.CombineTransaction()")
	}

	return res.TxID, res.SignedTransaction, nil
}

// SubmitTransaction calls SubmitTransaction API
// - signedTx is returned TxBlob by SignTransaction()
func (r *Ripple) SubmitTransaction(signedTx string) (*SentTx, uint64, error) {
	ctx := context.Background()
	req := &RequestSubmitTransaction{
		TxBlob: signedTx,
	}
	res, err := r.API.txClient.SubmitTransaction(ctx, req)
	if err != nil {
		return nil, 0, errors.Wrap(err, "fail to call client.SubmitTransaction()")
	}

	var sentTxJSON SentTx
	if err = json.Unmarshal([]byte(res.ResultJSONString), &sentTxJSON); err != nil {
		return nil, 0, errors.Wrap(err, "fail to call json.Unmarshal(sentTxJSON)")
	}

	// FIXME:
	// res.EarliestLedgerVersion may be useless because SentTxJSON includes `LastLedgerSequence` and it would be useful
	r.logger.Debug("response of submitTransaction",
		zap.String("res.ResultJSONString", res.ResultJSONString),
		zap.Uint64("res.EarliestLedgerVersion", res.EarliestLedgerVersion),
		zap.Uint64("sentTxJSON.TxJSON.LastLedgerSequence", sentTxJSON.TxJSON.LastLedgerSequence),
	)
	// res.EarliestLedgerVersion => for when calling GetTransaction()
	// sentTxJSON.TxJSON.LastLedgerSequence => for when calling WaitValidation()

	return &sentTxJSON, res.EarliestLedgerVersion, nil
	// return &sentTxJSON, sentTxJSON.TxJSON.LastLedgerSequence, nil
}

// WaitValidation calls WaitValidation API
// - handling server streaming
func (r *Ripple) WaitValidation(targetledgerVarsion uint64) (uint64, error) {
	ctx := context.Background()
	req := &emptypb.Empty{}
	resStream, err := r.API.txClient.WaitValidation(ctx, req)
	if err != nil {
		return 0, errors.Wrap(err, "fail to call client.WaitValidation()")
	}

	defer func() {
		r.logger.Debug("running in defer func()")
		if err := resStream.CloseSend(); err != nil {
			r.logger.Warn("fail to call resStream.CloseSend()")
		}
	}()

	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			r.logger.Warn("server is closed in WaitValidation()")
			return 0, errors.New("server is closed")
		} else if err != nil {
			if respErr, ok := status.FromError(err); ok {
				switch respErr.Code() {
				case codes.InvalidArgument:
					r.logger.Warn("parameter is invalid in WaitValidation()")
				case codes.DeadlineExceeded:
					r.logger.Warn("timeout in WaitValidation()")
				default:
					r.logger.Warn("gRPC error in WaitValidation()",
						zap.Uint32("code", uint32(respErr.Code())),
						zap.String("message", respErr.Message()),
					)
				}
			} else {
				r.logger.Warn("fail to call resStream.Recv()", zap.Error(err))
			}
			// break
			return 0, errors.Wrap(err, "fail to call resStream.Recv()")
		}
		// success
		r.logger.Info("response in WaitValidation()", zap.Uint64("LedgerVersion", res.LedgerVersion))
		if targetledgerVarsion <= res.LedgerVersion {
			// done
			return res.LedgerVersion, nil
		}
		// continue
	}
}

// GetTransaction calls GetTransaction API
func (r *Ripple) GetTransaction(txID string, targetLedgerVersion uint64) (*TxInfo, error) {
	ctx := context.Background()
	req := &RequestGetTransaction{
		TxID:             txID,
		MinLedgerVersion: targetLedgerVersion,
	}
	res, err := r.API.txClient.GetTransaction(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call client.GetTransaction()")
	}

	if res.ResultJSONString == "" {
		return nil, errors.Errorf("fail to get transaction info by %s", txID)
	}

	r.logger.Debug("response of getTransaction",
		zap.String("res.ResultJSONString", res.ResultJSONString),
	)

	var txInfo TxInfo
	if err = json.Unmarshal([]byte(res.ResultJSONString), &txInfo); err != nil {
		return nil, errors.Wrap(err, "fail to call json.Unmarshal(txInfo)")
	}
	// TODO: check
	// txInfo.Outcome.Result : tesSUCCESS
	return &txInfo, nil
}
