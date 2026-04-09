-- name: CreateInvoiceItem :one
INSERT INTO invoice_items (
    hospital_id, invoice_id, price_list_id, service_name, service_code,
    quantity, unit_price, discount_amount, tax_amount, total_amount,
    is_claimed, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetInvoiceItem :one
SELECT * FROM invoice_items WHERE id = $1 AND deleted_at IS NULL;

-- name: ListItemsByInvoice :many
SELECT * FROM invoice_items
WHERE invoice_id = $1 AND deleted_at IS NULL
ORDER BY created_at;

-- name: ListClaimedItems :many
SELECT * FROM invoice_items
WHERE hospital_id = $1 AND is_claimed = TRUE AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: MarkItemAsClaimed :one
UPDATE invoice_items
SET is_claimed = TRUE
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteInvoiceItem :exec
UPDATE invoice_items SET deleted_at = now() WHERE id = $1;
