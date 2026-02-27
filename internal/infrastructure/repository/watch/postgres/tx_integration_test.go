//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainBTC "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/btc"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	watchpostgres "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/postgres"
)

// ============================================================
// BTC Transaction (btc_tx) Repository Tests
// ============================================================

// TestBTCTxRepository_InsertAndGetOne verifies InsertUnsignedTx stores a record and
// GetOne retrieves it with correct field values.
func TestBTCTxRepository_InsertAndGetOne(t *testing.T) {
	db := openWatchDB(t)
	repo := watchpostgres.NewBTCTxRepositorySqlc(db, domainCoin.BTC)

	tx := &domainBTC.BTCTransaction{
		UnsignedHexTx:     "aabb0001",
		SignedHexTx:       "",
		SentHashTx:        "",
		TotalInputAmount:  "1.5000000000",
		TotalOutputAmount: "1.4000000000",
		Fee:               "0.1000000000",
		CurrentTxType:     domainTx.TxTypeUnsigned,
	}

	id, err := repo.InsertUnsignedTx(domainTx.ActionTypeDeposit, tx)
	require.NoError(t, err, "InsertUnsignedTx should succeed")
	require.Positive(t, id, "returned ID must be positive")
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), "DELETE FROM btc_tx WHERE id = $1", id)
	})

	got, err := repo.GetOne(id)
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, domainCoin.BTC, got.CoinTypeCode)
	assert.Equal(t, domainTx.ActionTypeDeposit, got.ActionType)
	assert.Equal(t, domainTx.TxTypeUnsigned, got.CurrentTxType)
	assert.Equal(t, "aabb0001", got.UnsignedHexTx)
}

// TestBTCTxRepository_DecimalPrecision verifies that NUMERIC(26,10) amounts are
// stored and retrieved with full 10-decimal-place precision.
func TestBTCTxRepository_DecimalPrecision(t *testing.T) {
	db := openWatchDB(t)
	repo := watchpostgres.NewBTCTxRepositorySqlc(db, domainCoin.BTC)

	// Use satoshi-level precision (8 decimal places) and beyond.
	const (
		totalIn  = "99999999.9999999999"
		totalOut = "0.0000000001"
		fee      = "0.0100000000"
	)

	tx := &domainBTC.BTCTransaction{
		UnsignedHexTx:     "decimal_precision_test",
		TotalInputAmount:  totalIn,
		TotalOutputAmount: totalOut,
		Fee:               fee,
		CurrentTxType:     domainTx.TxTypeUnsigned,
	}

	id, err := repo.InsertUnsignedTx(domainTx.ActionTypePayment, tx)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), "DELETE FROM btc_tx WHERE id = $1", id)
	})

	got, err := repo.GetOne(id)
	require.NoError(t, err)

	assert.Equal(t, totalIn, got.TotalInputAmount, "TotalInputAmount should preserve 10 decimal places")
	assert.Equal(t, totalOut, got.TotalOutputAmount, "TotalOutputAmount should preserve 10 decimal places")
	assert.Equal(t, fee, got.Fee, "Fee should preserve 10 decimal places")
}

// TestBTCTxRepository_UpdateTxType verifies that UpdateTxType changes the
// current_tx_type field and returns the correct row count.
func TestBTCTxRepository_UpdateTxType(t *testing.T) {
	db := openWatchDB(t)
	repo := watchpostgres.NewBTCTxRepositorySqlc(db, domainCoin.BTC)

	tx := &domainBTC.BTCTransaction{
		UnsignedHexTx:     "update_type_test",
		TotalInputAmount:  "1.0000000000",
		TotalOutputAmount: "0.9000000000",
		Fee:               "0.1000000000",
		CurrentTxType:     domainTx.TxTypeUnsigned,
	}
	id, err := repo.InsertUnsignedTx(domainTx.ActionTypeDeposit, tx)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), "DELETE FROM btc_tx WHERE id = $1", id)
	})

	rowsAffected, err := repo.UpdateTxType(id, domainTx.TxTypeSigned)
	require.NoError(t, err)
	assert.EqualValues(t, 1, rowsAffected, "exactly one row should be updated")

	got, err := repo.GetOne(id)
	require.NoError(t, err)
	assert.Equal(t, domainTx.TxTypeSigned, got.CurrentTxType, "tx type should be updated to signed")
}

