-- Create "address" table
CREATE TABLE `address` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT "ID",
  `coin` enum('btc','bch','eth','xrp','hyt') NOT NULL COMMENT "coin type code",
  `account` enum('client','deposit','payment','stored') NOT NULL COMMENT "account type",
  `wallet_address` varchar(255) NOT NULL COMMENT "wallet address",
  `is_allocated` bool NOT NULL DEFAULT 0 COMMENT "true: address is allocated(used)",
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date",
  PRIMARY KEY (`id`),
  INDEX `idx_account` (`account`),
  INDEX `idx_coin` (`coin`),
  UNIQUE INDEX `idx_wallet_address` (`wallet_address`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "table for account pubkey";
-- Create "btc_tx" table
CREATE TABLE `btc_tx` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT "transaction ID",
  `coin` enum('btc','bch') NOT NULL COMMENT "coin type code",
  `action` enum('deposit','payment','transfer') NOT NULL COMMENT "action type",
  `unsigned_hex_tx` text NOT NULL COMMENT "HEX string for unsigned transaction",
  `signed_hex_tx` text NOT NULL COMMENT "HEX string for signed transaction",
  `sent_hash_tx` text NOT NULL COMMENT "Hash for sent transaction",
  `total_input_amount` decimal(26,10) NOT NULL COMMENT "total amount of coin to send",
  `total_output_amount` decimal(26,10) NOT NULL COMMENT "total amount of coin to receive without fee",
  `fee` decimal(26,10) NOT NULL COMMENT "fee",
  `current_tx_type` tinyint NOT NULL DEFAULT 1 COMMENT "current transaction type",
  `unsigned_updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date for unsigned transaction created",
  `sent_updated_at` datetime NULL COMMENT "updated date for signed transaction sent",
  PRIMARY KEY (`id`),
  INDEX `idx_action` (`action`),
  INDEX `idx_coin` (`coin`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "table for btc transaction info";
-- Create "btc_tx_input" table
CREATE TABLE `btc_tx_input` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT "ID",
  `tx_id` bigint NOT NULL COMMENT "tx table ID",
  `input_txid` varchar(255) NOT NULL COMMENT "txid for input",
  `input_vout` mediumint unsigned NOT NULL COMMENT "vout for input",
  `input_address` varchar(255) NOT NULL COMMENT "sender address for input",
  `input_account` varchar(255) NOT NULL COMMENT "sender account for input",
  `input_amount` decimal(26,10) NOT NULL COMMENT "amount of coin to send for input",
  `input_confirmations` bigint unsigned NOT NULL COMMENT "block confirmations when unspent rpc returned",
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date",
  PRIMARY KEY (`id`),
  INDEX `idx_tx_id` (`tx_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "table for input transaction";
-- Create "btc_tx_output" table
CREATE TABLE `btc_tx_output` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT "ID",
  `tx_id` bigint NOT NULL COMMENT "tx table ID",
  `output_address` varchar(255) NOT NULL COMMENT "receiver address for output",
  `output_account` varchar(255) NOT NULL COMMENT "receiver account for output",
  `output_amount` decimal(26,10) NOT NULL COMMENT "amount of coin to receive",
  `is_change` bool NOT NULL DEFAULT 0 COMMENT "true: output is for fee",
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date",
  PRIMARY KEY (`id`),
  INDEX `idx_tx_id` (`tx_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "table for output transaction";
-- Create "eth_detail_tx" table
CREATE TABLE `eth_detail_tx` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT "ID",
  `tx_id` bigint NOT NULL COMMENT "eth_tx table ID",
  `uuid` varchar(36) NOT NULL COMMENT "UUID",
  `current_tx_type` tinyint NOT NULL DEFAULT 1 COMMENT "current transaction type",
  `sender_account` varchar(255) NOT NULL COMMENT "sender account",
  `sender_address` varchar(255) NOT NULL COMMENT "sender address",
  `receiver_account` varchar(255) NOT NULL COMMENT "receiver account",
  `receiver_address` varchar(255) NOT NULL COMMENT "receiver address",
  `amount` bigint unsigned NOT NULL COMMENT "amount of coin to receive",
  `fee` bigint unsigned NOT NULL COMMENT "fee",
  `gas_limit` mediumint unsigned NOT NULL COMMENT "gas limit",
  `nonce` bigint unsigned NOT NULL COMMENT "nonce",
  `unsigned_hex_tx` text NOT NULL COMMENT "HEX string for unsigned transaction",
  `signed_hex_tx` text NOT NULL COMMENT "HEX string for signed transaction",
  `sent_hash_tx` text NOT NULL COMMENT "Hash for sent transaction",
  `unsigned_updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date for unsigned transaction created",
  `sent_updated_at` datetime NULL COMMENT "updated date for signed transaction sent",
  PRIMARY KEY (`id`),
  INDEX `idx_receiver_account` (`receiver_account`),
  INDEX `idx_sender_account` (`sender_account`),
  INDEX `idx_txid` (`tx_id`),
  UNIQUE INDEX `idx_uuid` (`uuid`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "table for eth transaction detail";
-- Create "payment_request" table
CREATE TABLE `payment_request` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT "ID",
  `coin` enum('btc','bch','eth','xrp') NOT NULL COMMENT "coin type code",
  `payment_id` bigint NULL COMMENT "tx table ID for payment action",
  `sender_address` varchar(255) NOT NULL COMMENT "sender address",
  `sender_account` varchar(255) NOT NULL COMMENT "sender account",
  `receiver_address` varchar(255) NOT NULL COMMENT "receiver address",
  `amount` decimal(26,10) NOT NULL COMMENT "amount of coin to send",
  `is_done` bool NOT NULL DEFAULT 0 COMMENT "true: unsigned transaction is created",
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date",
  PRIMARY KEY (`id`),
  INDEX `idx_coin` (`coin`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "table for payment request";
-- Create "tx" table
CREATE TABLE `tx` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT "transaction ID",
  `coin` enum('eth','xrp','hyt') NOT NULL COMMENT "coin type code",
  `action` enum('deposit','payment','transfer') NOT NULL COMMENT "action type",
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date",
  PRIMARY KEY (`id`),
  INDEX `idx_action` (`action`),
  INDEX `idx_coin` (`coin`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "table for eth transaction info";
-- Create "xrp_detail_tx" table
CREATE TABLE `xrp_detail_tx` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT "ID",
  `tx_id` bigint NOT NULL COMMENT "xrp_tx table ID",
  `uuid` varchar(36) NOT NULL COMMENT "UUID",
  `current_tx_type` tinyint NOT NULL DEFAULT 1 COMMENT "current transaction type",
  `sender_account` varchar(255) NOT NULL COMMENT "sender account",
  `sender_address` varchar(255) NOT NULL COMMENT "sender address",
  `receiver_account` varchar(255) NOT NULL COMMENT "receiver account",
  `receiver_address` varchar(255) NOT NULL COMMENT "receiver address",
  `amount` varchar(255) NOT NULL COMMENT "amount of coin to receive",
  `xrp_tx_type` varchar(255) NOT NULL COMMENT "xrp tx type like `Payment`",
  `fee` varchar(255) NOT NULL COMMENT "tx fee",
  `flags` bigint unsigned NOT NULL COMMENT "tx flags",
  `last_ledger_sequence` bigint unsigned NOT NULL COMMENT "tx LastLedgerSequence",
  `sequence` bigint unsigned NOT NULL COMMENT "tx Sequence",
  `signing_pubkey` varchar(255) NOT NULL COMMENT "tx SigningPubKey",
  `txn_signature` varchar(255) NOT NULL COMMENT "tx TxnSignature",
  `hash` varchar(255) NOT NULL COMMENT "tx Hash",
  `earliest_ledger_version` bigint unsigned NOT NULL COMMENT "tx earliest_ledger_version after sending tx",
  `signed_tx_id` varchar(255) NOT NULL COMMENT "signed tx id",
  `tx_blob` text NOT NULL COMMENT "sent tx blob",
  `sent_updated_at` datetime NULL COMMENT "updated date for signed transaction sent",
  PRIMARY KEY (`id`),
  INDEX `idx_receiver_account` (`receiver_account`),
  INDEX `idx_sender_account` (`sender_account`),
  INDEX `idx_txid` (`tx_id`),
  UNIQUE INDEX `idx_uuid` (`uuid`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "table for xrp transaction detail";
