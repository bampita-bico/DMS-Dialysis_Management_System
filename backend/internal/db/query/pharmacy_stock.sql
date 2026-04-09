-- name: CreatePharmacyStock :one
INSERT INTO pharmacy_stock (
    hospital_id, medication_id, batch_number, quantity_current, unit_cost,
    expiry_date, received_date, supplier_name, storage_location, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetPharmacyStock :one
SELECT * FROM pharmacy_stock WHERE id = $1 AND deleted_at IS NULL;

-- name: ListStockByMedication :many
SELECT * FROM pharmacy_stock
WHERE medication_id = $1 AND deleted_at IS NULL
ORDER BY expiry_date;

-- name: ListLowStock :many
SELECT s.*, m.generic_name, m.reorder_level
FROM pharmacy_stock s
JOIN medications m ON s.medication_id = m.id
WHERE s.hospital_id = $1 AND s.is_low_stock = TRUE AND s.deleted_at IS NULL
ORDER BY m.generic_name;

-- name: ListExpiringStock :many
SELECT s.*, m.generic_name
FROM pharmacy_stock s
JOIN medications m ON s.medication_id = m.id
WHERE s.hospital_id = $1
  AND s.expiry_date <= $2
  AND s.quantity_current > 0
  AND s.deleted_at IS NULL
ORDER BY s.expiry_date;

-- name: UpdateStockQuantity :one
UPDATE pharmacy_stock
SET quantity_current = $2, quantity_available = $2 - quantity_reserved
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeletePharmacyStock :exec
UPDATE pharmacy_stock SET deleted_at = now() WHERE id = $1;
