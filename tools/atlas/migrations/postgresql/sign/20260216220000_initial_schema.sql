-- SECURITY NOTE:
-- This schema stores cryptographic keys and seeds. In production deployments:
-- 1. Use database-level encryption (PostgreSQL Transparent Data Encryption)
-- 2. Consider external Key Management Systems (KMS) or Hardware Security Modules (HSM)
-- 3. Implement application-level encryption before database persistence
-- 4. Restrict database access with role-based permissions
-- 5. Enable audit logging for sensitive table access
--
-- The offline wallet architecture (keygen/sign separated from watch) provides
-- air-gapped security, but database-level protections remain critical.


-- Create "auth_account_key" table
CREATE TABLE "auth_account_key" (
  "id" smallserial NOT NULL,
  "coin" text NOT NULL CHECK ("coin" IN ('btc','bch')),
  "key_type" varchar(20) NOT NULL DEFAULT 'bip44',
  "auth_account" varchar(20) NOT NULL,
  "account" varchar(20) NOT NULL DEFAULT 'deposit',
  "p2pkh_address" varchar(255) NOT NULL,
  "p2sh_segwit_address" varchar(255) NULL,
  "bech32_address" varchar(255) NULL,
  "taproot_address" varchar(255) NULL,
  "full_public_key" varchar(255) NOT NULL,
  "multisig_address" varchar(255) NOT NULL DEFAULT '',
  "redeem_script" varchar(255) NOT NULL DEFAULT '',
  "wallet_import_format" varchar(255) NOT NULL,
  "account_extended_privkey" varchar(255) NULL,
  "idx" bigint NOT NULL,
  "addr_status" smallint NOT NULL DEFAULT 0,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "idx_auth_account_key_coin_auth_account_account" ON "auth_account_key" ("coin", "auth_account", "account");
CREATE INDEX "idx_auth_account_key_auth_account" ON "auth_account_key" ("auth_account");
CREATE UNIQUE INDEX "idx_auth_account_key_bech32_address" ON "auth_account_key" ("bech32_address");
CREATE INDEX "idx_auth_account_key_coin" ON "auth_account_key" ("coin");
CREATE INDEX "idx_auth_account_key_key_type" ON "auth_account_key" ("key_type");
CREATE UNIQUE INDEX "idx_auth_account_key_p2pkh_address" ON "auth_account_key" ("p2pkh_address");
CREATE UNIQUE INDEX "idx_auth_account_key_p2sh_segwit_address" ON "auth_account_key" ("p2sh_segwit_address");
CREATE UNIQUE INDEX "idx_auth_account_key_wallet_import_format" ON "auth_account_key" ("wallet_import_format");
-- Create "musig2_nonces" table
CREATE TABLE "musig2_nonces" (
  "id" bigserial NOT NULL,
  "signer_id" varchar(255) NOT NULL,
  "transaction_id" varchar(255) NOT NULL,
  "public_nonce" bytea NOT NULL,
  "is_used" boolean NOT NULL DEFAULT false,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "used_at" timestamp NULL,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_musig2_nonces_created_at" ON "musig2_nonces" ("created_at");
CREATE INDEX "idx_musig2_nonces_is_used" ON "musig2_nonces" ("is_used");
CREATE UNIQUE INDEX "idx_musig2_nonces_signer_transaction" ON "musig2_nonces" ("signer_id", "transaction_id");
CREATE INDEX "idx_musig2_nonces_transaction_id" ON "musig2_nonces" ("transaction_id");
-- Create "seed" table
CREATE TABLE "seed" (
  "id" smallserial NOT NULL,
  "coin" text NOT NULL CHECK ("coin" IN ('btc','bch')),
  "seed" varchar(255) NOT NULL,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_seed_coin" ON "seed" ("coin");
-- Comments
COMMENT ON TABLE "auth_account_key" IS 'table for keys for auth account';
COMMENT ON COLUMN "auth_account_key"."id" IS 'ID';
COMMENT ON COLUMN "auth_account_key"."coin" IS 'coin type code';
COMMENT ON COLUMN "auth_account_key"."key_type" IS 'key type (bip44, bip49, bip84, bip86, musig2)';
COMMENT ON COLUMN "auth_account_key"."auth_account" IS 'auth type';
COMMENT ON COLUMN "auth_account_key"."account" IS 'multisig account type (deposit, payment, stored)';
COMMENT ON COLUMN "auth_account_key"."p2pkh_address" IS 'address as standard pubkey script that Pays To PubKey Hash (P2PKH)';
COMMENT ON COLUMN "auth_account_key"."p2sh_segwit_address" IS 'p2sh-segwit address';
COMMENT ON COLUMN "auth_account_key"."bech32_address" IS 'bech32 address';
COMMENT ON COLUMN "auth_account_key"."taproot_address" IS 'taproot address (BIP86)';
COMMENT ON COLUMN "auth_account_key"."full_public_key" IS 'full public key';
COMMENT ON COLUMN "auth_account_key"."multisig_address" IS 'multisig address';
COMMENT ON COLUMN "auth_account_key"."redeem_script" IS 'redeedScript after multisig address generated';
COMMENT ON COLUMN "auth_account_key"."wallet_import_format" IS 'WIF';
COMMENT ON COLUMN "auth_account_key"."account_extended_privkey" IS 'Account-level extended private key (xpriv) for BIP32 derivation';
COMMENT ON COLUMN "auth_account_key"."idx" IS 'index for hd wallet';
COMMENT ON COLUMN "auth_account_key"."addr_status" IS 'progress status for address generating';
COMMENT ON COLUMN "auth_account_key"."updated_at" IS 'updated date';
COMMENT ON TABLE "musig2_nonces" IS 'MuSig2 nonce commitments for secure storage';
COMMENT ON COLUMN "musig2_nonces"."id" IS 'ID';
COMMENT ON COLUMN "musig2_nonces"."signer_id" IS 'Signer identifier';
COMMENT ON COLUMN "musig2_nonces"."transaction_id" IS 'Transaction identifier';
COMMENT ON COLUMN "musig2_nonces"."public_nonce" IS 'Public nonce (66 bytes: two 33-byte compressed EC points R1||R2)';
COMMENT ON COLUMN "musig2_nonces"."is_used" IS 'true: nonce has been used in signing';
COMMENT ON COLUMN "musig2_nonces"."created_at" IS 'creation date';
COMMENT ON COLUMN "musig2_nonces"."used_at" IS 'date when nonce was marked as used';
COMMENT ON TABLE "seed" IS 'table for seed';
COMMENT ON COLUMN "seed"."id" IS 'ID';
COMMENT ON COLUMN "seed"."coin" IS 'coin type code';
COMMENT ON COLUMN "seed"."seed" IS 'seed';
COMMENT ON COLUMN "seed"."updated_at" IS 'updated date';
