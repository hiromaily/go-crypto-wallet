//go:build integration
// +build integration

package watchrepo_test

import (
	"testing"

	"github.com/ericlagergren/decimal"
	_ "github.com/go-sql-driver/mysql"

	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestPaymentRequestSqlc is integration test for PaymentRequestRepositorySqlc
func TestPaymentRequestSqlc(t *testing.T) {
	paymentRepo := testutil.NewPaymentRequestRepositorySqlc()

	// Delete all records
	if _, err := paymentRepo.DeleteAll(); err != nil {
		t.Fatalf("fail to call DeleteAll() %v", err)
	}

	// Create test payment requests
	amount1 := new(decimal.Big)
	amount1, _ = amount1.SetString("1.5")
	amount2 := new(decimal.Big)
	amount2, _ = amount2.SetString("2.5")

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
	if err := paymentRepo.InsertBulk(requests); err != nil {
		t.Fatalf("fail to call InsertBulk() %v", err)
	}

	// Get all
	allRequests, err := paymentRepo.GetAll()
	if err != nil {
		t.Fatalf("fail to call GetAll() %v", err)
	}
	if len(allRequests) < 2 {
		t.Errorf("GetAll() returned %d requests, want at least 2", len(allRequests))
		return
	}

	// Update payment ID
	paymentID := int64(12345)
	ids := []int64{allRequests[0].ID, allRequests[1].ID}
	rowsAffected, err := paymentRepo.UpdatePaymentID(paymentID, ids)
	if err != nil {
		t.Fatalf("fail to call UpdatePaymentID() %v", err)
	}
	if rowsAffected != 2 {
		t.Errorf("UpdatePaymentID() affected %d rows, want 2", rowsAffected)
		return
	}

	// Get all by payment ID
	requestsByPaymentID, err := paymentRepo.GetAllByPaymentID(paymentID)
	if err != nil {
		t.Fatalf("fail to call GetAllByPaymentID() %v", err)
	}
	if len(requestsByPaymentID) != 2 {
		t.Errorf("GetAllByPaymentID() returned %d requests, want 2", len(requestsByPaymentID))
		return
	}

	// Update is_done
	rowsAffected, err = paymentRepo.UpdateIsDone(paymentID)
	if err != nil {
		t.Fatalf("fail to call UpdateIsDone() %v", err)
	}
	if rowsAffected != 2 {
		t.Errorf("UpdateIsDone() affected %d rows, want 2", rowsAffected)
		return
	}

	// Verify is_done is true
	verifyRequests, err := paymentRepo.GetAllByPaymentID(paymentID)
	if err != nil {
		t.Fatalf("fail to call GetAllByPaymentID() after UpdateIsDone() %v", err)
	}
	for _, req := range verifyRequests {
		if !req.IsDone {
			t.Errorf("UpdateIsDone() did not set is_done to true for request ID %d", req.ID)
			return
		}
	}
}
