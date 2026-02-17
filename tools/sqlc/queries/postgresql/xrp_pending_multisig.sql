-- name: GetXRPPendingMultisigByUUID :one
SELECT * FROM xrp_pending_multisig
WHERE tx_uuid = $1;

-- name: GetXRPPendingMultisigByID :one
SELECT * FROM xrp_pending_multisig
WHERE id = $1;

-- name: GetXRPPendingMultisigsByAccountID :many
SELECT * FROM xrp_pending_multisig
WHERE account_id = $1 AND status = $2
ORDER BY created_at DESC;

-- name: GetXRPPendingMultisigsByStatus :many
SELECT * FROM xrp_pending_multisig
WHERE status = $1
ORDER BY created_at ASC;

-- name: GetExpiredXRPPendingMultisigs :many
SELECT * FROM xrp_pending_multisig
WHERE status = 'pending' AND expires_at IS NOT NULL AND expires_at < $1;

-- name: InsertXRPPendingMultisig :one
INSERT INTO xrp_pending_multisig (
  tx_uuid, account_id, unsigned_tx_json, xrp_tx_type,
  required_quorum, current_weight, status, expires_at, created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;

-- name: UpdateXRPPendingMultisigWeight :execresult
UPDATE xrp_pending_multisig
SET current_weight = $1, updated_at = $2
WHERE id = $3;

-- name: UpdateXRPPendingMultisigStatus :execresult
UPDATE xrp_pending_multisig
SET status = $1, updated_at = $2
WHERE id = $3;

-- name: UpdateXRPPendingMultisigReady :execresult
UPDATE xrp_pending_multisig
SET status = 'ready', combined_tx_blob = $1, updated_at = $2
WHERE id = $3;

-- name: UpdateXRPPendingMultisigSubmitted :execresult
UPDATE xrp_pending_multisig
SET status = 'submitted', submitted_tx_hash = $1, updated_at = $2
WHERE id = $3;

-- name: UpdateXRPPendingMultisigConfirmed :execresult
UPDATE xrp_pending_multisig
SET status = 'confirmed', updated_at = $1
WHERE id = $2;

-- name: UpdateXRPPendingMultisigFailed :execresult
UPDATE xrp_pending_multisig
SET status = 'failed', updated_at = $1
WHERE id = $2;

-- name: ExpireXRPPendingMultisigs :execresult
UPDATE xrp_pending_multisig
SET status = 'expired', updated_at = $1
WHERE status = 'pending' AND expires_at IS NOT NULL AND expires_at < $2;
