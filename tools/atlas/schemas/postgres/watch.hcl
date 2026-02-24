# Watch Schema Definition
# Online wallet data (addresses, transactions, payment requests)

# Schema declaration (PostgreSQL uses "public" schema within each database)
schema "public" {}

# Table: btc_tx
table "btc_tx" {
  schema  = schema.public
  comment = "table for btc transaction info"

  column "id" {
    type    = bigint
    null    = false
    comment = "transaction ID"
    identity {
      generated = BY_DEFAULT
    }
  }

  column "coin" {
    type    = varchar(10)
    null    = false
    comment = "coin type code"
  }

  column "action" {
    type    = varchar(20)
    null    = false
    comment = "action type"
  }

  column "unsigned_hex_tx" {
    type    = text
    null    = false
    comment = "HEX string for unsigned transaction"
  }

  column "signed_hex_tx" {
    type    = text
    null    = false
    comment = "HEX string for signed transaction"
  }

  column "sent_hash_tx" {
    type    = text
    null    = false
    comment = "Hash for sent transaction"
  }

  column "total_input_amount" {
    type    = numeric(26,10)
    null    = false
    comment = "total amount of coin to send"
  }

  column "total_output_amount" {
    type    = numeric(26,10)
    null    = false
    comment = "total amount of coin to receive without fee"
  }

  column "fee" {
    type    = numeric(26,10)
    null    = false
    comment = "fee"
  }

  column "current_tx_type" {
    type    = smallint
    null    = false
    default = 1
    comment = "current transaction type"
  }

  column "unsigned_updated_at" {
    type    = timestamptz
    null    = true
    default = sql("now()")
    comment = "updated date for unsigned transaction created"
  }

  column "sent_updated_at" {
    type    = timestamptz
    null    = true
    comment = "updated date for signed transaction sent"
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_btc_tx_coin" {
    columns = [column.coin]
  }

  index "idx_btc_tx_action" {
    columns = [column.action]
  }

  check "chk_btc_tx_coin" {
    expr = "coin IN ('btc', 'bch')"
  }

  check "chk_btc_tx_action" {
    expr = "action IN ('deposit', 'payment', 'transfer')"
  }
}

