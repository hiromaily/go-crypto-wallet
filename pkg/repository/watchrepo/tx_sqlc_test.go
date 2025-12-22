//go:build integration
// +build integration

package watchrepo_test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-crypto-wallet/pkg/action"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestTxSqlc is integration test for TxRepositorySqlc
func TestTxSqlc(t *testing.T) {
	txRepo := testutil.NewTxRepositorySqlc()

	// Delete all records
	if _, err := txRepo.DeleteAll(); err != nil {
		t.Fatalf("fail to call DeleteAll() %v", err)
	}

	// Insert unsigned tx
	actionType := action.ActionTypePayment
	id, err := txRepo.InsertUnsignedTx(actionType)
	if err != nil {
		t.Fatalf("fail to call InsertUnsignedTx() %v", err)
	}
	if id == 0 {
		t.Errorf("InsertUnsignedTx() returned id = 0, want non-zero")
		return
	}

	// Get one
	tx, err := txRepo.GetOne(id)
	if err != nil {
		t.Fatalf("fail to call GetOne() %v", err)
	}
	if tx.ID != id {
		t.Errorf("GetOne() returned id = %d, want %d", tx.ID, id)
		return
	}
	if tx.Action != actionType.String() {
		t.Errorf("GetOne() returned action = %s, want %s", tx.Action, actionType.String())
		return
	}

	// Get max ID
	maxID, err := txRepo.GetMaxID(actionType)
	if err != nil {
		t.Fatalf("fail to call GetMaxID() %v", err)
	}
	if maxID != id {
		t.Errorf("GetMaxID() = %d, want %d", maxID, id)
		return
	}

	// Insert another tx to test max ID
	id2, err := txRepo.InsertUnsignedTx(actionType)
	if err != nil {
		t.Fatalf("fail to call InsertUnsignedTx() second time %v", err)
	}

	// Get max ID again
	maxID2, err := txRepo.GetMaxID(actionType)
	if err != nil {
		t.Fatalf("fail to call GetMaxID() second time %v", err)
	}
	if maxID2 != id2 {
		t.Errorf("GetMaxID() = %d, want %d", maxID2, id2)
		return
	}
	if id2 <= id {
		t.Errorf("second InsertUnsignedTx() returned id = %d, want > %d", id2, id)
		return
	}

	// Update tx
	tx.Action = action.ActionTypeDeposit.String()
	rowsAffected, err := txRepo.Update(tx)
	if err != nil {
		t.Fatalf("fail to call Update() %v", err)
	}
	if rowsAffected != 1 {
		t.Errorf("Update() affected %d rows, want 1", rowsAffected)
		return
	}

	// Verify update
	updatedTx, err := txRepo.GetOne(id)
	if err != nil {
		t.Fatalf("fail to call GetOne() after update %v", err)
	}
	if updatedTx.Action != action.ActionTypeDeposit.String() {
		t.Errorf("Update() did not change action, got %s, want %s", updatedTx.Action, action.ActionTypeDeposit.String())
		return
	}
}
