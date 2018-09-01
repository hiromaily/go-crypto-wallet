-- MySQL dump 10.13  Distrib 5.7.21, for osx10.13 (x86_64)
--
-- Host: 127.0.0.1    Database: cold_wallet1
-- ------------------------------------------------------
-- Server version	5.7.21

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
-- DATABASE cold_wallet1
--
DROP DATABASE IF EXISTS `cold_wallet1`;

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `cold_wallet1` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `cold_wallet1`;


--
-- Table structure for table `seed`
--

DROP TABLE IF EXISTS `seed`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `seed` (
  `id`         tinyint(1) NOT NULL AUTO_INCREMENT COMMENT'ID',
  `seed`       VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'seed情報',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'更新日時',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='Seed情報テーブル';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `key_type`
--

DROP TABLE IF EXISTS `key_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `key_type` (
  `id`           tinyint(1) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT'ID',
  `purpose`      tinyint(1) UNSIGNED NOT NULL COMMENT'purpose',
  `coin_type`    tinyint(1) UNSIGNED NOT NULL COMMENT'コインの種類(Mainnet, Testnet)',
  `account_type` tinyint(1) UNSIGNED NOT NULL COMMENT'利用目的',
  `change_type`  tinyint(1) UNSIGNED NOT NULL COMMENT'受け取り階層',
  `description`  VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'説明',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'更新日時',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='Key種別テーブル';
/*!40101 SET character_set_client = @saved_cs_client */;


/*まずはcoin_type:1 => Testnet用として作成する*/
LOCK TABLES `key_type` WRITE;
/*!40000 ALTER TABLE `key_type` DISABLE KEYS */;
INSERT INTO `key_type` VALUES
  (1,44,0,0,0,'顧客用 MainNet環境のためのkey',now()),
  (2,44,0,1,0,'入金保管用 MainNet環境のためのkey',now()),
  (3,44,0,2,0,'支払い用 MainNet環境のためのkey',now()),
  (4,44,0,3,0,'承認用 MainNet環境のためのkey',now()),
  (5,44,1,0,0,'顧客用 TestNet環境のためのkey',now()),
  (6,44,1,1,0,'入金保管用 TestNet環境のためのkey',now()),
  (7,44,1,2,0,'支払い用 TestNet環境のためのkey',now()),
  (8,44,1,3,0,'承認用 TestNet環境のためのkey',now());
/*!40000 ALTER TABLE `key_type` ENABLE KEYS */;
UNLOCK TABLES;


--
-- Table structure for table `account_key_client`
--

DROP TABLE IF EXISTS `account_key_client`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_key_client` (
  `id`                      BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT'ID',
  `wallet_address`          VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'Walletアドレス',
  `p2sh_segwit_address`     VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'p2sh-segwitアドレス',
  `full_public_key`         VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'full public key',
  `wallet_multisig_address` VARCHAR(255) COLLATE utf8_unicode_ci DEFAULT '' NOT NULL COMMENT'multisigとしてのWalletアドレス',
  `redeem_script`           VARCHAR(255) COLLATE utf8_unicode_ci DEFAULT '' NOT NULL COMMENT'multisigアドレス生成後に渡されるredeedScript',
  `wallet_import_format`    VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'WIF',
  `account`                 VARCHAR(255) COLLATE utf8_unicode_ci DEFAULT '' NOT NULL COMMENT'アドレスに紐づくアカウント名',
  `key_type`                tinyint(1) UNSIGNED NOT NULL COMMENT'コインの種類',
  `idx`                     BIGINT(20) UNSIGNED NOT NULL COMMENT'HDウォレット生成時のindex',
  `key_status`              tinyint(1) UNSIGNED DEFAULT 0 NOT NULL COMMENT'keyの進捗ステータス',
  `updated_at`              datetime DEFAULT CURRENT_TIMESTAMP COMMENT'更新日時',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_wallet_address` (`wallet_address`),
  UNIQUE KEY `idx_p2sh_segwit_address` (`p2sh_segwit_address`),
  UNIQUE KEY `idx_wallet_import_format` (`wallet_import_format`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='受け取り用トランザクション情報Table';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `account_key_receipt`
--

DROP TABLE IF EXISTS `account_key_receipt`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_key_receipt` LIKE `account_key_client`;


--
-- Table structure for table `account_key_payment`
--

DROP TABLE IF EXISTS `account_key_payment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_key_payment` LIKE `account_key_client`;


--
-- Table structure for table `account_key_authorization`
--

DROP TABLE IF EXISTS `account_key_authorization`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_key_authorization` LIKE `account_key_client`;


--
-- Table structure for table `added_pubkey_history_receipt`
--  coldwallet2用

DROP TABLE IF EXISTS `added_pubkey_history_receipt`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `added_pubkey_history_receipt` (
  `id`                      BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT'ID',
  /*`wallet_address`          VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'Walletアドレス',*/
  `full_public_key`         VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'full public key',
  `auth_address1`           VARCHAR(255) COLLATE utf8_unicode_ci DEFAULT '' NOT NULL COMMENT'認証用Walletアドレス1',
  `auth_address2`           VARCHAR(255) COLLATE utf8_unicode_ci DEFAULT '' NOT NULL COMMENT'認証用Walletアドレス2',
  `wallet_multisig_address` VARCHAR(255) COLLATE utf8_unicode_ci DEFAULT '' NOT NULL COMMENT'multisigとしてのWalletアドレス',
  `redeem_script`           VARCHAR(255) COLLATE utf8_unicode_ci DEFAULT '' NOT NULL COMMENT'multisigアドレス生成後に渡されるredeedScript',
  `is_exported`             BOOL DEFAULT false COMMENT'CSV出力済かどうか',
  `updated_at`              datetime DEFAULT CURRENT_TIMESTAMP COMMENT'更新日時',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_full_public_key` (`full_public_key`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='受け取り用multisigアドレス情報Table';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `added_pubkey_history_payment`
--  coldwallet2用

DROP TABLE IF EXISTS `added_pubkey_history_payment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `added_pubkey_history_payment` LIKE `added_pubkey_history_receipt`;
