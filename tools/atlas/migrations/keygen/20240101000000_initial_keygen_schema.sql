-- Initial migration for keygen schema
-- Converted from docker/mysql/sqls/definition_keygen.sql

-- Table: seed
CREATE TABLE `seed` (
  `id`         tinyint(2) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `coin`       ENUM('btc', 'bch', 'eth', 'xrp', 'hyt') NOT NULL COMMENT 'coin type code',
  `seed`       VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'seed',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (`id`),
  INDEX idx_coin (`coin`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for seed';

-- Table: account_key
CREATE TABLE `account_key` (
  `id`                      BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `coin`                    ENUM('btc', 'bch', 'eth', 'xrp', 'hyt') NOT NULL COMMENT 'coin type code',
  `account`                 ENUM('client', 'deposit', 'payment', 'stored') NOT NULL COMMENT 'account type',
  `p2pkh_address`           VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'address as standard pubkey script that Pays To PubKey Hash (P2PKH)',
  `p2sh_segwit_address`     VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'p2sh-segwit address',
  `bech32_address`          VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'bech32 address',
  `taproot_address`         VARCHAR(255) COLLATE utf8_unicode_ci NULL DEFAULT NULL COMMENT 'taproot address (BIP86)',
  `full_public_key`         VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'full public key',
  `multisig_address`        VARCHAR(255) COLLATE utf8_unicode_ci DEFAULT '' NOT NULL COMMENT 'multisig address',
  `redeem_script`           VARCHAR(1000) COLLATE utf8_unicode_ci DEFAULT '' NOT NULL COMMENT 'redeedScript after multisig address generated',
  `wallet_import_format`    VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'WIF',
  `idx`                     BIGINT(20) NOT NULL COMMENT 'index for hd wallet',
  `addr_status`             tinyint(2) DEFAULT 0 NOT NULL COMMENT 'progress status for address generating',
  `updated_at`              datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_p2pkh_address` (`p2pkh_address`),
  UNIQUE KEY `idx_wallet_import_format` (`wallet_import_format`),
  INDEX idx_coin (`coin`),
  INDEX idx_account (`account`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for keys for any account';

-- Table: xrp_account_key
CREATE TABLE `xrp_account_key` (
  `id`                      BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `coin`                    ENUM('xrp') NOT NULL COMMENT 'coin type code',
  `account`                 ENUM('client', 'deposit', 'payment', 'stored') NOT NULL COMMENT 'account type',
  `account_id`              VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'account_id',
  `key_type`                tinyint(2) DEFAULT 0 NOT NULL COMMENT 'key_type',
  `master_key`              VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'master_key, DEPRECATED',
  `master_seed`             VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'master_seed',
  `master_seed_hex`         VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'master_seed_hex',
  `public_key`              VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'public_key',
  `public_key_hex`          VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'public_key_hex',
  `is_regular_key_pair`     BOOL NOT NULL DEFAULT false COMMENT 'true: this key is for regular key pair',
  `allocated_id`            BIGINT(20) DEFAULT 0 NOT NULL COMMENT 'index for hd wallet',
  `addr_status`             tinyint(2) DEFAULT 0 NOT NULL COMMENT 'progress status for address generating',
  `updated_at`              datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_account_id` (`account_id`),
  UNIQUE KEY `idx_master_seed` (`master_seed`),
  INDEX idx_coin (`coin`),
  INDEX idx_account (`account`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for xrp keys for any account';

-- Table: auth_fullpubkey
CREATE TABLE `auth_fullpubkey` (
  `id`                      SMALLINT(5) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `coin`                    ENUM('btc', 'bch') NOT NULL COMMENT 'coin type code',
  `auth_account`            VARCHAR(20)  COLLATE utf8_unicode_ci NOT NULL COMMENT 'auth type',
  `full_public_key`         VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'full public key',
  `updated_at`              datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idex_coin_auth_account` (`coin`, `auth_account`),
  UNIQUE KEY `idx_full_public_key` (`full_public_key`),
  INDEX idx_coin (`coin`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for auth key exported from sign db';

