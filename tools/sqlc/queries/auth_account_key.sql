-- name: GetAuthAccountKey :one
SELECT * FROM auth_account_key WHERE coin = ? AND auth_account = ? LIMIT 1;

-- name: InsertAuthAccountKey :execresult
INSERT INTO auth_account_key (
  coin, key_type, auth_account, p2pkh_address, p2sh_segwit_address, bech32_address, taproot_address,
  full_public_key, multisig_address, redeem_script, wallet_import_format, idx, addr_status
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateAuthAccountKeyAddrStatus :execresult
UPDATE auth_account_key SET addr_status = ?, updated_at = ?
WHERE coin = ? AND wallet_import_format = ?;
