-- MySQL dump 10.14  Distrib 5.7.28, for osx10.14 (x86_64)
--
-- Host: 127.0.0.1    Database: keygen
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
-- Table structure for table `seed`
--

DROP TABLE IF EXISTS `seed`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `seed` (
  `id`         tinyint(1) NOT NULL AUTO_INCREMENT COMMENT'ID',
  `coin`       ENUM('btc', 'bch') COMMENT'coin type code',
  `seed`       VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'seed',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date',
  PRIMARY KEY (`id`),
  INDEX idx_coin (`coin`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for seed';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `account_key`
--

DROP TABLE IF EXISTS `account_key`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_key` (
  `id`                      BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT'ID',
  `coin`                    ENUM('btc', 'bch') COMMENT'coin type code',
  `account`                 ENUM('client', 'receipt', 'payment', 'stored', 'fee', 'auth') COMMENT'account type',
  `wallet_address`          VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'address as standard pubkey script that Pays To PubKey Hash (P2PKH)',
  `p2sh_segwit_address`     VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'p2sh-segwit address',
  `full_public_key`         VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'full public key',
  `wallet_multisig_address` VARCHAR(255) COLLATE utf8_unicode_ci DEFAULT '' NOT NULL COMMENT'multisig address',
  `redeem_script`           VARCHAR(255) COLLATE utf8_unicode_ci DEFAULT '' NOT NULL COMMENT'redeedScript after multisig address generated',
  `wallet_import_format`    VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'WIF',
  `idx`                     BIGINT(20) UNSIGNED NOT NULL COMMENT'index for hd wallet',
  `addr_status`             tinyint(1) UNSIGNED DEFAULT 0 NOT NULL COMMENT'progress status for address generating',
  `updated_at`              datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_wallet_address` (`wallet_address`),
  UNIQUE KEY `idx_p2sh_segwit_address` (`p2sh_segwit_address`),
  UNIQUE KEY `idx_wallet_import_format` (`wallet_import_format`),
  INDEX idx_coin (`coin`),
  INDEX idx_account (`account`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for keys for any account';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `multisig_pubkey_history`
--  for sigunature database

DROP TABLE IF EXISTS `multisig_pubkey_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `multisig_pubkey_history` (
  `id`                      BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT'ID',
  `coin`                    ENUM('btc', 'bch') COMMENT'coin type code',
  `account`                 ENUM('receipt', 'payment', 'stored', 'fee') COMMENT'multisig account type',
  `full_public_key`         VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'full public key',
  `auth_address1`           VARCHAR(255) COLLATE utf8_unicode_ci DEFAULT '' NOT NULL COMMENT'1st address for authorization',
  `auth_address2`           VARCHAR(255) COLLATE utf8_unicode_ci DEFAULT '' NOT NULL COMMENT'2nd address for authorization',
  `wallet_multisig_address` VARCHAR(255) COLLATE utf8_unicode_ci DEFAULT '' NOT NULL COMMENT'multisig address',
  `redeem_script`           VARCHAR(255) COLLATE utf8_unicode_ci DEFAULT '' NOT NULL COMMENT'redeedScript after multisig address generated',
  `is_exported`             BOOL DEFAULT false COMMENT'true: csv file exported',
  `updated_at`              datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_full_public_key` (`full_public_key`),
  INDEX idx_coin (`coin`),
  INDEX idx_account (`account`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='multisig address info';
/*!40101 SET character_set_client = @saved_cs_client */;
