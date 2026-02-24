-- name: GetXRPSignerEntriesByListID :many
SELECT * FROM xrp_signer_entry
WHERE signer_list_id = $1
ORDER BY signer_weight DESC;

-- name: GetXRPSignerEntryByID :one
SELECT * FROM xrp_signer_entry
WHERE id = $1;

-- name: GetXRPSignerEntryByListAndAccount :one
SELECT * FROM xrp_signer_entry
WHERE signer_list_id = $1 AND signer_account = $2
LIMIT 1;

-- name: InsertXRPSignerEntry :one
INSERT INTO xrp_signer_entry (
  signer_list_id, signer_account, signer_weight, created_at
) VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: DeleteXRPSignerEntriesByListID :exec
DELETE FROM xrp_signer_entry
WHERE signer_list_id = $1;

-- name: GetTotalWeightByListID :one
SELECT COALESCE(SUM(signer_weight), 0) as total_weight
FROM xrp_signer_entry
WHERE signer_list_id = $1;
