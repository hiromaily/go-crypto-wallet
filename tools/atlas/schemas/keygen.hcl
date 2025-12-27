# Keygen Schema Definition
# Key generation data (seeds, account keys, full public keys)

schema "keygen" {
  charset   = "utf8mb4"
  collation = "utf8mb4_unicode_ci"
  comment   = "Keygen schema for key generation operations"
}

# Table: seed
table "seed" {
  schema = schema.keygen
  comment = "table for seed"

  column "id" {
    type     = tinyint
    null     = false
    auto_increment = true
    comment  = "ID"
  }

  column "coin" {
    type     = enum("btc", "bch", "eth", "xrp", "hyt")
    null     = false
    comment  = "coin type code"
  }

  column "seed" {
    type     = varchar(255)
    null     = false
    comment  = "seed"
  }

  column "updated_at" {
    type     = datetime
    null     = true
    default  = sql("CURRENT_TIMESTAMP")
    comment  = "updated date"
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_coin" {
    columns = [column.coin]
  }
}

# Table: account_key
table "account_key" {
  schema = schema.keygen
  comment = "table for keys for any account"

  column "id" {
    type     = bigint
    null     = false
    auto_increment = true
    comment  = "ID"
  }

  column "coin" {
    type     = enum("btc", "bch", "eth", "xrp", "hyt")
    null     = false
    comment  = "coin type code"
  }

  column "account" {
    type     = enum("client", "deposit", "payment", "stored")
    null     = false
    comment  = "account type"
  }

  column "p2pkh_address" {
    type     = varchar(255)
    null     = false
    comment  = "address as standard pubkey script that Pays To PubKey Hash (P2PKH)"
  }

  column "p2sh_segwit_address" {
    type     = varchar(255)
    null     = false
    comment  = "p2sh-segwit address"
  }

  column "bech32_address" {
    type     = varchar(255)
    null     = false
    comment  = "bech32 address"
  }

  column "taproot_address" {
    type     = varchar(255)
    null     = true
    comment  = "taproot address (BIP86)"
  }

  column "full_public_key" {
    type     = varchar(255)
    null     = false
    comment  = "full public key"
  }

  column "multisig_address" {
    type     = varchar(255)
    null     = false
    default  = ""
    comment  = "multisig address"
  }

  column "redeem_script" {
    type     = varchar(1000)
    null     = false
    default  = ""
    comment  = "redeedScript after multisig address generated"
  }

  column "wallet_import_format" {
    type     = varchar(255)
    null     = false
    comment  = "WIF"
  }

  column "idx" {
    type     = bigint
    null     = false
    comment  = "index for hd wallet"
  }

  column "addr_status" {
    type     = tinyint
    null     = false
    default  = 0
    comment  = "progress status for address generating"
  }

  column "updated_at" {
    type     = datetime
    null     = true
    default  = sql("CURRENT_TIMESTAMP")
    comment  = "updated date"
  }

  primary_key {
    columns = [column.id]
  }

  unique "idx_p2pkh_address" {
    columns = [column.p2pkh_address]
  }

  unique "idx_wallet_import_format" {
    columns = [column.wallet_import_format]
  }

  index "idx_coin" {
    columns = [column.coin]
  }

  index "idx_account" {
    columns = [column.account]
  }
}

# Table: xrp_account_key
table "xrp_account_key" {
  schema = schema.keygen
  comment = "table for xrp keys for any account"

  column "id" {
    type     = bigint
    null     = false
    auto_increment = true
    comment  = "ID"
  }

  column "coin" {
    type     = enum("xrp")
    null     = false
    comment  = "coin type code"
  }

  column "account" {
    type     = enum("client", "deposit", "payment", "stored")
    null     = false
    comment  = "account type"
  }

  column "account_id" {
    type     = varchar(255)
    null     = false
    comment  = "account_id"
  }

  column "key_type" {
    type     = tinyint
    null     = false
    default  = 0
    comment  = "key_type"
  }

  column "master_key" {
    type     = varchar(255)
    null     = false
    comment  = "master_key, DEPRECATED"
  }

  column "master_seed" {
    type     = varchar(255)
    null     = false
    comment  = "master_seed"
  }

  column "master_seed_hex" {
    type     = varchar(255)
    null     = false
    comment  = "master_seed_hex"
  }

  column "public_key" {
    type     = varchar(255)
    null     = false
    comment  = "public_key"
  }

  column "public_key_hex" {
    type     = varchar(255)
    null     = false
    comment  = "public_key_hex"
  }

  column "is_regular_key_pair" {
    type     = boolean
    null     = false
    default  = false
    comment  = "true: this key is for regular key pair"
  }

  column "allocated_id" {
    type     = bigint
    null     = false
    default  = 0
    comment  = "index for hd wallet"
  }

  column "addr_status" {
    type     = tinyint
    null     = false
    default  = 0
    comment  = "progress status for address generating"
  }

  column "updated_at" {
    type     = datetime
    null     = true
    default  = sql("CURRENT_TIMESTAMP")
    comment  = "updated date"
  }

  primary_key {
    columns = [column.id]
  }

  unique "idx_account_id" {
    columns = [column.account_id]
  }

  unique "idx_master_seed" {
    columns = [column.master_seed]
  }

  index "idx_coin" {
    columns = [column.coin]
  }

  index "idx_account" {
    columns = [column.account]
  }
}

# Table: auth_fullpubkey
table "auth_fullpubkey" {
  schema = schema.keygen
  comment = "table for auth key exported from sign db"

  column "id" {
    type     = smallint
    null     = false
    auto_increment = true
    comment  = "ID"
  }

  column "coin" {
    type     = enum("btc", "bch")
    null     = false
    comment  = "coin type code"
  }

  column "auth_account" {
    type     = varchar(20)
    null     = false
    comment  = "auth type"
  }

  column "full_public_key" {
    type     = varchar(255)
    null     = false
    comment  = "full public key"
  }

  column "updated_at" {
    type     = datetime
    null     = true
    default  = sql("CURRENT_TIMESTAMP")
    comment  = "updated date"
  }

  primary_key {
    columns = [column.id]
  }

  unique "idex_coin_auth_account" {
    columns = [column.coin, column.auth_account]
  }

  unique "idx_full_public_key" {
    columns = [column.full_public_key]
  }

  index "idx_coin" {
    columns = [column.coin]
  }
}

