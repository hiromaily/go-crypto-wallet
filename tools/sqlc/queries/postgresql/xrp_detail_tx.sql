-- name: GetXRPDetailTxByID :one
SELECT *
FROM xrp_detail_tx
WHERE id = $1;
-- name: GetXRPDetailTxsByTxID :many
SELECT *
FROM xrp_detail_tx
WHERE tx_id = $1;
-- name: GetXRPDetailTxBlobList :many
SELECT xrp_detail_tx.tx_blob
FROM xrp_detail_tx
  INNER JOIN tx ON tx.id = xrp_detail_tx.tx_id
WHERE tx.coin = $1
  AND xrp_detail_tx.current_tx_type = $2;
-- name: InsertXRPDetailTx :one
INSERT INTO xrp_detail_tx (
    tx_id,
    uuid,
    current_tx_type,
    sender_account,
    sender_address,
    receiver_account,
    receiver_address,
    amount,
    xrp_tx_type,
    fee,
    flags,
    last_ledger_sequence,
    sequence,
    signing_pubkey,
    txn_signature,
    hash,
    earliest_ledger_version,
    signed_tx_id,
    tx_blob,
    sent_updated_at
  )
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14,
    $15,
    $16,
    $17,
    $18,
    $19,
    $20
  )
RETURNING id;
-- name: UpdateXRPDetailTxAfterSent :execresult
UPDATE xrp_detail_tx
SET current_tx_type = $1,
  signed_tx_id = $2,
  tx_blob = $3,
  earliest_ledger_version = $4,
  sent_updated_at = $5
WHERE uuid = $6;
-- name: UpdateXRPDetailTxType :execresult
UPDATE xrp_detail_tx
SET current_tx_type = $1
WHERE id = $2;
-- name: UpdateXRPDetailTxTypeBySentHash :execresult
UPDATE xrp_detail_tx
SET current_tx_type = $1
WHERE tx_blob = $2;
