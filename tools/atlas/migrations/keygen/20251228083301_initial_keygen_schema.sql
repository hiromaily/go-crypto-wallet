-- Create "account_key" table
CREATE TABLE `account_key` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT "ID",
  `coin` enum('btc','bch','eth','xrp','hyt') NOT NULL COMMENT "coin type code",
  `account` enum('client','deposit','payment','stored') NOT NULL COMMENT "account type",
  `p2pkh_address` varchar(255) NOT NULL COMMENT "address as standard pubkey script that Pays To PubKey Hash (P2PKH)",
  `p2sh_segwit_address` varchar(255) NOT NULL COMMENT "p2sh-segwit address",
  `bech32_address` varchar(255) NOT NULL COMMENT "bech32 address",
  `taproot_address` varchar(255) NULL COMMENT "taproot address (BIP86)",
  `full_public_key` varchar(255) NOT NULL COMMENT "full public key",
  `multisig_address` varchar(255) NOT NULL DEFAULT "" COMMENT "multisig address",
  `redeem_script` varchar(1000) NOT NULL DEFAULT "" COMMENT "redeedScript after multisig address generated",
  `wallet_import_format` varchar(255) NOT NULL COMMENT "WIF",
  `idx` bigint NOT NULL COMMENT "index for hd wallet",
  `addr_status` tinyint NOT NULL DEFAULT 0 COMMENT "progress status for address generating",
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date",
  PRIMARY KEY (`id`),
  INDEX `idx_account` (`account`),
  INDEX `idx_coin` (`coin`),
  UNIQUE INDEX `idx_p2pkh_address` (`p2pkh_address`),
  UNIQUE INDEX `idx_wallet_import_format` (`wallet_import_format`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT "table for keys for any account";
-- Create "auth_fullpubkey" table
CREATE TABLE `auth_fullpubkey` (
  `id` smallint NOT NULL AUTO_INCREMENT COMMENT "ID",
  `coin` enum('btc','bch') NOT NULL COMMENT "coin type code",
  `auth_account` varchar(20) NOT NULL COMMENT "auth type",
  `full_public_key` varchar(255) NOT NULL COMMENT "full public key",
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date",
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idex_coin_auth_account` (`coin`, `auth_account`),
  INDEX `idx_coin` (`coin`),
  UNIQUE INDEX `idx_full_public_key` (`full_public_key`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT "table for auth key exported from sign db";
-- Create "seed" table
CREATE TABLE `seed` (
  `id` tinyint NOT NULL AUTO_INCREMENT COMMENT "ID",
  `coin` enum('btc','bch','eth','xrp','hyt') NOT NULL COMMENT "coin type code",
  `seed` varchar(255) NOT NULL COMMENT "seed",
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date",
  PRIMARY KEY (`id`),
  INDEX `idx_coin` (`coin`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT "table for seed";
-- Create "xrp_account_key" table
CREATE TABLE `xrp_account_key` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT "ID",
  `coin` enum('xrp') NOT NULL COMMENT "coin type code",
  `account` enum('client','deposit','payment','stored') NOT NULL COMMENT "account type",
  `account_id` varchar(255) NOT NULL COMMENT "account_id",
  `key_type` tinyint NOT NULL DEFAULT 0 COMMENT "key_type",
  `master_key` varchar(255) NOT NULL COMMENT "master_key, DEPRECATED",
  `master_seed` varchar(255) NOT NULL COMMENT "master_seed",
  `master_seed_hex` varchar(255) NOT NULL COMMENT "master_seed_hex",
  `public_key` varchar(255) NOT NULL COMMENT "public_key",
  `public_key_hex` varchar(255) NOT NULL COMMENT "public_key_hex",
  `is_regular_key_pair` bool NOT NULL DEFAULT 0 COMMENT "true: this key is for regular key pair",
  `allocated_id` bigint NOT NULL DEFAULT 0 COMMENT "index for hd wallet",
  `addr_status` tinyint NOT NULL DEFAULT 0 COMMENT "progress status for address generating",
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date",
  PRIMARY KEY (`id`),
  INDEX `idx_account` (`account`),
  UNIQUE INDEX `idx_account_id` (`account_id`),
  INDEX `idx_coin` (`coin`),
  UNIQUE INDEX `idx_master_seed` (`master_seed`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT "table for xrp keys for any account";
