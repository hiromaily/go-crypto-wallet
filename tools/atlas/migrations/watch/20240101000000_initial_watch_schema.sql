-- Initial migration for watch schema
-- Converted from docker/mysql/sqls/definition_watch.sql and payment_request.sql

-- Table: btc_tx
CREATE TABLE `btc_tx` (
  `id`                  BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'transaction ID',
  `coin`                ENUM('btc', 'bch') NOT NULL COMMENT 'coin type code',
  `action`              ENUM('deposit', 'payment', 'transfer') NOT NULL COMMENT 'action type',
  `unsigned_hex_tx`     TEXT COLLATE utf8_unicode_ci NOT NULL COMMENT 'HEX string for unsigned transaction',
  `signed_hex_tx`       TEXT COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT 'HEX string for signed transaction',
  `sent_hash_tx`        TEXT COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT 'Hash for sent transaction',
  `total_input_amount`  DECIMAL(26,10) NOT NULL COMMENT 'total amount of coin to send',
  `total_output_amount` DECIMAL(26,10) NOT NULL COMMENT 'total amount of coin to receive without fee',
  `fee`                 DECIMAL(26,10) NOT NULL COMMENT 'fee',
  `current_tx_type`     tinyint(2) NOT NULL DEFAULT 1 COMMENT 'current transaction type',
  `unsigned_updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date for unsigned transaction created',
  `sent_updated_at`     datetime DEFAULT NULL COMMENT 'updated date for signed transaction sent',
  PRIMARY KEY (`id`),
  INDEX idx_coin (`coin`),
  INDEX idx_action (`action`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for btc transaction info';

-- Table: btc_tx_input
CREATE TABLE `btc_tx_input` (
  `id`             BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `tx_id`          BIGINT(20) NOT NULL COMMENT 'tx table ID',
  `input_txid`     VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'txid for input',
  `input_vout`     MEDIUMINT(11) UNSIGNED NOT NULL COMMENT 'vout for input',
  `input_address`  VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'sender address for input',
  `input_account`  VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'sender account for input',
  `input_amount`   DECIMAL(26,10) NOT NULL COMMENT 'amount of coin to send for input',
  `input_confirmations` BIGINT(20) UNSIGNED NOT NULL COMMENT 'block confirmations when unspent rpc returned',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (`id`),
  INDEX idx_tx_id (`tx_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for input transaction';

-- Table: btc_tx_output
CREATE TABLE `btc_tx_output` (
  `id`             BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `tx_id`          BIGINT(20) NOT NULL COMMENT 'tx table ID',
  `output_address` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'receiver address for output',
  `output_account` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'receiver account for output',
  `output_amount`  DECIMAL(26,10) NOT NULL COMMENT 'amount of coin to receive',
  `is_change`      BOOL NOT NULL DEFAULT false COMMENT 'true: output is for fee',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (`id`),
  INDEX idx_tx_id (`tx_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for output transaction';

-- Table: tx
CREATE TABLE `tx` (
  `id`                  BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'transaction ID',
  `coin`                ENUM('eth','xrp','hyt') NOT NULL COMMENT 'coin type code',
  `action`              ENUM('deposit', 'payment', 'transfer') NOT NULL COMMENT 'action type',
  `updated_at`          datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (`id`),
  INDEX idx_coin (`coin`),
  INDEX idx_action (`action`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for eth transaction info';

-- Table: eth_detail_tx
CREATE TABLE `eth_detail_tx` (
  `id`               BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `tx_id`            BIGINT(20) NOT NULL COMMENT 'eth_tx table ID',
  `uuid`             VARCHAR(36) NOT NULL COMMENT 'UUID',
  `current_tx_type`  tinyint(2) NOT NULL DEFAULT 1 COMMENT 'current transaction type',
  `sender_account`   VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'sender account',
  `sender_address`   VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'sender address',
  `receiver_account` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'receiver account',
  `receiver_address` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'receiver address',
  `amount`           BIGINT(20) UNSIGNED NOT NULL COMMENT 'amount of coin to receive',
  `fee`              BIGINT(20) UNSIGNED NOT NULL COMMENT 'fee',
  `gas_limit`        MEDIUMINT(11) UNSIGNED NOT NULL COMMENT 'gas limit',
  `nonce`            BIGINT(20) UNSIGNED NOT NULL COMMENT 'nonce',
  `unsigned_hex_tx`  TEXT COLLATE utf8_unicode_ci NOT NULL COMMENT 'HEX string for unsigned transaction',
  `signed_hex_tx`    TEXT COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT 'HEX string for signed transaction',
  `sent_hash_tx`     TEXT COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT 'Hash for sent transaction',
  `unsigned_updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date for unsigned transaction created',
  `sent_updated_at`     datetime DEFAULT NULL COMMENT 'updated date for signed transaction sent',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uuid` (`uuid`),
  INDEX idx_txid (`tx_id`),
  INDEX idx_sender_account (`sender_account`),
  INDEX idx_receiver_account (`receiver_account`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for eth transaction detail';

-- Table: xrp_detail_tx
CREATE TABLE `xrp_detail_tx` (
  `id`                   BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `tx_id`                BIGINT(20) NOT NULL COMMENT 'xrp_tx table ID',
  `uuid`                 VARCHAR(36) NOT NULL COMMENT 'UUID',
  `current_tx_type`      tinyint(2) NOT NULL DEFAULT 1 COMMENT 'current transaction type',
  `sender_account`       VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'sender account',
  `sender_address`       VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'sender address',
  `receiver_account`     VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'receiver account',
  `receiver_address`     VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'receiver address',
  `amount`               VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'amount of coin to receive',
  `xrp_tx_type`          VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'xrp tx type like `Payment`',
  `fee`                  VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'tx fee',
  `flags`                BIGINT(20) UNSIGNED NOT NULL COMMENT 'tx flags',
  `last_ledger_sequence` BIGINT(20) UNSIGNED NOT NULL COMMENT 'tx LastLedgerSequence',
  `sequence`             BIGINT(20) UNSIGNED NOT NULL COMMENT 'tx Sequence',
  `signing_pubkey`       VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'tx SigningPubKey',
  `txn_signature`        VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'tx TxnSignature',
  `hash`                 VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'tx Hash',
  `earliest_ledger_version` BIGINT(20) UNSIGNED NOT NULL COMMENT 'tx earliest_ledger_version after sending tx',
  `signed_tx_id`         VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'signed tx id',
  `tx_blob`              TEXT COLLATE utf8_unicode_ci NOT NULL COMMENT 'sent tx blob',
  `sent_updated_at`      datetime DEFAULT NULL COMMENT 'updated date for signed transaction sent',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uuid` (`uuid`),
  INDEX idx_txid (`tx_id`),
  INDEX idx_sender_account (`sender_account`),
  INDEX idx_receiver_account (`receiver_account`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for xrp transaction detail';

-- Table: address
CREATE TABLE `address` (
  `id`                BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `coin`              ENUM('btc', 'bch', 'eth', 'xrp', 'hyt') NOT NULL COMMENT 'coin type code',
  `account`           ENUM('client', 'deposit', 'payment', 'stored') NOT NULL COMMENT 'account type',
  `wallet_address`    VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'wallet address',
  `is_allocated`      BOOL NOT NULL DEFAULT false COMMENT 'true: address is allocated(used)',
  `updated_at`        datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_wallet_address` (`wallet_address`),
  INDEX idx_coin (`coin`),
  INDEX idx_account (`account`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for account pubkey';

-- Table: payment_request
CREATE TABLE `payment_request` (
  `id`                BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `coin`              ENUM('btc', 'bch', 'eth', 'xrp') NOT NULL COMMENT 'coin type code',
  `payment_id`        BIGINT(20) DEFAULT NULL COMMENT 'tx table ID for payment action',
  `sender_address`    VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'sender address',
  `sender_account`    VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'sender account',
  `receiver_address`  VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'receiver address',
  `amount`            DECIMAL(26,10) NOT NULL COMMENT 'amount of coin to send',
  `is_done`           BOOL NOT NULL DEFAULT false COMMENT 'true: unsigned transaction is created',
  `updated_at`        datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (`id`),
  INDEX idx_coin (`coin`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for payment request';

