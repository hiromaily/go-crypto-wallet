-- Create "address" table
CREATE TABLE "address" (
  "id" bigserial NOT NULL,
  "coin" text NOT NULL CHECK ("coin" IN ('btc','bch','eth','xrp','hyt')),
  "account" text NOT NULL CHECK ("account" IN ('client','deposit','payment','stored')),
  "wallet_address" varchar(255) NOT NULL,
  "is_allocated" boolean NOT NULL DEFAULT false,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_address_account" ON "address" ("account");
CREATE INDEX "idx_address_coin" ON "address" ("coin");
CREATE UNIQUE INDEX "idx_address_wallet_address" ON "address" ("wallet_address");
-- Create "btc_tx" table
CREATE TABLE "btc_tx" (
  "id" bigserial NOT NULL,
  "coin" text NOT NULL CHECK ("coin" IN ('btc','bch')),
  "action" text NOT NULL CHECK ("action" IN ('deposit','payment','transfer')),
  "unsigned_hex_tx" text NOT NULL,
  "signed_hex_tx" text NOT NULL,
  "sent_hash_tx" text NOT NULL,
  "total_input_amount" numeric(26,10) NOT NULL,
  "total_output_amount" numeric(26,10) NOT NULL,
  "fee" numeric(26,10) NOT NULL,
  "current_tx_type" smallint NOT NULL DEFAULT 1,
  "unsigned_updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "sent_updated_at" timestamp NULL,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_btc_tx_action" ON "btc_tx" ("action");
CREATE INDEX "idx_btc_tx_coin" ON "btc_tx" ("coin");
-- Create "btc_tx_input" table
CREATE TABLE "btc_tx_input" (
  "id" bigserial NOT NULL,
  "tx_id" bigint NOT NULL,
  "input_txid" varchar(255) NOT NULL,
  "input_vout" integer NOT NULL,
  "input_address" varchar(255) NOT NULL,
  "input_account" varchar(255) NOT NULL,
  "input_amount" numeric(26,10) NOT NULL,
  "input_confirmations" bigint NOT NULL,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_btc_tx_input_tx_id" ON "btc_tx_input" ("tx_id");
-- Create "btc_tx_output" table
CREATE TABLE "btc_tx_output" (
  "id" bigserial NOT NULL,
  "tx_id" bigint NOT NULL,
  "output_address" varchar(255) NOT NULL,
  "output_account" varchar(255) NOT NULL,
  "output_amount" numeric(26,10) NOT NULL,
  "is_change" boolean NOT NULL DEFAULT false,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_btc_tx_output_tx_id" ON "btc_tx_output" ("tx_id");
-- Create "eth_detail_tx" table
CREATE TABLE "eth_detail_tx" (
  "id" bigserial NOT NULL,
  "tx_id" bigint NOT NULL,
  "uuid" varchar(36) NOT NULL,
  "current_tx_type" smallint NOT NULL DEFAULT 1,
  "sender_account" varchar(255) NOT NULL,
  "sender_address" varchar(255) NOT NULL,
  "receiver_account" varchar(255) NOT NULL,
  "receiver_address" varchar(255) NOT NULL,
  "amount" bigint NOT NULL,
  "fee" bigint NOT NULL,
  "gas_limit" integer NOT NULL,
  "nonce" bigint NOT NULL,
  "unsigned_hex_tx" text NOT NULL,
  "signed_hex_tx" text NOT NULL,
  "sent_hash_tx" text NOT NULL,
  "unsigned_updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "sent_updated_at" timestamp NULL,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_eth_detail_tx_receiver_account" ON "eth_detail_tx" ("receiver_account");
CREATE INDEX "idx_eth_detail_tx_sender_account" ON "eth_detail_tx" ("sender_account");
CREATE INDEX "idx_eth_detail_tx_tx_id" ON "eth_detail_tx" ("tx_id");
CREATE UNIQUE INDEX "idx_eth_detail_tx_uuid" ON "eth_detail_tx" ("uuid");
-- Create "payment_request" table
CREATE TABLE "payment_request" (
  "id" bigserial NOT NULL,
  "coin" text NOT NULL CHECK ("coin" IN ('btc','bch','eth','xrp')),
  "payment_id" bigint NULL,
  "sender_address" varchar(255) NOT NULL,
  "sender_account" varchar(255) NOT NULL,
  "receiver_address" varchar(255) NOT NULL,
  "amount" numeric(26,10) NOT NULL,
  "is_done" boolean NOT NULL DEFAULT false,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_payment_request_coin" ON "payment_request" ("coin");
-- Create "tx" table
CREATE TABLE "tx" (
  "id" bigserial NOT NULL,
  "coin" text NOT NULL CHECK ("coin" IN ('eth','xrp','hyt')),
  "action" text NOT NULL CHECK ("action" IN ('deposit','payment','transfer')),
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_tx_action" ON "tx" ("action");
CREATE INDEX "idx_tx_coin" ON "tx" ("coin");
-- Create "xrp_detail_tx" table
CREATE TABLE "xrp_detail_tx" (
  "id" bigserial NOT NULL,
  "tx_id" bigint NOT NULL,
  "uuid" varchar(36) NOT NULL,
  "current_tx_type" smallint NOT NULL DEFAULT 1,
  "sender_account" varchar(255) NOT NULL,
  "sender_address" varchar(255) NOT NULL,
  "receiver_account" varchar(255) NOT NULL,
  "receiver_address" varchar(255) NOT NULL,
  "amount" varchar(255) NOT NULL,
  "xrp_tx_type" varchar(255) NOT NULL,
  "fee" varchar(255) NOT NULL,
  "flags" bigint NOT NULL,
  "last_ledger_sequence" bigint NOT NULL,
  "sequence" bigint NOT NULL,
  "signing_pubkey" varchar(255) NOT NULL,
  "txn_signature" varchar(255) NOT NULL,
  "hash" varchar(255) NOT NULL,
  "earliest_ledger_version" bigint NOT NULL,
  "signed_tx_id" varchar(255) NOT NULL,
  "tx_blob" text NOT NULL,
  "sent_updated_at" timestamp NULL,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_xrp_detail_tx_receiver_account" ON "xrp_detail_tx" ("receiver_account");
CREATE INDEX "idx_xrp_detail_tx_sender_account" ON "xrp_detail_tx" ("sender_account");
CREATE INDEX "idx_xrp_detail_tx_tx_id" ON "xrp_detail_tx" ("tx_id");
CREATE UNIQUE INDEX "idx_xrp_detail_tx_uuid" ON "xrp_detail_tx" ("uuid");
-- Create "xrp_multisig_signature" table
CREATE TABLE "xrp_multisig_signature" (
  "id" bigserial NOT NULL,
  "pending_multisig_id" bigint NOT NULL,
  "signer_account" varchar(255) NOT NULL,
  "signed_tx_blob" text NOT NULL,
  "signer_weight" integer NOT NULL,
  "signed_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_xrp_multisig_signature_pending_multisig_id" ON "xrp_multisig_signature" ("pending_multisig_id");
CREATE UNIQUE INDEX "idx_xrp_multisig_signature_pending_signer" ON "xrp_multisig_signature" ("pending_multisig_id", "signer_account");
CREATE INDEX "idx_xrp_multisig_signature_signer_account" ON "xrp_multisig_signature" ("signer_account");
-- Create "xrp_pending_multisig" table
CREATE TABLE "xrp_pending_multisig" (
  "id" bigserial NOT NULL,
  "tx_uuid" varchar(36) NOT NULL,
  "account_id" varchar(255) NOT NULL,
  "unsigned_tx_json" text NOT NULL,
  "xrp_tx_type" varchar(50) NOT NULL,
  "required_quorum" integer NOT NULL,
  "current_weight" integer NOT NULL DEFAULT 0,
  "status" text NOT NULL CHECK ("status" IN ('pending','ready','submitted','confirmed','failed','expired')) DEFAULT 'pending',
  "combined_tx_blob" text NULL,
  "submitted_tx_hash" varchar(255) NULL,
  "expires_at" timestamp NULL,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_xrp_pending_multisig_account_id" ON "xrp_pending_multisig" ("account_id");
CREATE INDEX "idx_xrp_pending_multisig_account_status" ON "xrp_pending_multisig" ("account_id", "status");
CREATE INDEX "idx_xrp_pending_multisig_status" ON "xrp_pending_multisig" ("status");
CREATE UNIQUE INDEX "idx_xrp_pending_multisig_tx_uuid" ON "xrp_pending_multisig" ("tx_uuid");
-- Comments
COMMENT ON TABLE "address" IS 'table for account pubkey';
COMMENT ON COLUMN "address"."id" IS 'ID';
COMMENT ON COLUMN "address"."coin" IS 'coin type code';
COMMENT ON COLUMN "address"."account" IS 'account type';
COMMENT ON COLUMN "address"."wallet_address" IS 'wallet address';
COMMENT ON COLUMN "address"."is_allocated" IS 'true: address is allocated(used)';
COMMENT ON COLUMN "address"."updated_at" IS 'updated date';
COMMENT ON TABLE "btc_tx" IS 'table for btc transaction info';
COMMENT ON COLUMN "btc_tx"."id" IS 'transaction ID';
COMMENT ON COLUMN "btc_tx"."coin" IS 'coin type code';
COMMENT ON COLUMN "btc_tx"."action" IS 'action type';
COMMENT ON COLUMN "btc_tx"."unsigned_hex_tx" IS 'HEX string for unsigned transaction';
COMMENT ON COLUMN "btc_tx"."signed_hex_tx" IS 'HEX string for signed transaction';
COMMENT ON COLUMN "btc_tx"."sent_hash_tx" IS 'Hash for sent transaction';
COMMENT ON COLUMN "btc_tx"."total_input_amount" IS 'total amount of coin to send';
COMMENT ON COLUMN "btc_tx"."total_output_amount" IS 'total amount of coin to receive without fee';
COMMENT ON COLUMN "btc_tx"."fee" IS 'fee';
COMMENT ON COLUMN "btc_tx"."current_tx_type" IS 'current transaction type';
COMMENT ON COLUMN "btc_tx"."unsigned_updated_at" IS 'updated date for unsigned transaction created';
COMMENT ON COLUMN "btc_tx"."sent_updated_at" IS 'updated date for signed transaction sent';
COMMENT ON TABLE "btc_tx_input" IS 'table for input transaction';
COMMENT ON COLUMN "btc_tx_input"."id" IS 'ID';
COMMENT ON COLUMN "btc_tx_input"."tx_id" IS 'tx table ID';
COMMENT ON COLUMN "btc_tx_input"."input_txid" IS 'txid for input';
COMMENT ON COLUMN "btc_tx_input"."input_vout" IS 'vout for input';
COMMENT ON COLUMN "btc_tx_input"."input_address" IS 'sender address for input';
COMMENT ON COLUMN "btc_tx_input"."input_account" IS 'sender account for input';
COMMENT ON COLUMN "btc_tx_input"."input_amount" IS 'amount of coin to send for input';
COMMENT ON COLUMN "btc_tx_input"."input_confirmations" IS 'block confirmations when unspent rpc returned';
COMMENT ON COLUMN "btc_tx_input"."updated_at" IS 'updated date';
COMMENT ON TABLE "btc_tx_output" IS 'table for output transaction';
COMMENT ON COLUMN "btc_tx_output"."id" IS 'ID';
COMMENT ON COLUMN "btc_tx_output"."tx_id" IS 'tx table ID';
COMMENT ON COLUMN "btc_tx_output"."output_address" IS 'receiver address for output';
COMMENT ON COLUMN "btc_tx_output"."output_account" IS 'receiver account for output';
COMMENT ON COLUMN "btc_tx_output"."output_amount" IS 'amount of coin to receive';
COMMENT ON COLUMN "btc_tx_output"."is_change" IS 'true: output is for fee';
COMMENT ON COLUMN "btc_tx_output"."updated_at" IS 'updated date';
COMMENT ON TABLE "eth_detail_tx" IS 'table for eth transaction detail';
COMMENT ON COLUMN "eth_detail_tx"."id" IS 'ID';
COMMENT ON COLUMN "eth_detail_tx"."tx_id" IS 'eth_tx table ID';
COMMENT ON COLUMN "eth_detail_tx"."uuid" IS 'UUID';
COMMENT ON COLUMN "eth_detail_tx"."current_tx_type" IS 'current transaction type';
COMMENT ON COLUMN "eth_detail_tx"."sender_account" IS 'sender account';
COMMENT ON COLUMN "eth_detail_tx"."sender_address" IS 'sender address';
COMMENT ON COLUMN "eth_detail_tx"."receiver_account" IS 'receiver account';
COMMENT ON COLUMN "eth_detail_tx"."receiver_address" IS 'receiver address';
COMMENT ON COLUMN "eth_detail_tx"."amount" IS 'amount of coin to receive';
COMMENT ON COLUMN "eth_detail_tx"."fee" IS 'fee';
COMMENT ON COLUMN "eth_detail_tx"."gas_limit" IS 'gas limit';
COMMENT ON COLUMN "eth_detail_tx"."nonce" IS 'nonce';
COMMENT ON COLUMN "eth_detail_tx"."unsigned_hex_tx" IS 'HEX string for unsigned transaction';
COMMENT ON COLUMN "eth_detail_tx"."signed_hex_tx" IS 'HEX string for signed transaction';
COMMENT ON COLUMN "eth_detail_tx"."sent_hash_tx" IS 'Hash for sent transaction';
COMMENT ON COLUMN "eth_detail_tx"."unsigned_updated_at" IS 'updated date for unsigned transaction created';
COMMENT ON COLUMN "eth_detail_tx"."sent_updated_at" IS 'updated date for signed transaction sent';
COMMENT ON TABLE "payment_request" IS 'table for payment request';
COMMENT ON COLUMN "payment_request"."id" IS 'ID';
COMMENT ON COLUMN "payment_request"."coin" IS 'coin type code';
COMMENT ON COLUMN "payment_request"."payment_id" IS 'tx table ID for payment action';
COMMENT ON COLUMN "payment_request"."sender_address" IS 'sender address';
COMMENT ON COLUMN "payment_request"."sender_account" IS 'sender account';
COMMENT ON COLUMN "payment_request"."receiver_address" IS 'receiver address';
COMMENT ON COLUMN "payment_request"."amount" IS 'amount of coin to send';
COMMENT ON COLUMN "payment_request"."is_done" IS 'true: unsigned transaction is created';
COMMENT ON COLUMN "payment_request"."updated_at" IS 'updated date';
COMMENT ON TABLE "tx" IS 'table for eth transaction info';
COMMENT ON COLUMN "tx"."id" IS 'transaction ID';
COMMENT ON COLUMN "tx"."coin" IS 'coin type code';
COMMENT ON COLUMN "tx"."action" IS 'action type';
COMMENT ON COLUMN "tx"."updated_at" IS 'updated date';
COMMENT ON TABLE "xrp_detail_tx" IS 'table for xrp transaction detail';
COMMENT ON COLUMN "xrp_detail_tx"."id" IS 'ID';
COMMENT ON COLUMN "xrp_detail_tx"."tx_id" IS 'xrp_tx table ID';
COMMENT ON COLUMN "xrp_detail_tx"."uuid" IS 'UUID';
COMMENT ON COLUMN "xrp_detail_tx"."current_tx_type" IS 'current transaction type';
COMMENT ON COLUMN "xrp_detail_tx"."sender_account" IS 'sender account';
COMMENT ON COLUMN "xrp_detail_tx"."sender_address" IS 'sender address';
COMMENT ON COLUMN "xrp_detail_tx"."receiver_account" IS 'receiver account';
COMMENT ON COLUMN "xrp_detail_tx"."receiver_address" IS 'receiver address';
COMMENT ON COLUMN "xrp_detail_tx"."amount" IS 'amount of coin to receive';
COMMENT ON COLUMN "xrp_detail_tx"."xrp_tx_type" IS 'xrp tx type like "Payment"';
COMMENT ON COLUMN "xrp_detail_tx"."fee" IS 'tx fee';
COMMENT ON COLUMN "xrp_detail_tx"."flags" IS 'tx flags';
COMMENT ON COLUMN "xrp_detail_tx"."last_ledger_sequence" IS 'tx LastLedgerSequence';
COMMENT ON COLUMN "xrp_detail_tx"."sequence" IS 'tx Sequence';
COMMENT ON COLUMN "xrp_detail_tx"."signing_pubkey" IS 'tx SigningPubKey';
COMMENT ON COLUMN "xrp_detail_tx"."txn_signature" IS 'tx TxnSignature';
COMMENT ON COLUMN "xrp_detail_tx"."hash" IS 'tx Hash';
COMMENT ON COLUMN "xrp_detail_tx"."earliest_ledger_version" IS 'tx earliest_ledger_version after sending tx';
COMMENT ON COLUMN "xrp_detail_tx"."signed_tx_id" IS 'signed tx id';
COMMENT ON COLUMN "xrp_detail_tx"."tx_blob" IS 'sent tx blob';
COMMENT ON COLUMN "xrp_detail_tx"."sent_updated_at" IS 'updated date for signed transaction sent';
COMMENT ON TABLE "xrp_multisig_signature" IS 'Collected signatures for XRP multi-signature transactions';
COMMENT ON COLUMN "xrp_multisig_signature"."id" IS 'ID';
COMMENT ON COLUMN "xrp_multisig_signature"."pending_multisig_id" IS 'Reference to xrp_pending_multisig.id';
COMMENT ON COLUMN "xrp_multisig_signature"."signer_account" IS 'XRP address of the signer (r...)';
COMMENT ON COLUMN "xrp_multisig_signature"."signed_tx_blob" IS 'Signed transaction blob from this signer';
COMMENT ON COLUMN "xrp_multisig_signature"."signer_weight" IS 'Weight of this signature';
COMMENT ON COLUMN "xrp_multisig_signature"."signed_at" IS 'date when signature was collected';
COMMENT ON TABLE "xrp_pending_multisig" IS 'Pending XRP multi-signature transactions awaiting signatures';
COMMENT ON COLUMN "xrp_pending_multisig"."id" IS 'ID';
COMMENT ON COLUMN "xrp_pending_multisig"."tx_uuid" IS 'Unique identifier for the transaction';
COMMENT ON COLUMN "xrp_pending_multisig"."account_id" IS 'XRP account address (r...) that will sign this transaction';
COMMENT ON COLUMN "xrp_pending_multisig"."unsigned_tx_json" IS 'JSON representation of the unsigned transaction';
COMMENT ON COLUMN "xrp_pending_multisig"."xrp_tx_type" IS 'XRP transaction type (Payment, SignerListSet, etc.)';
COMMENT ON COLUMN "xrp_pending_multisig"."required_quorum" IS 'Total signature weight required for this transaction';
COMMENT ON COLUMN "xrp_pending_multisig"."current_weight" IS 'Current total weight of collected signatures';
COMMENT ON COLUMN "xrp_pending_multisig"."status" IS 'Current status of the multi-sig transaction';
COMMENT ON COLUMN "xrp_pending_multisig"."combined_tx_blob" IS 'Combined signed transaction blob (set when status is ready)';
COMMENT ON COLUMN "xrp_pending_multisig"."submitted_tx_hash" IS 'Transaction hash after submission to ledger';
COMMENT ON COLUMN "xrp_pending_multisig"."expires_at" IS 'Expiration time for this pending transaction';
COMMENT ON COLUMN "xrp_pending_multisig"."created_at" IS 'creation date';
COMMENT ON COLUMN "xrp_pending_multisig"."updated_at" IS 'last update date';
