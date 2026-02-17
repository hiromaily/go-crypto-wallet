-- name: GetXRPMultisigSignaturesByPendingID :many
SELECT * FROM xrp_multisig_signature
WHERE pending_multisig_id = $1
ORDER BY signed_at ASC;

-- name: GetXRPMultisigSignatureByID :one
SELECT * FROM xrp_multisig_signature
WHERE id = $1;

-- name: GetXRPMultisigSignatureByPendingAndSigner :one
SELECT * FROM xrp_multisig_signature
WHERE pending_multisig_id = $1 AND signer_account = $2
LIMIT 1;

-- name: InsertXRPMultisigSignature :one
INSERT INTO xrp_multisig_signature (
  pending_multisig_id, signer_account, signed_tx_blob, signer_weight, signed_at
) VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetSignatureCountByPendingID :one
SELECT COUNT(*) as signature_count
FROM xrp_multisig_signature
WHERE pending_multisig_id = $1;

-- name: GetTotalSignedWeightByPendingID :one
SELECT COALESCE(SUM(signer_weight), 0) as total_weight
FROM xrp_multisig_signature
WHERE pending_multisig_id = $1;

-- name: DeleteXRPMultisigSignaturesByPendingID :exec
DELETE FROM xrp_multisig_signature
WHERE pending_multisig_id = $1;
