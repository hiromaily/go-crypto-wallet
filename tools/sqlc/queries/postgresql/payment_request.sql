-- name: GetAllPaymentRequests :many
SELECT * FROM payment_request
WHERE coin = $1 AND payment_id IS NULL;

-- name: GetPaymentRequestsByPaymentID :many
SELECT * FROM payment_request
WHERE coin = $1 AND payment_id = $2;

-- name: InsertPaymentRequest :one
INSERT INTO payment_request (coin, payment_id, sender_address, sender_account, receiver_address, amount, is_done, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;

-- name: UpdatePaymentRequestPaymentID :execresult
UPDATE payment_request
SET payment_id = $1
WHERE id = $2;

-- name: UpdatePaymentRequestIsDone :execresult
UPDATE payment_request
SET is_done = $1
WHERE coin = $2 AND payment_id = $3;

-- name: DeleteAllPaymentRequests :execresult
DELETE FROM payment_request
WHERE coin = $1;
