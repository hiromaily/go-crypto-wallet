package xrp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	dtoRipple "github.com/hiromaily/go-crypto-wallet/internal/application/dto/ripple"
	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/protogen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// - Send XRP https://xrpl.org/send-xrp.html
// - Payment System Basics https://xrpl.org/payment-system-basics.html

// unquoteJSON attempts to unquote a JSON string. If unquoting fails,
// it returns the original string. This handles cases where the gRPC
// response may or may not include extra quotes around the JSON.
func unquoteJSON(s string) string {
	unquoted, err := strconv.Unquote(s)
	if err != nil {
		// If unquoting fails, assume the string is not quoted
		return s
	}
	return unquoted
}

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
func (r *Ripple) PrepareTransaction(
	ctx context.Context, senderAccount, receiverAccount string, amount float64, instructions *dtoRipple.Instructions,
) (*dtoRipple.TxInput, string, error) {
	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	req := protogen.RequestPrepareTransaction_builder{
		TxType:          protogen.EnumTransactionType_TX_PAYMENT,
		SenderAccount:   senderAccount,
		Amount:          amount,
		ReceiverAccount: receiverAccount,
		Instructions:    infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction(): %w", err)
	}
	logger.Debug("response",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput TxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(txJSON): %w", err)
	}

	// Convert infrastructure type to DTO
	return ToDTOTxInput(&txInput), unquotedJSON, nil
}

// SetRegularKeyTxInput is the transaction input for SetRegularKey
type SetRegularKeyTxInput struct {
	TransactionType    string `json:"TransactionType"`
	Account            string `json:"Account"`
	RegularKey         string `json:"RegularKey,omitempty"` // Omit to remove regular key
	Fee                string `json:"Fee"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence"`
	Sequence           uint64 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey,omitempty"`
	TxnSignature       string `json:"TxnSignature,omitempty"`
	Hash               string `json:"hash,omitempty"`
}

// PrepareSetRegularKeyTransaction prepares a SetRegularKey transaction
// - regularKey: the address to authorize as regular key, or empty to remove
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/setregularkey
func (r *Ripple) PrepareSetRegularKeyTransaction(
	ctx context.Context, senderAccount, regularKey string, instructions *dtoRipple.Instructions,
) (*SetRegularKeyTxInput, string, error) {
	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	req := protogen.RequestPrepareTransaction_builder{
		TxType:        protogen.EnumTransactionType_TX_SET_REGULAR_KEY,
		SenderAccount: senderAccount,
		RegularKey:    regularKey,
		Instructions:  infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction() for SetRegularKey: %w", err)
	}
	logger.Debug("response SetRegularKey",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput SetRegularKeyTxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(SetRegularKeyTxJSON): %w", err)
	}

	return &txInput, unquotedJSON, nil
}

// AccountSetTxInput is the transaction input for AccountSet
type AccountSetTxInput struct {
	TransactionType    string `json:"TransactionType"`
	Account            string `json:"Account"`
	SetFlag            uint32 `json:"SetFlag,omitempty"`
	ClearFlag          uint32 `json:"ClearFlag,omitempty"`
	Fee                string `json:"Fee"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence"`
	Sequence           uint64 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey,omitempty"`
	TxnSignature       string `json:"TxnSignature,omitempty"`
	Hash               string `json:"hash,omitempty"`
}

// AccountSet flag constants
const (
	// AsfDisableMaster disables the master key from signing transactions
	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/accountset#accountset-flags
	AsfDisableMaster uint32 = 4
)

// SignerListEntry represents a single signer in the SignerList
type SignerListEntry struct {
	SignerEntry struct {
		Account      string `json:"Account"`
		SignerWeight uint32 `json:"SignerWeight"`
	} `json:"SignerEntry"`
}

