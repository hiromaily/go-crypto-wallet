package xrppublic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	xrpkg "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/xrplgo"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// TxInput is the transaction input JSON type used for building Payment transactions.
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

// toInfraTxInput converts a DTO TxInput to the infrastructure TxInput.
func toInfraTxInput(dto *dtoxrp.TxInput) *TxInput {
	if dto == nil {
		return nil
	}
	return &TxInput{
		TransactionType:    dto.TransactionType,
		Account:            dto.Account,
		Amount:             dto.Amount,
		Destination:        dto.Destination,
		Fee:                dto.Fee,
		Flags:              dto.Flags,
		LastLedgerSequence: dto.LastLedgerSequence,
		Sequence:           dto.Sequence,
		SigningPubKey:      dto.SigningPubKey,
		TxnSignature:       dto.TxnSignature,
		Hash:               dto.Hash,
	}
}

// toDTOTxInput converts the infrastructure TxInput to a DTO TxInput.
func toDTOTxInput(infra *TxInput) *dtoxrp.TxInput {
	if infra == nil {
		return nil
	}
	return &dtoxrp.TxInput{
		TransactionType:    infra.TransactionType,
		Account:            infra.Account,
		Amount:             infra.Amount,
		Destination:        infra.Destination,
		Fee:                infra.Fee,
		Flags:              infra.Flags,
		LastLedgerSequence: infra.LastLedgerSequence,
		Sequence:           infra.Sequence,
		SigningPubKey:      infra.SigningPubKey,
		TxnSignature:       infra.TxnSignature,
		Hash:               infra.Hash,
	}
}

// PrepareTransaction builds an unsigned Payment transaction using WebSocket account_info.
func (p *PublicXRP) PrepareTransaction(
	ctx context.Context, senderAccount, receiverAccount string, amount float64, instructions *dtoxrp.Instructions,
) (*dtoxrp.TxInput, string, error) {
	accInfo, err := p.PublicRPC.AccountInfo(ctx, senderAccount)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction(): %w", err)
	}
	if accInfo.Error != "" {
		return nil, "", fmt.Errorf("fail to call client.PrepareTransaction(): %s", accInfo.Error)
	}

	sequence := uint64(accInfo.Result.AccountData.Sequence)
	lastLedgerSequence := uint64(accInfo.Result.LedgerCurrentIndex) + xrpkg.MaxLedgerVersionOffset
	fee := "12" // minimum fee in drops

	if instructions != nil {
		if instructions.Fee != "" {
			fee = instructions.Fee
		}
		if instructions.Sequence != 0 {
			sequence = instructions.Sequence
		}
		if instructions.MaxLedgerVersion != 0 {
			lastLedgerSequence = instructions.MaxLedgerVersion
		} else if instructions.MaxLedgerVersionOffset != 0 {
			lastLedgerSequence = uint64(accInfo.Result.LedgerCurrentIndex) + instructions.MaxLedgerVersionOffset
		}
	}

	amountDrops := strconv.FormatInt(int64(amount*float64(dropsPerXRP)), 10)

	txInput := &TxInput{
		TransactionType:    "Payment",
		Account:            senderAccount,
		Amount:             amountDrops,
		Destination:        receiverAccount,
		Fee:                fee,
		Flags:              0,
		LastLedgerSequence: lastLedgerSequence,
		Sequence:           sequence,
	}

	jsonBytes, err := json.Marshal(txInput)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call json.Marshal(txInput): %w", err)
	}

	logger.Debug("PrepareTransaction", "TxJSON", string(jsonBytes))

	return toDTOTxInput(txInput), string(jsonBytes), nil
}

