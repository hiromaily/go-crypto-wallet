package coldrepo_test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

//NewAccountKeyRepository

// TestTx is test for any data operation
func TestAccountKey(t *testing.T) {
	//boil.DebugMode = true
	akRepo := testutil.NewAccountKeyRepository()

	_, err := akRepo.UpdateAddrStatus(
		account.AccountTypeClient,
		address.AddrStatusAddressExported,
		[]string{"cND89yJJxG6KHxZrR7ZrwqQ3yrFSUjDRrGnHKiHFxNjJYmqUQRBu"})
	if err != nil {
		t.Fatalf("fail to call UpdateAddrStatus() %v", err)
	}
}
