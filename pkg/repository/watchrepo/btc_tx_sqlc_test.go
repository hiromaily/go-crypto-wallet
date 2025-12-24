//go:build integration
// +build integration

package watchrepo_test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/quagmt/udecimal"
	"github.com/stretchr/testify/require"

	"github.com/hiromaily/go-crypto-wallet/pkg/action"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
)

// TestBTCTxSqlc is integration test for BTCTxRepositorySqlc
func TestBTCTxSqlc(t *testing.T) {
	txRepo := testutil.NewBTCTxRepositorySqlc()

	// Delete records
	_, err := txRepo.DeleteAll()
	require.NoError(t, err, "fail to call DeleteAll()")

	// Insert
	inputAmt, _ := udecimal.Parse("0.100")
	outputAmt, _ := udecimal.Parse("0.090")
	feeAmt, _ := udecimal.Parse("0.010")

	hex := "unsigned-hex-sqlc"
	actionType := action.ActionTypePayment
	txItem := &models.BTCTX{
		Coin:              "btc",
		Action:            actionType.String(),
		UnsignedHexTX:     hex,
		TotalInputAmount:  inputAmt,
		TotalOutputAmount: outputAmt,
		Fee:               feeAmt,
	}
	id, err := txRepo.InsertUnsignedTx(actionType, txItem)
	require.NoError(t, err, "fail to call InsertUnsignedTx()")
	txItem.ID = id // Set the ID for later operations
	// check inserted record
	tmpTx, err := txRepo.GetOne(id)
	require.NoError(t, err, "fail to call GetOne()")
	require.Equal(t, hex, tmpTx.UnsignedHexTX, "InsertUnsignedTx() should insert correct hex")
	// check Count
	cnt, err := txRepo.GetCountByUnsignedHex(actionType, hex)
	require.NoError(t, err, "fail to call GetCount()")
	require.Equal(t, int64(1), cnt, "GetCount() should return 1")

	// Update only UnsignedHexTX
	hex2 := "unsigned-hex2-sqlc"
	txItem.UnsignedHexTX = hex2
	_, err = txRepo.Update(txItem)
	require.NoError(t, err, "fail to call UpdateTx()")
	// check updated unsigned hex tx
	tmpTx, err = txRepo.GetOne(txItem.ID)
	require.NoError(t, err, "fail to call GetOne()")
	require.Equal(t, hex2, tmpTx.UnsignedHexTX, "Update() should update UnsignedHexTX")

	// Update like after tx sent
	signedHex := "signed-hex-sqlc"
	sentHashTx := "sent-hash-tx-sqlc"
	_, err = txRepo.UpdateAfterTxSent(txItem.ID, tx.TxTypeSent, signedHex, sentHashTx)
	require.NoError(t, err, "fail to call UpdateTx()")
	// check updated record
	tmpTx, err = txRepo.GetOne(txItem.ID)
	require.NoError(t, err, "fail to call GetOne()")
	require.Equal(t, signedHex, tmpTx.SignedHexTX, "Update() should update SignedHexTX")
	// sent_hash_tx should be retrieved
	hashes, err := txRepo.GetSentHashTx(actionType, tx.TxTypeSent)
	require.NoError(t, err, "fail to call GetSentHashTx()")
	require.Len(t, hashes, 1, "GetSentHashTx() should return 1 hash")

	// update txType
	_, err = txRepo.UpdateTxTypeBySentHashTx(actionType, tx.TxTypeDone, sentHashTx)
	require.NoError(t, err, "fail to call UpdateTxTypeBySentHashTx()")
	// check updated record
	tmpTx, err = txRepo.GetOne(txItem.ID)
	require.NoError(t, err, "fail to call GetOne()")
	require.Equal(t, tx.TxTypeDone.Int8(), tmpTx.CurrentTXType, "UpdateTxTypeBySentHashTx() should update CurrentTXType to TxTypeDone")

	// update txType
	_, err = txRepo.UpdateTxType(txItem.ID, tx.TxTypeNotified)
	require.NoError(t, err, "fail to call UpdateTxType()")
	// check updated record
	tmpTx, err = txRepo.GetOne(txItem.ID)
	require.NoError(t, err, "fail to call GetOne()")
	require.Equal(t, tx.TxTypeNotified.Int8(), tmpTx.CurrentTXType, "UpdateTxType() should update CurrentTXType to TxTypeNotified")
}
