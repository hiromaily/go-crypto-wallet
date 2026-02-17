-- name: GetAllAddresses :many
SELECT * FROM address
WHERE coin = $1 AND account = $2;

-- name: GetAllAddressStrings :many
SELECT wallet_address FROM address
WHERE coin = $1 AND account = $2;

-- name: GetOneUnallocatedAddress :one
SELECT * FROM address
WHERE coin = $1 AND account = $2 AND is_allocated = false
LIMIT 1;

-- name: InsertAddress :one
INSERT INTO address (coin, account, wallet_address, is_allocated, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: UpdateAddressIsAllocated :execresult
UPDATE address
SET is_allocated = $1, updated_at = $2
WHERE coin = $3 AND wallet_address = $4;
