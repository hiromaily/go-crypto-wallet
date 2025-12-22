//go:build integration
// +build integration

package watchrepo_test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestAddressSqlc is integration test for AddressRepositorySqlc
func TestAddressSqlc(t *testing.T) {
	addressRepo := testutil.NewAddressRepositorySqlc()
	accountType := account.AccountTypeClient

	// Clean up any existing test data (address table has unique key on wallet_address)
	// Note: We can't access db directly from testutil, so we skip cleanup here
	// Tests should use unique addresses or be run with clean database

	// Insert bulk addresses
	addresses := []*models.Address{
		{
			Coin:          "btc",
			Account:       accountType.String(),
			WalletAddress: "address-sqlc-1",
			IsAllocated:   false,
		},
		{
			Coin:          "btc",
			Account:       accountType.String(),
			WalletAddress: "address-sqlc-2",
			IsAllocated:   false,
		},
		{
			Coin:          "btc",
			Account:       accountType.String(),
			WalletAddress: "address-sqlc-3",
			IsAllocated:   true,
		},
	}

	if err := addressRepo.InsertBulk(addresses); err != nil {
		t.Fatalf("fail to call InsertBulk() %v", err)
	}

	// Get all addresses
	allAddrs, err := addressRepo.GetAll(accountType)
	if err != nil {
		t.Fatalf("fail to call GetAll() %v", err)
	}
	if len(allAddrs) < 3 {
		t.Errorf("GetAll() returned %d addresses, want at least 3", len(allAddrs))
		return
	}

	// Get all address strings
	addrStrings, err := addressRepo.GetAllAddress(accountType)
	if err != nil {
		t.Fatalf("fail to call GetAllAddress() %v", err)
	}
	if len(addrStrings) < 3 {
		t.Errorf("GetAllAddress() returned %d addresses, want at least 3", len(addrStrings))
		return
	}

	// Get one unallocated address
	unallocAddr, err := addressRepo.GetOneUnAllocated(accountType)
	if err != nil {
		t.Fatalf("fail to call GetOneUnAllocated() %v", err)
	}
	if unallocAddr == nil {
		t.Fatal("GetOneUnAllocated() returned nil")
	}
	if unallocAddr.IsAllocated {
		t.Errorf("GetOneUnAllocated() returned allocated address")
		return
	}

	// Update is_allocated
	rowsAffected, err := addressRepo.UpdateIsAllocated(true, unallocAddr.WalletAddress)
	if err != nil {
		t.Fatalf("fail to call UpdateIsAllocated() %v", err)
	}
	if rowsAffected != 1 {
		t.Errorf("UpdateIsAllocated() affected %d rows, want 1", rowsAffected)
		return
	}

	// Verify the address is now allocated
	verifyAddr, err := addressRepo.GetOneUnAllocated(accountType)
	if err != nil {
		t.Fatalf("fail to call GetOneUnAllocated() after update %v", err)
	}
	if verifyAddr != nil && verifyAddr.WalletAddress == unallocAddr.WalletAddress {
		t.Errorf("GetOneUnAllocated() returned the same address after allocation")
		return
	}
}
