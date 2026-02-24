-- name: GetETHDetailTXByID :one
SELECT * FROM eth_detail_tx
WHERE id = $1;

-- name: GetETHDetailTXsByTxID :many
SELECT * FROM eth_detail_tx
WHERE tx_id = $1;

-- name: GetETHDetailTXSentHashList :many
SELECT eth_detail_tx.sent_hash_tx
FROM eth_detail_tx
INNER JOIN tx ON tx.id = eth_detail_tx.tx_id
WHERE tx.coin = $1 AND eth_detail_tx.current_tx_type = $2;

-- name: InsertETHDetailTX :one
INSERT INTO eth_detail_tx (
  tx_id, uuid, current_tx_type, sender_account, sender_address,
  receiver_account, receiver_address, amount, fee, gas_limit, nonce,
  unsigned_hex_tx, signed_hex_tx, sent_hash_tx, unsigned_updated_at, sent_updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING id;

-- name: UpdateETHDetailTXAfterSent :execresult
UPDATE eth_detail_tx
SET current_tx_type = $1, signed_hex_tx = $2, sent_hash_tx = $3, sent_updated_at = $4
WHERE uuid = $5;

-- name: UpdateETHDetailTXType :execresult
UPDATE eth_detail_tx
SET current_tx_type = $1
WHERE id = $2;

-- name: UpdateETHDetailTXTypeBySentHash :execresult
UPDATE eth_detail_tx
SET current_tx_type = $1
WHERE sent_hash_tx = $2;
