-- Add account_extended_privkey column to eth_account_key for BIP32 offline signing
ALTER TABLE `eth_account_key` ADD COLUMN `account_extended_privkey` varchar(255) NULL COMMENT 'BIP32 account-level extended private key for offline signing derivation';
