--
-- DATABASE wallet
--
DROP DATABASE IF EXISTS `watch`;

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `watch` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `watch`;

-- wallet table definition
source /sqls/definition_watch.sql
source /sqls/payment_request.sql
-- keygen table definition for sqlboiler, comment out later
source /sqls/definition_keygen.sql
-- sign table definition for sqlboiler, comment out later
source /sqls/definition_sign.sql
