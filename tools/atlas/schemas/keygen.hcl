# Keygen Schema Definition
# Key generation data (seeds, account keys, full public keys)

# Schema declaration (empty to avoid CREATE SCHEMA in migrations)
# Note: detailed attributes (charset, collation, comment) would generate CREATE SCHEMA statement
schema "keygen" {}

# Table: seed
table "seed" {
  schema  = schema.keygen
  comment = "table for seed"

  column "id" {
    type           = tinyint
    null           = false
    auto_increment = true
    comment        = "ID"
  }

  column "coin" {
    type    = enum("btc", "bch", "eth", "xrp", "hyt")
    null    = false
    comment = "coin type code"
  }

  column "seed" {
    type    = varchar(255)
    null    = false
    comment = "seed"
  }

  column "updated_at" {
    type    = datetime
    null    = true
    default = sql("CURRENT_TIMESTAMP")
    comment = "updated date"
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_coin" {
    columns = [column.coin]
  }
}

# Table: btc_account_key
# Bitcoin and Bitcoin Cash account keys with multiple address formats
table "btc_account_key" {
  schema  = schema.keygen
  comment = "table for BTC/BCH keys for any account"

  column "id" {
    type           = bigint
    null           = false
    auto_increment = true
    comment        = "ID"
  }

  column "coin" {
    type    = enum("btc", "bch")
    null    = false
    comment = "coin type code"
  }

  column "key_type" {
    type    = varchar(20)
    null    = false
    default = "bip44"
    comment = "key type (bip44, bip49, bip84, bip86, musig2)"
  }

  column "account" {
    type    = enum("client", "deposit", "payment", "stored")
    null    = false
    comment = "account type"
  }

  column "p2pkh_address" {
    type    = varchar(255)
    null    = false
    comment = "address as standard pubkey script that Pays To PubKey Hash (P2PKH)"
  }

  column "p2sh_segwit_address" {
    type    = varchar(255)
    null    = true
    comment = "p2sh-segwit address"
  }

  column "bech32_address" {
    type    = varchar(255)
    null    = true
    comment = "bech32 address"
  }

  column "taproot_address" {
    type    = varchar(255)
    null    = true
    comment = "taproot address (BIP86)"
  }

  column "full_public_key" {
    type    = varchar(255)
    null    = false
    comment = "full public key"
  }

  column "multisig_address" {
    type    = varchar(255)
    null    = false
    default = ""
    comment = "multisig address"
  }

  column "redeem_script" {
    type    = varchar(1000)
    null    = false
    default = ""
    comment = "redeedScript after multisig address generated"
  }

  column "wallet_import_format" {
    type    = varchar(255)
    null    = false
    comment = "WIF"
  }

  column "account_extended_privkey" {
    type    = varchar(255)
    null    = true
    comment = "Account-level extended private key (xpriv) for BIP32 derivation"
  }

  column "idx" {
    type    = bigint
    null    = false
    comment = "index for hd wallet"
  }

  column "addr_status" {
    type    = tinyint
    null    = false
    default = 0
    comment = "progress status for address generating"
  }

  column "updated_at" {
    type    = datetime
    null    = true
    default = sql("CURRENT_TIMESTAMP")
    comment = "updated date"
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_p2pkh_address" {
    unique  = true
    columns = [column.p2pkh_address]
  }

  index "idx_wallet_import_format" {
    unique  = true
    columns = [column.wallet_import_format]
  }

  index "idx_coin" {
    columns = [column.coin]
  }

  index "idx_key_type" {
    columns = [column.key_type]
  }

  index "idx_account" {
    columns = [column.account]
  }
}

