//go:build integration
// +build integration

package watchrepo_test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/config/account"
	models "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestAddressSqlc is integration test for AddressRepositorySqlc
func TestAddressSqlc(t *testing.T) {
	// Get db connection for cleanup
	db := testutil.GetDB()
	// Clean up any existing test data
	_, _ = db.Exec("DELETE FROM address WHERE wallet_address LIKE 'address-sqlc-%'")

	addressRepo := testutil.NewAddressRepositorySqlc()
	accountType := account.AccountTypeClient

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

	err := addressRepo.InsertBulk(addresses)
	require.NoError(t, err, "fail to call InsertBulk()")

	// Get all addresses
	allAddrs, err := addressRepo.GetAll(accountType)
	require.NoError(t, err, "fail to call GetAll()")
	require.GreaterOrEqual(t, len(allAddrs), 3, "GetAll() should return at least 3 addresses")

	// Get all address strings
	addrStrings, err := addressRepo.GetAllAddress(accountType)
	require.NoError(t, err, "fail to call GetAllAddress()")
	require.GreaterOrEqual(t, len(addrStrings), 3, "GetAllAddress() should return at least 3 addresses")

	// Get one unallocated address
	unallocAddr, err := addressRepo.GetOneUnAllocated(accountType)
	require.NoError(t, err, "fail to call GetOneUnAllocated()")
	require.NotNil(t, unallocAddr, "GetOneUnAllocated() returned nil")
	require.False(t, unallocAddr.IsAllocated, "GetOneUnAllocated() returned allocated address")

	// Update is_allocated
	rowsAffected, err := addressRepo.UpdateIsAllocated(true, unallocAddr.WalletAddress)
	require.NoError(t, err, "fail to call UpdateIsAllocated()")
	require.Equal(t, int64(1), rowsAffected, "UpdateIsAllocated() should affect 1 row")

	// Verify the address is now allocated
	verifyAddr, err := addressRepo.GetOneUnAllocated(accountType)
	require.NoError(t, err, "fail to call GetOneUnAllocated() after update")
	if verifyAddr != nil {
		assert.NotEqual(t, unallocAddr.WalletAddress, verifyAddr.WalletAddress,
			"GetOneUnAllocated() returned the same address after allocation")
	}
}
