# Sign Schema Definition
# Signing wallet data (auth account keys, seeds)
#
# Design Decision: varchar + CHECK constraints instead of native PostgreSQL enums
# - Better schema evolution flexibility (adding/removing values without ALTER TYPE locks)
# - sqlc generates proper `string` types in Go models (native enums produce `interface{}`)
# - Consistent with project design docs (docs/database/research.md)

# Schema declaration (PostgreSQL uses "public" schema within each database)
schema "public" {}

# Table: seed
table "seed" {
  schema  = schema.public
  comment = "table for seed"

  column "id" {
    type    = smallint
    null    = false
    comment = "ID"
    identity {
      generated = BY_DEFAULT
    }
  }

  column "coin" {
    type    = varchar(10)
    null    = false
    comment = "coin type code"
  }

  column "seed" {
    type    = varchar(255)
    null    = false
    comment = "seed"
  }

  column "updated_at" {
    type    = timestamptz
    null    = true
    default = sql("now()")
    comment = "updated date"
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_seed_coin" {
    columns = [column.coin]
  }

  check "chk_seed_coin" {
    expr = "coin IN ('btc', 'bch')"
  }
}

# Table: auth_account_key
table "auth_account_key" {
  schema  = schema.public
  comment = "table for keys for auth account"

  column "id" {
    type    = smallint
    null    = false
    comment = "ID"
    identity {
      generated = BY_DEFAULT
    }
  }

  column "coin" {
    type    = varchar(10)
    null    = false
    comment = "coin type code"
  }

  column "key_type" {
    type    = varchar(20)
    null    = false
    default = "bip44"
    comment = "key type (bip44, bip49, bip84, bip86, musig2)"
  }

  column "auth_account" {
    type    = varchar(20)
    null    = false
    comment = "auth type"
  }

  column "account" {
    type    = varchar(20)
    null    = false
    default = "deposit"
    comment = "multisig account type (deposit, payment, stored)"
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
    type    = varchar(255)
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
    type    = smallint
    null    = false
    default = 0
    comment = "progress status for address generating"
  }

  column "updated_at" {
    type    = timestamptz
    null    = true
    default = sql("now()")
    comment = "updated date"
  }

  primary_key {
    columns = [column.id]
  }

  index "idex_coin_auth_account_account" {
    unique  = true
    columns = [column.coin, column.auth_account, column.account]
  }

  index "idx_p2pkh_address" {
    unique  = true
    columns = [column.p2pkh_address]
  }

  index "idx_p2sh_segwit_address" {
    unique  = true
    columns = [column.p2sh_segwit_address]
  }

  index "idx_bech32_address" {
    unique  = true
    columns = [column.bech32_address]
  }

  index "idx_wallet_import_format" {
    unique  = true
    columns = [column.wallet_import_format]
  }

  index "idx_auth_account_key_coin" {
    columns = [column.coin]
  }

  index "idx_auth_account_key_key_type" {
    columns = [column.key_type]
  }

  index "idx_auth_account" {
    columns = [column.auth_account]
  }

  check "chk_auth_account_key_coin" {
    expr = "coin IN ('btc', 'bch')"
  }
}

# Table: musig2_nonces
# MuSig2 nonce commitments for secure storage
# [SHARED] This table is also defined in: keygen.hcl
# When modifying, update all schemas for consistency!
table "musig2_nonces" {
  schema  = schema.public
  comment = "MuSig2 nonce commitments for secure storage"

  column "id" {
    type    = bigint
    null    = false
    comment = "ID"
    identity {
      generated = BY_DEFAULT
    }
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
    type    = bytea
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
    type    = timestamptz
    null    = true
    default = sql("now()")
    comment = "creation date"
  }

  column "used_at" {
    type    = timestamptz
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