// TestBTCTxRepository_GetCountByUnsignedHex verifies query filtering by coin, action,
// and unsigned_hex_tx.
func TestBTCTxRepository_GetCountByUnsignedHex(t *testing.T) {
	db := openWatchDB(t)
	repo := watchpostgres.NewBTCTxRepositorySqlc(db, domainCoin.BTC)

	hexTx := fmt.Sprintf("unique_hex_%d", time.Now().UnixNano())
	tx := &domainBTC.BTCTransaction{
		UnsignedHexTx:     hexTx,
		TotalInputAmount:  "1.0000000000",
		TotalOutputAmount: "0.9000000000",
		Fee:               "0.1000000000",
		CurrentTxType:     domainTx.TxTypeUnsigned,
	}
	id, err := repo.InsertUnsignedTx(domainTx.ActionTypeDeposit, tx)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), "DELETE FROM btc_tx WHERE id = $1", id)
	})

	count, err := repo.GetCountByUnsignedHex(domainTx.ActionTypeDeposit, hexTx)
	require.NoError(t, err)
	assert.EqualValues(t, 1, count, "should find exactly one record with the unique hex")

	// Different action type should return 0.
	count, err = repo.GetCountByUnsignedHex(domainTx.ActionTypePayment, hexTx)
	require.NoError(t, err)
	assert.EqualValues(t, 0, count, "different action should return 0")
}

// TestBTCTxRepository_UpdateAfterTxSent verifies that UpdateAfterTxSent advances the
// transaction state to sent and persists the signed hex and sent hash.
func TestBTCTxRepository_UpdateAfterTxSent(t *testing.T) {
	db := openWatchDB(t)
	repo := watchpostgres.NewBTCTxRepositorySqlc(db, domainCoin.BTC)

	tx := &domainBTC.BTCTransaction{
		UnsignedHexTx:     "after_sent_test",
		TotalInputAmount:  "2.0000000000",
		TotalOutputAmount: "1.9000000000",
		Fee:               "0.1000000000",
		CurrentTxType:     domainTx.TxTypeUnsigned,
	}
	id, err := repo.InsertUnsignedTx(domainTx.ActionTypePayment, tx)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), "DELETE FROM btc_tx WHERE id = $1", id)
	})

	sentHash := fmt.Sprintf("senthash_%d", time.Now().UnixNano())
	rowsAffected, err := repo.UpdateAfterTxSent(id, domainTx.TxTypeSent, "signed_hex_data", sentHash)
	require.NoError(t, err)
	assert.EqualValues(t, 1, rowsAffected)

	got, err := repo.GetOne(id)
	require.NoError(t, err)
	assert.Equal(t, domainTx.TxTypeSent, got.CurrentTxType)
	assert.Equal(t, "signed_hex_data", got.SignedHexTx)
	assert.Equal(t, sentHash, got.SentHashTx)
	assert.NotNil(t, got.SentUpdatedAt, "SentUpdatedAt should be set after tx sent")
}

// ============================================================
// BTC TX Input (btc_tx_input) Repository Tests
// ============================================================

// TestBTCTxInputRepository_InsertAndGetOne verifies Insert stores an input record and
// GetOne retrieves it with correct field values.
func TestBTCTxInputRepository_InsertAndGetOne(t *testing.T) {
	db := openWatchDB(t)
	repo := watchpostgres.NewBTCTxInputRepositorySqlc(db, domainCoin.BTC)

	inputTxid := fmt.Sprintf("inputtxid_%d", time.Now().UnixNano())
	input := &domainBTC.BTCTxInput{
		TxID:               1000,
		InputTxid:          inputTxid,
		InputVout:          0,
		InputAddress:       "1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2",
		InputAccount:       "deposit",
		InputAmount:        "0.5000000000",
		InputConfirmations: 6,
	}

	err := repo.Insert(input)
	require.NoError(t, err, "Insert should succeed for valid input")
	t.Cleanup(func() {
		_, _ = db.ExecContext(
			context.Background(),
			"DELETE FROM btc_tx_input WHERE input_txid = $1",
			inputTxid,
		)
	})

	inputs, err := repo.GetAllByTxID(1000)
	require.NoError(t, err)

	var found *domainBTC.BTCTxInput
	for _, inp := range inputs {
		if inp.InputTxid == inputTxid {
			found = inp
			break
		}
	}
	require.NotNil(t, found, "inserted input should be retrievable via GetAllByTxID")

	assert.Equal(t, int64(1000), found.TxID)
	assert.Equal(t, inputTxid, found.InputTxid)
	assert.Equal(t, uint32(0), found.InputVout)
	assert.Equal(t, uint64(6), found.InputConfirmations)
}

