-- name: CreateLabPanel :one
INSERT INTO lab_panels (
    hospital_id, code, name, description, tests, cost_amount
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetLabPanel :one
SELECT * FROM lab_panels WHERE id = $1 AND deleted_at IS NULL;

-- name: ListLabPanelsByHospital :many
SELECT * FROM lab_panels
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY name;

-- name: ListActiveLabPanels :many
SELECT * FROM lab_panels
WHERE hospital_id = $1 AND is_active = TRUE AND deleted_at IS NULL
ORDER BY name;

-- name: GetLabPanelByCode :one
SELECT * FROM lab_panels
WHERE hospital_id = $1 AND code = $2 AND deleted_at IS NULL;

-- name: UpdateLabPanel :one
UPDATE lab_panels
SET code = $2, name = $3, description = $4, tests = $5,
    cost_amount = $6, is_active = $7
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteLabPanel :exec
UPDATE lab_panels SET deleted_at = now() WHERE id = $1;
