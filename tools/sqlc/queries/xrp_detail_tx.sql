-- name: GetXrpDetailTxByID :one
SELECT * FROM xrp_detail_tx
WHERE id = ?;

-- name: GetXrpDetailTxsByTxID :many
SELECT * FROM xrp_detail_tx
WHERE tx_id = ?;

-- name: GetXrpDetailTxBlobList :many
SELECT xrp_detail_tx.tx_blob
FROM xrp_detail_tx
INNER JOIN tx ON tx.id = xrp_detail_tx.tx_id
WHERE tx.coin = ? AND xrp_detail_tx.current_tx_type = ?;

-- name: InsertXrpDetailTx :execresult
INSERT INTO xrp_detail_tx (
  tx_id, uuid, current_tx_type, sender_account, sender_address,
  receiver_account, receiver_address, amount, xrp_tx_type, fee,
  flags, last_ledger_sequence, sequence, signing_pubkey, txn_signature,
  hash, earliest_ledger_version, signed_tx_id, tx_blob, sent_updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateXrpDetailTxAfterSent :execresult
UPDATE xrp_detail_tx
SET current_tx_type = ?, signed_tx_id = ?, tx_blob = ?,
    earliest_ledger_version = ?, sent_updated_at = ?
WHERE uuid = ?;

-- name: UpdateXrpDetailTxType :execresult
UPDATE xrp_detail_tx
SET current_tx_type = ?
WHERE id = ?;

-- name: UpdateXrpDetailTxTypeBySentHash :execresult
UPDATE xrp_detail_tx
SET current_tx_type = ?
WHERE tx_blob = ?;
