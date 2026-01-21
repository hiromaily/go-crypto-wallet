-- name: GetXRPMultisigSignaturesByPendingID :many
SELECT * FROM xrp_multisig_signature
WHERE pending_multisig_id = ?
ORDER BY signed_at ASC;

-- name: GetXRPMultisigSignatureByID :one
SELECT * FROM xrp_multisig_signature
WHERE id = ?;

-- name: GetXRPMultisigSignatureByPendingAndSigner :one
SELECT * FROM xrp_multisig_signature
WHERE pending_multisig_id = ? AND signer_account = ?
LIMIT 1;

-- name: InsertXRPMultisigSignature :execresult
INSERT INTO xrp_multisig_signature (
  pending_multisig_id, signer_account, signed_tx_blob, signer_weight, signed_at
) VALUES (?, ?, ?, ?, ?);

-- name: GetSignatureCountByPendingID :one
SELECT COUNT(*) as signature_count
FROM xrp_multisig_signature
WHERE pending_multisig_id = ?;

-- name: GetTotalSignedWeightByPendingID :one
SELECT COALESCE(SUM(signer_weight), 0) as total_weight
FROM xrp_multisig_signature
WHERE pending_multisig_id = ?;

-- name: DeleteXRPMultisigSignaturesByPendingID :exec
DELETE FROM xrp_multisig_signature
WHERE pending_multisig_id = ?;
