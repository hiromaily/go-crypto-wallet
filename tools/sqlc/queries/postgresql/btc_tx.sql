-- name: GetBtcTxByID :one
SELECT * FROM btc_tx
WHERE id = $1;

-- name: GetBtcTxCountByUnsignedHex :one
SELECT COUNT(*) as count FROM btc_tx
WHERE coin = $1 AND action = $2 AND unsigned_hex_tx = $3;

-- name: GetBtcTxIDBySentHash :one
SELECT id FROM btc_tx
WHERE coin = $1 AND action = $2 AND sent_hash_tx = $3;

-- name: GetBtcTxIDByUnsignedHex :one
SELECT id FROM btc_tx
WHERE coin = $1 AND action = $2 AND unsigned_hex_tx = $3;

-- name: GetBtcTxSentHashList :many
SELECT sent_hash_tx FROM btc_tx
WHERE coin = $1 AND action = $2 AND current_tx_type = $3;

-- name: InsertBtcTx :one
INSERT INTO btc_tx (
  coin, action, unsigned_hex_tx, signed_hex_tx, sent_hash_tx,
  total_input_amount, total_output_amount, fee, current_tx_type,
  unsigned_updated_at, sent_updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id;

-- name: UpdateBtcTx :exec
UPDATE btc_tx
SET coin = $1, action = $2, unsigned_hex_tx = $3, signed_hex_tx = $4, sent_hash_tx = $5,
    total_input_amount = $6, total_output_amount = $7, fee = $8, current_tx_type = $9,
    unsigned_updated_at = $10, sent_updated_at = $11
WHERE id = $12;

-- name: UpdateBtcTxAfterSent :execresult
UPDATE btc_tx
SET current_tx_type = $1, signed_hex_tx = $2, sent_hash_tx = $3, sent_updated_at = $4
WHERE id = $5;

-- name: UpdateBtcTxType :execresult
UPDATE btc_tx
SET current_tx_type = $1
WHERE id = $2;

-- name: UpdateBtcTxTypeBySentHash :execresult
UPDATE btc_tx
SET current_tx_type = $1
WHERE coin = $2 AND action = $3 AND sent_hash_tx = $4;

-- name: DeleteAllBtcTx :execresult
DELETE FROM btc_tx;
