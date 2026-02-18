-- name: GetAuthAccountKey :one
SELECT
  id, coin, key_type, auth_account, account, p2pkh_address,
  COALESCE(p2sh_segwit_address, '') as p2sh_segwit_address,
  COALESCE(bech32_address, '') as bech32_address,
  taproot_address, full_public_key, multisig_address, redeem_script,
  wallet_import_format, account_extended_privkey, idx, addr_status, updated_at
FROM auth_account_key WHERE coin = $1 AND auth_account = $2 LIMIT 1;

-- name: GetAuthAccountKeyByAccount :one
SELECT
  id, coin, key_type, auth_account, account, p2pkh_address,
  COALESCE(p2sh_segwit_address, '') as p2sh_segwit_address,
  COALESCE(bech32_address, '') as bech32_address,
  taproot_address, full_public_key, multisig_address, redeem_script,
  wallet_import_format, account_extended_privkey, idx, addr_status, updated_at
FROM auth_account_key WHERE coin = $1 AND auth_account = $2 AND account = $3 LIMIT 1;

-- name: InsertAuthAccountKey :one
INSERT INTO auth_account_key (
  coin, key_type, auth_account, account, p2pkh_address, p2sh_segwit_address, bech32_address, taproot_address,
  full_public_key, multisig_address, redeem_script, wallet_import_format, account_extended_privkey, idx, addr_status
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING id;

-- name: UpdateAuthAccountKeyAddrStatus :execresult
UPDATE auth_account_key SET addr_status = $1, updated_at = $2
WHERE coin = $3 AND wallet_import_format = $4;
