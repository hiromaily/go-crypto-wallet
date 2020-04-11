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
-- DATABASE signature
--
DROP DATABASE IF EXISTS `signature`;

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `signature` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `signature`;


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
-- Table structure for table `account_key_quoine`
--

DROP TABLE IF EXISTS `account_key_quoine`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_key_quoine` LIKE `account_key_client`;


--
-- Table structure for table `account_key_fee`
--

DROP TABLE IF EXISTS `account_key_fee`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_key_fee` LIKE `account_key_client`;


--
-- Table structure for table `account_key_stored`
--

DROP TABLE IF EXISTS `account_key_stored`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_key_stored` LIKE `account_key_client`;


--
-- Table structure for table `account_key_authorization`
--  for sigunature database

DROP TABLE IF EXISTS `account_key_authorization`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_key_authorization` LIKE `account_key_client`;


--
-- Table structure for table `added_pubkey_history_receipt`
--  for sigunature database

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
--  for sigunature database

DROP TABLE IF EXISTS `added_pubkey_history_payment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `added_pubkey_history_payment` LIKE `added_pubkey_history_receipt`;


--
-- Table structure for table `added_pubkey_history_quoine`
--  for sigunature database

DROP TABLE IF EXISTS `added_pubkey_history_quoine`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `added_pubkey_history_quoine` LIKE `added_pubkey_history_receipt`;


--
-- Table structure for table `added_pubkey_history_fee`
--  for sigunature database

DROP TABLE IF EXISTS `added_pubkey_history_fee`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `added_pubkey_history_fee` LIKE `added_pubkey_history_receipt`;


--
-- Table structure for table `added_pubkey_history_stored`
--  for sigunature database

DROP TABLE IF EXISTS `added_pubkey_history_stored`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `added_pubkey_history_stored` LIKE `added_pubkey_history_receipt`;
