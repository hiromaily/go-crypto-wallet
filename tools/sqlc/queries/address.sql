-- name: GetAllAddresses :many
SELECT * FROM address
WHERE coin = ? AND account = ?;

-- name: GetAllAddressStrings :many
SELECT wallet_address FROM address
WHERE coin = ? AND account = ?;

-- name: GetOneUnallocatedAddress :one
SELECT * FROM address
WHERE coin = ? AND account = ? AND is_allocated = false
LIMIT 1;

-- name: InsertAddress :execresult
INSERT INTO address (coin, account, wallet_address, is_allocated, updated_at)
VALUES (?, ?, ?, ?, ?);

-- name: UpdateAddressIsAllocated :execresult
UPDATE address
SET is_allocated = ?, updated_at = ?
WHERE coin = ? AND wallet_address = ?;
