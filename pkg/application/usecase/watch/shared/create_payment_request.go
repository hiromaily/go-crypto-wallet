package shared

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch"
	sharedwatchsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/watch/shared"
)

type createPaymentRequestUseCase struct {
	paymentRequestCreator *sharedwatchsrv.PaymentRequestCreate
}

// NewCreatePaymentRequestUseCase creates a new CreatePaymentRequestUseCase
func NewCreatePaymentRequestUseCase(paymentRequestCreator *sharedwatchsrv.PaymentRequestCreate) watch.CreatePaymentRequestUseCase {
	return &createPaymentRequestUseCase{
		paymentRequestCreator: paymentRequestCreator,
	}
}

func (u *createPaymentRequestUseCase) Execute(ctx context.Context, input watch.CreatePaymentRequestInput) error {
	if err := u.paymentRequestCreator.CreatePaymentRequest(input.AmountList); err != nil {
		return fmt.Errorf("failed to create payment request: %w", err)
	}
	return nil
}
