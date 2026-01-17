-- SQLite schema for watch wallet
-- Converted from MySQL schema

CREATE TABLE address (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  coin TEXT NOT NULL CHECK(coin IN ('btc','bch','eth','xrp','hyt')),
  account TEXT NOT NULL CHECK(account IN ('client','deposit','payment','stored')),
  wallet_address TEXT NOT NULL UNIQUE,
  is_allocated INTEGER NOT NULL DEFAULT 0,
  updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_address_account ON address(account);
CREATE INDEX idx_address_coin ON address(coin);

CREATE TABLE btc_tx (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  coin TEXT NOT NULL CHECK(coin IN ('btc','bch')),
  action TEXT NOT NULL CHECK(action IN ('deposit','payment','transfer')),
  unsigned_hex_tx TEXT NOT NULL,
  signed_hex_tx TEXT NOT NULL,
  sent_hash_tx TEXT NOT NULL,
  total_input_amount TEXT NOT NULL,
  total_output_amount TEXT NOT NULL,
  fee TEXT NOT NULL,
  current_tx_type INTEGER NOT NULL DEFAULT 1,
  unsigned_updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
  sent_updated_at TEXT DEFAULT NULL
);

CREATE INDEX idx_btc_tx_action ON btc_tx(action);
CREATE INDEX idx_btc_tx_coin ON btc_tx(coin);

CREATE TABLE btc_tx_input (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  tx_id INTEGER NOT NULL,
  input_txid TEXT NOT NULL,
  input_vout INTEGER NOT NULL,
  input_address TEXT NOT NULL,
  input_account TEXT NOT NULL,
  input_amount TEXT NOT NULL,
  input_confirmations INTEGER NOT NULL,
  updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_btc_tx_input_tx_id ON btc_tx_input(tx_id);

CREATE TABLE btc_tx_output (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  tx_id INTEGER NOT NULL,
  output_address TEXT NOT NULL,
  output_account TEXT NOT NULL,
  output_amount TEXT NOT NULL,
  is_change INTEGER NOT NULL DEFAULT 0,
  updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_btc_tx_output_tx_id ON btc_tx_output(tx_id);

CREATE TABLE eth_detail_tx (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  tx_id INTEGER NOT NULL,
  uuid TEXT NOT NULL UNIQUE,
  current_tx_type INTEGER NOT NULL DEFAULT 1,
  sender_account TEXT NOT NULL,
  sender_address TEXT NOT NULL,
  receiver_account TEXT NOT NULL,
  receiver_address TEXT NOT NULL,
  amount INTEGER NOT NULL,
  fee INTEGER NOT NULL,
  gas_limit INTEGER NOT NULL,
  nonce INTEGER NOT NULL,
  unsigned_hex_tx TEXT NOT NULL,
  signed_hex_tx TEXT NOT NULL,
  sent_hash_tx TEXT NOT NULL,
  unsigned_updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
  sent_updated_at TEXT DEFAULT NULL
);

CREATE INDEX idx_eth_detail_tx_receiver_account ON eth_detail_tx(receiver_account);
CREATE INDEX idx_eth_detail_tx_sender_account ON eth_detail_tx(sender_account);
CREATE INDEX idx_eth_detail_tx_txid ON eth_detail_tx(tx_id);

CREATE TABLE payment_request (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  coin TEXT NOT NULL CHECK(coin IN ('btc','bch','eth','xrp')),
  payment_id INTEGER DEFAULT NULL,
  sender_address TEXT NOT NULL,
  sender_account TEXT NOT NULL,
  receiver_address TEXT NOT NULL,
  amount TEXT NOT NULL,
  is_done INTEGER NOT NULL DEFAULT 0,
  updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payment_request_coin ON payment_request(coin);

CREATE TABLE tx (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  coin TEXT NOT NULL CHECK(coin IN ('eth','xrp','hyt')),
  action TEXT NOT NULL CHECK(action IN ('deposit','payment','transfer')),
  updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tx_action ON tx(action);
CREATE INDEX idx_tx_coin ON tx(coin);

CREATE TABLE xrp_detail_tx (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  tx_id INTEGER NOT NULL,
  uuid TEXT NOT NULL UNIQUE,
  current_tx_type INTEGER NOT NULL DEFAULT 1,
  sender_account TEXT NOT NULL,
  sender_address TEXT NOT NULL,
  receiver_account TEXT NOT NULL,
  receiver_address TEXT NOT NULL,
  amount TEXT NOT NULL,
  xrp_tx_type TEXT NOT NULL,
  fee TEXT NOT NULL,
  flags INTEGER NOT NULL,
  last_ledger_sequence INTEGER NOT NULL,
  sequence INTEGER NOT NULL,
  signing_pubkey TEXT NOT NULL,
  txn_signature TEXT NOT NULL,
  hash TEXT NOT NULL,
  earliest_ledger_version INTEGER NOT NULL,
  signed_tx_id TEXT NOT NULL,
  tx_blob TEXT NOT NULL,
  sent_updated_at TEXT DEFAULT NULL
);

CREATE INDEX idx_xrp_detail_tx_receiver_account ON xrp_detail_tx(receiver_account);
CREATE INDEX idx_xrp_detail_tx_sender_account ON xrp_detail_tx(sender_account);
CREATE INDEX idx_xrp_detail_tx_txid ON xrp_detail_tx(tx_id);
