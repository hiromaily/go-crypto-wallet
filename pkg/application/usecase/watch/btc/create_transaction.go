package btc

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	btcwatchsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/watch/btc"
)

// TxCreator interface defines the methods needed from the BTC TxCreate service
type TxCreator interface {
	CreateDepositTx(adjustmentFee float64) (string, string, error)
	CreatePaymentTx(adjustmentFee float64) (string, string, error)
	CreateTransferTx(
		sender, receiver domainAccount.AccountType,
		floatAmount, adjustmentFee float64,
	) (string, string, error)
}

type createTransactionUseCase struct {
	txCreator *btcwatchsrv.TxCreate
}

// NewCreateTransactionUseCase creates a new CreateTransactionUseCase
func NewCreateTransactionUseCase(txCreator *btcwatchsrv.TxCreate) watch.CreateTransactionUseCase {
	return &createTransactionUseCase{
		txCreator: txCreator,
	}
}

func (u *createTransactionUseCase) Execute(
	ctx context.Context,
	input watch.CreateTransactionInput,
) (watch.CreateTransactionOutput, error) {
	// Convert action type string to domain type
	actionType := domainTx.ActionType(input.ActionType)
	if !domainTx.ValidateActionType(input.ActionType) {
		return watch.CreateTransactionOutput{}, fmt.Errorf("invalid action type: %s", input.ActionType)
	}

	var hex, fileName string
	var execErr error

	switch actionType {
	case domainTx.ActionTypeDeposit:
		hex, fileName, execErr = u.txCreator.CreateDepositTx(input.AdjustmentFee)
	case domainTx.ActionTypePayment:
		hex, fileName, execErr = u.txCreator.CreatePaymentTx(input.AdjustmentFee)
	case domainTx.ActionTypeTransfer:
		hex, fileName, execErr = u.txCreator.CreateTransferTx(
			input.SenderAccount,
			input.ReceiverAccount,
			input.Amount,
			input.AdjustmentFee,
		)
	default:
		return watch.CreateTransactionOutput{}, fmt.Errorf("unsupported action type: %s", input.ActionType)
	}

	if execErr != nil {
		return watch.CreateTransactionOutput{}, fmt.Errorf("failed to create transaction: %w", execErr)
	}

	return watch.CreateTransactionOutput{
		TransactionHex: hex,
		FileName:       fileName,
	}, nil
}
