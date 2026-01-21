-- name: GetXRPSignerListByAccountID :one
SELECT * FROM xrp_signer_list
WHERE account_id = ? AND is_active = TRUE
LIMIT 1;

-- name: GetXRPSignerListByID :one
SELECT * FROM xrp_signer_list
WHERE id = ?;

-- name: GetXRPSignerListHistoryByAccountID :many
SELECT * FROM xrp_signer_list
WHERE account_id = ?
ORDER BY created_at DESC;

-- name: InsertXRPSignerList :execresult
INSERT INTO xrp_signer_list (
  account_id, signer_quorum, is_active, set_tx_hash, created_at
) VALUES (?, ?, ?, ?, ?);

-- name: UpdateXRPSignerListStatus :execresult
UPDATE xrp_signer_list
SET is_active = ?, updated_at = ?
WHERE id = ?;

-- name: UpdateXRPSignerListTxHash :execresult
UPDATE xrp_signer_list
SET set_tx_hash = ?, updated_at = ?
WHERE id = ?;

-- name: DeactivateXRPSignerListByAccountID :execresult
UPDATE xrp_signer_list
SET is_active = FALSE, updated_at = ?
WHERE account_id = ? AND is_active = TRUE;
