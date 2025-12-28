--
-- Consolidated Database Initialization Script
-- Creates all three schemas (watch, keygen, sign) in a single MySQL instance
--

-- Create watch schema
DROP DATABASE IF EXISTS `watch`;
CREATE DATABASE /*!32312 IF NOT EXISTS*/ `watch` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;
USE `watch`;
-- watch table definitions
source /sqls/definition_watch.sql;
source /sqls/payment_request.sql;

-- Create keygen schema
DROP DATABASE IF EXISTS `keygen`;
CREATE DATABASE /*!32312 IF NOT EXISTS*/ `keygen` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;
USE `keygen`;
-- keygen table definitions
source /sqls/definition_keygen.sql;

-- Create sign schema
DROP DATABASE IF EXISTS `sign`;
CREATE DATABASE /*!32312 IF NOT EXISTS*/ `sign` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;
USE `sign`;
-- sign table definitions
source /sqls/definition_sign.sql;
