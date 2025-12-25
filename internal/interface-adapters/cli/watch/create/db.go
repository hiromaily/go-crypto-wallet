package create

import (
	"context"
	"fmt"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/di"
)

// runDB creates payment_request table with dummy data
//
// Deprecated: Use query with shell script instead of go code.
func runDB(container di.Container, tableName string) error {
	fmt.Println("-table: " + tableName)

	// validator
	if tableName == "" {
		tableName = "payment_request"
	}
	if tableName == "payment_request" {
		// Get use case from container
		useCase := container.NewWatchCreatePaymentRequestUseCase()

		// create payment_request table
		// Note: AmountList can be empty for dummy data generation
		if err := useCase.Execute(context.Background(), watchusecase.CreatePaymentRequestInput{
			AmountList: []float64{}, // Empty for default dummy data
		}); err != nil {
			return fmt.Errorf("fail to create payment request: %w", err)
		}
	}

	return nil
}