// SignerListSetTxInput is the transaction input for SignerListSet
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/signerlistset
type SignerListSetTxInput struct {
	TransactionType    string            `json:"TransactionType"`
	Account            string            `json:"Account"`
	SignerQuorum       uint32            `json:"SignerQuorum"`
	SignerEntries      []SignerListEntry `json:"SignerEntries,omitempty"`
	Fee                string            `json:"Fee"`
	Flags              uint64            `json:"Flags"`
	LastLedgerSequence uint64            `json:"LastLedgerSequence"`
	Sequence           uint64            `json:"Sequence"`
	SigningPubKey      string            `json:"SigningPubKey,omitempty"`
	TxnSignature       string            `json:"TxnSignature,omitempty"`
	Hash               string            `json:"hash,omitempty"`
}

// PrepareAccountSetTransaction prepares an AccountSet transaction
// - setFlag: flag to set (e.g., AsfDisableMaster = 4)
// - clearFlag: flag to clear
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/accountset
func (r *Ripple) PrepareAccountSetTransaction(
	ctx context.Context, senderAccount string, setFlag, clearFlag uint32, instructions *dtoRipple.Instructions,
) (*AccountSetTxInput, string, error) {
	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	req := protogen.RequestPrepareTransaction_builder{
		TxType:        protogen.EnumTransactionType_TX_ACCOUNT_SET,
		SenderAccount: senderAccount,
		SetFlag:       setFlag,
		ClearFlag:     clearFlag,
		Instructions:  infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction() for AccountSet: %w", err)
	}
	logger.Debug("response AccountSet",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput AccountSetTxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(AccountSetTxJSON): %w", err)
	}

	return &txInput, unquotedJSON, nil
}

// PrepareSignerListSetTransaction prepares a SignerListSet transaction
// - signerQuorum: minimum total weight of signatures required (0 to remove signer list)
// - signerEntries: list of signers with their weights (must be empty if signerQuorum is 0)
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/signerlistset
func (r *Ripple) PrepareSignerListSetTransaction(
	ctx context.Context,
	senderAccount string,
	signerQuorum uint32,
	signerEntries []apixrp.SignerEntryInput,
	instructions *dtoRipple.Instructions,
) (*dtoRipple.SignerListSetTxInput, string, error) {
	// Convert DTO to infrastructure type
	infraInstructions := ToInfraInstructions(instructions)

	// Convert SignerEntryInput to protogen.SignerEntry
	protoSignerEntries := make([]*protogen.SignerEntry, len(signerEntries))
	for i, entry := range signerEntries {
		protoSignerEntries[i] = protogen.SignerEntry_builder{
			Account: entry.Account,
			Weight:  entry.Weight,
		}.Build()
	}

	req := protogen.RequestPrepareTransaction_builder{
		TxType:        protogen.EnumTransactionType_TX_SINGER_LIST_SET,
		SenderAccount: senderAccount,
		SignerQuorum:  signerQuorum,
		SignerEntries: protoSignerEntries,
		Instructions:  infraInstructions,
	}.Build()

	res, err := r.API.txClient.PrepareTransaction(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction() for SignerListSet: %w", err)
	}
	logger.Debug("response SignerListSet",
		"TxJSON", res.GetTxJSON(),
		"Instructions", res.GetInstructions(),
	)

	var txInput SignerListSetTxInput
	unquotedJSON := unquoteJSON(res.GetTxJSON())
	if err = json.Unmarshal([]byte(unquotedJSON), &txInput); err != nil {
		return nil, "", fmt.Errorf("fail to call json.Unmarshal(SignerListSetTxJSON): %w", err)
	}

	// Convert infrastructure type to DTO
	return ToDTOSignerListSetTxInput(&txInput), unquotedJSON, nil
}

// SignSignerListSetTransaction signs a SignerListSet transaction
func (r *Ripple) SignSignerListSetTransaction(
	ctx context.Context, txInput *SignerListSetTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "SignerListSet")
}

