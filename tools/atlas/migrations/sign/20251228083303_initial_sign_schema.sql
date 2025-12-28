-- Create "auth_account_key" table
CREATE TABLE `auth_account_key` (
  `id` smallint NOT NULL AUTO_INCREMENT COMMENT "ID",
  `coin` enum('btc','bch') NOT NULL COMMENT "coin type code",
  `auth_account` varchar(20) NOT NULL COMMENT "auth type",
  `p2pkh_address` varchar(255) NOT NULL COMMENT "address as standard pubkey script that Pays To PubKey Hash (P2PKH)",
  `p2sh_segwit_address` varchar(255) NOT NULL COMMENT "p2sh-segwit address",
  `bech32_address` varchar(255) NOT NULL COMMENT "bech32 address",
  `taproot_address` varchar(255) NULL COMMENT "taproot address (BIP86)",
  `full_public_key` varchar(255) NOT NULL COMMENT "full public key",
  `multisig_address` varchar(255) NOT NULL DEFAULT "" COMMENT "multisig address",
  `redeem_script` varchar(255) NOT NULL DEFAULT "" COMMENT "redeedScript after multisig address generated",
  `wallet_import_format` varchar(255) NOT NULL COMMENT "WIF",
  `idx` bigint NOT NULL COMMENT "index for hd wallet",
  `addr_status` tinyint NOT NULL DEFAULT 0 COMMENT "progress status for address generating",
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date",
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idex_coin_auth_account` (`coin`, `auth_account`),
  INDEX `idx_auth_account` (`auth_account`),
  UNIQUE INDEX `idx_bech32_address` (`bech32_address`),
  INDEX `idx_coin` (`coin`),
  UNIQUE INDEX `idx_p2pkh_address` (`p2pkh_address`),
  UNIQUE INDEX `idx_p2sh_segwit_address` (`p2sh_segwit_address`),
  UNIQUE INDEX `idx_wallet_import_format` (`wallet_import_format`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT "table for keys for auth account";
-- Create "seed" table
CREATE TABLE `seed` (
  `id` tinyint NOT NULL AUTO_INCREMENT COMMENT "ID",
  `coin` enum('btc','bch') NOT NULL COMMENT "coin type code",
  `seed` varchar(255) NOT NULL COMMENT "seed",
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date",
  PRIMARY KEY (`id`),
  INDEX `idx_coin` (`coin`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT "table for seed";
