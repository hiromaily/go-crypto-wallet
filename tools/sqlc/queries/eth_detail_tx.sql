-- name: GetEthDetailTxByID :one
SELECT * FROM eth_detail_tx
WHERE id = ?;

-- name: GetEthDetailTxsByTxID :many
SELECT * FROM eth_detail_tx
WHERE tx_id = ?;

-- name: GetEthDetailTxSentHashList :many
SELECT eth_detail_tx.sent_hash_tx
FROM eth_detail_tx
INNER JOIN tx ON tx.id = eth_detail_tx.tx_id
WHERE tx.coin = ? AND eth_detail_tx.current_tx_type = ?;

-- name: InsertEthDetailTx :execresult
INSERT INTO eth_detail_tx (
  tx_id, uuid, current_tx_type, sender_account, sender_address,
  receiver_account, receiver_address, amount, fee, gas_limit, nonce,
  unsigned_hex_tx, signed_hex_tx, sent_hash_tx, unsigned_updated_at, sent_updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateEthDetailTxAfterSent :execresult
UPDATE eth_detail_tx
SET current_tx_type = ?, signed_hex_tx = ?, sent_hash_tx = ?, sent_updated_at = ?
WHERE uuid = ?;

-- name: UpdateEthDetailTxType :execresult
UPDATE eth_detail_tx
SET current_tx_type = ?
WHERE id = ?;

-- name: UpdateEthDetailTxTypeBySentHash :execresult
UPDATE eth_detail_tx
SET current_tx_type = ?
WHERE sent_hash_tx = ?;