// signTransactionJSON is a generic helper that signs any transaction type.
// It marshals the input to JSON and calls the gRPC SignTransaction API.
func (r *Ripple) signTransactionJSON(
	ctx context.Context, txInput any, secret, txTypeName string,
) (string, string, error) {
	strJSON, err := json.Marshal(txInput)
	if err != nil {
		return "", "", fmt.Errorf("fail to call json.Marshal(%sTxInput): %w", txTypeName, err)
	}
	req := protogen.RequestSignTransaction_builder{
		TxJSON: string(strJSON),
		Secret: secret,
	}.Build()

	res, err := r.API.txClient.SignTransaction(ctx, req)
	if err != nil {
		return "", "", fmt.Errorf("fail to call client.SignTransaction() for %s: %w", txTypeName, err)
	}

	return res.GetTxID(), res.GetTxBlob(), nil
}

// SignSetRegularKeyTransaction signs a SetRegularKey transaction
func (r *Ripple) SignSetRegularKeyTransaction(
	ctx context.Context, txInput *SetRegularKeyTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "SetRegularKey")
}

// SignAccountSetTransaction signs an AccountSet transaction
func (r *Ripple) SignAccountSetTransaction(
	ctx context.Context, txInput *AccountSetTxInput, secret string,
) (string, string, error) {
	return r.signTransactionJSON(ctx, txInput, secret, "AccountSet")
}

// SignTransaction calls SignTransaction API
// Offline functionality
// - https://xrpl.org/rippleapi-reference.html#offline-functionality
func (r *Ripple) SignTransaction(
	ctx context.Context, txInput *dtoRipple.TxInput, secret string,
) (string, string, error) {
	// Convert DTO to infrastructure type
	infraTxInput := ToInfraTxInput(txInput)

	strJSON, err := json.Marshal(infraTxInput)
	if err != nil {
		return "", "", fmt.Errorf("fail to call json.Marshal(txJSON): %w", err)
	}
	req := protogen.RequestSignTransaction_builder{
		TxJSON: string(strJSON),
		Secret: secret,
	}.Build()

	res, err := r.API.txClient.SignTransaction(ctx, req)
	if err != nil {
		return "", "", fmt.Errorf("fail to call client.SignTransaction(): %w", err)
	}

	return res.GetTxID(), res.GetTxBlob(), nil
}

// CombineTransaction combines signed transactions from multiple accounts for a multisignature transaction.
// - The signed transaction must subsequently be submitted.
func (r *Ripple) CombineTransaction(ctx context.Context, signedTxs []string) (string, string, error) {
	req := protogen.RequestCombineTransaction_builder{
		SignedTransactions: signedTxs,
	}.Build()

	res, err := r.API.txClient.CombineTransaction(ctx, req)
	if err != nil {
		return "", "", fmt.Errorf("fail to call client.CombineTransaction(): %w", err)
	}

	return res.GetTxID(), res.GetSignedTransaction(), nil
}

// SubmitTransaction calls SubmitTransaction API
// - signedTx is returned TxBlob by SignTransaction()
func (r *Ripple) SubmitTransaction(ctx context.Context, signedTx string) (*dtoRipple.SentTx, uint64, error) {
	req := protogen.RequestSubmitTransaction_builder{
		TxBlob: signedTx,
	}.Build()
	res, err := r.API.txClient.SubmitTransaction(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to call client.SubmitTransaction(): %w", err)
	}

	var sentTxJSON SentTx
	if err = json.Unmarshal([]byte(res.GetResultJSONString()), &sentTxJSON); err != nil {
		return nil, 0, fmt.Errorf("fail to call json.Unmarshal(sentTxJSON): %w", err)
	}

	// FIXME:
	// res.EarliestLedgerVersion may be useless because SentTxJSON includes `LastLedgerSequence` and it would be useful
	logger.Debug("response of submitTransaction",
		"res.ResultJSONString", res.GetResultJSONString(),
		"res.EarliestLedgerVersion", res.GetEarliestLedgerVersion(),
		"sentTxJSON.TxJSON.LastLedgerSequence", sentTxJSON.TxJSON.LastLedgerSequence,
	)
	// res.EarliestLedgerVersion => for when calling GetTransaction()
	// sentTxJSON.TxJSON.LastLedgerSequence => for when calling WaitValidation()

	// Convert infrastructure type to DTO
	return ToDTOSentTx(&sentTxJSON), res.GetEarliestLedgerVersion(), nil
	// return ToDTOSentTx(&sentTxJSON), sentTxJSON.TxJSON.LastLedgerSequence, nil
}

