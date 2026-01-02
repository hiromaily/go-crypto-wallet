package transaction

import (
	"time"

	"github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// Transaction represents a base cryptocurrency transaction in the domain layer.
// This entity follows Clean Architecture principles and is independent of any infrastructure concerns.
//
// Transaction is a high-level record that tracks the lifecycle of a transaction
// (unsigned, signed, sent, done, notified, canceled). Coin-specific transaction details
// are stored in separate entities (BTCTransaction, ETHDetailTx, XRPDetailTx).
type Transaction struct {
	// ID is the unique identifier for this transaction record
	ID int64

	// CoinTypeCode identifies the cryptocurrency for this transaction (BTC, BCH, ETH, XRP)
	CoinTypeCode coin.CoinTypeCode

	// ActionType identifies the purpose of this transaction (deposit, payment, transfer)
	ActionType ActionType

	// UpdatedAt is the timestamp of the last update to this transaction record
	// Nil if never updated after creation
	UpdatedAt *time.Time
}

// NewTransaction creates a new Transaction entity.
func NewTransaction(
	coinTypeCode coin.CoinTypeCode,
	actionType ActionType,
) *Transaction {
	return &Transaction{
		CoinTypeCode: coinTypeCode,
		ActionType:   actionType,
	}
}

// Touch updates the UpdatedAt timestamp to the current time.
func (t *Transaction) Touch() {
	now := time.Now()
	t.UpdatedAt = &now
}
