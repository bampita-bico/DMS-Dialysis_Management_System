-- name: CreateInvoice :one
INSERT INTO invoices (
    hospital_id, account_id, patient_id, session_id, invoice_number,
    invoice_date, due_date, total_amount, discount_amount, tax_amount,
    net_amount, status, issued_by, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: GetInvoice :one
SELECT * FROM invoices WHERE id = $1 AND deleted_at IS NULL;

-- name: GetInvoiceByNumber :one
SELECT * FROM invoices
WHERE hospital_id = $1 AND invoice_number = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: ListInvoicesByPatient :many
SELECT * FROM invoices
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY invoice_date DESC
LIMIT $2 OFFSET $3;

-- name: ListInvoicesByAccount :many
SELECT * FROM invoices
WHERE account_id = $1 AND deleted_at IS NULL
ORDER BY invoice_date DESC;

-- name: ListInvoicesByStatus :many
SELECT * FROM invoices
WHERE hospital_id = $1 AND status = $2 AND deleted_at IS NULL
ORDER BY invoice_date DESC;

-- name: ListOverdueInvoices :many
SELECT * FROM invoices
WHERE hospital_id = $1
  AND status IN ('issued', 'partially_paid', 'overdue')
  AND due_date < CURRENT_DATE
  AND deleted_at IS NULL
ORDER BY due_date ASC;

-- name: UpdateInvoiceStatus :one
UPDATE invoices
SET status = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateInvoicePayment :one
UPDATE invoices
SET paid_amount = $2, balance_due = net_amount - $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteInvoice :exec
UPDATE invoices SET deleted_at = now() WHERE id = $1;
