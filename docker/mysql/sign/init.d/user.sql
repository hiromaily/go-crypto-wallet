-- Create users if they don't exist (MySQL 8.4+ syntax)
-- MySQL 8.4 uses caching_sha2_password by default (mysql_native_password was removed)
CREATE USER IF NOT EXISTS 'root'@'%' IDENTIFIED BY 'root';
CREATE USER IF NOT EXISTS 'hiromaily'@'%' IDENTIFIED BY 'hiromaily';

-- Grant privileges
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON *.* TO 'hiromaily'@'%' WITH GRANT OPTION;

-- Apply changes
FLUSH PRIVILEGES;
