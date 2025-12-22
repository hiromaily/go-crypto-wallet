-- name: GetMaxAccountKeyIndex :one
SELECT COALESCE(MAX(idx), 0) as max_idx FROM account_key WHERE coin = ? AND account = ?;

-- name: GetOneAccountKeyByMaxID :one
SELECT * FROM account_key WHERE coin = ? AND account = ? ORDER BY id DESC LIMIT 1;

-- name: GetAccountKeysByAddrStatus :many
SELECT * FROM account_key WHERE coin = ? AND account = ? AND addr_status = ?;

-- name: GetAccountKeysByMultisigAddresses :many
SELECT * FROM account_key WHERE coin = ? AND account = ? AND multisig_address IN (sqlc.slice('addrs'));

-- name: InsertAccountKey :execresult
INSERT INTO account_key (
  coin, account, p2pkh_address, p2sh_segwit_address, bech32_address,
  full_public_key, multisig_address, redeem_script, wallet_import_format, idx, addr_status
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateAccountKeyAddress :execresult
UPDATE account_key SET p2pkh_address = ?, updated_at = ?
WHERE coin = ? AND account = ? AND p2sh_segwit_address = ?;

-- name: UpdateAccountKeyAddrStatus :execresult
UPDATE account_key SET addr_status = ?, updated_at = ?
WHERE coin = ? AND account = ? AND wallet_import_format = ?;

-- name: UpdateAccountKeyMultisigAddr :execresult
UPDATE account_key
SET multisig_address = ?, redeem_script = ?, addr_status = ?, updated_at = ?
WHERE coin = ? AND account = ? AND full_public_key = ?;
