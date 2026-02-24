-- MuSig2 nonce management queries
-- SECURITY CRITICAL: These queries enforce nonce uniqueness to prevent private key leakage

-- name: SaveNonce :one
INSERT INTO musig2_nonces (signer_id, transaction_id, public_nonce, is_used, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetNonceBySignerAndTx :one
SELECT * FROM musig2_nonces
WHERE signer_id = $1 AND transaction_id = $2;

-- name: GetAllNoncesForTransaction :many
SELECT * FROM musig2_nonces
WHERE transaction_id = $1;

-- name: GetUnusedNoncesForTransaction :many
SELECT * FROM musig2_nonces
WHERE transaction_id = $1 AND is_used = false;

-- name: MarkNonceUsed :execresult
UPDATE musig2_nonces
SET is_used = $1, used_at = $2
WHERE signer_id = $3 AND transaction_id = $4;

-- name: DeleteNoncesForTransaction :execresult
DELETE FROM musig2_nonces
WHERE transaction_id = $1;

-- name: CleanupOldUnusedNonces :execresult
DELETE FROM musig2_nonces
WHERE is_used = false AND created_at < $1;
