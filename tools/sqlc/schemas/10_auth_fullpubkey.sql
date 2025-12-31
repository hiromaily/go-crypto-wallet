-- Table structure for table `auth_fullpubkey`

CREATE TABLE `auth_fullpubkey` (
  `id`                      SMALLINT(5) NOT NULL AUTO_INCREMENT COMMENT'ID',
  `coin`                    ENUM('btc', 'bch') NOT NULL COMMENT'coin type code',
  `auth_account`            VARCHAR(20) NOT NULL COMMENT'auth type',
  `full_public_key`         VARCHAR(255) NOT NULL COMMENT'full public key',
  `fingerprint`             VARCHAR(8) DEFAULT NULL COMMENT'BIP32 master key fingerprint (8 hex chars)',
  `updated_at`              datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idex_coin_auth_account` (`coin`, `auth_account`),
  UNIQUE KEY `idx_full_public_key` (`full_public_key`),
  INDEX idx_coin (`coin`),
  INDEX idx_fingerprint (`fingerprint`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='table for auth key exported from sign db';
