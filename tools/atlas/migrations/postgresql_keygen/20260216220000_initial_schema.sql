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
  "id" smallserial NOT NULL COMMENT "ID",
  "coin" text NOT NULL CHECK ("coin" IN ('btc','bch')) COMMENT "coin type code",
  "auth_account" varchar(20) NOT NULL COMMENT "auth type",
  "purpose" smallint NOT NULL DEFAULT 49 COMMENT "BIP purpose (44, 49, 84, 86) - default 49 for backward compatibility",
  "full_public_key" varchar(255) NOT NULL COMMENT "full public key (legacy: compressed pubkey, new: may be empty if using extended_pubkey)",
  "extended_pubkey" varchar(255) NULL COMMENT "BIP32 extended public key (xpub/tpub format)",
  "fingerprint" varchar(8) NULL COMMENT "BIP32 master key fingerprint (8 hex chars)",
  "derivation_path" varchar(50) NULL COMMENT "BIP32 derivation path (e.g., m/49'/1'/0')",
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date",
  PRIMARY KEY ("id"),
  UNIQUE INDEX "idex_coin_auth_account_purpose" ("coin", "auth_account", "purpose") COMMENT "unique constraint for coin, auth_account, and purpose combination",
  INDEX "idx_coin" ("coin"),
  INDEX "idx_fingerprint" ("fingerprint"),
  INDEX "idx_purpose" ("purpose")
) COMMENT "table for auth key exported from sign db";
-- Create "btc_account_key" table
CREATE TABLE "btc_account_key" (
  "id" bigserial NOT NULL COMMENT "ID",
  "coin" text NOT NULL CHECK ("coin" IN ('btc','bch')) COMMENT "coin type code",
  "key_type" varchar(20) NOT NULL DEFAULT "bip44" COMMENT "key type (bip44, bip49, bip84, bip86, musig2)",
  "account" text NOT NULL CHECK ("account" IN ('client','deposit','payment','stored')) COMMENT "account type",
  "p2pkh_address" varchar(255) NOT NULL COMMENT "address as standard pubkey script that Pays To PubKey Hash (P2PKH)",
  "p2sh_segwit_address" varchar(255) NULL COMMENT "p2sh-segwit address",
  "bech32_address" varchar(255) NULL COMMENT "bech32 address",
  "taproot_address" varchar(255) NULL COMMENT "taproot address (BIP86)",
  "full_public_key" varchar(255) NOT NULL COMMENT "full public key",
  "multisig_address" varchar(255) NOT NULL DEFAULT "" COMMENT "multisig address",
  "redeem_script" varchar(1000) NOT NULL DEFAULT "" COMMENT "redeedScript after multisig address generated",
  "wallet_import_format" varchar(255) NOT NULL COMMENT "WIF",
  "account_extended_privkey" varchar(255) NULL COMMENT "Account-level extended private key (xpriv) for BIP32 derivation",
  "idx" bigint NOT NULL COMMENT "index for hd wallet",
  "addr_status" smallint NOT NULL DEFAULT false COMMENT "progress status for address generating",
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date",
  PRIMARY KEY ("id"),
  INDEX "idx_account" ("account"),
  INDEX "idx_coin" ("coin"),
  INDEX "idx_key_type" ("key_type"),
  UNIQUE INDEX "idx_p2pkh_address" ("p2pkh_address"),
  UNIQUE INDEX "idx_wallet_import_format" ("wallet_import_format")
) COMMENT "table for BTC/BCH keys for any account";
-- Create "eth_account_key" table
CREATE TABLE "eth_account_key" (
  "id" bigserial NOT NULL COMMENT "ID",
  "account" text NOT NULL CHECK ("account" IN ('client','deposit','payment','stored')) COMMENT "account type",
  "address" varchar(42) NOT NULL COMMENT "Ethereum address (0x...)",
  "full_public_key" varchar(130) NOT NULL COMMENT "full public key (uncompressed, 65 bytes hex)",
  "private_key" varchar(64) NOT NULL COMMENT "private key (hex encoded)",
  "idx" bigint NOT NULL COMMENT "index for hd wallet",
  "addr_status" smallint NOT NULL DEFAULT false COMMENT "progress status for address generating",
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date",
  PRIMARY KEY ("id"),
  INDEX "idx_account" ("account"),
  UNIQUE INDEX "idx_address" ("address"),
  UNIQUE INDEX "idx_private_key" ("private_key")
) COMMENT "table for ETH keys for any account";
-- Create "musig2_nonces" table
CREATE TABLE "musig2_nonces" (
  "id" bigserial NOT NULL COMMENT "ID",
  "signer_id" varchar(255) NOT NULL COMMENT "Signer identifier",
  "transaction_id" varchar(255) NOT NULL COMMENT "Transaction identifier",
  "public_nonce" binary(66) NOT NULL COMMENT "Public nonce (66 bytes: two 33-byte compressed EC points R1||R2)",
  "is_used" boolean NOT NULL DEFAULT false COMMENT "true: nonce has been used in signing",
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT "creation date",
  "used_at" timestamp NULL COMMENT "date when nonce was marked as used",
  PRIMARY KEY ("id"),
  INDEX "idx_created_at" ("created_at"),
  INDEX "idx_is_used" ("is_used"),
  UNIQUE INDEX "idx_signer_transaction" ("signer_id", "transaction_id") COMMENT "Prevent duplicate nonces per signer-tx pair (CRITICAL for security)",
  INDEX "idx_transaction_id" ("transaction_id")
) COMMENT "MuSig2 nonce commitments for secure storage";
-- Create "seed" table
CREATE TABLE "seed" (
  "id" smallserial NOT NULL COMMENT "ID",
  "coin" text NOT NULL CHECK ("coin" IN ('btc','bch','eth','xrp','hyt')) COMMENT "coin type code",
  "seed" varchar(255) NOT NULL COMMENT "seed",
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date",
  PRIMARY KEY ("id"),
  INDEX "idx_coin" ("coin")
) COMMENT "table for seed";
-- Create "xrp_account_key" table
CREATE TABLE "xrp_account_key" (
  "id" bigserial NOT NULL COMMENT "ID",
  "coin" text NOT NULL CHECK ("coin" IN ('xrp')) COMMENT "coin type code",
  "account" text NOT NULL CHECK ("account" IN ('client','deposit','payment','stored')) COMMENT "account type",
  "account_id" varchar(255) NOT NULL COMMENT "account_id",
  "key_type" smallint NOT NULL DEFAULT false COMMENT "key_type",
  "master_key" varchar(255) NOT NULL COMMENT "master_key, DEPRECATED",
  "master_seed" varchar(255) NOT NULL COMMENT "master_seed",
  "master_seed_hex" varchar(255) NOT NULL COMMENT "master_seed_hex",
  "public_key" varchar(255) NOT NULL COMMENT "public_key",
  "public_key_hex" varchar(255) NOT NULL COMMENT "public_key_hex",
  "is_regular_key_pair" boolean NOT NULL DEFAULT false COMMENT "true: this key is for regular key pair",
  "allocated_id" bigint NOT NULL DEFAULT false COMMENT "index for hd wallet",
  "addr_status" smallint NOT NULL DEFAULT false COMMENT "progress status for address generating",
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date",
  PRIMARY KEY ("id"),
  INDEX "idx_account" ("account"),
  UNIQUE INDEX "idx_account_id" ("account_id"),
  INDEX "idx_coin" ("coin"),
  UNIQUE INDEX "idx_master_seed" ("master_seed")
) COMMENT "table for xrp keys for any account";
-- Create "xrp_regular_key" table
CREATE TABLE "xrp_regular_key" (
  "id" bigserial NOT NULL COMMENT "ID",
  "account_id" varchar(255) NOT NULL COMMENT "XRP account address (r...) that owns this regular key",
  "regular_key_address" varchar(255) NOT NULL COMMENT "Regular key address (r...) authorized to sign for account",
  "public_key" varchar(255) NOT NULL COMMENT "Regular key public key",
  "public_key_hex" varchar(255) NOT NULL COMMENT "Regular key public key in hex format",
  "is_active" boolean NOT NULL DEFAULT true COMMENT "true: this regular key is currently active for signing",
  "set_tx_hash" varchar(255) NULL COMMENT "Transaction hash of SetRegularKey that activated this key",
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT "creation date",
  "rotated_at" timestamp NULL COMMENT "date when this key was rotated out (set inactive)",
  PRIMARY KEY ("id"),
  INDEX "idx_account_active" ("account_id", "is_active") COMMENT "Composite index for finding active regular key by account",
  INDEX "idx_account_id" ("account_id"),
  INDEX "idx_is_active" ("is_active"),
  UNIQUE INDEX "idx_regular_key_address" ("regular_key_address")
) COMMENT "table for XRP regular key assignments";
-- Create "xrp_signer_entry" table
CREATE TABLE "xrp_signer_entry" (
  "id" bigserial NOT NULL COMMENT "ID",
  "signer_list_id" bigint NOT NULL COMMENT "Reference to xrp_signer_list.id",
  "signer_account" varchar(255) NOT NULL COMMENT "XRP address of the authorized signer (r...)",
  "signer_weight" integer NOT NULL COMMENT "Weight of this signer (contributes to quorum)",
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT "creation date",
  PRIMARY KEY ("id"),
  UNIQUE INDEX "idx_list_signer" ("signer_list_id", "signer_account") COMMENT "Each signer can only appear once per signer list",
  INDEX "idx_signer_account" ("signer_account"),
  INDEX "idx_signer_list_id" ("signer_list_id")
) COMMENT "Individual signer entries within an XRP signer list";
-- Create "xrp_signer_list" table
CREATE TABLE "xrp_signer_list" (
  "id" bigserial NOT NULL COMMENT "ID",
  "account_id" varchar(255) NOT NULL COMMENT "XRP account address (r...) that owns this signer list",
  "signer_quorum" integer NOT NULL COMMENT "Minimum total weight of signatures required to authorize a transaction",
  "is_active" boolean NOT NULL DEFAULT true COMMENT "true: this signer list is currently active on the ledger",
  "set_tx_hash" varchar(255) NULL COMMENT "Transaction hash of SignerListSet that created/updated this list",
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT "creation date",
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT "last update date",
  PRIMARY KEY ("id"),
  UNIQUE INDEX "idx_account_active" ("account_id", "is_active") COMMENT "Only one active signer list per account",
  INDEX "idx_account_id" ("account_id")
) COMMENT "XRP signer list configuration for multi-signature accounts";
