-- name: GetBtcTxOutputByID :one
SELECT * FROM btc_tx_output
WHERE id = $1;

-- name: GetBtcTxOutputsByTxID :many
SELECT * FROM btc_tx_output
WHERE tx_id = $1;

-- name: InsertBtcTxOutput :one
INSERT INTO btc_tx_output (
  tx_id, output_address, output_account, output_amount, is_change, updated_at
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
