-- name: GetXRPRegularKeyByAccountID :one
SELECT * FROM xrp_regular_key
WHERE account_id = ? AND is_active = TRUE
LIMIT 1;

-- name: GetXRPRegularKeysByAccountID :many
SELECT * FROM xrp_regular_key
WHERE account_id = ?
ORDER BY created_at DESC;

-- name: GetActiveXRPRegularKeys :many
SELECT * FROM xrp_regular_key
WHERE is_active = TRUE;

-- name: GetXRPRegularKeyByAddress :one
SELECT * FROM xrp_regular_key
WHERE regular_key_address = ?
LIMIT 1;

-- name: InsertXRPRegularKey :execresult
INSERT INTO xrp_regular_key (
  account_id, regular_key_address, public_key, public_key_hex,
  is_active, set_tx_hash, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateXRPRegularKeyStatus :execresult
UPDATE xrp_regular_key
SET is_active = ?, rotated_at = ?
WHERE id = ?;

-- name: DeactivateXRPRegularKeyByAccountID :execresult
UPDATE xrp_regular_key
SET is_active = FALSE, rotated_at = ?
WHERE account_id = ? AND is_active = TRUE;

-- name: UpdateXRPRegularKeyTxHash :execresult
UPDATE xrp_regular_key
SET set_tx_hash = ?
WHERE id = ?;
