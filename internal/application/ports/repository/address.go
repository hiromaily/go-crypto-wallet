package repository

import (
	"context"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
)

// AddressRepositorier is AddressRepository interface
type AddressRepositorier interface {
	GetAll(accountType domainAccount.AccountType) ([]*domainAddress.Address, error)
	GetAllAddress(accountType domainAccount.AccountType) ([]string, error)
	GetOneUnAllocated(accountType domainAccount.AccountType) (*domainAddress.Address, error)
	InsertBulk(ctx context.Context, items []*domainAddress.Address) error
	UpdateIsAllocated(isAllocated bool, Address string) (int64, error)
}
