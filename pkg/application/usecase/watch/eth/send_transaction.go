package eth

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch"
	ethwatchsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/watch/eth"
)

type sendTransactionUseCase struct {
	txSender *ethwatchsrv.TxSend
}

// NewSendTransactionUseCase creates a new SendTransactionUseCase
func NewSendTransactionUseCase(txSender *ethwatchsrv.TxSend) watch.SendTransactionUseCase {
	return &sendTransactionUseCase{
		txSender: txSender,
	}
}

func (u *sendTransactionUseCase) Execute(
	ctx context.Context,
	input watch.SendTransactionInput,
) (watch.SendTransactionOutput, error) {
	txID, err := u.txSender.SendTx(input.FilePath)
	if err != nil {
		return watch.SendTransactionOutput{}, fmt.Errorf("failed to send transaction: %w", err)
	}

	return watch.SendTransactionOutput{
		TxID: txID,
	}, nil
}
