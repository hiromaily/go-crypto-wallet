-- name: GetXRPRegularKeyByAccountID :one
SELECT * FROM xrp_regular_key
WHERE account_id = $1 AND is_active = TRUE
LIMIT 1;

-- name: GetXRPRegularKeysByAccountID :many
SELECT * FROM xrp_regular_key
WHERE account_id = $1
ORDER BY created_at DESC;

-- name: GetActiveXRPRegularKeys :many
SELECT * FROM xrp_regular_key
WHERE is_active = TRUE;

-- name: GetXRPRegularKeyByAddress :one
SELECT * FROM xrp_regular_key
WHERE regular_key_address = $1
LIMIT 1;

-- name: InsertXRPRegularKey :one
INSERT INTO xrp_regular_key (
  account_id, regular_key_address, public_key, public_key_hex,
  is_active, set_tx_hash, created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: UpdateXRPRegularKeyStatus :execresult
UPDATE xrp_regular_key
SET is_active = $1, rotated_at = $2
WHERE id = $3;

-- name: DeactivateXRPRegularKeyByAccountID :execresult
UPDATE xrp_regular_key
SET is_active = FALSE, rotated_at = $1
WHERE account_id = $2 AND is_active = TRUE;

-- name: UpdateXRPRegularKeyTxHash :execresult
UPDATE xrp_regular_key
SET set_tx_hash = $1
WHERE id = $2;
