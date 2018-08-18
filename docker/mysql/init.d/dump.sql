-- MySQL dump 10.13  Distrib 5.7.21, for osx10.13 (x86_64)
--
-- Host: 127.0.0.1    Database: wallet
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
-- DATABASE wallet
--
DROP DATABASE IF EXISTS `wallet`;

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `wallet` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `wallet`;


--
-- Table structure for table `block_checked`
--

DROP TABLE IF EXISTS `block_checked`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `block_checked` (
  `id`         tinyint(1) NOT NULL AUTO_INCREMENT COMMENT'ID',
  `count`      int(11) NOT NULL COMMENT'現在のチェックしたブロック数(番号)',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'更新日時',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='ブロック情報テーブル';
/*!40101 SET character_set_client = @saved_cs_client */;


LOCK TABLES `block_checked` WRITE;
/*!40000 ALTER TABLE `block_checked` DISABLE KEYS */;
INSERT INTO `block_checked` VALUES
  (1,10000,now());
/*!40000 ALTER TABLE `block_checked` ENABLE KEYS */;
UNLOCK TABLES;


--
-- Table structure for table `tx_type`
--

DROP TABLE IF EXISTS `tx_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tx_type` (
  `id`         tinyint(1) NOT NULL AUTO_INCREMENT COMMENT'ID',
  `type`       VARCHAR(20) COLLATE utf8_unicode_ci NOT NULL COMMENT'トランザクション種別',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'更新日時',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='トランザクション種別テーブル';
/*!40101 SET character_set_client = @saved_cs_client */;


LOCK TABLES `tx_type` WRITE;
/*!40000 ALTER TABLE `tx_type` DISABLE KEYS */;
INSERT INTO `tx_type` VALUES
  (1,'unsigned',now()),
  (2,'signed',now()),
  (3,'sent',now()),
  (4,'done',now()),
  (5,'cancel',now());
/*!40000 ALTER TABLE `tx_type` ENABLE KEYS */;
UNLOCK TABLES;


--
-- Table structure for table `tx_receipt`
--

DROP TABLE IF EXISTS `tx_receipt`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tx_receipt` (
  `id`     BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT'ID',
  `unsigned_hex_tx`     TEXT COLLATE utf8_unicode_ci NOT NULL COMMENT'未署名トランザクションのHex',
  `signed_hex_tx`       TEXT COLLATE utf8_unicode_ci DEFAULT NULL COMMENT'署名済トランザクションのHex',
  `sent_hash_tx`         TEXT COLLATE utf8_unicode_ci DEFAULT NULL COMMENT'送信済トランザクションのHash',
  `total_amount`        DECIMAL(26,10) NOT NULL COMMENT'手数料込みの送信される金額(実際は手数料が引かれるので、これより少ないAmoutが送金される)',
  `fee`                 DECIMAL(26,10) NOT NULL COMMENT'手数料',
  `receiver_address`    VARCHAR(40) COLLATE utf8_unicode_ci NOT NULL COMMENT'受取先アドレス(固定だがlogも兼ねるので念の為保持する)',
  `current_tx_type`     tinyint(1) NOT NULL COMMENT'現在のtx_type(ステータス)',
  `unsigned_updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'未署名トランザクション 更新日時',
  `signed_updated_at`   datetime DEFAULT NULL COMMENT'署名済トランザクション 更新日時',
  `sent_updated_at`     datetime DEFAULT NULL COMMENT'送信済トランザクション 更新日時',
  PRIMARY KEY (`id`)
  /*INDEX idx_unsigned_hex (`unsigned_hex_tx(255)`),*/
  /*INDEX idx_signed_hex (`signed_hex_tx(255)`),*/
  /*INDEX idx_sent_hex (`sent_hex_tx(255)`)*/
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='受け取り用トランザクション情報Table';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `tx_receipt_detail`
--

DROP TABLE IF EXISTS `tx_receipt_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tx_receipt_detail` (
  `id`             BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT'ID',
  `receipt_id`     BIGINT(20) UNSIGNED NOT NULL COMMENT'tx_receipt ID',
  `input_txid`     VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'inputに利用されるtxid',
  `input_vout`     INT(11) NOT NULL COMMENT'inputに利用されるvout',
  `input_address`  VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'inputに利用されるaddress(入金した人)',
  `input_account`  VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'inputに利用されるaccount(入金した人)',
  `input_amount`   DECIMAL(26,10) NOT NULL COMMENT'inputに利用されるamount(入金金額)',
  `input_confirmations` BIGINT(20) UNSIGNED NOT NULL COMMENT'unspent取得時に確定していたconfirmations数',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'更新日時',
  PRIMARY KEY (`id`),
  INDEX idx_receipt_id (`receipt_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='受け取り用トランザクション情報詳細Table';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `tx_payment`
--

DROP TABLE IF EXISTS `tx_payment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tx_payment` LIKE `tx_receipt`;


--
-- Table structure for table `tx_payment`
--

DROP TABLE IF EXISTS `tx_payment_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tx_payment_detail` LIKE `tx_receipt_detail`;
