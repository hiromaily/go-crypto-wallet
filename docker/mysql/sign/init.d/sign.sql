--
-- DATABASE signature
--
DROP DATABASE IF EXISTS `sign`;

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `sign` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;

USE `sign`;

-- signature table definition
source /sqls/definition_sign.sql
