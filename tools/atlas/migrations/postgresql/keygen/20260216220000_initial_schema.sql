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


-- Create "auth_fullpubkey" table
CREATE TABLE "auth_fullpubkey" (
  "id" smallserial NOT NULL,
  "coin" text NOT NULL CHECK ("coin" IN ('btc','bch')),
  "auth_account" varchar(20) NOT NULL,
  "purpose" smallint NOT NULL DEFAULT 49,
  "full_public_key" varchar(255) NOT NULL,
  "extended_pubkey" varchar(255) NULL,
  "fingerprint" varchar(8) NULL,
  "derivation_path" varchar(50) NULL,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "idx_auth_fullpubkey_coin_auth_account_purpose" ON "auth_fullpubkey" ("coin", "auth_account", "purpose");
CREATE INDEX "idx_auth_fullpubkey_coin" ON "auth_fullpubkey" ("coin");
CREATE INDEX "idx_auth_fullpubkey_fingerprint" ON "auth_fullpubkey" ("fingerprint");
CREATE INDEX "idx_auth_fullpubkey_purpose" ON "auth_fullpubkey" ("purpose");
-- Create "btc_account_key" table
CREATE TABLE "btc_account_key" (
  "id" bigserial NOT NULL,
  "coin" text NOT NULL CHECK ("coin" IN ('btc','bch')),
  "key_type" varchar(20) NOT NULL DEFAULT 'bip44',
  "account" text NOT NULL CHECK ("account" IN ('client','deposit','payment','stored')),
  "p2pkh_address" varchar(255) NOT NULL,
  "p2sh_segwit_address" varchar(255) NULL,
  "bech32_address" varchar(255) NULL,
  "taproot_address" varchar(255) NULL,
  "full_public_key" varchar(255) NOT NULL,
  "multisig_address" varchar(255) NOT NULL DEFAULT '',
  "redeem_script" varchar(1000) NOT NULL DEFAULT '',
  "wallet_import_format" varchar(255) NOT NULL,
  "account_extended_privkey" varchar(255) NULL,
  "idx" bigint NOT NULL,
  "addr_status" smallint NOT NULL DEFAULT 0,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_btc_account_key_account" ON "btc_account_key" ("account");
CREATE INDEX "idx_btc_account_key_coin" ON "btc_account_key" ("coin");
CREATE INDEX "idx_btc_account_key_key_type" ON "btc_account_key" ("key_type");
CREATE UNIQUE INDEX "idx_btc_account_key_p2pkh_address" ON "btc_account_key" ("p2pkh_address");
CREATE UNIQUE INDEX "idx_btc_account_key_wallet_import_format" ON "btc_account_key" ("wallet_import_format");
-- Create "eth_account_key" table
CREATE TABLE "eth_account_key" (
  "id" bigserial NOT NULL,
  "account" text NOT NULL CHECK ("account" IN ('client','deposit','payment','stored')),
  "address" varchar(42) NOT NULL,
  "full_public_key" varchar(130) NOT NULL,
  "private_key" varchar(64) NOT NULL,
  "idx" bigint NOT NULL,
  "addr_status" smallint NOT NULL DEFAULT 0,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_eth_account_key_account" ON "eth_account_key" ("account");
CREATE UNIQUE INDEX "idx_eth_account_key_address" ON "eth_account_key" ("address");
CREATE UNIQUE INDEX "idx_eth_account_key_private_key" ON "eth_account_key" ("private_key");
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
  "coin" text NOT NULL CHECK ("coin" IN ('btc','bch','eth','xrp','hyt')),
  "seed" varchar(255) NOT NULL,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_seed_coin" ON "seed" ("coin");
-- Create "xrp_account_key" table
CREATE TABLE "xrp_account_key" (
  "id" bigserial NOT NULL,
  "coin" text NOT NULL CHECK ("coin" IN ('xrp')),
  "account" text NOT NULL CHECK ("account" IN ('client','deposit','payment','stored')),
  "account_id" varchar(255) NOT NULL,
  "key_type" smallint NOT NULL DEFAULT 0,
  "master_key" varchar(255) NOT NULL,
  "master_seed" varchar(255) NOT NULL,
  "master_seed_hex" varchar(255) NOT NULL,
  "public_key" varchar(255) NOT NULL,
  "public_key_hex" varchar(255) NOT NULL,
  "is_regular_key_pair" boolean NOT NULL DEFAULT false,
  "allocated_id" bigint NOT NULL DEFAULT 0,
  "addr_status" smallint NOT NULL DEFAULT 0,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_xrp_account_key_account" ON "xrp_account_key" ("account");
CREATE UNIQUE INDEX "idx_xrp_account_key_account_id" ON "xrp_account_key" ("account_id");
CREATE INDEX "idx_xrp_account_key_coin" ON "xrp_account_key" ("coin");
CREATE UNIQUE INDEX "idx_xrp_account_key_master_seed" ON "xrp_account_key" ("master_seed");
-- Create "xrp_regular_key" table
CREATE TABLE "xrp_regular_key" (
  "id" bigserial NOT NULL,
  "account_id" varchar(255) NOT NULL,
  "regular_key_address" varchar(255) NOT NULL,
  "public_key" varchar(255) NOT NULL,
  "public_key_hex" varchar(255) NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "set_tx_hash" varchar(255) NULL,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "rotated_at" timestamp NULL,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_xrp_regular_key_account_active" ON "xrp_regular_key" ("account_id", "is_active");
CREATE INDEX "idx_xrp_regular_key_account_id" ON "xrp_regular_key" ("account_id");
CREATE INDEX "idx_xrp_regular_key_is_active" ON "xrp_regular_key" ("is_active");
CREATE UNIQUE INDEX "idx_xrp_regular_key_regular_key_address" ON "xrp_regular_key" ("regular_key_address");
-- Create "xrp_signer_entry" table
CREATE TABLE "xrp_signer_entry" (
  "id" bigserial NOT NULL,
  "signer_list_id" bigint NOT NULL,
  "signer_account" varchar(255) NOT NULL,
  "signer_weight" integer NOT NULL,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "idx_xrp_signer_entry_list_signer" ON "xrp_signer_entry" ("signer_list_id", "signer_account");
CREATE INDEX "idx_xrp_signer_entry_signer_account" ON "xrp_signer_entry" ("signer_account");
CREATE INDEX "idx_xrp_signer_entry_signer_list_id" ON "xrp_signer_entry" ("signer_list_id");
-- Create "xrp_signer_list" table
CREATE TABLE "xrp_signer_list" (
  "id" bigserial NOT NULL,
  "account_id" varchar(255) NOT NULL,
  "signer_quorum" integer NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "set_tx_hash" varchar(255) NULL,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "idx_xrp_signer_list_account_active" ON "xrp_signer_list" ("account_id", "is_active");
CREATE INDEX "idx_xrp_signer_list_account_id" ON "xrp_signer_list" ("account_id");
-- Comments
COMMENT ON TABLE "auth_fullpubkey" IS 'table for auth key exported from sign db';
COMMENT ON COLUMN "auth_fullpubkey"."id" IS 'ID';
COMMENT ON COLUMN "auth_fullpubkey"."coin" IS 'coin type code';
COMMENT ON COLUMN "auth_fullpubkey"."auth_account" IS 'auth type';
COMMENT ON COLUMN "auth_fullpubkey"."purpose" IS 'BIP purpose (44, 49, 84, 86) - default 49 for backward compatibility';
COMMENT ON COLUMN "auth_fullpubkey"."full_public_key" IS 'full public key (legacy: compressed pubkey, new: may be empty if using extended_pubkey)';
COMMENT ON COLUMN "auth_fullpubkey"."extended_pubkey" IS 'BIP32 extended public key (xpub/tpub format)';
COMMENT ON COLUMN "auth_fullpubkey"."fingerprint" IS 'BIP32 master key fingerprint (8 hex chars)';
COMMENT ON COLUMN "auth_fullpubkey"."derivation_path" IS 'BIP32 derivation path (e.g., m/49''/1''/0'')';
COMMENT ON COLUMN "auth_fullpubkey"."updated_at" IS 'updated date';
COMMENT ON TABLE "btc_account_key" IS 'table for BTC/BCH keys for any account';
COMMENT ON COLUMN "btc_account_key"."id" IS 'ID';
COMMENT ON COLUMN "btc_account_key"."coin" IS 'coin type code';
COMMENT ON COLUMN "btc_account_key"."key_type" IS 'key type (bip44, bip49, bip84, bip86, musig2)';
COMMENT ON COLUMN "btc_account_key"."account" IS 'account type';
COMMENT ON COLUMN "btc_account_key"."p2pkh_address" IS 'address as standard pubkey script that Pays To PubKey Hash (P2PKH)';
COMMENT ON COLUMN "btc_account_key"."p2sh_segwit_address" IS 'p2sh-segwit address';
COMMENT ON COLUMN "btc_account_key"."bech32_address" IS 'bech32 address';
COMMENT ON COLUMN "btc_account_key"."taproot_address" IS 'taproot address (BIP86)';
COMMENT ON COLUMN "btc_account_key"."full_public_key" IS 'full public key';
COMMENT ON COLUMN "btc_account_key"."multisig_address" IS 'multisig address';
COMMENT ON COLUMN "btc_account_key"."redeem_script" IS 'redeedScript after multisig address generated';
COMMENT ON COLUMN "btc_account_key"."wallet_import_format" IS 'WIF';
COMMENT ON COLUMN "btc_account_key"."account_extended_privkey" IS 'Account-level extended private key (xpriv) for BIP32 derivation';
COMMENT ON COLUMN "btc_account_key"."idx" IS 'index for hd wallet';
COMMENT ON COLUMN "btc_account_key"."addr_status" IS 'progress status for address generating';
COMMENT ON COLUMN "btc_account_key"."updated_at" IS 'updated date';
COMMENT ON TABLE "eth_account_key" IS 'table for ETH keys for any account';
COMMENT ON COLUMN "eth_account_key"."id" IS 'ID';
COMMENT ON COLUMN "eth_account_key"."account" IS 'account type';
COMMENT ON COLUMN "eth_account_key"."address" IS 'Ethereum address (0x...)';
COMMENT ON COLUMN "eth_account_key"."full_public_key" IS 'full public key (uncompressed, 65 bytes hex)';
COMMENT ON COLUMN "eth_account_key"."private_key" IS 'private key (hex encoded)';
COMMENT ON COLUMN "eth_account_key"."idx" IS 'index for hd wallet';
COMMENT ON COLUMN "eth_account_key"."addr_status" IS 'progress status for address generating';
COMMENT ON COLUMN "eth_account_key"."updated_at" IS 'updated date';
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
COMMENT ON TABLE "xrp_account_key" IS 'table for xrp keys for any account';
COMMENT ON COLUMN "xrp_account_key"."id" IS 'ID';
COMMENT ON COLUMN "xrp_account_key"."coin" IS 'coin type code';
COMMENT ON COLUMN "xrp_account_key"."account" IS 'account type';
COMMENT ON COLUMN "xrp_account_key"."account_id" IS 'account_id';
COMMENT ON COLUMN "xrp_account_key"."key_type" IS 'key_type';
COMMENT ON COLUMN "xrp_account_key"."master_key" IS 'master_key, DEPRECATED';
COMMENT ON COLUMN "xrp_account_key"."master_seed" IS 'master_seed';
COMMENT ON COLUMN "xrp_account_key"."master_seed_hex" IS 'master_seed_hex';
COMMENT ON COLUMN "xrp_account_key"."public_key" IS 'public_key';
COMMENT ON COLUMN "xrp_account_key"."public_key_hex" IS 'public_key_hex';
COMMENT ON COLUMN "xrp_account_key"."is_regular_key_pair" IS 'true: this key is for regular key pair';
COMMENT ON COLUMN "xrp_account_key"."allocated_id" IS 'index for hd wallet';
COMMENT ON COLUMN "xrp_account_key"."addr_status" IS 'progress status for address generating';
COMMENT ON COLUMN "xrp_account_key"."updated_at" IS 'updated date';
COMMENT ON TABLE "xrp_regular_key" IS 'table for XRP regular key assignments';
COMMENT ON COLUMN "xrp_regular_key"."id" IS 'ID';
COMMENT ON COLUMN "xrp_regular_key"."account_id" IS 'XRP account address (r...) that owns this regular key';
COMMENT ON COLUMN "xrp_regular_key"."regular_key_address" IS 'Regular key address (r...) authorized to sign for account';
COMMENT ON COLUMN "xrp_regular_key"."public_key" IS 'Regular key public key';
COMMENT ON COLUMN "xrp_regular_key"."public_key_hex" IS 'Regular key public key in hex format';
COMMENT ON COLUMN "xrp_regular_key"."is_active" IS 'true: this regular key is currently active for signing';
COMMENT ON COLUMN "xrp_regular_key"."set_tx_hash" IS 'Transaction hash of SetRegularKey that activated this key';
COMMENT ON COLUMN "xrp_regular_key"."created_at" IS 'creation date';
COMMENT ON COLUMN "xrp_regular_key"."rotated_at" IS 'date when this key was rotated out (set inactive)';
COMMENT ON TABLE "xrp_signer_entry" IS 'Individual signer entries within an XRP signer list';
COMMENT ON COLUMN "xrp_signer_entry"."id" IS 'ID';
COMMENT ON COLUMN "xrp_signer_entry"."signer_list_id" IS 'Reference to xrp_signer_list.id';
COMMENT ON COLUMN "xrp_signer_entry"."signer_account" IS 'XRP address of the authorized signer (r...)';
COMMENT ON COLUMN "xrp_signer_entry"."signer_weight" IS 'Weight of this signer (contributes to quorum)';
COMMENT ON COLUMN "xrp_signer_entry"."created_at" IS 'creation date';
COMMENT ON TABLE "xrp_signer_list" IS 'XRP signer list configuration for multi-signature accounts';
COMMENT ON COLUMN "xrp_signer_list"."id" IS 'ID';
COMMENT ON COLUMN "xrp_signer_list"."account_id" IS 'XRP account address (r...) that owns this signer list';
COMMENT ON COLUMN "xrp_signer_list"."signer_quorum" IS 'Minimum total weight of signatures required to authorize a transaction';
COMMENT ON COLUMN "xrp_signer_list"."is_active" IS 'true: this signer list is currently active on the ledger';
COMMENT ON COLUMN "xrp_signer_list"."set_tx_hash" IS 'Transaction hash of SignerListSet that created/updated this list';
COMMENT ON COLUMN "xrp_signer_list"."created_at" IS 'creation date';
COMMENT ON COLUMN "xrp_signer_list"."updated_at" IS 'last update date';
