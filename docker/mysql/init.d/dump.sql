-- MySQL dump 10.13  Distrib 5.7.21, for osx10.13 (x86_64)
--
-- Host: 127.0.0.1    Database: hiromaily
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
  `id`         int(11) NOT NULL AUTO_INCREMENT COMMENT'ID',
  `count`      int(11) COLLATE utf8_unicode_ci NOT NULL COMMENT'現在のチェックしたブロック数',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='Checked Block Table';
/*!40101 SET character_set_client = @saved_cs_client */;


LOCK TABLES `block_checked` WRITE;
/*!40000 ALTER TABLE `block_checked` DISABLE KEYS */;
INSERT INTO `block_checked` VALUES
  (1,10000,now());
/*!40000 ALTER TABLE `block_checked` ENABLE KEYS */;
UNLOCK TABLES;


--
-- Table structure for table `transaction_unsigned`
--

DROP TABLE IF EXISTS `transaction_unsigned`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `transaction_unsigned` (
  `id`     int(11) NOT NULL AUTO_INCREMENT COMMENT'ID',
  `tx_hex` varchar(80) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT'CreatedRawTransactionによって算出されたHex',
  `total`  int(11) COLLATE utf8_unicode_ci NOT NULL COMMENT'手数料込みの送信される金額',
  `fee`    int(11) COLLATE utf8_unicode_ci NOT NULL COMMENT'手数料',
  `to_address` varchar(80) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT'受取先アドレス',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date',
  PRIMARY KEY (`id`),
  INDEX hex (`tx_hex`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='Unsigned Transaction Table';
/*!40101 SET character_set_client = @saved_cs_client */;

