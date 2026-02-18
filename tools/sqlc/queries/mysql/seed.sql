-- name: GetSeed :one
SELECT * FROM seed WHERE coin = ? LIMIT 1;

-- name: InsertSeed :execresult
INSERT INTO seed (coin, seed) VALUES (?, ?);
