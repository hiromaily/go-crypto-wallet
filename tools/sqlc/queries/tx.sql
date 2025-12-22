-- name: GetTxByID :one
SELECT * FROM tx
WHERE id = ?;

-- name: GetMaxTxID :one
SELECT MAX(id) as max_id FROM tx
WHERE coin = ? AND action = ?;

-- name: InsertTx :execresult
INSERT INTO tx (coin, action, updated_at)
VALUES (?, ?, CURRENT_TIMESTAMP);

-- name: UpdateTx :exec
UPDATE tx
SET coin = ?, action = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteAllTx :execresult
DELETE FROM tx;

-- name: GetAllTx :many
SELECT * FROM tx;
