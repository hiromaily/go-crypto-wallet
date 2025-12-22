//go:build integration
// +build integration

package watchrepo_test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"

	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
)

// TestXrpDetailTxSqlc is integration test for XrpDetailTxInputRepositorySqlc
func TestXrpDetailTxSqlc(t *testing.T) {
	xrpDetailTxRepo := testutil.NewXrpDetailTxRepositorySqlc()

	// Create test xrp detail tx
	uuid := "xrp-uuid-sqlc-test"
	xrpTx := &models.XRPDetailTX{
		TXID:                  1,
		UUID:                  uuid,
		CurrentTXType:         tx.TxTypeUnsigned.Int8(),
		SenderAccount:         "deposit",
		SenderAddress:         "rSender-sqlc",
		ReceiverAccount:       "client",
		ReceiverAddress:       "rReceiver-sqlc",
		Amount:                "1000000",
		XRPTXType:             "Payment",
		Fee:                   "12",
		Flags:                 0,
		LastLedgerSequence:    12345,
		Sequence:              1,
		SigningPubkey:         "pubkey-sqlc",
		TXNSignature:          "",
		Hash:                  "",
		EarliestLedgerVersion: 0,
		SignedTXID:            "",
		TXBlob:                "",
	}

	// Insert
	if err := xrpDetailTxRepo.Insert(xrpTx); err != nil {
		t.Fatalf("fail to call Insert() %v", err)
	}

	// Get all by tx ID
	xrpTxs, err := xrpDetailTxRepo.GetAllByTxID(1)
	if err != nil {
		t.Fatalf("fail to call GetAllByTxID() %v", err)
	}
	if len(xrpTxs) < 1 {
		t.Errorf("GetAllByTxID() returned %d records, want at least 1", len(xrpTxs))
		return
	}

	// Get one
	retrievedTx, err := xrpDetailTxRepo.GetOne(xrpTxs[0].ID)
	if err != nil {
		t.Fatalf("fail to call GetOne() %v", err)
	}
	if retrievedTx.UUID != uuid {
		t.Errorf("GetOne() returned UUID = %s, want %s", retrievedTx.UUID, uuid)
		return
	}

	// Update after tx sent
	signedTxID := "signed-txid-sqlc"
	txBlob := "tx-blob-sqlc"
	earliestLedgerVersion := uint64(12340)
	rowsAffected, err := xrpDetailTxRepo.UpdateAfterTxSent(uuid, tx.TxTypeSent, signedTxID, txBlob, earliestLedgerVersion)
	if err != nil {
		t.Fatalf("fail to call UpdateAfterTxSent() %v", err)
	}
	if rowsAffected < 1 {
		t.Errorf("UpdateAfterTxSent() affected %d rows, want at least 1", rowsAffected)
		return
	}

	// Verify update
	updatedTx, err := xrpDetailTxRepo.GetOne(retrievedTx.ID)
	if err != nil {
		t.Fatalf("fail to call GetOne() after update %v", err)
	}
	if updatedTx.SignedTXID != signedTxID {
		t.Errorf("UpdateAfterTxSent() did not update SignedTXID, got %s, want %s", updatedTx.SignedTXID, signedTxID)
		return
	}
	if updatedTx.TXBlob != txBlob {
		t.Errorf("UpdateAfterTxSent() did not update TXBlob, got %s, want %s", updatedTx.TXBlob, txBlob)
		return
	}
	if updatedTx.CurrentTXType != tx.TxTypeSent.Int8() {
		t.Errorf("UpdateAfterTxSent() did not update CurrentTXType, got %d, want %d", updatedTx.CurrentTXType, tx.TxTypeSent.Int8())
		return
	}
	if updatedTx.EarliestLedgerVersion != earliestLedgerVersion {
		t.Errorf("UpdateAfterTxSent() did not update EarliestLedgerVersion, got %d, want %d", updatedTx.EarliestLedgerVersion, earliestLedgerVersion)
		return
	}

	// Get sent hash tx (for XRP, this is tx_blob)
	blobs, err := xrpDetailTxRepo.GetSentHashTx(tx.TxTypeSent)
	if err != nil {
		t.Fatalf("fail to call GetSentHashTx() %v", err)
	}
	if len(blobs) < 1 {
		t.Errorf("GetSentHashTx() returned %d blobs, want at least 1", len(blobs))
		return
	}

	// Update tx type by sent hash tx (tx_blob)
	rowsAffected, err = xrpDetailTxRepo.UpdateTxTypeBySentHashTx(tx.TxTypeDone, txBlob)
	if err != nil {
		t.Fatalf("fail to call UpdateTxTypeBySentHashTx() %v", err)
	}
	if rowsAffected < 1 {
		t.Errorf("UpdateTxTypeBySentHashTx() affected %d rows, want at least 1", rowsAffected)
		return
	}

	// Verify tx type update
	verifyTx, err := xrpDetailTxRepo.GetOne(retrievedTx.ID)
	if err != nil {
		t.Fatalf("fail to call GetOne() after UpdateTxTypeBySentHashTx() %v", err)
	}
	if verifyTx.CurrentTXType != tx.TxTypeDone.Int8() {
		t.Errorf("UpdateTxTypeBySentHashTx() did not update CurrentTXType, got %d, want %d", verifyTx.CurrentTXType, tx.TxTypeDone.Int8())
		return
	}

	// Update tx type by ID
	rowsAffected, err = xrpDetailTxRepo.UpdateTxType(retrievedTx.ID, tx.TxTypeNotified)
	if err != nil {
		t.Fatalf("fail to call UpdateTxType() %v", err)
	}
	if rowsAffected != 1 {
		t.Errorf("UpdateTxType() affected %d rows, want 1", rowsAffected)
		return
	}

	// Verify final tx type
	finalTx, err := xrpDetailTxRepo.GetOne(retrievedTx.ID)
	if err != nil {
		t.Fatalf("fail to call GetOne() after UpdateTxType() %v", err)
	}
	if finalTx.CurrentTXType != tx.TxTypeNotified.Int8() {
		t.Errorf("UpdateTxType() did not update CurrentTXType, got %d, want %d", finalTx.CurrentTXType, tx.TxTypeNotified.Int8())
		return
	}

	// Test InsertBulk
	bulkTxs := []*models.XRPDetailTX{
		{
			TXID:                  2,
			UUID:                  "xrp-uuid-bulk-1",
			CurrentTXType:         tx.TxTypeUnsigned.Int8(),
			SenderAccount:         "deposit",
			SenderAddress:         "rSender-bulk-1",
			ReceiverAccount:       "client",
			ReceiverAddress:       "rReceiver-bulk-1",
			Amount:                "2000000",
			XRPTXType:             "Payment",
			Fee:                   "12",
			Flags:                 0,
			LastLedgerSequence:    12346,
			Sequence:              2,
			SigningPubkey:         "pubkey-bulk-1",
			TXNSignature:          "",
			Hash:                  "",
			EarliestLedgerVersion: 0,
			SignedTXID:            "",
			TXBlob:                "",
		},
		{
			TXID:                  2,
			UUID:                  "xrp-uuid-bulk-2",
			CurrentTXType:         tx.TxTypeUnsigned.Int8(),
			SenderAccount:         "deposit",
			SenderAddress:         "rSender-bulk-2",
			ReceiverAccount:       "client",
			ReceiverAddress:       "rReceiver-bulk-2",
			Amount:                "3000000",
			XRPTXType:             "Payment",
			Fee:                   "12",
			Flags:                 0,
			LastLedgerSequence:    12347,
			Sequence:              3,
			SigningPubkey:         "pubkey-bulk-2",
			TXNSignature:          "",
			Hash:                  "",
			EarliestLedgerVersion: 0,
			SignedTXID:            "",
			TXBlob:                "",
		},
	}

	if err := xrpDetailTxRepo.InsertBulk(bulkTxs); err != nil {
		t.Fatalf("fail to call InsertBulk() %v", err)
	}

	// Verify bulk insert
	bulkRetrieved, err := xrpDetailTxRepo.GetAllByTxID(2)
	if err != nil {
		t.Fatalf("fail to call GetAllByTxID() after InsertBulk() %v", err)
	}
	if len(bulkRetrieved) != 2 {
		t.Errorf("InsertBulk() inserted %d records, want 2", len(bulkRetrieved))
		return
	}
}
