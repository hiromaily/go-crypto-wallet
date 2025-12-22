-- name: GetAuthFullPubkey :one
SELECT * FROM auth_fullpubkey WHERE coin = ? AND auth_account = ? LIMIT 1;

-- name: InsertAuthFullPubkey :execresult
INSERT INTO auth_fullpubkey (coin, auth_account, full_public_key) VALUES (?, ?, ?);
