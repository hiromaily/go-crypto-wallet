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
-- Table structure for table `account_type`
--

DROP TABLE IF EXISTS `account_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_type` (
  `id`          tinyint(1) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT'ID',
  `type`        VARCHAR(20) COLLATE utf8_unicode_ci NOT NULL COMMENT'アカウント種別',
  `description` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'説明',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'更新日時',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='アカウント種別テーブル';
/*!40101 SET character_set_client = @saved_cs_client */;


LOCK TABLES `account_type` WRITE;
/*!40000 ALTER TABLE `account_type` DISABLE KEYS */;
INSERT INTO `account_type` VALUES
  (0,'client','顧客',now()),
  (1,'receipt','入金保管用',now()),
  (2,'payment','支払い用',now()),
  (3,'authorization','Multisigのための承認用',now());
/*!40000 ALTER TABLE `account_type` ENABLE KEYS */;
UNLOCK TABLES;


--
-- Table structure for table `coin_type`
--

DROP TABLE IF EXISTS `coin_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `coin_type` (
  `id`          tinyint(1) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT'ID',
  `type`        VARCHAR(20) COLLATE utf8_unicode_ci NOT NULL COMMENT'コイン種別',
  `description` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'説明',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'更新日時',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='コイン種別テーブル';
/*!40101 SET character_set_client = @saved_cs_client */;


LOCK TABLES `coin_type` WRITE;
/*!40000 ALTER TABLE `coin_type` DISABLE KEYS */;
INSERT INTO `coin_type` VALUES
  (0,'mainnet','Bitcoin Mainnet',now()),
  (1,'testnet3','Bitcoin Testnet3',now());
/*!40000 ALTER TABLE `coin_type` ENABLE KEYS */;
UNLOCK TABLES;


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
  (1,44,1,0,0,'顧客用 TestNet環境のためのkey',now()),
  (2,44,1,1,0,'入金保管用 TestNet環境のためのkey',now()),
  (3,44,1,2,0,'支払い用 TestNet環境のためのkey',now());
/*!40000 ALTER TABLE `key_type` ENABLE KEYS */;
UNLOCK TABLES;


--
-- Table structure for table `account_key_client`
--

DROP TABLE IF EXISTS `account_key_client`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_key_client` (
  `id`     BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT'ID',
  `wallet_address`          VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'Walletアドレス',
  `wallet_multisig_address` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'multisigとしてのWalletアドレス',
  `wallet_import_format`    VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'WIF',
  `account`                 VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'アドレスに紐づくアカウント名',
  `key_type`                tinyint(1) UNSIGNED NOT NULL COMMENT'コインの種類',
  `idx`    BIGINT(20) UNSIGNED NOT NULL COMMENT'HDウォレット生成時のindex',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'更新日時',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_wallet_address` (`wallet_address`),
  UNIQUE KEY `idx_wallet_multisig_address` (`wallet_multisig_address`),
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
CREATE TABLE `account_key_authorization` LIKE `account_key_authorization`;
