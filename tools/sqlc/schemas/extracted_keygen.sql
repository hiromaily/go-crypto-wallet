-- Extracted from dump_keygen.sql

CREATE TABLE account_key (
  id bigint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  coin enum('btc','bch','eth','xrp','hyt') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'coin type code',
  key_type varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'bip44' COMMENT 'key type (bip44, bip49, bip84, bip86, musig2)',
  account enum('client','deposit','payment','stored') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'account type',
  p2pkh_address varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'address as standard pubkey script that Pays To PubKey Hash (P2PKH)',
  p2sh_segwit_address varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'p2sh-segwit address',
  bech32_address varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'bech32 address',
  taproot_address varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'taproot address (BIP86)',
  full_public_key varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'full public key',
  multisig_address varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'multisig address',
  redeem_script varchar(1000) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'redeedScript after multisig address generated',
  wallet_import_format varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'WIF',
  idx bigint NOT NULL COMMENT 'index for hd wallet',
  addr_status tinyint NOT NULL DEFAULT '0' COMMENT 'progress status for address generating',
  updated_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (id),
  UNIQUE KEY idx_p2pkh_address (p2pkh_address),
  UNIQUE KEY idx_wallet_import_format (wallet_import_format),
  KEY idx_account (account),
  KEY idx_coin (coin),
  KEY idx_key_type (key_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='table for keys for any account';


CREATE TABLE auth_fullpubkey (
  id smallint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  coin enum('btc','bch') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'coin type code',
  auth_account varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'auth type',
  full_public_key varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'full public key',
  fingerprint varchar(8) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'BIP32 master key fingerprint (8 hex chars)',
  updated_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (id),
  UNIQUE KEY idex_coin_auth_account (coin,auth_account),
  UNIQUE KEY idx_full_public_key (full_public_key),
  KEY idx_coin (coin),
  KEY idx_fingerprint (fingerprint)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='table for auth key exported from sign db';


CREATE TABLE musig2_nonces (
  id bigint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  signer_id varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Signer identifier',
  transaction_id varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Transaction identifier',
  public_nonce binary(66) NOT NULL COMMENT 'Public nonce (66 bytes: two 33-byte compressed EC points R1||R2)',
  is_used tinyint(1) NOT NULL DEFAULT '0' COMMENT 'true: nonce has been used in signing',
  created_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'creation date',
  used_at datetime DEFAULT NULL COMMENT 'date when nonce was marked as used',
  PRIMARY KEY (id),
  UNIQUE KEY idx_signer_transaction (signer_id,transaction_id) COMMENT 'Prevent duplicate nonces per signer-tx pair (CRITICAL for security)',
  KEY idx_created_at (created_at),
  KEY idx_is_used (is_used),
  KEY idx_transaction_id (transaction_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='MuSig2 nonce commitments for secure storage';


CREATE TABLE seed (
  id tinyint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  coin enum('btc','bch','eth','xrp','hyt') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'coin type code',
  seed varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'seed',
  updated_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (id),
  KEY idx_coin (coin)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='table for seed';


CREATE TABLE xrp_account_key (
  id bigint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  coin enum('xrp') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'coin type code',
  account enum('client','deposit','payment','stored') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'account type',
  account_id varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'account_id',
  key_type tinyint NOT NULL DEFAULT '0' COMMENT 'key_type',
  master_key varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'master_key, DEPRECATED',
  master_seed varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'master_seed',
  master_seed_hex varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'master_seed_hex',
  public_key varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'public_key',
  public_key_hex varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'public_key_hex',
  is_regular_key_pair tinyint(1) NOT NULL DEFAULT '0' COMMENT 'true: this key is for regular key pair',
  allocated_id bigint NOT NULL DEFAULT '0' COMMENT 'index for hd wallet',
  addr_status tinyint NOT NULL DEFAULT '0' COMMENT 'progress status for address generating',
  updated_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (id),
  UNIQUE KEY idx_account_id (account_id),
  UNIQUE KEY idx_master_seed (master_seed),
  KEY idx_account (account),
  KEY idx_coin (coin)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='table for xrp keys for any account';


