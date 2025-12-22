-- Watch database: Ethereum transaction details

CREATE TABLE eth_detail_tx (
  id                  BIGINT NOT NULL AUTO_INCREMENT COMMENT 'ID',
  tx_id               BIGINT NOT NULL COMMENT 'eth_tx table ID',
  uuid                VARCHAR(36) NOT NULL COMMENT 'UUID',
  current_tx_type     TINYINT NOT NULL DEFAULT 1 COMMENT 'current transaction type',
  sender_account      VARCHAR(255) NOT NULL COMMENT 'sender account',
  sender_address      VARCHAR(255) NOT NULL COMMENT 'sender address',
  receiver_account    VARCHAR(255) NOT NULL COMMENT 'receiver account',
  receiver_address    VARCHAR(255) NOT NULL COMMENT 'receiver address',
  amount              BIGINT UNSIGNED NOT NULL COMMENT 'amount of coin to receive',
  fee                 BIGINT UNSIGNED NOT NULL COMMENT 'fee',
  gas_limit           MEDIUMINT UNSIGNED NOT NULL COMMENT 'gas limit',
  nonce               BIGINT UNSIGNED NOT NULL COMMENT 'nonce',
  unsigned_hex_tx     TEXT NOT NULL COMMENT 'HEX string for unsigned transaction',
  signed_hex_tx       TEXT NOT NULL DEFAULT '' COMMENT 'HEX string for signed transaction',
  sent_hash_tx        TEXT NOT NULL DEFAULT '' COMMENT 'Hash for sent transaction',
  unsigned_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date for unsigned transaction created',
  sent_updated_at     DATETIME DEFAULT NULL COMMENT 'updated date for signed transaction sent',
  PRIMARY KEY (id),
  UNIQUE KEY idx_uuid (uuid),
  INDEX idx_txid (tx_id),
  INDEX idx_sender_account (sender_account),
  INDEX idx_receiver_account (receiver_account)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for eth transaction detail';
