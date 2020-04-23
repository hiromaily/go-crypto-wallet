--
-- DATABASE wallet
--
DROP DATABASE IF EXISTS `wallet`;

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `wallet` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `wallet`;

-- wallet table definition
source /sqls/definition_wallet.sql
source /sqls/payment_request.sql
-- keygen table definition for sqlboiler
source /sqls/definition_keygen.sql

