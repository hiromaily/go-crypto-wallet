-- MuSig2 nonce management queries
-- SECURITY CRITICAL: These queries enforce nonce uniqueness to prevent private key leakage

-- name: SaveNonce :execresult
INSERT INTO musig2_nonces (signer_id, transaction_id, public_nonce, is_used, created_at)
VALUES (?, ?, ?, ?, ?);

-- name: GetNonceBySignerAndTx :one
SELECT * FROM musig2_nonces
WHERE signer_id = ? AND transaction_id = ?;

-- name: GetAllNoncesForTransaction :many
SELECT * FROM musig2_nonces
WHERE transaction_id = ?;

-- name: GetUnusedNoncesForTransaction :many
SELECT * FROM musig2_nonces
WHERE transaction_id = ? AND is_used = false;

-- name: MarkNonceUsed :execresult
UPDATE musig2_nonces
SET is_used = ?, used_at = ?
WHERE signer_id = ? AND transaction_id = ?;

-- name: DeleteNoncesForTransaction :execresult
DELETE FROM musig2_nonces
WHERE transaction_id = ?;

-- name: CleanupOldUnusedNonces :execresult
DELETE FROM musig2_nonces
WHERE is_used = false AND created_at < ?;
