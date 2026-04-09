-- name: CreatePayment :one
INSERT INTO payments (
    hospital_id, invoice_id, account_id, patient_id, payment_date,
    payment_time, amount, payment_method, reference_number,
    mobile_money_number, bank_name, cheque_number, card_last_four,
    received_by, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING *;

-- name: GetPayment :one
SELECT * FROM payments WHERE id = $1 AND deleted_at IS NULL;

-- name: ListPaymentsByInvoice :many
SELECT * FROM payments
WHERE invoice_id = $1 AND deleted_at IS NULL
ORDER BY payment_date DESC, payment_time DESC;

-- name: ListPaymentsByPatient :many
SELECT * FROM payments
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY payment_date DESC, payment_time DESC
LIMIT $2 OFFSET $3;

-- name: ListPaymentsByDate :many
SELECT * FROM payments
WHERE hospital_id = $1
  AND payment_date BETWEEN $2 AND $3
  AND deleted_at IS NULL
ORDER BY payment_date DESC, payment_time DESC;

-- name: ListPaymentsByMethod :many
SELECT * FROM payments
WHERE hospital_id = $1 AND payment_method = $2 AND deleted_at IS NULL
ORDER BY payment_date DESC;

-- name: GetPaymentTotal :one
SELECT COALESCE(SUM(amount), 0) AS total
FROM payments
WHERE hospital_id = $1
  AND payment_date BETWEEN $2 AND $3
  AND deleted_at IS NULL;

-- name: DeletePayment :exec
UPDATE payments SET deleted_at = now() WHERE id = $1;
