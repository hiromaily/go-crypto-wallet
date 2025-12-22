-- Watch database: Address table

CREATE TABLE address (
  id             BIGINT NOT NULL AUTO_INCREMENT COMMENT 'ID',
  coin           ENUM('btc', 'bch', 'eth', 'xrp', 'hyt') NOT NULL COMMENT 'coin type code',
  account        ENUM('client', 'deposit', 'payment', 'stored') NOT NULL COMMENT 'account type',
  wallet_address VARCHAR(255) NOT NULL COMMENT 'wallet address',
  is_allocated   BOOL NOT NULL DEFAULT false COMMENT 'true: address is allocated(used)',
  updated_at     DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (id),
  UNIQUE KEY idx_wallet_address (wallet_address),
  INDEX idx_coin (coin),
  INDEX idx_account (account)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for account pubkey';
