//go:build integration
// +build integration

package watchrepo_test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/quagmt/udecimal"
	"github.com/stretchr/testify/require"

	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestPaymentRequestSqlc is integration test for PaymentRequestRepositorySqlc
func TestPaymentRequestSqlc(t *testing.T) {
	paymentRepo := testutil.NewPaymentRequestRepositorySqlc()

	// Delete all records
	_, err := paymentRepo.DeleteAll()
	require.NoError(t, err, "fail to call DeleteAll()")

	// Create test payment requests
	amount1, _ := udecimal.Parse("1.5")
	amount2, _ := udecimal.Parse("2.5")

	requests := []*models.PaymentRequest{
		{
			Coin:            "btc",
			SenderAddress:   "sender-sqlc-1",
			ReceiverAddress: "receiver-sqlc-1",
			Amount:          amount1,
			IsDone:          false,
		},
		{
			Coin:            "btc",
			SenderAddress:   "sender-sqlc-2",
			ReceiverAddress: "receiver-sqlc-2",
			Amount:          amount2,
			IsDone:          false,
		},
	}

	// Insert bulk
	err = paymentRepo.InsertBulk(requests)
	require.NoError(t, err, "fail to call InsertBulk()")

	// Get all
	allRequests, err := paymentRepo.GetAll()
	require.NoError(t, err, "fail to call GetAll()")
	require.GreaterOrEqual(t, len(allRequests), 2, "GetAll() should return at least 2 requests")

	// Update payment ID
	paymentID := int64(12345)
	ids := []int64{allRequests[0].ID, allRequests[1].ID}
	rowsAffected, err := paymentRepo.UpdatePaymentID(paymentID, ids)
	require.NoError(t, err, "fail to call UpdatePaymentID()")
	require.Equal(t, int64(2), rowsAffected, "UpdatePaymentID() should affect 2 rows")

	// Get all by payment ID
	requestsByPaymentID, err := paymentRepo.GetAllByPaymentID(paymentID)
	require.NoError(t, err, "fail to call GetAllByPaymentID()")
	require.Len(t, requestsByPaymentID, 2, "GetAllByPaymentID() should return 2 requests")

	// Update is_done
	rowsAffected, err = paymentRepo.UpdateIsDone(paymentID)
	require.NoError(t, err, "fail to call UpdateIsDone()")
	require.Equal(t, int64(2), rowsAffected, "UpdateIsDone() should affect 2 rows")

	// Verify is_done is true
	verifyRequests, err := paymentRepo.GetAllByPaymentID(paymentID)
	require.NoError(t, err, "fail to call GetAllByPaymentID() after UpdateIsDone()")
	for _, req := range verifyRequests {
		require.True(t, req.IsDone, "UpdateIsDone() should set is_done to true for request ID %d", req.ID)
	}
}
