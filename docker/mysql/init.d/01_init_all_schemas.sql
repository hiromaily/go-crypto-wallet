--
-- Consolidated Database Initialization Script
-- Creates all three schemas (watch, keygen, sign) in a single MySQL instance
-- and grants permissions to the application user
--
-- Note: This script only creates empty schemas. Table definitions are managed
-- by Atlas migrations. After the database container starts, run:
--   make atlas-migrate-docker
-- to apply schema migrations.
--

-- Create watch schema
CREATE DATABASE /*!32312 IF NOT EXISTS*/ `watch` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;

-- Create keygen schema
CREATE DATABASE /*!32312 IF NOT EXISTS*/ `keygen` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;

-- Create sign schema
CREATE DATABASE /*!32312 IF NOT EXISTS*/ `sign` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;

-- Grant permissions to application user (created by MYSQL_USER environment variable)
-- Note: User must be created first by MySQL's built-in initialization
-- Make sure by `docker compose exec wallet-db mysql -uhiromaily -phiromaily watch -e "SELECT DATABASE(), USER();"`
GRANT ALL PRIVILEGES ON `watch`.* TO 'hiromaily'@'%';
GRANT ALL PRIVILEGES ON `keygen`.* TO 'hiromaily'@'%';
GRANT ALL PRIVILEGES ON `sign`.* TO 'hiromaily'@'%';
FLUSH PRIVILEGES;
