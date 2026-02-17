-- BTC/BCH Account Key Queries

-- name: GetMaxBtcAccountKeyIndex :one
SELECT COALESCE(MAX(idx), 0) as max_idx FROM btc_account_key WHERE coin = $1 AND account = $2;

-- name: GetOneBtcAccountKeyByMaxID :one
SELECT
  id, coin, key_type, account, p2pkh_address,
  COALESCE(p2sh_segwit_address, '') as p2sh_segwit_address,
  COALESCE(bech32_address, '') as bech32_address,
  taproot_address, full_public_key, multisig_address, redeem_script,
  wallet_import_format, account_extended_privkey, idx, addr_status, updated_at
FROM btc_account_key WHERE coin = $1 AND account = $2 ORDER BY id DESC LIMIT 1;

-- name: GetBtcAccountKeysByAddrStatus :many
SELECT
  id, coin, key_type, account, p2pkh_address,
  COALESCE(p2sh_segwit_address, '') as p2sh_segwit_address,
  COALESCE(bech32_address, '') as bech32_address,
  taproot_address, full_public_key, multisig_address, redeem_script,
  wallet_import_format, account_extended_privkey, idx, addr_status, updated_at
FROM btc_account_key WHERE coin = $1 AND account = $2 AND addr_status = $3;

-- name: GetBtcAccountKeysByMultisigAddresses :many
SELECT
  id, coin, key_type, account, p2pkh_address,
  COALESCE(p2sh_segwit_address, '') as p2sh_segwit_address,
  COALESCE(bech32_address, '') as bech32_address,
  taproot_address, full_public_key, multisig_address, redeem_script,
  wallet_import_format, account_extended_privkey, idx, addr_status, updated_at
FROM btc_account_key WHERE coin = $1 AND account = $2 AND multisig_address = ANY(sqlc.arg('addrs')::text[]);

-- name: InsertBtcAccountKey :one
INSERT INTO btc_account_key (
  coin, key_type, account, p2pkh_address, p2sh_segwit_address, bech32_address, taproot_address,
  full_public_key, multisig_address, redeem_script, wallet_import_format, account_extended_privkey, idx, addr_status
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING id;

-- name: UpdateBtcAccountKeyAddress :execresult
UPDATE btc_account_key SET p2pkh_address = $1, updated_at = $2
WHERE coin = $3 AND account = $4 AND p2sh_segwit_address = $5;

-- name: UpdateBtcAccountKeyAddrStatus :execresult
UPDATE btc_account_key SET addr_status = $1, updated_at = $2
WHERE coin = $3 AND account = $4 AND wallet_import_format = $5;

-- name: UpdateBtcAccountKeyMultisigAddr :execresult
UPDATE btc_account_key
SET multisig_address = $1, redeem_script = $2, addr_status = $3, updated_at = $4
WHERE coin = $5 AND account = $6 AND full_public_key = $7;
