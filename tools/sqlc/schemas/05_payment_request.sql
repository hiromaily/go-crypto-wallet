-- Watch database: Payment request table

CREATE TABLE payment_request (
  id               BIGINT NOT NULL AUTO_INCREMENT COMMENT 'ID',
  coin             ENUM('btc', 'bch', 'eth', 'xrp') NOT NULL COMMENT 'coin type code',
  payment_id       BIGINT DEFAULT NULL COMMENT 'tx table ID for payment action',
  sender_address   VARCHAR(255) NOT NULL COMMENT 'sender address',
  sender_account   VARCHAR(255) NOT NULL COMMENT 'sender account',
  receiver_address VARCHAR(255) NOT NULL COMMENT 'receiver address',
  amount           DECIMAL(26,10) NOT NULL COMMENT 'amount of coin to send',
  is_done          BOOL NOT NULL DEFAULT false COMMENT 'true: unsigned transaction is created',
  updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (id),
  INDEX idx_coin (coin)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for payment request';
