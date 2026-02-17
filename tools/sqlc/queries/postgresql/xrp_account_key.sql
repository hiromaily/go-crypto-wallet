-- name: GetXRPAccountKeysByAddrStatus :many
SELECT * FROM xrp_account_key WHERE coin = $1 AND account = $2 AND addr_status = $3;

-- name: GetXRPAccountKeySecret :one
SELECT master_seed FROM xrp_account_key WHERE coin = $1 AND account = $2 AND account_id = $3 LIMIT 1;

-- name: InsertXRPAccountKey :one
INSERT INTO xrp_account_key (
  coin, account, account_id, key_type, master_key, master_seed, master_seed_hex,
  public_key, public_key_hex, is_regular_key_pair, allocated_id, addr_status
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id;

-- name: UpdateXRPAccountKeyAddrStatus :execresult
UPDATE xrp_account_key SET addr_status = $1, updated_at = $2
WHERE coin = $3 AND account = $4 AND account_id = $5;
