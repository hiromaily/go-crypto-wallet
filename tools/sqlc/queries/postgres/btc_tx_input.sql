-- name: GetBtcTxInputByID :one
SELECT * FROM btc_tx_input
WHERE id = $1;

-- name: GetBtcTxInputsByTxID :many
SELECT * FROM btc_tx_input
WHERE tx_id = $1;

-- name: InsertBtcTxInput :one
INSERT INTO btc_tx_input (
  tx_id, input_txid, input_vout, input_address, input_account,
  input_amount, input_confirmations, updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;