// CreateRawTransaction creates a raw unsigned Payment transaction.
// https://xrpl.org/ja/send-xrp.html
func (p *PublicXRP) CreateRawTransaction(
	ctx context.Context, senderAccount, receiverAccount string, amount float64, instructions *dtoxrp.Instructions,
) (*dtoxrp.TxInput, string, error) {
	if senderAccount == "" {
		return nil, "", errors.New("senderAccount is empty")
	}
	if receiverAccount == "" {
		return nil, "", errors.New("receiverAccount is empty")
	}

	accountInfo, err := p.GetAccountInfo(ctx, senderAccount)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call GetAccountInfo(): %w", err)
	}
	if amount != 0 && (xrpkg.ToFloat64(accountInfo.XrpBalance)-xrpkg.MinimumReserve) <= amount {
		return nil, "", fmt.Errorf("balance is short to send %s", accountInfo.XrpBalance)
	}

	txJSON, stringJSON, err := p.PrepareTransaction(ctx, senderAccount, receiverAccount, amount, instructions)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call PrepareTransaction(): %w", err)
	}
	feeDrops := xrpkg.XRPToDrops(xrpkg.ToFloat64(txJSON.Fee))
	calculatedAmount := xrpkg.ToFloat64(accountInfo.XrpBalance) - xrpkg.MinimumReserve - feeDrops
	if amount == 0 {
		if calculatedAmount <= 0 {
			return nil, "", fmt.Errorf("balance is short to send %s", accountInfo.XrpBalance)
		}
		txJSON, stringJSON, err = p.PrepareTransaction(
			ctx, senderAccount, receiverAccount, calculatedAmount, instructions,
		)
		if err != nil {
			return nil, "", fmt.Errorf("fail to call PrepareTransaction(): %w", err)
		}
	} else if calculatedAmount < amount {
		return nil, "", fmt.Errorf("balance is short to send %s", accountInfo.XrpBalance)
	}

	return txJSON, stringJSON, nil
}

// SubmitTransaction submits a signed transaction blob via WebSocket.
func (p *PublicXRP) SubmitTransaction(ctx context.Context, signedTx string) (*xrplgo.SentTx, uint64, error) {
	res, err := p.PublicRPC.Submit(ctx, signedTx)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to call client.SubmitTransaction(): %w", err)
	}
	if res.Error != "" {
		return nil, 0, fmt.Errorf("fail to call client.SubmitTransaction(): %s", res.Error)
	}

	logger.Debug("response of submitTransaction",
		"engine_result", res.Result.EngineResult,
		"engine_result_code", res.Result.EngineResultCode,
		"LastLedgerSequence", res.Result.TxJSON.LastLedgerSequence,
		"validated_ledger_index", res.Result.ValidatedLedgerIndex,
	)

	return &xrplgo.SentTx{
		ResultCode:          res.Result.EngineResult,
		ResultMessage:       res.Result.EngineResultMessage,
		EngineResult:        res.Result.EngineResult,
		EngineResultCode:    res.Result.EngineResultCode,
		EngineResultMessage: res.Result.EngineResultMessage,
		TxBlob:              res.Result.TxBlob,
		TxJSON: xrplgo.TxInput{
			TransactionType:    res.Result.TxJSON.TransactionType,
			Account:            res.Result.TxJSON.Account,
			Amount:             res.Result.TxJSON.Amount,
			Destination:        res.Result.TxJSON.Destination,
			Fee:                res.Result.TxJSON.Fee,
			Flags:              res.Result.TxJSON.Flags,
			LastLedgerSequence: res.Result.TxJSON.LastLedgerSequence,
			Sequence:           res.Result.TxJSON.Sequence,
			SigningPubKey:      res.Result.TxJSON.SigningPubKey,
			TxnSignature:       res.Result.TxJSON.TxnSignature,
			Hash:               res.Result.TxJSON.Hash,
		},
	}, res.Result.ValidatedLedgerIndex, nil
}

// GetTransaction retrieves a validated transaction by hash via WebSocket.
func (p *PublicXRP) GetTransaction(
	ctx context.Context, txID string, targetLedgerVersion uint64,
) (*xrplgo.TxInfo, error) {
	res, err := p.PublicRPC.GetTx(ctx, txID, targetLedgerVersion)
	if err != nil {
		return nil, fmt.Errorf("fail to call client.GetTransaction(): %w", err)
	}
	if res.Error != "" {
		return nil, fmt.Errorf("fail to get transaction info by %s: %s", txID, res.Error)
	}

	logger.Debug("response of getTransaction",
		"hash", res.Result.Hash,
		"ledger_index", res.Result.LedgerIndex,
		"validated", res.Result.Validated,
		"TransactionResult", res.Result.Meta.TransactionResult,
	)

	return &xrplgo.TxInfo{
		ID: res.Result.Hash,
		Outcome: xrplgo.TxOutcome{
			Result:        res.Result.Meta.TransactionResult,
			LedgerVersion: int(res.Result.LedgerIndex),
		},
	}, nil
}

// ToInfraTxInput converts a DTO TxInput to the infrastructure TxInput.
// Exported for use by sign wallets that work with TxInput directly.
func ToInfraTxInput(dto *dtoxrp.TxInput) *TxInput {
	return toInfraTxInput(dto)
}

// ToDTOTxInput converts the infrastructure TxInput to a DTO TxInput.
// Exported for use by callers that receive a TxInput and need to pass it upstream.
func ToDTOTxInput(infra *TxInput) *dtoxrp.TxInput {
	return toDTOTxInput(infra)
}
