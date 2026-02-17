-- name: GetBtcTxByID :one
SELECT * FROM btc_tx
WHERE id = ?;

-- name: GetBtcTxCountByUnsignedHex :one
SELECT COUNT(*) as count FROM btc_tx
WHERE coin = ? AND action = ? AND unsigned_hex_tx = ?;

-- name: GetBtcTxIDBySentHash :one
SELECT id FROM btc_tx
WHERE coin = ? AND action = ? AND sent_hash_tx = ?;

-- name: GetBtcTxIDByUnsignedHex :one
SELECT id FROM btc_tx
WHERE coin = ? AND action = ? AND unsigned_hex_tx = ?;

-- name: GetBtcTxSentHashList :many
SELECT sent_hash_tx FROM btc_tx
WHERE coin = ? AND action = ? AND current_tx_type = ?;

-- name: InsertBtcTx :execresult
INSERT INTO btc_tx (
  coin, action, unsigned_hex_tx, signed_hex_tx, sent_hash_tx,
  total_input_amount, total_output_amount, fee, current_tx_type,
  unsigned_updated_at, sent_updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateBtcTx :exec
UPDATE btc_tx
SET coin = ?, action = ?, unsigned_hex_tx = ?, signed_hex_tx = ?, sent_hash_tx = ?,
    total_input_amount = ?, total_output_amount = ?, fee = ?, current_tx_type = ?,
    unsigned_updated_at = ?, sent_updated_at = ?
WHERE id = ?;

-- name: UpdateBtcTxAfterSent :execresult
UPDATE btc_tx
SET current_tx_type = ?, signed_hex_tx = ?, sent_hash_tx = ?, sent_updated_at = ?
WHERE id = ?;

-- name: UpdateBtcTxType :execresult
UPDATE btc_tx
SET current_tx_type = ?
WHERE id = ?;

-- name: UpdateBtcTxTypeBySentHash :execresult
UPDATE btc_tx
SET current_tx_type = ?
WHERE coin = ? AND action = ? AND sent_hash_tx = ?;

-- name: DeleteAllBtcTx :execresult
DELETE FROM btc_tx;
