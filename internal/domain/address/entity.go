package address

import (
	"errors"
	"time"

	"github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	"github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// Address represents a cryptocurrency wallet address in the domain layer.
// This entity follows Clean Architecture principles and is independent of any infrastructure concerns.
//
// An address is associated with a specific cryptocurrency (coin) and account type,
// and tracks whether it has been allocated for use (e.g., for payment requests).
type Address struct {
	// ID is the unique identifier for this address record
	ID int64

	// CoinTypeCode identifies the cryptocurrency this address belongs to (BTC, BCH, ETH, XRP)
	CoinTypeCode coin.CoinTypeCode

	// AccountType identifies the purpose/type of this address (client, deposit, payment, etc.)
	AccountType account.AccountType

	// WalletAddress is the actual blockchain address string (e.g., "1A1zP1eP5QGeFi...")
	WalletAddress string

	// IsAllocated indicates whether this address has been assigned for use
	// (e.g., allocated to a payment request)
	IsAllocated bool

	// UpdatedAt is the timestamp of the last update to this address record
	// Nil if never updated after creation
	UpdatedAt *time.Time
}

// NewAddress creates a new Address entity with validation.
func NewAddress(
	coinTypeCode coin.CoinTypeCode,
	accountType account.AccountType,
	walletAddress string,
	isAllocated bool,
) (*Address, error) {
	if walletAddress == "" {
		return nil, errors.New("wallet address cannot be empty")
	}

	return &Address{
		CoinTypeCode:  coinTypeCode,
		AccountType:   accountType,
		WalletAddress: walletAddress,
		IsAllocated:   isAllocated,
	}, nil
}

// Allocate marks this address as allocated.
func (a *Address) Allocate() {
	a.IsAllocated = true
	now := time.Now()
	a.UpdatedAt = &now
}

// Deallocate marks this address as not allocated.
func (a *Address) Deallocate() {
	a.IsAllocated = false
	now := time.Now()
	a.UpdatedAt = &now
}
