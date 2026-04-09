-- name: CreateLabTest :one
INSERT INTO lab_test_catalog (
    hospital_id, code, name, category, unit, turnaround_hrs,
    specimen_type, specimen_volume_ml, requires_fasting, cost_amount, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetLabTest :one
SELECT * FROM lab_test_catalog WHERE id = $1 AND deleted_at IS NULL;

-- name: ListLabTestsByHospital :many
SELECT * FROM lab_test_catalog
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY category, name;

-- name: ListActiveLabTests :many
SELECT * FROM lab_test_catalog
WHERE hospital_id = $1 AND is_active = TRUE AND deleted_at IS NULL
ORDER BY category, name;

-- name: ListLabTestsByCategory :many
SELECT * FROM lab_test_catalog
WHERE hospital_id = $1 AND category = $2 AND deleted_at IS NULL
ORDER BY name;

-- name: GetLabTestByCode :one
SELECT * FROM lab_test_catalog
WHERE hospital_id = $1 AND code = $2 AND deleted_at IS NULL;

-- name: UpdateLabTest :one
UPDATE lab_test_catalog
SET code = $2, name = $3, category = $4, unit = $5, turnaround_hrs = $6,
    specimen_type = $7, specimen_volume_ml = $8, requires_fasting = $9,
    cost_amount = $10, is_active = $11, notes = $12
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteLabTest :exec
UPDATE lab_test_catalog SET deleted_at = now() WHERE id = $1;
