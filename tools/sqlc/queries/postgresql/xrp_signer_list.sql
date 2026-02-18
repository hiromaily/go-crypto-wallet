-- name: GetXRPSignerListByAccountID :one
SELECT * FROM xrp_signer_list
WHERE account_id = $1 AND is_active = TRUE
LIMIT 1;

-- name: GetXRPSignerListByID :one
SELECT * FROM xrp_signer_list
WHERE id = $1;

-- name: GetXRPSignerListHistoryByAccountID :many
SELECT * FROM xrp_signer_list
WHERE account_id = $1
ORDER BY created_at DESC;

-- name: InsertXRPSignerList :one
INSERT INTO xrp_signer_list (
  account_id, signer_quorum, is_active, set_tx_hash, created_at
) VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: UpdateXRPSignerListStatus :execresult
UPDATE xrp_signer_list
SET is_active = $1, updated_at = $2
WHERE id = $3;

-- name: UpdateXRPSignerListTxHash :execresult
UPDATE xrp_signer_list
SET set_tx_hash = $1, updated_at = $2
WHERE id = $3;

-- name: DeactivateXRPSignerListByAccountID :execresult
UPDATE xrp_signer_list
SET is_active = FALSE, updated_at = $1
WHERE account_id = $2 AND is_active = TRUE;
