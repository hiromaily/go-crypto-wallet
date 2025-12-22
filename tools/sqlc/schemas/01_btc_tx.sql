-- Watch database: Bitcoin transaction tables

CREATE TABLE btc_tx (
  id                  BIGINT NOT NULL AUTO_INCREMENT COMMENT 'transaction ID',
  coin                ENUM('btc', 'bch') NOT NULL COMMENT 'coin type code',
  action              ENUM('deposit', 'payment', 'transfer') NOT NULL COMMENT 'action type',
  unsigned_hex_tx     TEXT NOT NULL COMMENT 'HEX string for unsigned transaction',
  signed_hex_tx       TEXT NOT NULL DEFAULT '' COMMENT 'HEX string for signed transaction',
  sent_hash_tx        TEXT NOT NULL DEFAULT '' COMMENT 'Hash for sent transaction',
  total_input_amount  DECIMAL(26,10) NOT NULL COMMENT 'total amount of coin to send',
  total_output_amount DECIMAL(26,10) NOT NULL COMMENT 'total amount of coin to receive without fee',
  fee                 DECIMAL(26,10) NOT NULL COMMENT 'fee',
  current_tx_type     TINYINT NOT NULL DEFAULT 1 COMMENT 'current transaction type',
  unsigned_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date for unsigned transaction created',
  sent_updated_at     DATETIME DEFAULT NULL COMMENT 'updated date for signed transaction sent',
  PRIMARY KEY (id),
  INDEX idx_coin (coin),
  INDEX idx_action (action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for btc transaction info';

CREATE TABLE btc_tx_input (
  id                  BIGINT NOT NULL AUTO_INCREMENT COMMENT 'ID',
  tx_id               BIGINT NOT NULL COMMENT 'tx table ID',
  input_txid          VARCHAR(255) NOT NULL COMMENT 'txid for input',
  input_vout          MEDIUMINT UNSIGNED NOT NULL COMMENT 'vout for input',
  input_address       VARCHAR(255) NOT NULL COMMENT 'sender address for input',
  input_account       VARCHAR(255) NOT NULL COMMENT 'sender account for input',
  input_amount        DECIMAL(26,10) NOT NULL COMMENT 'amount of coin to send for input',
  input_confirmations BIGINT UNSIGNED NOT NULL COMMENT 'block confirmations when unspent rpc returned',
  updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (id),
  INDEX idx_tx_id (tx_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for input transaction';

CREATE TABLE btc_tx_output (
  id             BIGINT NOT NULL AUTO_INCREMENT COMMENT 'ID',
  tx_id          BIGINT NOT NULL COMMENT 'tx table ID',
  output_address VARCHAR(255) NOT NULL COMMENT 'receiver address for output',
  output_account VARCHAR(255) NOT NULL COMMENT 'receiver account for output',
  output_amount  DECIMAL(26,10) NOT NULL COMMENT 'amount of coin to receive',
  is_change      BOOL NOT NULL DEFAULT false COMMENT 'true: output is for fee',
  updated_at     DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (id),
  INDEX idx_tx_id (tx_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for output transaction';
