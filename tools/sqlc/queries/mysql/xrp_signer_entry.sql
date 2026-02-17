-- name: GetXRPSignerEntriesByListID :many
SELECT * FROM xrp_signer_entry
WHERE signer_list_id = ?
ORDER BY signer_weight DESC;

-- name: GetXRPSignerEntryByID :one
SELECT * FROM xrp_signer_entry
WHERE id = ?;

-- name: GetXRPSignerEntryByListAndAccount :one
SELECT * FROM xrp_signer_entry
WHERE signer_list_id = ? AND signer_account = ?
LIMIT 1;

-- name: InsertXRPSignerEntry :execresult
INSERT INTO xrp_signer_entry (
  signer_list_id, signer_account, signer_weight, created_at
) VALUES (?, ?, ?, ?);

-- name: DeleteXRPSignerEntriesByListID :exec
DELETE FROM xrp_signer_entry
WHERE signer_list_id = ?;

-- name: GetTotalWeightByListID :one
SELECT COALESCE(SUM(signer_weight), 0) as total_weight
FROM xrp_signer_entry
WHERE signer_list_id = ?;
