-- name: CreateEquipment :one
INSERT INTO equipment (
    hospital_id, name, category, serial_number, model, manufacturer,
    purchase_date, purchase_cost, warranty_expiry_date, status,
    location, department_id, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetEquipment :one
SELECT * FROM equipment WHERE id = $1 AND deleted_at IS NULL;

-- name: ListEquipmentByHospital :many
SELECT * FROM equipment
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY name;

-- name: ListEquipmentByCategory :many
SELECT * FROM equipment
WHERE hospital_id = $1 AND category = $2 AND deleted_at IS NULL
ORDER BY name;

-- name: ListEquipmentByStatus :many
SELECT * FROM equipment
WHERE hospital_id = $1 AND status = $2 AND deleted_at IS NULL
ORDER BY name;

-- name: ListEquipmentByDepartment :many
SELECT * FROM equipment
WHERE department_id = $1 AND deleted_at IS NULL
ORDER BY name;

-- name: UpdateEquipmentStatus :one
UPDATE equipment
SET status = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteEquipment :exec
UPDATE equipment SET deleted_at = now() WHERE id = $1;
