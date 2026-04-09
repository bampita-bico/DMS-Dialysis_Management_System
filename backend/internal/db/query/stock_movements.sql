-- name: CreateStockMovement :one
INSERT INTO stock_movements (
    hospital_id, stock_id, movement_type, quantity, reference_type,
    reference_id, performed_by, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetStockMovement :one
SELECT * FROM stock_movements WHERE id = $1 AND deleted_at IS NULL;

-- name: ListMovementsByStock :many
SELECT * FROM stock_movements
WHERE stock_id = $1 AND deleted_at IS NULL
ORDER BY movement_date DESC
LIMIT $2 OFFSET $3;

-- name: ListMovementsByHospital :many
SELECT * FROM stock_movements
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY movement_date DESC
LIMIT $2 OFFSET $3;

-- name: ListMovementsByType :many
SELECT * FROM stock_movements
WHERE hospital_id = $1 AND movement_type = $2 AND deleted_at IS NULL
ORDER BY movement_date DESC
LIMIT $3 OFFSET $4;

-- name: GetMovementsByReference :many
SELECT * FROM stock_movements
WHERE reference_type = $1 AND reference_id = $2 AND deleted_at IS NULL
ORDER BY movement_date;

-- name: DeleteStockMovement :exec
UPDATE stock_movements SET deleted_at = now() WHERE id = $1;
