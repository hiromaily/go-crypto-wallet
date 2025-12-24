package xrp

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	xrpwatchsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/watch/xrp"
)

// TxCreator interface defines the methods needed from the XRP TxCreate service
type TxCreator interface {
	CreateDepositTx() (string, string, error)
	CreatePaymentTx() (string, string, error)
	CreateTransferTx(sender, receiver domainAccount.AccountType, floatValue float64) (string, string, error)
}

type createTransactionUseCase struct {
	txCreator *xrpwatchsrv.TxCreate
}

// NewCreateTransactionUseCase creates a new CreateTransactionUseCase
func NewCreateTransactionUseCase(txCreator *xrpwatchsrv.TxCreate) watch.CreateTransactionUseCase {
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
		// Note: XRP CreateDepositTx does not take adjustmentFee parameter
		hex, fileName, execErr = u.txCreator.CreateDepositTx()
	case domainTx.ActionTypePayment:
		// Note: XRP CreatePaymentTx does not take adjustmentFee parameter
		hex, fileName, execErr = u.txCreator.CreatePaymentTx()
	case domainTx.ActionTypeTransfer:
		// Note: XRP CreateTransferTx does not take adjustmentFee parameter
		hex, fileName, execErr = u.txCreator.CreateTransferTx(
			input.SenderAccount,
			input.ReceiverAccount,
			input.Amount,
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
