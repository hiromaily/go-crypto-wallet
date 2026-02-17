-- name: GetXRPPendingMultisigByUUID :one
SELECT * FROM xrp_pending_multisig
WHERE tx_uuid = ?;

-- name: GetXRPPendingMultisigByID :one
SELECT * FROM xrp_pending_multisig
WHERE id = ?;

-- name: GetXRPPendingMultisigsByAccountID :many
SELECT * FROM xrp_pending_multisig
WHERE account_id = ? AND status = ?
ORDER BY created_at DESC;

-- name: GetXRPPendingMultisigsByStatus :many
SELECT * FROM xrp_pending_multisig
WHERE status = ?
ORDER BY created_at ASC;

-- name: GetExpiredXRPPendingMultisigs :many
SELECT * FROM xrp_pending_multisig
WHERE status = 'pending' AND expires_at IS NOT NULL AND expires_at < ?;

-- name: InsertXRPPendingMultisig :execresult
INSERT INTO xrp_pending_multisig (
  tx_uuid, account_id, unsigned_tx_json, xrp_tx_type,
  required_quorum, current_weight, status, expires_at, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateXRPPendingMultisigWeight :execresult
UPDATE xrp_pending_multisig
SET current_weight = ?, updated_at = ?
WHERE id = ?;

-- name: UpdateXRPPendingMultisigStatus :execresult
UPDATE xrp_pending_multisig
SET status = ?, updated_at = ?
WHERE id = ?;

-- name: UpdateXRPPendingMultisigReady :execresult
UPDATE xrp_pending_multisig
SET status = 'ready', combined_tx_blob = ?, updated_at = ?
WHERE id = ?;

-- name: UpdateXRPPendingMultisigSubmitted :execresult
UPDATE xrp_pending_multisig
SET status = 'submitted', submitted_tx_hash = ?, updated_at = ?
WHERE id = ?;

-- name: UpdateXRPPendingMultisigConfirmed :execresult
UPDATE xrp_pending_multisig
SET status = 'confirmed', updated_at = ?
WHERE id = ?;

-- name: UpdateXRPPendingMultisigFailed :execresult
UPDATE xrp_pending_multisig
SET status = 'failed', updated_at = ?
WHERE id = ?;

-- name: ExpireXRPPendingMultisigs :execresult
UPDATE xrp_pending_multisig
SET status = 'expired', updated_at = ?
WHERE status = 'pending' AND expires_at IS NOT NULL AND expires_at < ?;
