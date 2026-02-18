-- name: GetBtcTxOutputByID :one
SELECT * FROM btc_tx_output
WHERE id = ?;

-- name: GetBtcTxOutputsByTxID :many
SELECT * FROM btc_tx_output
WHERE tx_id = ?;

-- name: InsertBtcTxOutput :execresult
INSERT INTO btc_tx_output (
  tx_id, output_address, output_account, output_amount, is_change, updated_at
) VALUES (?, ?, ?, ?, ?, ?);
