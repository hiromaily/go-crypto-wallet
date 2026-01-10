-- name: GetAuthFullPubkey :one
SELECT id, coin, auth_account, full_public_key, extended_pubkey, fingerprint, derivation_path, updated_at
FROM auth_fullpubkey WHERE coin = ? AND auth_account = ? LIMIT 1;

-- name: GetAuthFullPubkeyByFingerprint :one
SELECT id, coin, auth_account, full_public_key, extended_pubkey, fingerprint, derivation_path, updated_at
FROM auth_fullpubkey WHERE fingerprint = ? LIMIT 1;

-- name: InsertAuthFullPubkey :execresult
INSERT INTO auth_fullpubkey (coin, auth_account, full_public_key, extended_pubkey, fingerprint, derivation_path)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateAuthFullPubkeyFingerprint :exec
UPDATE auth_fullpubkey SET fingerprint = ?, updated_at = CURRENT_TIMESTAMP WHERE coin = ? AND auth_account = ?;