# Table: eth_account_key
# Ethereum account keys (simplified structure)
table "eth_account_key" {
  schema  = schema.keygen
  comment = "table for ETH keys for any account"

  column "id" {
    type           = bigint
    null           = false
    auto_increment = true
    comment        = "ID"
  }

  column "account" {
    type    = enum("client", "deposit", "payment", "stored")
    null    = false
    comment = "account type"
  }

  column "address" {
    type    = varchar(42)
    null    = false
    comment = "Ethereum address (0x...)"
  }

  column "full_public_key" {
    type    = varchar(130)
    null    = false
    comment = "full public key (uncompressed, 65 bytes hex)"
  }

  column "private_key" {
    type    = varchar(64)
    null    = false
    comment = "private key (hex encoded)"
  }

  column "idx" {
    type    = bigint
    null    = false
    comment = "index for hd wallet"
  }

  column "addr_status" {
    type    = tinyint
    null    = false
    default = 0
    comment = "progress status for address generating"
  }

  column "updated_at" {
    type    = datetime
    null    = true
    default = sql("CURRENT_TIMESTAMP")
    comment = "updated date"
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_address" {
    unique  = true
    columns = [column.address]
  }

  index "idx_private_key" {
    unique  = true
    columns = [column.private_key]
  }

  index "idx_account" {
    columns = [column.account]
  }
}

# Table: xrp_account_key
table "xrp_account_key" {
  schema  = schema.keygen
  comment = "table for xrp keys for any account"

  column "id" {
    type           = bigint
    null           = false
    auto_increment = true
    comment        = "ID"
  }

  column "coin" {
    type    = enum("xrp")
    null    = false
    comment = "coin type code"
  }

  column "account" {
    type    = enum("client", "deposit", "payment", "stored")
    null    = false
    comment = "account type"
  }

  column "account_id" {
    type    = varchar(255)
    null    = false
    comment = "account_id"
  }

  column "key_type" {
    type    = tinyint
    null    = false
    default = 0
    comment = "key_type"
  }

  column "master_key" {
    type    = varchar(255)
    null    = false
    comment = "master_key, DEPRECATED"
  }

  column "master_seed" {
    type    = varchar(255)
    null    = false
    comment = "master_seed"
  }

  column "master_seed_hex" {
    type    = varchar(255)
    null    = false
    comment = "master_seed_hex"
  }

  column "public_key" {
    type    = varchar(255)
    null    = false
    comment = "public_key"
  }

  column "public_key_hex" {
    type    = varchar(255)
    null    = false
    comment = "public_key_hex"
  }

  column "is_regular_key_pair" {
    type    = boolean
    null    = false
    default = false
    comment = "true: this key is for regular key pair"
  }

  column "allocated_id" {
    type    = bigint
    null    = false
    default = 0
    comment = "index for hd wallet"
  }

  column "addr_status" {
    type    = tinyint
    null    = false
    default = 0
    comment = "progress status for address generating"
  }

  column "updated_at" {
    type    = datetime
    null    = true
    default = sql("CURRENT_TIMESTAMP")
    comment = "updated date"
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_account_id" {
    unique  = true
    columns = [column.account_id]
  }

  index "idx_master_seed" {
    unique  = true
    columns = [column.master_seed]
  }

  index "idx_coin" {
    columns = [column.coin]
  }

  index "idx_account" {
    columns = [column.account]
  }
}

# Table: xrp_regular_key
# Tracks regular key assignments to XRP accounts for enhanced security
# Regular keys allow accounts to sign transactions without using the master key
# Reference: https://xrpl.org/docs/concepts/accounts/cryptographic-keys#regular-key-pair
table "xrp_regular_key" {
  schema  = schema.keygen
  comment = "table for XRP regular key assignments"

  column "id" {
    type           = bigint
    null           = false
    auto_increment = true
    comment        = "ID"
  }

  column "account_id" {
    type    = varchar(255)
    null    = false
    comment = "XRP account address (r...) that owns this regular key"
  }

  column "regular_key_address" {
    type    = varchar(255)
    null    = false
    comment = "Regular key address (r...) authorized to sign for account"
  }

  column "public_key" {
    type    = varchar(255)
    null    = false
    comment = "Regular key public key"
  }

  column "public_key_hex" {
    type    = varchar(255)
    null    = false
    comment = "Regular key public key in hex format"
  }

  column "is_active" {
    type    = boolean
    null    = false
    default = true
    comment = "true: this regular key is currently active for signing"
  }

  column "set_tx_hash" {
    type    = varchar(255)
    null    = true
    comment = "Transaction hash of SetRegularKey that activated this key"
  }

  column "created_at" {
    type    = datetime
    null    = true
    default = sql("CURRENT_TIMESTAMP")
    comment = "creation date"
  }

  column "rotated_at" {
    type    = datetime
    null    = true
    comment = "date when this key was rotated out (set inactive)"
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_account_id" {
    columns = [column.account_id]
  }

  index "idx_regular_key_address" {
    unique  = true
    columns = [column.regular_key_address]
  }

  index "idx_is_active" {
    columns = [column.is_active]
  }

  index "idx_account_active" {
    columns = [column.account_id, column.is_active]
    comment = "Composite index for finding active regular key by account"
  }
}

