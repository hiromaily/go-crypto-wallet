-- MySQL dump 10.14  Distrib 5.7.28, for osx10.14 (x86_64)
--
-- Host: 127.0.0.1    Database: sign
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
  `id`         tinyint(2) NOT NULL AUTO_INCREMENT COMMENT'ID',
  `coin`       ENUM('btc', 'bch') NOT NULL COMMENT'coin type code',
  `seed`       VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'seed',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date',
  PRIMARY KEY (`id`),
  INDEX idx_coin (`coin`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for seed';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `auth_account_key`
--

DROP TABLE IF EXISTS `auth_account_key`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `auth_account_key` (
  `id`                      SMALLINT(5) NOT NULL AUTO_INCREMENT COMMENT'ID',
  `coin`                    ENUM('btc', 'bch') NOT NULL COMMENT'coin type code',
  `auth_account`            VARCHAR(20)  COLLATE utf8_unicode_ci NOT NULL COMMENT'auth type',
  `p2pkh_address`           VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'address as standard pubkey script that Pays To PubKey Hash (P2PKH)',
  `p2sh_segwit_address`     VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'p2sh-segwit address',
  `full_public_key`         VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'full public key',
  `multisig_address`        VARCHAR(255) COLLATE utf8_unicode_ci DEFAULT '' NOT NULL COMMENT'multisig address',
  `redeem_script`           VARCHAR(255) COLLATE utf8_unicode_ci DEFAULT '' NOT NULL COMMENT'redeedScript after multisig address generated',
  `wallet_import_format`    VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'WIF',
  `idx`                     BIGINT(20) NOT NULL COMMENT'index for hd wallet',
  `addr_status`             tinyint(2) DEFAULT 0 NOT NULL COMMENT'progress status for address generating',
  `updated_at`              datetime DEFAULT CURRENT_TIMESTAMP COMMENT'updated date',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idex_coin_auth_account` (`coin`, `auth_account`),
  UNIQUE KEY `idx_p2pkh_address` (`p2pkh_address`),
  UNIQUE KEY `idx_p2sh_segwit_address` (`p2sh_segwit_address`),
  UNIQUE KEY `idx_wallet_import_format` (`wallet_import_format`),
  INDEX idx_coin (`coin`),
  INDEX idx_auth_account (`auth_account`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table for keys for auth account';
/*!40101 SET character_set_client = @saved_cs_client */;