// TestBTCTxInputRepository_DecimalPrecision verifies that input_amount (NUMERIC 26,10)
// preserves full precision through a store-and-retrieve round-trip.
func TestBTCTxInputRepository_DecimalPrecision(t *testing.T) {
	db := openWatchDB(t)
	repo := watchpostgres.NewBTCTxInputRepositorySqlc(db, domainCoin.BTC)

	const preciseAmount = "0.0000000001"
	inputTxid := fmt.Sprintf("inputdecimal_%d", time.Now().UnixNano())
	input := &domainBTC.BTCTxInput{
		TxID:               2000,
		InputTxid:          inputTxid,
		InputVout:          0,
		InputAddress:       "1BTCPrecisionTestAddress",
		InputAccount:       "client",
		InputAmount:        preciseAmount,
		InputConfirmations: 1,
	}

	err := repo.Insert(input)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = db.ExecContext(
			context.Background(),
			"DELETE FROM btc_tx_input WHERE input_txid = $1",
			inputTxid,
		)
	})

	inputs, err := repo.GetAllByTxID(2000)
	require.NoError(t, err)
	for _, inp := range inputs {
		if inp.InputTxid == inputTxid {
			assert.Equal(t, preciseAmount, inp.InputAmount,
				"input_amount NUMERIC(26,10) should preserve full precision")
			return
		}
	}
	t.Errorf("inserted input %q not found", inputTxid)
}

// TestBTCTxInputRepository_GetAllByTxID verifies that InsertBulk stores multiple inputs
// and GetAllByTxID returns all of them filtered by tx_id.
func TestBTCTxInputRepository_GetAllByTxID(t *testing.T) {
	db := openWatchDB(t)
	repo := watchpostgres.NewBTCTxInputRepositorySqlc(db, domainCoin.BTC)

	txID := int64(3000)
	inputTxid1 := fmt.Sprintf("bulk_input1_%d", time.Now().UnixNano())
	inputTxid2 := fmt.Sprintf("bulk_input2_%d", time.Now().UnixNano())
	t.Cleanup(func() {
		_, _ = db.ExecContext(
			context.Background(),
			"DELETE FROM btc_tx_input WHERE input_txid IN ($1, $2)",
			inputTxid1, inputTxid2,
		)
	})

	inputs := []*domainBTC.BTCTxInput{
		{
			TxID: txID, InputTxid: inputTxid1, InputVout: 0,
			InputAddress: "addr1", InputAccount: "deposit",
			InputAmount: "1.0000000000", InputConfirmations: 3,
		},
		{
			TxID: txID, InputTxid: inputTxid2, InputVout: 1,
			InputAddress: "addr2", InputAccount: "payment",
			InputAmount: "2.0000000000", InputConfirmations: 5,
		},
	}

	err := repo.InsertBulk(inputs)
	require.NoError(t, err)

	got, err := repo.GetAllByTxID(txID)
	require.NoError(t, err)

	foundTxids := make(map[string]bool, len(got))
	for _, inp := range got {
		foundTxids[inp.InputTxid] = true
	}
	assert.True(t, foundTxids[inputTxid1], "first bulk input should be retrievable")
	assert.True(t, foundTxids[inputTxid2], "second bulk input should be retrievable")
}

// ============================================================
// BTC TX Output (btc_tx_output) Repository Tests
// ============================================================

// TestBTCTxOutputRepository_InsertAndGetAllByTxID verifies Insert stores an output
// and GetAllByTxID retrieves it with correct values including NUMERIC precision.
func TestBTCTxOutputRepository_InsertAndGetAllByTxID(t *testing.T) {
	db := openWatchDB(t)
	repo := watchpostgres.NewBTCTxOutputRepositorySqlc(db, domainCoin.BTC)

	txID := int64(4000)
	outputAddr := fmt.Sprintf("1OutputTest%d", time.Now().UnixNano())
	const outputAmount = "1.2345678901"

	output := &domainBTC.BTCTxOutput{
		TxID:          txID,
		OutputAddress: outputAddr,
		OutputAccount: "payment",
		OutputAmount:  outputAmount,
		IsChange:      false,
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(
			context.Background(),
			"DELETE FROM btc_tx_output WHERE output_address = $1",
			outputAddr,
		)
	})

	err := repo.Insert(output)
	require.NoError(t, err, "Insert should succeed")

	outputs, err := repo.GetAllByTxID(txID)
	require.NoError(t, err)

	var found *domainBTC.BTCTxOutput
	for _, out := range outputs {
		if out.OutputAddress == outputAddr {
			found = out
			break
		}
	}
	require.NotNil(t, found, "inserted output should be retrievable")
	assert.Equal(t, outputAmount, found.OutputAmount, "output_amount precision should be preserved")
	assert.False(t, found.IsChange)
}

