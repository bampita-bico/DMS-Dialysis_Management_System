-- name: CreateConsumable :one
INSERT INTO consumables (
    hospital_id, name, category, unit, manufacturer, model,
    is_reusable, max_reuse_count, min_stock_level, cost_per_unit, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetConsumable :one
SELECT * FROM consumables WHERE id = $1 AND deleted_at IS NULL;

-- name: ListConsumablesByHospital :many
SELECT * FROM consumables
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY name;

-- name: ListConsumablesByCategory :many
SELECT * FROM consumables
WHERE hospital_id = $1 AND category = $2 AND deleted_at IS NULL
ORDER BY name;

-- name: ListActiveConsumables :many
SELECT * FROM consumables
WHERE hospital_id = $1 AND is_active = TRUE AND deleted_at IS NULL
ORDER BY name;

-- name: ListReusableConsumables :many
SELECT * FROM consumables
WHERE hospital_id = $1 AND is_reusable = TRUE AND deleted_at IS NULL
ORDER BY name;

-- name: DeleteConsumable :exec
UPDATE consumables SET deleted_at = now() WHERE id = $1;
