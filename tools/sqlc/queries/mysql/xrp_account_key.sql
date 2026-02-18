-- name: GetXRPAccountKeysByAddrStatus :many
SELECT * FROM xrp_account_key WHERE coin = ? AND account = ? AND addr_status = ?;

-- name: GetXRPAccountKeySecret :one
SELECT master_seed FROM xrp_account_key WHERE coin = ? AND account = ? AND account_id = ? LIMIT 1;

-- name: InsertXRPAccountKey :execresult
INSERT INTO xrp_account_key (
  coin, account, account_id, key_type, master_key, master_seed, master_seed_hex,
  public_key, public_key_hex, is_regular_key_pair, allocated_id, addr_status
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateXRPAccountKeyAddrStatus :execresult
UPDATE xrp_account_key SET addr_status = ?, updated_at = ?
WHERE coin = ? AND account = ? AND account_id = ?;
