-- Add account_extended_privkey column to eth_account_key for BIP32 offline signing
ALTER TABLE "eth_account_key" ADD COLUMN "account_extended_privkey" character varying(255);
