# Watch Schema Definition
# Online wallet data (addresses, transactions, payment requests)

schema "watch" {
  charset   = "utf8mb4"
  collation = "utf8mb4_unicode_ci"
  comment   = "Watch schema for online wallet operations"
}

# Table: btc_tx
table "btc_tx" {
  schema = schema.watch
  comment = "table for btc transaction info"

  column "id" {
    type     = bigint
    null     = false
    auto_increment = true
    comment  = "transaction ID"
  }

  column "coin" {
    type     = enum("btc", "bch")
    null     = false
    comment  = "coin type code"
  }

  column "action" {
    type     = enum("deposit", "payment", "transfer")
    null     = false
    comment  = "action type"
  }

  column "unsigned_hex_tx" {
    type     = text
    null     = false
    comment  = "HEX string for unsigned transaction"
  }

  column "signed_hex_tx" {
    type     = text
    null     = false
    default  = ""
    comment  = "HEX string for signed transaction"
  }

  column "sent_hash_tx" {
    type     = text
    null     = false
    default  = ""
    comment  = "Hash for sent transaction"
  }

  column "total_input_amount" {
    type     = decimal(26, 10)
    null     = false
    comment  = "total amount of coin to send"
  }

  column "total_output_amount" {
    type     = decimal(26, 10)
    null     = false
    comment  = "total amount of coin to receive without fee"
  }

  column "fee" {
    type     = decimal(26, 10)
    null     = false
    comment  = "fee"
  }

  column "current_tx_type" {
    type     = tinyint
    null     = false
    default  = 1
    comment  = "current transaction type"
  }

  column "unsigned_updated_at" {
    type     = datetime
    null     = true
    default  = sql("CURRENT_TIMESTAMP")
    comment  = "updated date for unsigned transaction created"
  }

  column "sent_updated_at" {
    type     = datetime
    null     = true
    comment  = "updated date for signed transaction sent"
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_coin" {
    columns = [column.coin]
  }

  index "idx_action" {
    columns = [column.action]
  }
}

