-- BTC/BCH Account Key Queries
-- Note: This file was renamed from account_key.sql to btc_account_key.sql

-- name: GetMaxBtcAccountKeyIndex :one
SELECT COALESCE(MAX(idx), 0) as max_idx FROM btc_account_key WHERE coin = ? AND account = ?;

-- name: GetOneBtcAccountKeyByMaxID :one
SELECT * FROM btc_account_key WHERE coin = ? AND account = ? ORDER BY id DESC LIMIT 1;

-- name: GetBtcAccountKeysByAddrStatus :many
SELECT * FROM btc_account_key WHERE coin = ? AND account = ? AND addr_status = ?;

-- name: GetBtcAccountKeysByMultisigAddresses :many
SELECT * FROM btc_account_key WHERE coin = ? AND account = ? AND multisig_address IN (sqlc.slice('addrs'));

-- name: InsertBtcAccountKey :execresult
INSERT INTO btc_account_key (
  coin, key_type, account, p2pkh_address, p2sh_segwit_address, bech32_address, taproot_address,
  full_public_key, multisig_address, redeem_script, wallet_import_format, idx, addr_status
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateBtcAccountKeyAddress :execresult
UPDATE btc_account_key SET p2pkh_address = ?, updated_at = ?
WHERE coin = ? AND account = ? AND p2sh_segwit_address = ?;

-- name: UpdateBtcAccountKeyAddrStatus :execresult
UPDATE btc_account_key SET addr_status = ?, updated_at = ?
WHERE coin = ? AND account = ? AND wallet_import_format = ?;

-- name: UpdateBtcAccountKeyMultisigAddr :execresult
UPDATE btc_account_key
SET multisig_address = ?, redeem_script = ?, addr_status = ?, updated_at = ?
WHERE coin = ? AND account = ? AND full_public_key = ?;
