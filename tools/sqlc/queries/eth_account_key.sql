-- ETH Account Key Queries
-- name: GetMaxETHAccountKeyIndex :one
SELECT COALESCE(MAX(idx), 0) as max_idx
FROM eth_account_key
WHERE account = ?;
-- name: GetOneETHAccountKeyByMaxID :one
SELECT *
FROM eth_account_key
WHERE account = ?
ORDER BY id DESC
LIMIT 1;
-- name: GetETHAccountKeysByAddrStatus :many
SELECT *
FROM eth_account_key
WHERE account = ?
  AND addr_status = ?;
-- name: GetETHAccountKeyByAddress :one
SELECT *
FROM eth_account_key
WHERE address = ?
LIMIT 1;
-- name: InsertETHAccountKey :execresult
INSERT INTO eth_account_key (
    account,
    address,
    full_public_key,
    private_key,
    idx,
    addr_status
  )
VALUES (?, ?, ?, ?, ?, ?);
-- name: UpdateETHAccountKeyAddrStatus :execresult
UPDATE eth_account_key
SET addr_status = ?,
  updated_at = ?
WHERE account = ?
  AND private_key = ?;
