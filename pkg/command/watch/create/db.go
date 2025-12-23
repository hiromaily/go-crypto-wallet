package create

import (
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// runDB creates payment_request table with dummy data
//
// Deprecated: Use query with shell script instead of go code.
func runDB(wallet wallets.Watcher, tableName string) error {
	fmt.Println("-table: " + tableName)

	// validator
	if tableName == "" {
		tableName = "payment_request"
	}
	if tableName == "payment_request" {
		// create payment_request table
		if err := wallet.CreatePaymentRequest(); err != nil {
			return fmt.Errorf("fail to call CreatePaymentRequest() %w", err)
		}
	}

	return nil
}
