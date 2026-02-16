-- Create "auth_account_key" table
CREATE TABLE "auth_account_key" (
  "id" smallint NOT NULL AUTO_INCREMENT COMMENT "ID",
  "coin" text NOT NULL CHECK ("coin" IN ('btc','bch')) COMMENT "coin type code",
  "key_type" varchar(20) NOT NULL DEFAULT "bip44" COMMENT "key type (bip44, bip49, bip84, bip86, musig2)",
  "auth_account" varchar(20) NOT NULL COMMENT "auth type",
  "account" varchar(20) NOT NULL DEFAULT "deposit" COMMENT "multisig account type (deposit, payment, stored)",
  "p2pkh_address" varchar(255) NOT NULL COMMENT "address as standard pubkey script that Pays To PubKey Hash (P2PKH)",
  "p2sh_segwit_address" varchar(255) NULL COMMENT "p2sh-segwit address",
  "bech32_address" varchar(255) NULL COMMENT "bech32 address",
  "taproot_address" varchar(255) NULL COMMENT "taproot address (BIP86)",
  "full_public_key" varchar(255) NOT NULL COMMENT "full public key",
  "multisig_address" varchar(255) NOT NULL DEFAULT "" COMMENT "multisig address",
  "redeem_script" varchar(255) NOT NULL DEFAULT "" COMMENT "redeedScript after multisig address generated",
  "wallet_import_format" varchar(255) NOT NULL COMMENT "WIF",
  "account_extended_privkey" varchar(255) NULL COMMENT "Account-level extended private key (xpriv) for BIP32 derivation",
  "idx" bigint NOT NULL COMMENT "index for hd wallet",
  "addr_status" smallint NOT NULL DEFAULT false COMMENT "progress status for address generating",
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date",
  PRIMARY KEY ("id"),
  UNIQUE INDEX "idex_coin_auth_account_account" ("coin", "auth_account", "account"),
  INDEX "idx_auth_account" ("auth_account"),
  UNIQUE INDEX "idx_bech32_address" ("bech32_address"),
  INDEX "idx_coin" ("coin"),
  INDEX "idx_key_type" ("key_type"),
  UNIQUE INDEX "idx_p2pkh_address" ("p2pkh_address"),
  UNIQUE INDEX "idx_p2sh_segwit_address" ("p2sh_segwit_address"),
  UNIQUE INDEX "idx_wallet_import_format" ("wallet_import_format")
) COMMENT "table for keys for auth account";
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
  "id" smallint NOT NULL AUTO_INCREMENT COMMENT "ID",
  "coin" text NOT NULL CHECK ("coin" IN ('btc','bch')) COMMENT "coin type code",
  "seed" varchar(255) NOT NULL COMMENT "seed",
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT "updated date",
  PRIMARY KEY ("id"),
  INDEX "idx_coin" ("coin")
) COMMENT "table for seed";
