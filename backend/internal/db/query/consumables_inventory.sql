-- name: CreateConsumablesInventory :one
INSERT INTO consumables_inventory (
    hospital_id, consumable_id, batch_number, quantity_current,
    quantity_reserved, quantity_available, unit_cost, expiry_date,
    received_date, supplier_name, storage_location, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetConsumablesInventory :one
SELECT * FROM consumables_inventory WHERE id = $1 AND deleted_at IS NULL;

-- name: ListInventoryByConsumable :many
SELECT * FROM consumables_inventory
WHERE consumable_id = $1 AND deleted_at IS NULL
ORDER BY expiry_date ASC NULLS LAST;

-- name: ListLowStockInventory :many
SELECT i.*, c.name AS consumable_name, c.min_stock_level
FROM consumables_inventory i
JOIN consumables c ON i.consumable_id = c.id
WHERE i.hospital_id = $1 AND i.is_low_stock = TRUE AND i.deleted_at IS NULL
ORDER BY c.name;

-- name: ListExpiringInventory :many
SELECT i.*, c.name AS consumable_name
FROM consumables_inventory i
JOIN consumables c ON i.consumable_id = c.id
WHERE i.hospital_id = $1
  AND i.expiry_date <= $2
  AND i.quantity_current > 0
  AND i.deleted_at IS NULL
ORDER BY i.expiry_date ASC;

-- name: UpdateInventoryQuantity :one
UPDATE consumables_inventory
SET quantity_current = $2, quantity_available = $2 - quantity_reserved
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: GetAvailableInventoryBatch :one
SELECT * FROM consumables_inventory
WHERE consumable_id = $1
  AND hospital_id = $2
  AND quantity_available > 0
  AND deleted_at IS NULL
ORDER BY expiry_date ASC NULLS LAST
LIMIT 1;

-- name: DeductInventory :one
UPDATE consumables_inventory
SET quantity_current = quantity_current - $2,
    quantity_available = quantity_available - $2,
    is_low_stock = CASE WHEN (quantity_current - $2) <= (
        SELECT COALESCE(c.min_stock_level, 0) FROM consumables c WHERE c.id = consumables_inventory.consumable_id
    ) THEN TRUE ELSE FALSE END
WHERE consumables_inventory.id = $1 AND quantity_available >= $2 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteConsumablesInventory :exec
UPDATE consumables_inventory SET deleted_at = now() WHERE id = $1;
