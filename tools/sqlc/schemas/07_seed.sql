-- Table structure for table `seed`

CREATE TABLE `seed` (
  `id`         tinyint(2) NOT NULL AUTO_INCREMENT COMMENT'ID',
  `coin`       ENUM('btc', 'bch', 'eth', 'xrp', 'hyt') NOT NULL COMMENT'coin type code',
  `seed`       VARCHAR(255) NOT NULL COMMENT'seed',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date',
  PRIMARY KEY (`id`),
  INDEX idx_coin (`coin`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='table for seed';
