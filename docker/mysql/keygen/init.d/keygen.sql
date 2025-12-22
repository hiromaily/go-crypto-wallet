--
-- DATABASE keygen
--
DROP DATABASE IF EXISTS `keygen`;

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `keygen` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;

USE `keygen`;

-- keygen table definition
source /sqls/definition_keygen.sql
