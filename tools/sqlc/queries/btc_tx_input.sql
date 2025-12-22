-- name: GetBtcTxInputByID :one
SELECT * FROM btc_tx_input
WHERE id = ?;

-- name: GetBtcTxInputsByTxID :many
SELECT * FROM btc_tx_input
WHERE tx_id = ?;

-- name: InsertBtcTxInput :execresult
INSERT INTO btc_tx_input (
  tx_id, input_txid, input_vout, input_address, input_account,
  input_amount, input_confirmations, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?);
