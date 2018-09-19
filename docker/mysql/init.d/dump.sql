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
-- Table structure for table `tx_receipt`
--

DROP TABLE IF EXISTS `tx_receipt`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tx_receipt` (
  `id`     BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT'ID',
  `unsigned_hex_tx`     TEXT COLLATE utf8_unicode_ci NOT NULL COMMENT'未署名トランザクションのHex',
  `signed_hex_tx`       TEXT COLLATE utf8_unicode_ci DEFAULT NULL COMMENT'署名済トランザクションのHex',
  `sent_hash_tx`        TEXT COLLATE utf8_unicode_ci DEFAULT NULL COMMENT'送信済トランザクションのHash',
  `total_input_amount`  DECIMAL(26,10) NOT NULL COMMENT'送信される金額合計',
  `total_output_amount` DECIMAL(26,10) NOT NULL COMMENT'受信される金額合計(手数料は含まない',
  `fee`                 DECIMAL(26,10) NOT NULL COMMENT'手数料',
  `current_tx_type`     tinyint(1) NOT NULL COMMENT'現在のtx_type(ステータス)',
  `unsigned_updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'未署名トランザクション 更新日時',
  `signed_updated_at`   datetime DEFAULT NULL COMMENT'署名済トランザクション 更新日時',
  `sent_updated_at`     datetime DEFAULT NULL COMMENT'送信済トランザクション 更新日時',
  PRIMARY KEY (`id`)
  /*UNIQUE KEY `idx_unsigned_hex` (`unsigned_hex_tx`)*/
  /*INDEX idx_unsigned_hex (`unsigned_hex_tx(255)`),*/
  /*INDEX idx_signed_hex (`signed_hex_tx(255)`),*/
  /*INDEX idx_sent_hash (`sent_hash_tx(255)`)*/
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='受け取り用トランザクション情報Table';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `tx_receipt_input`
--

DROP TABLE IF EXISTS `tx_receipt_input`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tx_receipt_input` (
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
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='受け取り用トランザクションに紐づくinputトランザクション情報Table';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `tx_receipt_output`
--

DROP TABLE IF EXISTS `tx_receipt_output`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tx_receipt_output` (
  `id`             BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT'ID',
  `receipt_id`     BIGINT(20) UNSIGNED NOT NULL COMMENT'tx_receipt ID',
  `output_address` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'outputに利用されるaddress(受け取る人)',
  `output_account` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'outputに利用されるaccount(受け取る人)',
  `output_amount`  DECIMAL(26,10) NOT NULL COMMENT'outputに利用されるamount(入金金額)',
  `is_change`      BOOL DEFAULT false COMMENT'お釣り用のoutputであればtrue',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'更新日時',
  PRIMARY KEY (`id`),
  INDEX idx_receipt_id (`receipt_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='受け取り用トランザクション情報に紐づくoutput情報のTable';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `tx_payment`
--

DROP TABLE IF EXISTS `tx_payment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tx_payment` LIKE `tx_receipt`;


--
-- Table structure for table `tx_payment_input`
--

DROP TABLE IF EXISTS `tx_payment_input`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tx_payment_input` LIKE `tx_receipt_input`;


--
-- Table structure for table `tx_payment_output`
--

DROP TABLE IF EXISTS `tx_payment_output`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tx_payment_output` LIKE `tx_receipt_output`;


--
-- Table structure for table `tx_transfer`
--

DROP TABLE IF EXISTS `tx_transfer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tx_transfer` LIKE `tx_receipt`;


--
-- Table structure for table `tx_transfer_input`
--

DROP TABLE IF EXISTS `tx_transfer_input`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tx_transfer_input` LIKE `tx_receipt_input`;


--
-- Table structure for table `tx_transfer_output`
--

DROP TABLE IF EXISTS `tx_transfer_output`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tx_transfer_output` LIKE `tx_receipt_output`;


--
-- Table structure for table `payment_request`
--

DROP TABLE IF EXISTS `payment_request`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `payment_request` (
  `id`           BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT'ID',
  `payment_id`   BIGINT(20) UNSIGNED DEFAULT NULL COMMENT'tx_paymentのID',
  `address_from` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'送信者アドレス',
  `account_from` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'送信者アカウント名',
  `address_to`   VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'受け取り者アドレス',
  `amount`       DECIMAL(26,10) NOT NULL COMMENT'送金されるamount(出金金額)',
  `is_done`      BOOL DEFAULT false COMMENT'出金手続き(トランザクション作成)が完了済かどうか',
  `updated_at`   datetime DEFAULT CURRENT_TIMESTAMP COMMENT'更新日時',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='出金依頼情報のTable';


/*テストデータのため本番では削除すること*/
/*
LOCK TABLES `payment_request` WRITE;
INSERT INTO `payment_request` VALUES
  (1,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','yasui','2N33pRYgyuHn6K2xCrrq9dPzuW6ZAvFJfVz',0.1,false,now()),
  (2,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','yasui','2NFd6TEUgSpy8LvttBgVrLB6ZBA5X9BSUSz',0.2,false,now()),
  (3,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','yasui','2MucBdUqkP5XqNFVTCj35H6WQPC5u2a2BKV',0.25,false,now()),
  (4,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','yasui','2MucBdUqkP5XqNFVTCj35H6WQPC5u2a2BKV',0.3,false,now()),
  (5,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','yasui','2N7WsiDc4yK7PoUL9saGE5ZGsbRQ8R9NafS',0.4,false,now());
UNLOCK TABLES;
*/


--
-- Table structure for table `account_pubkey_client`
--

DROP TABLE IF EXISTS `account_pubkey_client`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_pubkey_client` (
  `id`     BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT'ID',
  `wallet_address`    VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'Walletアドレス',
  `account`           VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'アドレスに紐づくアカウント名',
  `is_allocated`      BOOL DEFAULT false COMMENT'アドレスが割り当てられたかどうか(未使用かどうか)',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'更新日時',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_wallet_address` (`wallet_address`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='顧客用Publicキー情報Table';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `account_pubkey_receipt`
--

DROP TABLE IF EXISTS `account_pubkey_receipt`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_pubkey_receipt` LIKE `account_pubkey_client`;


--
-- Table structure for table `account_pubkey_payment`
--

DROP TABLE IF EXISTS `account_pubkey_payment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_pubkey_payment` LIKE `account_pubkey_client`;


--
-- Table structure for table `account_pubkey_quoine`
--

DROP TABLE IF EXISTS `account_pubkey_quoine`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_pubkey_quoine` LIKE `account_pubkey_client`;


--
-- Table structure for table `account_pubkey_fee`
--

DROP TABLE IF EXISTS `account_pubkey_fee`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_pubkey_fee` LIKE `account_pubkey_client`;


--
-- Table structure for table `account_pubkey_stored`
--

DROP TABLE IF EXISTS `account_pubkey_stored`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_pubkey_stored` LIKE `account_pubkey_client`;