# Table: auth_fullpubkey
table "auth_fullpubkey" {
  schema  = schema.keygen
  comment = "table for auth key exported from sign db"

  column "id" {
    type           = smallint
    null           = false
    auto_increment = true
    comment        = "ID"
  }

  column "coin" {
    type    = enum("btc", "bch")
    null    = false
    comment = "coin type code"
  }

  column "auth_account" {
    type    = varchar(20)
    null    = false
    comment = "auth type"
  }

  column "purpose" {
    type    = tinyint
    null    = false
    default = 49
    comment = "BIP purpose (44, 49, 84, 86) - default 49 for backward compatibility"
  }

  column "full_public_key" {
    type    = varchar(255)
    null    = false
    comment = "full public key (legacy: compressed pubkey, new: may be empty if using extended_pubkey)"
  }

  column "extended_pubkey" {
    type    = varchar(255)
    null    = true
    comment = "BIP32 extended public key (xpub/tpub format)"
  }

  column "fingerprint" {
    type    = varchar(8)
    null    = true
    comment = "BIP32 master key fingerprint (8 hex chars)"
  }

  column "derivation_path" {
    type    = varchar(50)
    null    = true
    comment = "BIP32 derivation path (e.g., m/49'/1'/0')"
  }

  column "updated_at" {
    type    = datetime
    null    = true
    default = sql("CURRENT_TIMESTAMP")
    comment = "updated date"
  }

  primary_key {
    columns = [column.id]
  }

  index "idex_coin_auth_account_purpose" {
    unique  = true
    columns = [column.coin, column.auth_account, column.purpose]
    comment = "unique constraint for coin, auth_account, and purpose combination"
  }

  index "idx_coin" {
    columns = [column.coin]
  }

  index "idx_fingerprint" {
    columns = [column.fingerprint]
  }

  index "idx_purpose" {
    columns = [column.purpose]
  }
}

# Table: musig2_nonces
# MuSig2 nonce commitments for secure storage
# [SHARED] This table is also defined in: sign.hcl
# When modifying, update all schemas for consistency!
table "musig2_nonces" {
  schema  = schema.keygen
  comment = "MuSig2 nonce commitments for secure storage"

  column "id" {
    type           = bigint
    null           = false
    auto_increment = true
    comment        = "ID"
  }

  column "signer_id" {
    type    = varchar(255)
    null    = false
    comment = "Signer identifier"
  }

  column "transaction_id" {
    type    = varchar(255)
    null    = false
    comment = "Transaction identifier"
  }

  column "public_nonce" {
    type    = binary(66)
    null    = false
    comment = "Public nonce (66 bytes: two 33-byte compressed EC points R1||R2)"
  }

  column "is_used" {
    type    = boolean
    null    = false
    default = false
    comment = "true: nonce has been used in signing"
  }

  column "created_at" {
    type    = datetime
    null    = true
    default = sql("CURRENT_TIMESTAMP")
    comment = "creation date"
  }

  column "used_at" {
    type    = datetime
    null    = true
    comment = "date when nonce was marked as used"
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_signer_transaction" {
    unique  = true
    columns = [column.signer_id, column.transaction_id]
    comment = "Prevent duplicate nonces per signer-tx pair (CRITICAL for security)"
  }

  index "idx_transaction_id" {
    columns = [column.transaction_id]
  }

  index "idx_is_used" {
    columns = [column.is_used]
  }

  index "idx_created_at" {
    columns = [column.created_at]
  }
}

