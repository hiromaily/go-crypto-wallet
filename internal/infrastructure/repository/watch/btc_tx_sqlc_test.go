//go:build integration
// +build integration

package watch_test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
	watchTestutil "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/testutil"
)

// TestBTCTxSqlc is integration test for BTCTxRepositorySqlc
func TestBTCTxSqlc(t *testing.T) {
	txRepo := watchTestutil.NewBTCTxRepositorySqlc()

	// Delete records
	_, err := txRepo.DeleteAll()
	require.NoError(t, err, "fail to call DeleteAll()")

	// Insert
	hex := "unsigned-hex-sqlc"
	actionType := domainTx.ActionTypePayment
	txItem := &sqlcgen.BtcTx{
		Coin:              sqlcgen.BtcTxCoinBtc,
		Action:            sqlcgen.BtcTxActionPayment,
		UnsignedHexTx:     hex,
		TotalInputAmount:  "0.100",
		TotalOutputAmount: "0.090",
		Fee:               "0.010",
	}
	id, err := txRepo.InsertUnsignedTx(actionType, txItem)
	require.NoError(t, err, "fail to call InsertUnsignedTx()")
	txItem.ID = id // Set the ID for later operations
	// check inserted record
	tmpTx, err := txRepo.GetOne(id)
	require.NoError(t, err, "fail to call GetOne()")
	require.Equal(t, hex, tmpTx.UnsignedHexTx, "InsertUnsignedTx() should insert correct hex")
	// check Count
	cnt, err := txRepo.GetCountByUnsignedHex(actionType, hex)
	require.NoError(t, err, "fail to call GetCount()")
	require.Equal(t, int64(1), cnt, "GetCount() should return 1")

	// Update only UnsignedHexTx
	hex2 := "unsigned-hex2-sqlc"
	txItem.UnsignedHexTx = hex2
	_, err = txRepo.Update(txItem)
	require.NoError(t, err, "fail to call UpdateTx()")
	// check updated unsigned hex tx
	tmpTx, err = txRepo.GetOne(txItem.ID)
	require.NoError(t, err, "fail to call GetOne()")
	require.Equal(t, hex2, tmpTx.UnsignedHexTx, "Update() should update UnsignedHexTx")

	// Update like after tx sent
	signedHex := "signed-hex-sqlc"
	sentHashTx := "sent-hash-tx-sqlc"
	_, err = txRepo.UpdateAfterTxSent(txItem.ID, domainTx.TxTypeSent, signedHex, sentHashTx)
	require.NoError(t, err, "fail to call UpdateTx()")
	// check updated record
	tmpTx, err = txRepo.GetOne(txItem.ID)
	require.NoError(t, err, "fail to call GetOne()")
	require.Equal(t, signedHex, tmpTx.SignedHexTx, "Update() should update SignedHexTx")
	// sent_hash_tx should be retrieved
	hashes, err := txRepo.GetSentHashTx(actionType, domainTx.TxTypeSent)
	require.NoError(t, err, "fail to call GetSentHashTx()")
	require.Len(t, hashes, 1, "GetSentHashTx() should return 1 hash")

	// update txType
	_, err = txRepo.UpdateTxTypeBySentHashTx(actionType, domainTx.TxTypeDone, sentHashTx)
	require.NoError(t, err, "fail to call UpdateTxTypeBySentHashTx()")
	// check updated record
	tmpTx, err = txRepo.GetOne(txItem.ID)
	require.NoError(t, err, "fail to call GetOne()")
	require.Equal(
		t,
		domainTx.TxTypeDone.Int8(),
		tmpTx.CurrentTxType,
		"UpdateTxTypeBySentHashTx() should update CurrentTxType to TxTypeDone",
	)

	// update txType
	_, err = txRepo.UpdateTxType(txItem.ID, domainTx.TxTypeNotified)
	require.NoError(t, err, "fail to call UpdateTxType()")
	// check updated record
	tmpTx, err = txRepo.GetOne(txItem.ID)
	require.NoError(t, err, "fail to call GetOne()")
	require.Equal(
		t,
		domainTx.TxTypeNotified.Int8(),
		tmpTx.CurrentTxType,
		"UpdateTxType() should update CurrentTxType to TxTypeNotified",
	)
}
