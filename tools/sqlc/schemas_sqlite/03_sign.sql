-- SQLite schema for sign wallet
-- Converted from MySQL schema

CREATE TABLE auth_account_key (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  coin TEXT NOT NULL CHECK(coin IN ('btc','bch')),
  key_type TEXT NOT NULL DEFAULT 'bip44',
  auth_account TEXT NOT NULL,
  account TEXT NOT NULL DEFAULT 'deposit',
  p2pkh_address TEXT NOT NULL UNIQUE,
  p2sh_segwit_address TEXT NOT NULL UNIQUE,
  bech32_address TEXT NOT NULL UNIQUE,
  taproot_address TEXT DEFAULT NULL,
  full_public_key TEXT NOT NULL,
  multisig_address TEXT NOT NULL DEFAULT '',
  redeem_script TEXT NOT NULL DEFAULT '',
  wallet_import_format TEXT NOT NULL UNIQUE,
  account_extended_privkey TEXT DEFAULT NULL,
  idx INTEGER NOT NULL,
  addr_status INTEGER NOT NULL DEFAULT 0,
  updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(coin, auth_account, account)
);

CREATE INDEX idx_auth_account_key_auth_account ON auth_account_key(auth_account);
CREATE INDEX idx_auth_account_key_coin ON auth_account_key(coin);
CREATE INDEX idx_auth_account_key_key_type ON auth_account_key(key_type);

CREATE TABLE seed (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  coin TEXT NOT NULL CHECK(coin IN ('btc','bch','eth','xrp','hyt')),
  seed TEXT NOT NULL,
  updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_seed_coin ON seed(coin);
