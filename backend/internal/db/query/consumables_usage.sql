-- name: CreateConsumablesUsage :one
INSERT INTO consumables_usage (
    hospital_id, session_id, consumable_id, inventory_id,
    quantity_used, reuse_number, recorded_by, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetConsumablesUsage :one
SELECT * FROM consumables_usage WHERE id = $1 AND deleted_at IS NULL;

-- name: ListUsageBySession :many
SELECT u.*, c.name AS consumable_name, c.category
FROM consumables_usage u
JOIN consumables c ON u.consumable_id = c.id
WHERE u.session_id = $1 AND u.deleted_at IS NULL
ORDER BY c.name;

-- name: ListUsageByConsumable :many
SELECT * FROM consumables_usage
WHERE consumable_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetTotalUsageByConsumable :one
SELECT consumable_id, SUM(quantity_used) AS total_quantity
FROM consumables_usage
WHERE consumable_id = $1 AND deleted_at IS NULL
GROUP BY consumable_id;

-- name: DeleteConsumablesUsage :exec
UPDATE consumables_usage SET deleted_at = now() WHERE id = $1;
