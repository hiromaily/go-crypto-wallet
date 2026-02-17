-- name: GetAuthFullPubkey :one
SELECT id, coin, auth_account, purpose, full_public_key, extended_pubkey, fingerprint, derivation_path, updated_at
FROM auth_fullpubkey WHERE coin = $1 AND auth_account = $2 AND purpose = 49 LIMIT 1;

-- name: GetAuthFullPubkeyByPurpose :one
SELECT id, coin, auth_account, purpose, full_public_key, extended_pubkey, fingerprint, derivation_path, updated_at
FROM auth_fullpubkey WHERE coin = $1 AND auth_account = $2 AND purpose = $3 LIMIT 1;

-- name: GetAuthFullPubkeyByFingerprint :one
SELECT id, coin, auth_account, purpose, full_public_key, extended_pubkey, fingerprint, derivation_path, updated_at
FROM auth_fullpubkey WHERE fingerprint = $1 LIMIT 1;

-- name: InsertAuthFullPubkey :one
INSERT INTO auth_fullpubkey (coin, auth_account, purpose, full_public_key, extended_pubkey, fingerprint, derivation_path)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: UpdateAuthFullPubkeyFingerprint :exec
UPDATE auth_fullpubkey SET fingerprint = $1, updated_at = CURRENT_TIMESTAMP WHERE coin = $2 AND auth_account = $3 AND purpose = $4;
