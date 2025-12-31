-- Modify "auth_fullpubkey" table
ALTER TABLE `auth_fullpubkey` ADD COLUMN `fingerprint` varchar(8) NULL COMMENT "BIP32 master key fingerprint (8 hex chars)" AFTER `full_public_key`;
-- Create index "idx_fingerprint" to table: "auth_fullpubkey"
CREATE INDEX `idx_fingerprint` ON `auth_fullpubkey` (`fingerprint`);
