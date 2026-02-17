-- name: GetSeed :one
SELECT * FROM seed WHERE coin = $1 LIMIT 1;

-- name: InsertSeed :one
INSERT INTO seed (coin, seed) VALUES ($1, $2)
RETURNING id;
