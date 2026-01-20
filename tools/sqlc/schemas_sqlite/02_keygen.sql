-- SQLite schema for keygen wallet
-- Converted from MySQL schema

CREATE TABLE auth_fullpubkey (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  coin TEXT NOT NULL CHECK(coin IN ('btc','bch')),
  auth_account TEXT NOT NULL,
  purpose INTEGER NOT NULL DEFAULT 49,
  full_public_key TEXT NOT NULL,
  extended_pubkey TEXT DEFAULT NULL,
  fingerprint TEXT DEFAULT NULL,
  derivation_path TEXT DEFAULT NULL,
  updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(coin, auth_account, purpose)
);

CREATE INDEX idx_auth_fullpubkey_coin ON auth_fullpubkey(coin);
CREATE INDEX idx_auth_fullpubkey_fingerprint ON auth_fullpubkey(fingerprint);
CREATE INDEX idx_auth_fullpubkey_purpose ON auth_fullpubkey(purpose);

CREATE TABLE btc_account_key (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  coin TEXT NOT NULL CHECK(coin IN ('btc','bch')),
  key_type TEXT NOT NULL DEFAULT 'bip44',
  account TEXT NOT NULL CHECK(account IN ('client','deposit','payment','stored')),
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
  updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_btc_account_key_p2pkh_address ON btc_account_key(p2pkh_address);
CREATE UNIQUE INDEX idx_btc_account_key_wallet_import_format ON btc_account_key(wallet_import_format);
CREATE INDEX idx_btc_account_key_account ON btc_account_key(account);
CREATE INDEX idx_btc_account_key_coin ON btc_account_key(coin);
CREATE INDEX idx_btc_account_key_key_type ON btc_account_key(key_type);

CREATE TABLE eth_account_key (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  account TEXT NOT NULL CHECK(account IN ('client','deposit','payment','stored')),
  address TEXT NOT NULL UNIQUE,
  full_public_key TEXT NOT NULL,
  private_key TEXT NOT NULL UNIQUE,
  idx INTEGER NOT NULL,
  addr_status INTEGER NOT NULL DEFAULT 0,
  updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_eth_account_key_account ON eth_account_key(account);

CREATE TABLE musig2_nonces (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  signer_id TEXT NOT NULL,
  transaction_id TEXT NOT NULL,
  public_nonce BLOB NOT NULL,
  is_used INTEGER NOT NULL DEFAULT 0,
  created_at TEXT DEFAULT CURRENT_TIMESTAMP,
  used_at TEXT DEFAULT NULL,
  UNIQUE(signer_id, transaction_id)
);

CREATE INDEX idx_musig2_nonces_created_at ON musig2_nonces(created_at);
CREATE INDEX idx_musig2_nonces_is_used ON musig2_nonces(is_used);
CREATE INDEX idx_musig2_nonces_transaction_id ON musig2_nonces(transaction_id);

CREATE TABLE seed (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  coin TEXT NOT NULL CHECK(coin IN ('btc','bch','eth','xrp','hyt')),
  seed TEXT NOT NULL,
  updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_seed_coin ON seed(coin);

CREATE TABLE xrp_account_key (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  coin TEXT NOT NULL CHECK(coin IN ('xrp')),
  account TEXT NOT NULL CHECK(account IN ('client','deposit','payment','stored')),
  account_id TEXT NOT NULL UNIQUE,
  key_type INTEGER NOT NULL DEFAULT 0,
  master_key TEXT NOT NULL,
  master_seed TEXT NOT NULL UNIQUE,
  master_seed_hex TEXT NOT NULL,
  public_key TEXT NOT NULL,
  public_key_hex TEXT NOT NULL,
  is_regular_key_pair INTEGER NOT NULL DEFAULT 0,
  allocated_id INTEGER NOT NULL DEFAULT 0,
  addr_status INTEGER NOT NULL DEFAULT 0,
  updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_xrp_account_key_account ON xrp_account_key(account);
CREATE INDEX idx_xrp_account_key_coin ON xrp_account_key(coin);