// ============================================================
// ETH Detail TX (eth_detail_tx) Repository Tests
// ============================================================

// newTestETHDetailTx creates a minimal valid ETHDetailTx for testing.
func newTestETHDetailTx(txID int64, uuid string) *domainETH.ETHDetailTx {
	return &domainETH.ETHDetailTx{
		TxID:            txID,
		UUID:            uuid,
		CurrentTxType:   domainTx.TxTypeUnsigned,
		SenderAccount:   "payment",
		SenderAddress:   "0xSenderAddr",
		ReceiverAccount: "client",
		ReceiverAddress: "0xReceiverAddr",
		Amount:          1000000000000000000, // 1 ETH in wei
		Fee:             21000,
		GasLimit:        21000,
		Nonce:           1,
		UnsignedHexTx:   "0xabcdef",
		SignedHexTx:     "",
		SentHashTx:      "",
	}
}

// TestETHDetailTxRepository_InsertAndGetOne verifies Insert stores an ETH detail tx
// and GetOne retrieves it with correct field values.
func TestETHDetailTxRepository_InsertAndGetOne(t *testing.T) {
	db := openWatchDB(t)
	repo := watchpostgres.NewETHDetailTXInputRepositorySqlc(db, domainCoin.ETH)

	uuid := fmt.Sprintf("eth-uuid-%d", time.Now().UnixNano())
	tx := newTestETHDetailTx(5000, uuid)

	t.Cleanup(func() {
		_, _ = db.ExecContext(
			context.Background(),
			"DELETE FROM eth_detail_tx WHERE uuid = $1",
			uuid,
		)
	})

	err := repo.Insert(tx)
	require.NoError(t, err, "Insert should succeed for valid ETH detail tx")

	// GetOne requires an ID; retrieve via GetAllByTxID first.
	all, err := repo.GetAllByTxID(5000)
	require.NoError(t, err)

	var found *domainETH.ETHDetailTx
	for _, item := range all {
		if item.UUID == uuid {
			found = item
			break
		}
	}
	require.NotNil(t, found, "inserted ETH detail tx should be retrievable by tx_id")

	// Now verify GetOne with the discovered ID.
	got, err := repo.GetOne(found.ID)
	require.NoError(t, err)
	assert.Equal(t, uuid, got.UUID)
	assert.Equal(t, domainTx.TxTypeUnsigned, got.CurrentTxType)
	assert.Equal(t, uint64(1000000000000000000), got.Amount)
	assert.Equal(t, uint64(21000), got.Fee)
	assert.Equal(t, uint32(21000), got.GasLimit)
	assert.Equal(t, uint64(1), got.Nonce)
}

// TestETHDetailTxRepository_GetAllByTxID verifies InsertBulk stores multiple ETH txs
// and GetAllByTxID returns all records for a given tx_id.
func TestETHDetailTxRepository_GetAllByTxID(t *testing.T) {
	db := openWatchDB(t)
	repo := watchpostgres.NewETHDetailTXInputRepositorySqlc(db, domainCoin.ETH)

	txID := int64(6000)
	uuid1 := fmt.Sprintf("eth-bulk-1-%d", time.Now().UnixNano())
	uuid2 := fmt.Sprintf("eth-bulk-2-%d", time.Now().UnixNano())
	t.Cleanup(func() {
		_, _ = db.ExecContext(
			context.Background(),
			"DELETE FROM eth_detail_tx WHERE uuid IN ($1, $2)",
			uuid1, uuid2,
		)
	})

	err := repo.InsertBulk([]*domainETH.ETHDetailTx{
		newTestETHDetailTx(txID, uuid1),
		newTestETHDetailTx(txID, uuid2),
	})
	require.NoError(t, err)

	got, err := repo.GetAllByTxID(txID)
	require.NoError(t, err)

	uuids := make(map[string]bool, len(got))
	for _, item := range got {
		uuids[item.UUID] = true
	}
	assert.True(t, uuids[uuid1], "first bulk ETH tx should be retrievable")
	assert.True(t, uuids[uuid2], "second bulk ETH tx should be retrievable")
}

