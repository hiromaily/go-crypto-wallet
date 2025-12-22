-- Watch database: Generic transaction table for ETH/XRP

CREATE TABLE tx (
  id         BIGINT NOT NULL AUTO_INCREMENT COMMENT 'transaction ID',
  coin       ENUM('eth','xrp','hyt') NOT NULL COMMENT 'coin type code',
  action     ENUM('deposit', 'payment', 'transfer') NOT NULL COMMENT 'action type',
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (id),
  INDEX idx_coin (coin),
  INDEX idx_action (action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for eth/xrp transaction info';