# Table: btc_tx_input
table "btc_tx_input" {
  schema  = schema.public
  comment = "table for input transaction"

  column "id" {
    type    = bigint
    null    = false
    comment = "ID"
    identity {
      generated = BY_DEFAULT
    }
  }

  column "tx_id" {
    type    = bigint
    null    = false
    comment = "tx table ID"
  }

  column "input_txid" {
    type    = varchar(255)
    null    = false
    comment = "txid for input"
  }

  column "input_vout" {
    type    = integer
    null    = false
    comment = "vout for input"
  }

  column "input_address" {
    type    = varchar(255)
    null    = false
    comment = "sender address for input"
  }

  column "input_account" {
    type    = varchar(255)
    null    = false
    comment = "sender account for input"
  }

  column "input_amount" {
    type    = numeric(26,10)
    null    = false
    comment = "amount of coin to send for input"
  }

  column "input_confirmations" {
    type    = bigint
    null    = false
    comment = "block confirmations when unspent rpc returned"
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

  index "idx_btc_tx_input_tx_id" {
    columns = [column.tx_id]
  }
}

# Table: btc_tx_output
table "btc_tx_output" {
  schema  = schema.public
  comment = "table for output transaction"

  column "id" {
    type    = bigint
    null    = false
    comment = "ID"
    identity {
      generated = BY_DEFAULT
    }
  }

  column "tx_id" {
    type    = bigint
    null    = false
    comment = "tx table ID"
  }

  column "output_address" {
    type    = varchar(255)
    null    = false
    comment = "receiver address for output"
  }

  column "output_account" {
    type    = varchar(255)
    null    = false
    comment = "receiver account for output"
  }

  column "output_amount" {
    type    = numeric(26,10)
    null    = false
    comment = "amount of coin to receive"
  }

  column "is_change" {
    type    = boolean
    null    = false
    default = false
    comment = "true: output is for fee"
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

  index "idx_btc_tx_output_tx_id" {
    columns = [column.tx_id]
  }
}

# Table: tx
table "tx" {
  schema  = schema.public
  comment = "table for eth transaction info"

  column "id" {
    type    = bigint
    null    = false
    comment = "transaction ID"
    identity {
      generated = BY_DEFAULT
    }
  }

  column "coin" {
    type    = varchar(10)
    null    = false
    comment = "coin type code"
  }

  column "action" {
    type    = varchar(20)
    null    = false
    comment = "action type"
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

  index "idx_tx_coin" {
    columns = [column.coin]
  }

  index "idx_tx_action" {
    columns = [column.action]
  }

  check "chk_tx_coin" {
    expr = "coin IN ('eth', 'xrp', 'hyt')"
  }

  check "chk_tx_action" {
    expr = "action IN ('deposit', 'payment', 'transfer')"
  }
}

# Table: eth_detail_tx
table "eth_detail_tx" {
  schema  = schema.public
  comment = "table for eth transaction detail"

  column "id" {
    type    = bigint
    null    = false
    comment = "ID"
    identity {
      generated = BY_DEFAULT
    }
  }

  column "tx_id" {
    type    = bigint
    null    = false
    comment = "eth_tx table ID"
  }

  column "uuid" {
    type    = varchar(36)
    null    = false
    comment = "UUID"
  }

  column "current_tx_type" {
    type    = smallint
    null    = false
    default = 1
    comment = "current transaction type"
  }

  column "sender_account" {
    type    = varchar(255)
    null    = false
    comment = "sender account"
  }

  column "sender_address" {
    type    = varchar(255)
    null    = false
    comment = "sender address"
  }

  column "receiver_account" {
    type    = varchar(255)
    null    = false
    comment = "receiver account"
  }

  column "receiver_address" {
    type    = varchar(255)
    null    = false
    comment = "receiver address"
  }

  column "amount" {
    type    = bigint
    null    = false
    comment = "amount of coin to receive"
  }

  column "fee" {
    type    = bigint
    null    = false
    comment = "fee"
  }

  column "gas_limit" {
    type    = integer
    null    = false
    comment = "gas limit"
  }

  column "nonce" {
    type    = bigint
    null    = false
    comment = "nonce"
  }

  column "unsigned_hex_tx" {
    type    = text
    null    = false
    comment = "HEX string for unsigned transaction"
  }

  column "signed_hex_tx" {
    type    = text
    null    = false
    comment = "HEX string for signed transaction"
  }

  column "sent_hash_tx" {
    type    = text
    null    = false
    comment = "Hash for sent transaction"
  }

  column "unsigned_updated_at" {
    type    = timestamptz
    null    = true
    default = sql("now()")
    comment = "updated date for unsigned transaction created"
  }

  column "sent_updated_at" {
    type    = timestamptz
    null    = true
    comment = "updated date for signed transaction sent"
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_eth_detail_tx_uuid" {
    unique  = true
    columns = [column.uuid]
  }

  index "idx_eth_detail_tx_txid" {
    columns = [column.tx_id]
  }

  index "idx_eth_detail_tx_sender_account" {
    columns = [column.sender_account]
  }

  index "idx_eth_detail_tx_receiver_account" {
    columns = [column.receiver_account]
  }
}

# Table: xrp_detail_tx
table "xrp_detail_tx" {
  schema  = schema.public
  comment = "table for xrp transaction detail"

  column "id" {
    type    = bigint
    null    = false
    comment = "ID"
    identity {
      generated = BY_DEFAULT
    }
  }

  column "tx_id" {
    type    = bigint
    null    = false
    comment = "xrp_tx table ID"
  }

  column "uuid" {
    type    = varchar(36)
    null    = false
    comment = "UUID"
  }

  column "current_tx_type" {
    type    = smallint
    null    = false
    default = 1
    comment = "current transaction type"
  }

  column "sender_account" {
    type    = varchar(255)
    null    = false
    comment = "sender account"
  }

  column "sender_address" {
    type    = varchar(255)
    null    = false
    comment = "sender address"
  }

  column "receiver_account" {
    type    = varchar(255)
    null    = false
    comment = "receiver account"
  }

  column "receiver_address" {
    type    = varchar(255)
    null    = false
    comment = "receiver address"
  }

  column "amount" {
    type    = varchar(255)
    null    = false
    comment = "amount of coin to receive"
  }

  column "xrp_tx_type" {
    type    = varchar(255)
    null    = false
    comment = "xrp tx type like `Payment`"
  }

  column "fee" {
    type    = varchar(255)
    null    = false
    comment = "tx fee"
  }

  column "flags" {
    type    = bigint
    null    = false
    comment = "tx flags"
  }

  column "last_ledger_sequence" {
    type    = bigint
    null    = false
    comment = "tx LastLedgerSequence"
  }

  column "sequence" {
    type    = bigint
    null    = false
    comment = "tx Sequence"
  }

  column "signing_pubkey" {
    type    = varchar(255)
    null    = false
    comment = "tx SigningPubKey"
  }

  column "txn_signature" {
    type    = varchar(255)
    null    = false
    comment = "tx TxnSignature"
  }

  column "hash" {
    type    = varchar(255)
    null    = false
    comment = "tx Hash"
  }

  column "earliest_ledger_version" {
    type    = bigint
    null    = false
    comment = "tx earliest_ledger_version after sending tx"
  }

  column "signed_tx_id" {
    type    = varchar(255)
    null    = false
    comment = "signed tx id"
  }

  column "tx_blob" {
    type    = text
    null    = false
    comment = "sent tx blob"
  }

  column "sent_updated_at" {
    type    = timestamptz
    null    = true
    comment = "updated date for signed transaction sent"
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_xrp_detail_tx_uuid" {
    unique  = true
    columns = [column.uuid]
  }

  index "idx_xrp_detail_tx_txid" {
    columns = [column.tx_id]
  }

  index "idx_xrp_detail_tx_sender_account" {
    columns = [column.sender_account]
  }

  index "idx_xrp_detail_tx_receiver_account" {
    columns = [column.receiver_account]
  }
}

# Table: xrp_pending_multisig
# Pending multi-signature transactions waiting for signature collection
# Reference: https://xrpl.org/docs/concepts/accounts/multi-signing
table "xrp_pending_multisig" {
  schema  = schema.public
  comment = "Pending XRP multi-signature transactions awaiting signatures"

  column "id" {
    type    = bigint
    null    = false
    comment = "ID"
    identity {
      generated = BY_DEFAULT
    }
  }

  column "tx_uuid" {
    type    = varchar(36)
    null    = false
    comment = "Unique identifier for the transaction"
  }

  column "account_id" {
    type    = varchar(255)
    null    = false
    comment = "XRP account address (r...) that will sign this transaction"
  }

  column "unsigned_tx_json" {
    type    = text
    null    = false
    comment = "JSON representation of the unsigned transaction"
  }

  column "xrp_tx_type" {
    type    = varchar(50)
    null    = false
    comment = "XRP transaction type (Payment, SignerListSet, etc.)"
  }

  column "required_quorum" {
    type    = integer
    null    = false
    comment = "Total signature weight required for this transaction"
  }

  column "current_weight" {
    type    = integer
    null    = false
    default = 0
    comment = "Current total weight of collected signatures"
  }

  column "status" {
    type    = varchar(20)
    null    = false
    default = "pending"
    comment = "Current status of the multi-sig transaction"
  }

  column "combined_tx_blob" {
    type    = text
    null    = true
    comment = "Combined signed transaction blob (set when status is ready)"
  }

  column "submitted_tx_hash" {
    type    = varchar(255)
    null    = true
    comment = "Transaction hash after submission to ledger"
  }

  column "expires_at" {
    type    = timestamptz
    null    = true
    comment = "Expiration time for this pending transaction"
  }

  column "created_at" {
    type    = timestamptz
    null    = true
    default = sql("now()")
    comment = "creation date"
  }

  column "updated_at" {
    type    = timestamptz
    null    = true
    default = sql("now()")
    comment = "last update date"
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_tx_uuid" {
    unique  = true
    columns = [column.tx_uuid]
  }

  index "idx_account_id" {
    columns = [column.account_id]
  }

  index "idx_status" {
    columns = [column.status]
  }

  index "idx_account_status" {
    columns = [column.account_id, column.status]
    comment = "Find pending transactions for an account"
  }

  check "chk_xrp_pending_multisig_status" {
    expr = "status IN ('pending', 'ready', 'submitted', 'confirmed', 'failed', 'expired')"
  }
}

# Table: xrp_multisig_signature
# Collected signatures for pending multi-signature transactions
table "xrp_multisig_signature" {
  schema  = schema.public
  comment = "Collected signatures for XRP multi-signature transactions"

  column "id" {
    type    = bigint
    null    = false
    comment = "ID"
    identity {
      generated = BY_DEFAULT
    }
  }

  column "pending_multisig_id" {
    type    = bigint
    null    = false
    comment = "Reference to xrp_pending_multisig.id"
  }

  column "signer_account" {
    type    = varchar(255)
    null    = false
    comment = "XRP address of the signer (r...)"
  }

  column "signed_tx_blob" {
    type    = text
    null    = false
    comment = "Signed transaction blob from this signer"
  }

  column "signer_weight" {
    type    = integer
    null    = false
    comment = "Weight of this signature"
  }

  column "signed_at" {
    type    = timestamptz
    null    = true
    default = sql("now()")
    comment = "date when signature was collected"
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_pending_multisig_id" {
    columns = [column.pending_multisig_id]
  }

  index "idx_signer_account" {
    columns = [column.signer_account]
  }

  index "idx_pending_signer" {
    unique  = true
    columns = [column.pending_multisig_id, column.signer_account]
    comment = "Each signer can only sign once per transaction"
  }
}

# Table: address
table "address" {
  schema  = schema.public
  comment = "table for account pubkey"

  column "id" {
    type    = bigint
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

  column "account" {
    type    = varchar(20)
    null    = false
    comment = "account type"
  }

  column "wallet_address" {
    type    = varchar(255)
    null    = false
    comment = "wallet address"
  }

  column "is_allocated" {
    type    = boolean
    null    = false
    default = false
    comment = "true: address is allocated(used)"
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

  index "idx_wallet_address" {
    unique  = true
    columns = [column.wallet_address]
  }

  index "idx_address_coin" {
    columns = [column.coin]
  }

  index "idx_address_account" {
    columns = [column.account]
  }

  check "chk_address_coin" {
    expr = "coin IN ('btc', 'bch', 'eth', 'xrp', 'hyt')"
  }

  check "chk_address_account" {
    expr = "account IN ('client', 'deposit', 'payment', 'stored')"
  }
}

# Table: payment_request
table "payment_request" {
  schema  = schema.public
  comment = "table for payment request"

  column "id" {
    type    = bigint
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

  column "payment_id" {
    type    = bigint
    null    = true
    comment = "tx table ID for payment action"
  }

  column "sender_address" {
    type    = varchar(255)
    null    = false
    comment = "sender address"
  }

  column "sender_account" {
    type    = varchar(255)
    null    = false
    comment = "sender account"
  }

  column "receiver_address" {
    type    = varchar(255)
    null    = false
    comment = "receiver address"
  }

  column "amount" {
    type    = numeric(26,10)
    null    = false
    comment = "amount of coin to send"
  }

  column "is_done" {
    type    = boolean
    null    = false
    default = false
    comment = "true: unsigned transaction is created"
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

  index "idx_payment_request_coin" {
    columns = [column.coin]
  }

  check "chk_payment_request_coin" {
    expr = "coin IN ('btc', 'bch', 'eth', 'xrp')"
  }
}