// TestETHDetailTxRepository_UpdateAfterTxSent verifies that UpdateAfterTxSent
// advances the ETH tx state and persists the sent hash.
func TestETHDetailTxRepository_UpdateAfterTxSent(t *testing.T) {
	db := openWatchDB(t)
	repo := watchpostgres.NewETHDetailTXInputRepositorySqlc(db, domainCoin.ETH)

	uuid := fmt.Sprintf("eth-sent-%d", time.Now().UnixNano())
	t.Cleanup(func() {
		_, _ = db.ExecContext(
			context.Background(),
			"DELETE FROM eth_detail_tx WHERE uuid = $1",
			uuid,
		)
	})

	err := repo.Insert(newTestETHDetailTx(7000, uuid))
	require.NoError(t, err)

	sentHash := fmt.Sprintf("0xETHSentHash%d", time.Now().UnixNano())
	rowsAffected, err := repo.UpdateAfterTxSent(uuid, domainTx.TxTypeSent, "0xSignedHex", sentHash)
	require.NoError(t, err)
	assert.EqualValues(t, 1, rowsAffected)

	all, err := repo.GetAllByTxID(7000)
	require.NoError(t, err)
	for _, item := range all {
		if item.UUID == uuid {
			assert.Equal(t, domainTx.TxTypeSent, item.CurrentTxType)
			assert.Equal(t, sentHash, item.SentHashTx)
			return
		}
	}
	t.Errorf("ETH tx with uuid %q not found after update", uuid)
}

// ============================================================
// XRP Detail TX (xrp_detail_tx) Repository Tests
// ============================================================

// newTestXRPDetailTx creates a minimal valid XRPDetailTx for testing.
func newTestXRPDetailTx(txID int64, uuid string) *domainXRP.XRPDetailTx {
	return &domainXRP.XRPDetailTx{
		TxID:               txID,
		UUID:               uuid,
		CurrentTxType:      domainTx.TxTypeUnsigned,
		SenderAccount:      "payment",
		SenderAddress:      "rXRPSenderAddr",
		ReceiverAccount:    "client",
		ReceiverAddress:    "rXRPReceiverAddr",
		Amount:             "1000000", // 1 XRP in drops
		XrpTxType:          "Payment",
		Fee:                "12",
		Flags:              0,
		LastLedgerSequence: 1000,
		Sequence:           42,
		SigningPubkey:      "",
		TxnSignature:       "",
		Hash:               "",
	}
}

// TestXRPDetailTxRepository_InsertAndGetOne verifies Insert stores an XRP detail tx
// and GetOne retrieves it with correct field values.
func TestXRPDetailTxRepository_InsertAndGetOne(t *testing.T) {
	db := openWatchDB(t)
	repo := watchpostgres.NewXRPDetailTxInputRepositorySqlc(db, domainCoin.XRP)

	uuid := fmt.Sprintf("xrp-uuid-%d", time.Now().UnixNano())
	tx := newTestXRPDetailTx(8000, uuid)
	t.Cleanup(func() {
		_, _ = db.ExecContext(
			context.Background(),
			"DELETE FROM xrp_detail_tx WHERE uuid = $1",
			uuid,
		)
	})

	err := repo.Insert(tx)
	require.NoError(t, err, "Insert should succeed for valid XRP detail tx")

	all, err := repo.GetAllByTxID(8000)
	require.NoError(t, err)

	var found *domainXRP.XRPDetailTx
	for _, item := range all {
		if item.UUID == uuid {
			found = item
			break
		}
	}
	require.NotNil(t, found, "inserted XRP detail tx should be retrievable by tx_id")

	got, err := repo.GetOne(found.ID)
	require.NoError(t, err)
	assert.Equal(t, uuid, got.UUID)
	assert.Equal(t, "1000000", got.Amount)
	assert.Equal(t, "Payment", got.XrpTxType)
	assert.Equal(t, uint64(42), got.Sequence)
	assert.Equal(t, domainTx.TxTypeUnsigned, got.CurrentTxType)
}