# Table: btc_tx_input
table "btc_tx_input" {
  schema = schema.watch
  comment = "table for input transaction"

  column "id" {
    type     = bigint
    null     = false
    auto_increment = true
    comment  = "ID"
  }

  column "tx_id" {
    type     = bigint
    null     = false
    comment  = "tx table ID"
  }

  column "input_txid" {
    type     = varchar(255)
    null     = false
    comment  = "txid for input"
  }

  column "input_vout" {
    type     = mediumint
    unsigned = true
    null     = false
    comment  = "vout for input"
  }

  column "input_address" {
    type     = varchar(255)
    null     = false
    comment  = "sender address for input"
  }

  column "input_account" {
    type     = varchar(255)
    null     = false
    comment  = "sender account for input"
  }

  column "input_amount" {
    type     = decimal(26, 10)
    null     = false
    comment  = "amount of coin to send for input"
  }

  column "input_confirmations" {
    type     = bigint
    unsigned = true
    null     = false
    comment  = "block confirmations when unspent rpc returned"
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

  index "idx_tx_id" {
    columns = [column.tx_id]
  }
}

# Table: btc_tx_output
table "btc_tx_output" {
  schema = schema.watch
  comment = "table for output transaction"

  column "id" {
    type     = bigint
    null     = false
    auto_increment = true
    comment  = "ID"
  }

  column "tx_id" {
    type     = bigint
    null     = false
    comment  = "tx table ID"
  }

  column "output_address" {
    type     = varchar(255)
    null     = false
    comment  = "receiver address for output"
  }

  column "output_account" {
    type     = varchar(255)
    null     = false
    comment  = "receiver account for output"
  }

  column "output_amount" {
    type     = decimal(26, 10)
    null     = false
    comment  = "amount of coin to receive"
  }

  column "is_change" {
    type     = boolean
    null     = false
    default  = false
    comment  = "true: output is for fee"
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

  index "idx_tx_id" {
    columns = [column.tx_id]
  }
}

# Table: tx
table "tx" {
  schema = schema.watch
  comment = "table for eth transaction info"

  column "id" {
    type     = bigint
    null     = false
    auto_increment = true
    comment  = "transaction ID"
  }

  column "coin" {
    type     = enum("eth", "xrp", "hyt")
    null     = false
    comment  = "coin type code"
  }

  column "action" {
    type     = enum("deposit", "payment", "transfer")
    null     = false
    comment  = "action type"
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

  index "idx_action" {
    columns = [column.action]
  }
}

# Table: eth_detail_tx
table "eth_detail_tx" {
  schema = schema.watch
  comment = "table for eth transaction detail"

  column "id" {
    type     = bigint
    null     = false
    auto_increment = true
    comment  = "ID"
  }

  column "tx_id" {
    type     = bigint
    null     = false
    comment  = "eth_tx table ID"
  }

  column "uuid" {
    type     = varchar(36)
    null     = false
    comment  = "UUID"
  }

  column "current_tx_type" {
    type     = tinyint
    null     = false
    default  = 1
    comment  = "current transaction type"
  }

  column "sender_account" {
    type     = varchar(255)
    null     = false
    comment  = "sender account"
  }

  column "sender_address" {
    type     = varchar(255)
    null     = false
    comment  = "sender address"
  }

  column "receiver_account" {
    type     = varchar(255)
    null     = false
    comment  = "receiver account"
  }

  column "receiver_address" {
    type     = varchar(255)
    null     = false
    comment  = "receiver address"
  }

  column "amount" {
    type     = bigint
    unsigned = true
    null     = false
    comment  = "amount of coin to receive"
  }

  column "fee" {
    type     = bigint
    unsigned = true
    null     = false
    comment  = "fee"
  }

  column "gas_limit" {
    type     = mediumint
    unsigned = true
    null     = false
    comment  = "gas limit"
  }

  column "nonce" {
    type     = bigint
    unsigned = true
    null     = false
    comment  = "nonce"
  }

  column "unsigned_hex_tx" {
    type     = text
    null     = false
    comment  = "HEX string for unsigned transaction"
  }

  column "signed_hex_tx" {
    type     = text
    null     = false
    default  = ""
    comment  = "HEX string for signed transaction"
  }

  column "sent_hash_tx" {
    type     = text
    null     = false
    default  = ""
    comment  = "Hash for sent transaction"
  }

  column "unsigned_updated_at" {
    type     = datetime
    null     = true
    default  = sql("CURRENT_TIMESTAMP")
    comment  = "updated date for unsigned transaction created"
  }

  column "sent_updated_at" {
    type     = datetime
    null     = true
    comment  = "updated date for signed transaction sent"
  }

  primary_key {
    columns = [column.id]
  }

  unique "idx_uuid" {
    columns = [column.uuid]
  }

  index "idx_txid" {
    columns = [column.tx_id]
  }

  index "idx_sender_account" {
    columns = [column.sender_account]
  }

  index "idx_receiver_account" {
    columns = [column.receiver_account]
  }
}

# Table: xrp_detail_tx
table "xrp_detail_tx" {
  schema = schema.watch
  comment = "table for xrp transaction detail"

  column "id" {
    type     = bigint
    null     = false
    auto_increment = true
    comment  = "ID"
  }

  column "tx_id" {
    type     = bigint
    null     = false
    comment  = "xrp_tx table ID"
  }

  column "uuid" {
    type     = varchar(36)
    null     = false
    comment  = "UUID"
  }

  column "current_tx_type" {
    type     = tinyint
    null     = false
    default  = 1
    comment  = "current transaction type"
  }

  column "sender_account" {
    type     = varchar(255)
    null     = false
    comment  = "sender account"
  }

  column "sender_address" {
    type     = varchar(255)
    null     = false
    comment  = "sender address"
  }

  column "receiver_account" {
    type     = varchar(255)
    null     = false
    comment  = "receiver account"
  }

  column "receiver_address" {
    type     = varchar(255)
    null     = false
    comment  = "receiver address"
  }

  column "amount" {
    type     = varchar(255)
    null     = false
    comment  = "amount of coin to receive"
  }

  column "xrp_tx_type" {
    type     = varchar(255)
    null     = false
    comment  = "xrp tx type like `Payment`"
  }

  column "fee" {
    type     = varchar(255)
    null     = false
    comment  = "tx fee"
  }

  column "flags" {
    type     = bigint
    unsigned = true
    null     = false
    comment  = "tx flags"
  }

  column "last_ledger_sequence" {
    type     = bigint
    unsigned = true
    null     = false
    comment  = "tx LastLedgerSequence"
  }

  column "sequence" {
    type     = bigint
    unsigned = true
    null     = false
    comment  = "tx Sequence"
  }

  column "signing_pubkey" {
    type     = varchar(255)
    null     = false
    comment  = "tx SigningPubKey"
  }

  column "txn_signature" {
    type     = varchar(255)
    null     = false
    comment  = "tx TxnSignature"
  }

  column "hash" {
    type     = varchar(255)
    null     = false
    comment  = "tx Hash"
  }

  column "earliest_ledger_version" {
    type     = bigint
    unsigned = true
    null     = false
    comment  = "tx earliest_ledger_version after sending tx"
  }

  column "signed_tx_id" {
    type     = varchar(255)
    null     = false
    comment  = "signed tx id"
  }

  column "tx_blob" {
    type     = text
    null     = false
    comment  = "sent tx blob"
  }

  column "sent_updated_at" {
    type     = datetime
    null     = true
    comment  = "updated date for signed transaction sent"
  }

  primary_key {
    columns = [column.id]
  }

  unique "idx_uuid" {
    columns = [column.uuid]
  }

  index "idx_txid" {
    columns = [column.tx_id]
  }

  index "idx_sender_account" {
    columns = [column.sender_account]
  }

  index "idx_receiver_account" {
    columns = [column.receiver_account]
  }
}

# Table: address
table "address" {
  schema = schema.watch
  comment = "table for account pubkey"

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

  column "wallet_address" {
    type     = varchar(255)
    null     = false
    comment  = "wallet address"
  }

  column "is_allocated" {
    type     = boolean
    null     = false
    default  = false
    comment  = "true: address is allocated(used)"
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

  unique "idx_wallet_address" {
    columns = [column.wallet_address]
  }

  index "idx_coin" {
    columns = [column.coin]
  }

  index "idx_account" {
    columns = [column.account]
  }
}

# Table: payment_request
table "payment_request" {
  schema = schema.watch
  comment = "table for payment request"

  column "id" {
    type     = bigint
    null     = false
    auto_increment = true
    comment  = "ID"
  }

  column "coin" {
    type     = enum("btc", "bch", "eth", "xrp")
    null     = false
    comment  = "coin type code"
  }

  column "payment_id" {
    type     = bigint
    null     = true
    comment  = "tx table ID for payment action"
  }

  column "sender_address" {
    type     = varchar(255)
    null     = false
    comment  = "sender address"
  }

  column "sender_account" {
    type     = varchar(255)
    null     = false
    comment  = "sender account"
  }

  column "receiver_address" {
    type     = varchar(255)
    null     = false
    comment  = "receiver address"
  }

  column "amount" {
    type     = decimal(26, 10)
    null     = false
    comment  = "amount of coin to send"
  }

  column "is_done" {
    type     = boolean
    null     = false
    default  = false
    comment  = "true: unsigned transaction is created"
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