// WaitValidation calls WaitValidation API
// - handling server streaming
func (r *Ripple) WaitValidation(ctx context.Context, targetledgerVarsion uint64) (uint64, error) {
	req := &emptypb.Empty{}
	resStream, err := r.API.txClient.WaitValidation(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("fail to call client.WaitValidation(): %w", err)
	}

	defer func() {
		logger.Debug("running in defer func()")
		if closeErr := resStream.CloseSend(); closeErr != nil {
			logger.Warn("fail to call resStream.CloseSend()")
		}
	}()

	for {
		var res *protogen.ResponseWaitValidation
		res, err = resStream.Recv()
		if err == io.EOF {
			logger.Warn("server is closed in WaitValidation()")
			return 0, errors.New("server is closed")
		} else if err != nil {
			if respErr, ok := status.FromError(err); ok {
				switch respErr.Code() {
				case codes.InvalidArgument:
					logger.Warn("parameter is invalid in WaitValidation()")
				case codes.DeadlineExceeded:
					logger.Warn("timeout in WaitValidation()")
				case codes.OK, codes.Canceled, codes.Unknown, codes.NotFound, codes.AlreadyExists,
					codes.PermissionDenied, codes.ResourceExhausted, codes.FailedPrecondition,
					codes.Aborted, codes.OutOfRange, codes.Unimplemented, codes.Internal,
					codes.Unavailable, codes.DataLoss, codes.Unauthenticated:
					logger.Warn("gRPC error in WaitValidation()",
						"code", uint32(respErr.Code()),
						"message", respErr.Message(),
					)
				default:
					logger.Warn("gRPC error in WaitValidation()",
						"code", uint32(respErr.Code()),
						"message", respErr.Message(),
					)
				}
			} else {
				logger.Warn("fail to call resStream.Recv()", "error", err)
			}
			// break
			return 0, fmt.Errorf("fail to call resStream.Recv(): %w", err)
		}
		// success
		logger.Info("response in WaitValidation()", "LedgerVersion", res.GetLedgerVersion())
		if targetledgerVarsion <= res.GetLedgerVersion() {
			// done
			return res.GetLedgerVersion(), nil
		}
		// continue
	}
}

// GetTransaction calls GetTransaction API
func (r *Ripple) GetTransaction(
	ctx context.Context, txID string, targetLedgerVersion uint64,
) (*dtoRipple.TxInfo, error) {
	req := protogen.RequestGetTransaction_builder{
		TxID:             txID,
		MinLedgerVersion: targetLedgerVersion,
	}.Build()
	res, err := r.API.txClient.GetTransaction(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("fail to call client.GetTransaction(): %w", err)
	}

	if res.GetResultJSONString() == "" {
		return nil, fmt.Errorf("fail to get transaction info by %s", txID)
	}

	logger.Debug("response of getTransaction",
		"res.ResultJSONString", res.GetResultJSONString(),
	)

	var txInfo TxInfo
	if err = json.Unmarshal([]byte(res.GetResultJSONString()), &txInfo); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(txInfo): %w", err)
	}
	// TODO: check
	// txInfo.Outcome.Result : tesSUCCESS

	// Convert infrastructure type to DTO
	return ToDTOTxInfo(&txInfo), nil
}