// TestXRPDetailTxRepository_GetAllByTxID verifies that InsertBulk stores multiple XRP
// txs and GetAllByTxID returns all of them for a given tx_id.
func TestXRPDetailTxRepository_GetAllByTxID(t *testing.T) {
	db := openWatchDB(t)
	repo := watchpostgres.NewXRPDetailTxInputRepositorySqlc(db, domainCoin.XRP)

	txID := int64(9000)
	uuid1 := fmt.Sprintf("xrp-bulk-1-%d", time.Now().UnixNano())
	uuid2 := fmt.Sprintf("xrp-bulk-2-%d", time.Now().UnixNano())
	t.Cleanup(func() {
		_, _ = db.ExecContext(
			context.Background(),
			"DELETE FROM xrp_detail_tx WHERE uuid IN ($1, $2)",
			uuid1, uuid2,
		)
	})

	err := repo.InsertBulk([]*domainXRP.XRPDetailTx{
		newTestXRPDetailTx(txID, uuid1),
		newTestXRPDetailTx(txID, uuid2),
	})
	require.NoError(t, err)

	got, err := repo.GetAllByTxID(txID)
	require.NoError(t, err)

	uuids := make(map[string]bool, len(got))
	for _, item := range got {
		uuids[item.UUID] = true
	}
	assert.True(t, uuids[uuid1], "first bulk XRP tx should be retrievable")
	assert.True(t, uuids[uuid2], "second bulk XRP tx should be retrievable")
}

// TestXRPDetailTxRepository_UpdateAfterTxSent verifies UpdateAfterTxSent advances the
// XRP tx state and persists signed tx_id and tx_blob.
func TestXRPDetailTxRepository_UpdateAfterTxSent(t *testing.T) {
	db := openWatchDB(t)
	repo := watchpostgres.NewXRPDetailTxInputRepositorySqlc(db, domainCoin.XRP)

	uuid := fmt.Sprintf("xrp-sent-%d", time.Now().UnixNano())
	t.Cleanup(func() {
		_, _ = db.ExecContext(
			context.Background(),
			"DELETE FROM xrp_detail_tx WHERE uuid = $1",
			uuid,
		)
	})

	err := repo.Insert(newTestXRPDetailTx(10000, uuid))
	require.NoError(t, err)

	signedTxID := fmt.Sprintf("XRPSignedTxID%d", time.Now().UnixNano())
	txBlob := "AABBCCDD"
	rowsAffected, err := repo.UpdateAfterTxSent(uuid, domainTx.TxTypeSent, signedTxID, txBlob, 999)
	require.NoError(t, err)
	assert.EqualValues(t, 1, rowsAffected)

	all, err := repo.GetAllByTxID(10000)
	require.NoError(t, err)
	for _, item := range all {
		if item.UUID == uuid {
			assert.Equal(t, domainTx.TxTypeSent, item.CurrentTxType)
			assert.Equal(t, signedTxID, item.SignedTxID)
			assert.Equal(t, txBlob, item.TxBlob)
			assert.Equal(t, uint64(999), item.EarliestLedgerVersion)
			return
		}
	}
	t.Errorf("XRP tx with uuid %q not found after update", uuid)
}

// ============================================================
// Concurrent Access Test
// ============================================================

// TestBTCTxRepository_ConcurrentAccess verifies that multiple goroutines can insert
// BTC transactions concurrently without data races or errors.
func TestBTCTxRepository_ConcurrentAccess(t *testing.T) {
	db := openWatchDB(t)

	const n = 5
	type result struct {
		id  int64
		err error
	}
	results := make(chan result, n)

	for i := range n {
		go func(idx int) {
			repo := watchpostgres.NewBTCTxRepositorySqlc(db, domainCoin.BTC)
			tx := &domainBTC.BTCTransaction{
				UnsignedHexTx:     fmt.Sprintf("concurrent_%d_%d", idx, time.Now().UnixNano()),
				TotalInputAmount:  "1.0000000000",
				TotalOutputAmount: "0.9000000000",
				Fee:               "0.1000000000",
				CurrentTxType:     domainTx.TxTypeUnsigned,
			}
			id, err := repo.InsertUnsignedTx(domainTx.ActionTypeDeposit, tx)
			results <- result{id: id, err: err}
		}(i)
	}

	ids := make([]int64, 0, n)
	for range n {
		r := <-results
		assert.NoError(t, r.err, "concurrent insert should not error")
		if r.id > 0 {
			ids = append(ids, r.id)
		}
	}

	t.Cleanup(func() {
		for _, id := range ids {
			_, _ = db.ExecContext(
				context.Background(),
				"DELETE FROM btc_tx WHERE id = $1",
				id,
			)
		}
	})

	assert.Len(t, ids, n, "all %d concurrent inserts should succeed", n)
}
