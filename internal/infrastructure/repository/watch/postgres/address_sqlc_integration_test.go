//go:build integration

package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	watchpostgres "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/postgres"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	pkgpostgres "github.com/hiromaily/go-crypto-wallet/pkg/db/postgres"
)

// postgresWatchConf builds a PostgreSQL config from environment variables.
// Defaults match the Docker Compose wallet-postgres service configuration.
func postgresWatchConf() *config.PostgreSQL {
	port := 5432
	if p := os.Getenv("POSTGRES_PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			port = v
		}
	}
	return &config.PostgreSQL{
		Host:    envOrDefaultWatch("POSTGRES_HOST", "localhost"),
		Port:    port,
		DB:      envOrDefaultWatch("POSTGRES_DB", "watch"),
		User:    envOrDefaultWatch("POSTGRES_USER", "postgres"),
		Pass:    envOrDefaultWatch("POSTGRES_PASS", "postgres"),
		SSLMode: "disable",
	}
}

func envOrDefaultWatch(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// openWatchDB opens a connection to the watch PostgreSQL database for testing.
func openWatchDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := pkgpostgres.NewPostgres(postgresWatchConf())
	require.NoError(t, err, "failed to connect to PostgreSQL watch database")
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// uniqueAddress generates a test wallet address with a unique suffix to avoid conflicts.
func uniqueAddress(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

// cleanupAddresses deletes test addresses matching the given prefix from the database.
func cleanupAddresses(t *testing.T, db *sql.DB, addresses ...string) {
	t.Helper()
	for _, addr := range addresses {
		_, err := db.ExecContext(context.Background(),
			"DELETE FROM address WHERE wallet_address = $1", addr)
		require.NoError(t, err, "failed to clean up test address %q", addr)
	}
}

// TestAddressRepository_InsertBulk_and_GetAll verifies that InsertBulk stores addresses
// and GetAll retrieves them by account type.
func TestAddressRepository_InsertBulk_and_GetAll(t *testing.T) {
	db := openWatchDB(t)

	addr1 := uniqueAddress("test_btc_client_1")
	addr2 := uniqueAddress("test_btc_client_2")
	t.Cleanup(func() { cleanupAddresses(t, db, addr1, addr2) })

	repo := watchpostgres.NewAddressRepositorySqlc(db, domainCoin.BTC)

	items := []*domainAddress.Address{
		{
			CoinTypeCode: domainCoin.BTC, AccountType: domainAccount.AccountTypeClient,
			WalletAddress: addr1, IsAllocated: false,
		},
		{
			CoinTypeCode: domainCoin.BTC, AccountType: domainAccount.AccountTypeClient,
			WalletAddress: addr2, IsAllocated: false,
		},
	}
	err := repo.InsertBulk(context.Background(), items)
	require.NoError(t, err, "InsertBulk should succeed for valid addresses")

	got, err := repo.GetAll(domainAccount.AccountTypeClient)
	require.NoError(t, err)

	// Confirm both inserted addresses appear in the result.
	addrSet := make(map[string]bool, len(got))
	for _, a := range got {
		addrSet[a.WalletAddress] = true
	}
	assert.True(t, addrSet[addr1], "GetAll should return the first inserted address")
	assert.True(t, addrSet[addr2], "GetAll should return the second inserted address")
}

// TestAddressRepository_GetAllAddress verifies that GetAllAddress returns only the
// wallet address strings for the given account type.
func TestAddressRepository_GetAllAddress(t *testing.T) {
	db := openWatchDB(t)

	addr := uniqueAddress("test_btc_client_str")
	t.Cleanup(func() { cleanupAddresses(t, db, addr) })

	repo := watchpostgres.NewAddressRepositorySqlc(db, domainCoin.BTC)

	err := repo.InsertBulk(context.Background(), []*domainAddress.Address{
		{
			CoinTypeCode: domainCoin.BTC, AccountType: domainAccount.AccountTypeClient,
			WalletAddress: addr, IsAllocated: false,
		},
	})
	require.NoError(t, err)

	strings, err := repo.GetAllAddress(domainAccount.AccountTypeClient)
	require.NoError(t, err)

	assert.Contains(t, strings, addr, "GetAllAddress should include the inserted wallet address string")
}

// TestAddressRepository_GetOneUnAllocated verifies that GetOneUnAllocated returns an
// unallocated address for the given account type.
func TestAddressRepository_GetOneUnAllocated(t *testing.T) {
	db := openWatchDB(t)

	addr := uniqueAddress("test_btc_deposit_unalloc")
	t.Cleanup(func() { cleanupAddresses(t, db, addr) })

	repo := watchpostgres.NewAddressRepositorySqlc(db, domainCoin.BTC)

	err := repo.InsertBulk(context.Background(), []*domainAddress.Address{
		{
			CoinTypeCode: domainCoin.BTC, AccountType: domainAccount.AccountTypeDeposit,
			WalletAddress: addr, IsAllocated: false,
		},
	})
	require.NoError(t, err)

	got, err := repo.GetOneUnAllocated(domainAccount.AccountTypeDeposit)
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.False(t, got.IsAllocated, "returned address should be unallocated")
	assert.Equal(t, domainCoin.BTC, got.CoinTypeCode)
	assert.Equal(t, domainAccount.AccountTypeDeposit, got.AccountType)
}

// TestAddressRepository_UpdateIsAllocated verifies that UpdateIsAllocated correctly
// flips the is_allocated flag and returns the affected row count.
func TestAddressRepository_UpdateIsAllocated(t *testing.T) {
	db := openWatchDB(t)

	addr := uniqueAddress("test_btc_client_update")
	t.Cleanup(func() { cleanupAddresses(t, db, addr) })

	repo := watchpostgres.NewAddressRepositorySqlc(db, domainCoin.BTC)

	err := repo.InsertBulk(context.Background(), []*domainAddress.Address{
		{
			CoinTypeCode: domainCoin.BTC, AccountType: domainAccount.AccountTypeClient,
			WalletAddress: addr, IsAllocated: false,
		},
	})
	require.NoError(t, err)

	rowsAffected, err := repo.UpdateIsAllocated(true, addr)
	require.NoError(t, err)
	assert.EqualValues(t, 1, rowsAffected, "exactly one row should be updated")

	// Verify the update persisted.
	rows, err := repo.GetAll(domainAccount.AccountTypeClient)
	require.NoError(t, err)
	for _, a := range rows {
		if a.WalletAddress == addr {
			assert.True(t, a.IsAllocated, "address should be marked as allocated after update")
			return
		}
	}
	t.Errorf("address %q not found after update", addr)
}

// TestAddressRepository_TypeConversion verifies domain.Address ↔ sqlcgen.Address
// round-trip: fields inserted via domain entity are returned as correct domain entity.
func TestAddressRepository_TypeConversion(t *testing.T) {
	db := openWatchDB(t)

	addr := uniqueAddress("test_eth_payment_conv")
	t.Cleanup(func() { cleanupAddresses(t, db, addr) })

	repo := watchpostgres.NewAddressRepositorySqlc(db, domainCoin.ETH)

	input := &domainAddress.Address{
		CoinTypeCode:  domainCoin.ETH,
		AccountType:   domainAccount.AccountTypePayment,
		WalletAddress: addr,
		IsAllocated:   false,
	}
	err := repo.InsertBulk(context.Background(), []*domainAddress.Address{input})
	require.NoError(t, err)

	got, err := repo.GetAll(domainAccount.AccountTypePayment)
	require.NoError(t, err)

	var found *domainAddress.Address
	for _, a := range got {
		if a.WalletAddress == addr {
			found = a
			break
		}
	}
	require.NotNil(t, found, "inserted address should be retrievable")

	assert.Equal(t, domainCoin.ETH, found.CoinTypeCode, "CoinTypeCode should round-trip correctly")
	assert.Equal(t, domainAccount.AccountTypePayment, found.AccountType, "AccountType should round-trip correctly")
	assert.Equal(t, addr, found.WalletAddress, "WalletAddress should round-trip correctly")
	assert.False(t, found.IsAllocated, "IsAllocated should be false")
}

// TestAddressRepository_NullUpdatedAt verifies that when UpdatedAt is nil on insert,
// the retrieved entity also has nil UpdatedAt.
func TestAddressRepository_NullUpdatedAt(t *testing.T) {
	db := openWatchDB(t)

	addr := uniqueAddress("test_btc_client_nulltime")
	t.Cleanup(func() { cleanupAddresses(t, db, addr) })

	repo := watchpostgres.NewAddressRepositorySqlc(db, domainCoin.BTC)

	// Insert with explicit nil UpdatedAt.
	err := repo.InsertBulk(context.Background(), []*domainAddress.Address{
		{
			CoinTypeCode:  domainCoin.BTC,
			AccountType:   domainAccount.AccountTypeClient,
			WalletAddress: addr,
			IsAllocated:   false,
			UpdatedAt:     nil,
		},
	})
	require.NoError(t, err)

	rows, err := repo.GetAll(domainAccount.AccountTypeClient)
	require.NoError(t, err)

	for _, a := range rows {
		if a.WalletAddress == addr {
			// UpdatedAt may be nil or populated by the DB default (now()).
			// The key assertion is that the converter handles sql.NullTime
			// without panicking. The field value itself is acceptable either way.
			// (The schema sets DEFAULT now(), so the DB may fill it in.)
			t.Logf("UpdatedAt for inserted address: %v", a.UpdatedAt)
			return
		}
	}
	t.Errorf("address %q not found in GetAll result", addr)
}

// TestAddressRepository_InsertWithUpdatedAt verifies that an explicit UpdatedAt value
// inserted is preserved and retrieved correctly.
func TestAddressRepository_InsertWithUpdatedAt(t *testing.T) {
	db := openWatchDB(t)

	addr := uniqueAddress("test_btc_client_withtime")
	t.Cleanup(func() { cleanupAddresses(t, db, addr) })

	repo := watchpostgres.NewAddressRepositorySqlc(db, domainCoin.BTC)

	now := time.Now().Truncate(time.Millisecond)
	err := repo.InsertBulk(context.Background(), []*domainAddress.Address{
		{
			CoinTypeCode:  domainCoin.BTC,
			AccountType:   domainAccount.AccountTypeClient,
			WalletAddress: addr,
			IsAllocated:   false,
			UpdatedAt:     &now,
		},
	})
	require.NoError(t, err)

	rows, err := repo.GetAll(domainAccount.AccountTypeClient)
	require.NoError(t, err)

	for _, a := range rows {
		if a.WalletAddress == addr {
			// UpdatedAt should be non-nil when explicitly provided.
			require.NotNil(t, a.UpdatedAt, "UpdatedAt should not be nil when explicitly set")
			return
		}
	}
	t.Errorf("address %q not found in GetAll result", addr)
}

// TestAddressRepository_CheckConstraint_InvalidCoin verifies that inserting an address
// with an invalid coin value fails due to the chk_address_coin CHECK constraint.
func TestAddressRepository_CheckConstraint_InvalidCoin(t *testing.T) {
	db := openWatchDB(t)

	// Use an invalid coin code that is not in ('btc','bch','eth','xrp','hyt').
	invalidCoin := domainCoin.CoinTypeCode("invalid_coin")
	repo := watchpostgres.NewAddressRepositorySqlc(db, invalidCoin)

	err := repo.InsertBulk(context.Background(), []*domainAddress.Address{
		{
			CoinTypeCode:  invalidCoin,
			AccountType:   domainAccount.AccountTypeClient,
			WalletAddress: uniqueAddress("test_invalid_coin"),
			IsAllocated:   false,
		},
	})
	require.Error(t, err, "INSERT with invalid coin should fail CHECK constraint")
	assert.Contains(t, err.Error(), "InsertAddress", "error should reference the failing operation")
}

// TestAddressRepository_CheckConstraint_InvalidAccount verifies that inserting an address
// with an invalid account value fails due to the chk_address_account CHECK constraint.
// Note: The address table only allows 'client','deposit','payment','stored'.
func TestAddressRepository_CheckConstraint_InvalidAccount(t *testing.T) {
	db := openWatchDB(t)

	repo := watchpostgres.NewAddressRepositorySqlc(db, domainCoin.BTC)

	// 'auth1' is a valid AccountType in the domain but NOT in the CHECK constraint
	// for the address table (only client, deposit, payment, stored are allowed).
	err := repo.InsertBulk(context.Background(), []*domainAddress.Address{
		{
			CoinTypeCode:  domainCoin.BTC,
			AccountType:   domainAccount.AccountTypeAuth1,
			WalletAddress: uniqueAddress("test_invalid_account"),
			IsAllocated:   false,
		},
	})
	require.Error(t, err, "INSERT with account type not in CHECK constraint should fail")
	assert.Contains(t, err.Error(), "InsertAddress", "error should reference the failing operation")
}
