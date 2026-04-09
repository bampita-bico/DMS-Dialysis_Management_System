-- name: CreateMedication :one
INSERT INTO medications (
    hospital_id, generic_name, brand_names, drug_class, form, strength, unit,
    is_controlled, requires_prescription, is_essential_who, storage_conditions,
    cost_per_unit, reorder_level, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: GetMedication :one
SELECT * FROM medications WHERE id = $1 AND deleted_at IS NULL;

-- name: ListMedicationsByHospital :many
SELECT * FROM medications
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY generic_name;

-- name: ListActiveMedications :many
SELECT * FROM medications
WHERE hospital_id = $1 AND is_active = TRUE AND deleted_at IS NULL
ORDER BY generic_name;

-- name: ListMedicationsByClass :many
SELECT * FROM medications
WHERE hospital_id = $1 AND drug_class = $2 AND deleted_at IS NULL
ORDER BY generic_name;

-- name: SearchMedications :many
SELECT * FROM medications
WHERE hospital_id = $1
  AND (generic_name ILIKE '%' || $2 || '%' OR drug_class ILIKE '%' || $2 || '%')
  AND deleted_at IS NULL
ORDER BY generic_name
LIMIT 50;

-- name: UpdateMedication :one
UPDATE medications
SET generic_name = $2, brand_names = $3, drug_class = $4, form = $5,
    strength = $6, unit = $7, is_controlled = $8, requires_prescription = $9,
    is_essential_who = $10, storage_conditions = $11, cost_per_unit = $12,
    reorder_level = $13, is_active = $14, notes = $15
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteMedication :exec
UPDATE medications SET deleted_at = now() WHERE id = $1;
