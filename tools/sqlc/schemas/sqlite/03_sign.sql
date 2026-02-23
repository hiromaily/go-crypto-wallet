-- SQLite schema for sign wallet
-- Converted from MySQL schema

CREATE TABLE auth_account_key (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  coin TEXT NOT NULL CHECK(coin IN ('btc','bch')),
  key_type TEXT NOT NULL DEFAULT 'bip44',
  auth_account TEXT NOT NULL,
  account TEXT NOT NULL DEFAULT 'deposit',
  p2pkh_address TEXT NOT NULL,
  p2sh_segwit_address TEXT NULL,
  bech32_address TEXT NULL,
  taproot_address TEXT NULL,
  full_public_key TEXT NOT NULL,
  multisig_address TEXT NOT NULL DEFAULT '',
  redeem_script TEXT NOT NULL DEFAULT '',
  wallet_import_format TEXT NOT NULL,
  account_extended_privkey TEXT DEFAULT NULL,
  idx INTEGER NOT NULL,
  addr_status INTEGER NOT NULL DEFAULT 0,
  updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(coin, auth_account, account)
);

CREATE UNIQUE INDEX idx_auth_account_key_p2pkh_address ON auth_account_key(p2pkh_address);
CREATE UNIQUE INDEX idx_auth_account_key_p2sh_segwit_address ON auth_account_key(p2sh_segwit_address) WHERE p2sh_segwit_address IS NOT NULL;
CREATE UNIQUE INDEX idx_auth_account_key_bech32_address ON auth_account_key(bech32_address) WHERE bech32_address IS NOT NULL;
CREATE UNIQUE INDEX idx_auth_account_key_wallet_import_format ON auth_account_key(wallet_import_format);
CREATE INDEX idx_auth_account_key_auth_account ON auth_account_key(auth_account);
CREATE INDEX idx_auth_account_key_coin ON auth_account_key(coin);
CREATE INDEX idx_auth_account_key_key_type ON auth_account_key(key_type);
