-- name: GetAllPaymentRequests :many
SELECT * FROM payment_request
WHERE coin = ? AND payment_id IS NULL;

-- name: GetPaymentRequestsByPaymentID :many
SELECT * FROM payment_request
WHERE coin = ? AND payment_id = ?;

-- name: InsertPaymentRequest :execresult
INSERT INTO payment_request (coin, payment_id, sender_address, sender_account, receiver_address, amount, is_done, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdatePaymentRequestPaymentID :execresult
UPDATE payment_request
SET payment_id = ?
WHERE id = ?;

-- name: UpdatePaymentRequestIsDone :execresult
UPDATE payment_request
SET is_done = ?
WHERE coin = ? AND payment_id = ?;

-- name: DeleteAllPaymentRequests :execresult
DELETE FROM payment_request
WHERE coin = ?;
