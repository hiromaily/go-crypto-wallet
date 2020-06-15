-- MySQL dump 10.14  Distrib 5.7.28, for osx10.14 (x86_64)
--
-- Host: 127.0.0.1    Database: wallet
-- ------------------------------------------------------
-- Server version	5.7.28

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


--
-- Table structure for table `btc_tx`
--

DROP TABLE IF EXISTS `btc_tx`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `btc_tx` (
  /*`id`                  BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT'transaction ID',*/
  `id`                  BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT'transaction ID',
  `coin`                ENUM('btc', 'bch') NOT NULL COMMENT'coin type code',
  `action`              ENUM('deposit', 'payment', 'transfer') NOT NULL COMMENT'action type',
  `unsigned_hex_tx`     TEXT COLLATE utf8_unicode_ci NOT NULL COMMENT'HEX string for unsigned transaction',
  `signed_hex_tx`       TEXT COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT'HEX string for signed transaction',
  `sent_hash_tx`        TEXT COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT'Hash for sent transaction',
  `total_input_amount`  DECIMAL(26,10) NOT NULL COMMENT'total amount of coin to send',
  `total_output_amount` DECIMAL(26,10) NOT NULL COMMENT'total amount of coin to receive without fee',
  `fee`                 DECIMAL(26,10) NOT NULL COMMENT'fee',
  `current_tx_type`     tinyint(2) NOT NULL DEFAULT 1 COMMENT'current transaction type',
  `unsigned_updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date for unsigned transaction created',
  `sent_updated_at`     datetime DEFAULT NULL COMMENT'updated date for signed transaction sent',
  PRIMARY KEY (`id`),
  INDEX idx_coin (`coin`),
  INDEX idx_action (`action`)
  /*UNIQUE KEY `idx_unsigned_hex` (`unsigned_hex_tx`)*/
  /*INDEX idx_unsigned_hex (`unsigned_hex_tx(255)`),*/
  /*INDEX idx_signed_hex (`signed_hex_tx(255)`),*/
  /*INDEX idx_sent_hash (`sent_hash_tx(255)`)*/
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for btc transaction info';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `btc_tx_input`
--

DROP TABLE IF EXISTS `btc_tx_input`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `btc_tx_input` (
  `id`             BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT'ID',
  `tx_id`          BIGINT(20) NOT NULL COMMENT'tx table ID',
  `input_txid`     VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'txid for input',
  `input_vout`     MEDIUMINT(11) UNSIGNED NOT NULL COMMENT'vout for input',
  `input_address`  VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'sender address for input',
  `input_account`  VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'sender account for input',
  `input_amount`   DECIMAL(26,10) NOT NULL COMMENT'amount of coin to send for input',
  `input_confirmations` BIGINT(20) UNSIGNED NOT NULL COMMENT'block confirmations when unspent rpc returned',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date',
  PRIMARY KEY (`id`),
  INDEX idx_tx_id (`tx_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for input transaction';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `btc_tx_output`
--

DROP TABLE IF EXISTS `btc_tx_output`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `btc_tx_output` (
  `id`             BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT'ID',
  `tx_id`          BIGINT(20) NOT NULL COMMENT'tx table ID',
  `output_address` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'receiver address for output',
  `output_account` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'receiver account for output',
  `output_amount`  DECIMAL(26,10) NOT NULL COMMENT'amount of coin to receive',
  `is_change`      BOOL NOT NULL DEFAULT false COMMENT'true: output is for fee',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date',
  PRIMARY KEY (`id`),
  INDEX idx_tx_id (`tx_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for output transaction';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `tx`
--

DROP TABLE IF EXISTS `tx`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tx` (
  `id`                  BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT'transaction ID',
  `coin`                ENUM('eth','xrp') NOT NULL COMMENT'coin type code',
  `action`              ENUM('deposit', 'payment', 'transfer') NOT NULL COMMENT'action type',
  `updated_at`          datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date',
  PRIMARY KEY (`id`),
  INDEX idx_coin (`coin`),
  INDEX idx_action (`action`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for eth transaction info';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `eth_detail_tx`
--

DROP TABLE IF EXISTS `eth_detail_tx`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `eth_detail_tx` (
  `id`               BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT'ID',
  `tx_id`            BIGINT(20) NOT NULL COMMENT'eth_tx table ID',
  `uuid`             VARCHAR(36) NOT NULL COMMENT'UUID',
  `current_tx_type`  tinyint(2) NOT NULL DEFAULT 1 COMMENT'current transaction type',
  `sender_account`   VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'sender account',
  `sender_address`   VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'sender address',
  `receiver_account` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'receiver account',
  `receiver_address` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'receiver address',
  `amount`           BIGINT(20) UNSIGNED NOT NULL COMMENT'amount of coin to receive',
  `fee`              BIGINT(20) UNSIGNED NOT NULL COMMENT'fee',
  `gas_limit`        MEDIUMINT(11) UNSIGNED NOT NULL COMMENT'gas limit',
  `nonce`            BIGINT(20) UNSIGNED NOT NULL COMMENT'nonce',
  `unsigned_hex_tx`  TEXT COLLATE utf8_unicode_ci NOT NULL COMMENT'HEX string for unsigned transaction',
  `signed_hex_tx`    TEXT COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT'HEX string for signed transaction',
  `sent_hash_tx`     TEXT COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT'Hash for sent transaction',
  `unsigned_updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date for unsigned transaction created',
  `sent_updated_at`     datetime DEFAULT NULL COMMENT'updated date for signed transaction sent',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uuid` (`uuid`),
  INDEX idx_txid (`tx_id`),
  INDEX idx_sender_account (`sender_account`),
  INDEX idx_receiver_account (`receiver_account`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for eth transaction detail';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `xrp_detail_tx`
--

DROP TABLE IF EXISTS `xrp_detail_tx`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `xrp_detail_tx` (
  `id`                   BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT'ID',
  `tx_id`                BIGINT(20) NOT NULL COMMENT'xrp_tx table ID',
  `uuid`                 VARCHAR(36) NOT NULL COMMENT'UUID',
  `current_tx_type`      tinyint(2) NOT NULL DEFAULT 1 COMMENT'current transaction type',
  `sender_account`       VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'sender account',
  `sender_address`       VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'sender address',
  `receiver_account`     VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'receiver account',
  `receiver_address`     VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'receiver address',
  `amount`               VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'amount of coin to receive',
  `xrp_tx_type`          VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'xrp tx type like `Payment`',
  `fee`                  VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'tx fee',
  `flags`                BIGINT(20) UNSIGNED NOT NULL COMMENT'tx flags',
  `last_ledger_sequence` BIGINT(20) UNSIGNED NOT NULL COMMENT'tx LastLedgerSequence',
  `sequence`             BIGINT(20) UNSIGNED NOT NULL COMMENT'tx Sequence',
  `signing_pubkey`       VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'tx SigningPubKey',
  `txn_signature`        VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'tx TxnSignature',
  `hash`                 VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'tx Hash',
  `earliest_ledger_version` BIGINT(20) UNSIGNED NOT NULL COMMENT'tx earliest_ledger_version after sending tx',
  `signed_tx_id`         VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'signed tx id',
  `signed_tx_blob`       VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'signed tx blob',
  `sent_tx_blob`         TEXT COLLATE utf8_unicode_ci NOT NULL COMMENT'sent tx blob',
  `unsigned_updated_at`  datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date for unsigned transaction created',
  `sent_updated_at`      datetime DEFAULT NULL COMMENT'updated date for signed transaction sent',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uuid` (`uuid`),
  INDEX idx_txid (`tx_id`),
  INDEX idx_sender_account (`sender_account`),
  INDEX idx_receiver_account (`receiver_account`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for eth transaction detail';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `payment_request`
--

DROP TABLE IF EXISTS `payment_request`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `payment_request` (
  `id`                BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT'ID',
  `coin`              ENUM('btc', 'bch', 'eth', 'xrp') NOT NULL COMMENT'coin type code',
  `payment_id`        BIGINT(20) DEFAULT NULL COMMENT'tx table ID for payment action',
  `sender_address`    VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'sender address',
  `sender_account`    VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'sender account',
  `receiver_address`  VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'receiver address',
  `amount`            DECIMAL(26,10) NOT NULL COMMENT'amount of coin to send',
  `is_done`           BOOL NOT NULL DEFAULT false COMMENT'true: unsigned transaction is created',
  `updated_at`        datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date',
  PRIMARY KEY (`id`),
  INDEX idx_coin (`coin`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for payment request';


/*this is test data for development*/
/*
LOCK TABLES `payment_request` WRITE;
INSERT INTO `payment_request` VALUES
  (1,'btc',NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','tom1','2N33pRYgyuHn6K2xCrrq9dPzuW6ZAvFJfVz',0.001,false,now()),
  (2,'btc',NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','tom2','2NFd6TEUgSpy8LvttBgVrLB6ZBA5X9BSUSz',0.002,false,now()),
  (3,'btc',NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','tom3','2MucBdUqkP5XqNFVTCj35H6WQPC5u2a2BKV',0.0025,false,now()),
  (4,'btc',NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','tom4','2MucBdUqkP5XqNFVTCj35H6WQPC5u2a2BKV',0.0015,false,now()),
  (5,'btc',NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','tom5','2N7WsiDc4yK7PoUL9saGE5ZGsbRQ8R9NafS',0.0022,false,now());
UNLOCK TABLES;
*/

--
-- Table structure for table `address`
--

DROP TABLE IF EXISTS `address`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `address` (
  `id`                BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT'ID',
  `coin`              ENUM('btc', 'bch', 'eth', 'xrp') NOT NULL COMMENT'coin type code',
  `account`           ENUM('client', 'deposit', 'payment', 'stored') NOT NULL COMMENT'account type',
  `wallet_address`    VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'wallet address',
  `is_allocated`      BOOL NOT NULL DEFAULT false COMMENT'true: address is allocated(used)',
  `updated_at`        datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_wallet_address` (`wallet_address`),
  INDEX idx_coin (`coin`),
  INDEX idx_account (`account`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for account pubkey';
/*!40101 SET character_set_client = @saved_cs_client */;
