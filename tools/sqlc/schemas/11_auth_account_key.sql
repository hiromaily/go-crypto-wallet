-- Table structure for table `auth_account_key`

CREATE TABLE `auth_account_key` (
  `id`                      SMALLINT(5) NOT NULL AUTO_INCREMENT COMMENT'ID',
  `coin`                    ENUM('btc', 'bch') NOT NULL COMMENT'coin type code',
  `key_type`                VARCHAR(20) DEFAULT 'bip44' NOT NULL COMMENT 'key type (bip44, bip49, bip84, bip86, musig2)',
  `auth_account`            VARCHAR(20) NOT NULL COMMENT'auth type',
  `p2pkh_address`           VARCHAR(255) NOT NULL COMMENT'address as standard pubkey script that Pays To PubKey Hash (P2PKH)',
  `p2sh_segwit_address`     VARCHAR(255) NOT NULL COMMENT'p2sh-segwit address',
  `bech32_address`          VARCHAR(255) NOT NULL COMMENT'bech32 address',
  `full_public_key`         VARCHAR(255) NOT NULL COMMENT'full public key',
  `multisig_address`        VARCHAR(255) DEFAULT '' NOT NULL COMMENT'multisig address',
  `redeem_script`           VARCHAR(255) DEFAULT '' NOT NULL COMMENT'redeedScript after multisig address generated',
  `wallet_import_format`    VARCHAR(255) NOT NULL COMMENT'WIF',
  `idx`                     BIGINT(20) NOT NULL COMMENT'index for hd wallet',
  `addr_status`             tinyint(2) DEFAULT 0 NOT NULL COMMENT'progress status for address generating',
  `updated_at`              datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idex_coin_auth_account` (`coin`, `auth_account`),
  UNIQUE KEY `idx_p2pkh_address` (`p2pkh_address`),
  UNIQUE KEY `idx_p2sh_segwit_address` (`p2sh_segwit_address`),
  UNIQUE KEY `idx_bech32_address` (`bech32_address`),
  UNIQUE KEY `idx_wallet_import_format` (`wallet_import_format`),
  INDEX idx_coin (`coin`),
  INDEX idx_key_type (`key_type`),
  INDEX idx_auth_account (`auth_account`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='table for keys for auth account';
