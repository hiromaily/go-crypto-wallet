-- name: GetTxByID :one
SELECT * FROM tx
WHERE id = $1;

-- name: GetMaxTxID :one
SELECT MAX(id) as max_id FROM tx
WHERE coin = $1 AND action = $2;

-- name: InsertTx :one
INSERT INTO tx (coin, action, updated_at)
VALUES ($1, $2, CURRENT_TIMESTAMP)
RETURNING id;

-- name: UpdateTx :exec
UPDATE tx
SET coin = $1, action = $2, updated_at = $3
WHERE id = $4;

-- name: DeleteAllTx :execresult
DELETE FROM tx;

-- name: GetAllTx :many
SELECT * FROM tx;
