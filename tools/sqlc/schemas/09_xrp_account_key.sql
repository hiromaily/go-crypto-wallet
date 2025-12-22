-- Table structure for table `xrp_account_key`

CREATE TABLE `xrp_account_key` (
  `id`                      BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT'ID',
  `coin`                    ENUM('xrp') NOT NULL COMMENT'coin type code',
  `account`                 ENUM('client', 'deposit', 'payment', 'stored') NOT NULL COMMENT'account type',
  `account_id`              VARCHAR(255) NOT NULL COMMENT'account_id',
  `key_type`                tinyint(2) DEFAULT 0 NOT NULL COMMENT'key_type',
  `master_key`              VARCHAR(255) NOT NULL COMMENT'master_key, DEPRECATED',
  `master_seed`             VARCHAR(255) NOT NULL COMMENT'master_seed',
  `master_seed_hex`         VARCHAR(255) NOT NULL COMMENT'master_seed_hex',
  `public_key`              VARCHAR(255) NOT NULL COMMENT'public_key',
  `public_key_hex`          VARCHAR(255) NOT NULL COMMENT'public_key_hex',
  `is_regular_key_pair`     BOOL NOT NULL DEFAULT false COMMENT'true: this key is for regular key pair',
  `allocated_id`            BIGINT(20) DEFAULT 0 NOT NULL COMMENT'index for hd wallet',
  `addr_status`             tinyint(2) DEFAULT 0 NOT NULL COMMENT'progress status for address generating',
  `updated_at`              datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_account_id` (`account_id`),
  UNIQUE KEY `idx_master_seed` (`master_seed`),
  INDEX idx_coin (`coin`),
  INDEX idx_account (`account`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='table for xrp keys for any account';
