-- ETH Account Key Queries
-- name: GetMaxETHAccountKeyIndex :one
SELECT COALESCE(MAX(idx), 0) as max_idx
FROM eth_account_key
WHERE account = $1;
-- name: GetOneETHAccountKeyByMaxID :one
SELECT *
FROM eth_account_key
WHERE account = $1
ORDER BY id DESC
LIMIT 1;
-- name: GetETHAccountKeysByAddrStatus :many
SELECT *
FROM eth_account_key
WHERE account = $1
  AND addr_status = $2;
-- name: GetETHAccountKeyByAddress :one
SELECT *
FROM eth_account_key
WHERE address = $1
LIMIT 1;
-- name: InsertETHAccountKey :one
INSERT INTO eth_account_key (
    account,
    address,
    full_public_key,
    private_key,
    idx,
    addr_status
  )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
-- name: UpdateETHAccountKeyAddrStatus :execresult
UPDATE eth_account_key
SET addr_status = $1,
  updated_at = $2
WHERE account = $3
  AND private_key = $4;
