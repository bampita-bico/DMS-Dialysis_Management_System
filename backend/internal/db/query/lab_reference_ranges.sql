-- name: CreateLabReferenceRange :one
INSERT INTO lab_reference_ranges (
    hospital_id, test_id, age_group, sex, min_value, max_value,
    critical_low, critical_high, reference_text, is_default
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetLabReferenceRange :one
SELECT * FROM lab_reference_ranges WHERE id = $1 AND deleted_at IS NULL;

-- name: ListReferenceRangesByTest :many
SELECT * FROM lab_reference_ranges
WHERE test_id = $1 AND deleted_at IS NULL
ORDER BY age_group, sex;

-- name: GetReferenceRangeForDemographics :one
SELECT * FROM lab_reference_ranges
WHERE test_id = $1 AND age_group = $2 AND sex = $3 AND deleted_at IS NULL
LIMIT 1;

-- name: GetDefaultReferenceRange :one
SELECT * FROM lab_reference_ranges
WHERE test_id = $1 AND is_default = TRUE AND deleted_at IS NULL
LIMIT 1;

-- name: UpdateLabReferenceRange :one
UPDATE lab_reference_ranges
SET age_group = $2, sex = $3, min_value = $4, max_value = $5,
    critical_low = $6, critical_high = $7, reference_text = $8, is_default = $9
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteLabReferenceRange :exec
UPDATE lab_reference_ranges SET deleted_at = now() WHERE id = $1;
