-- name: CreateAllergy :one
INSERT INTO allergies (
    hospital_id, patient_id, allergen, category, reaction, severity, onset_date, notes, recorded_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetActiveAllergies :many
SELECT * FROM allergies
WHERE patient_id = $1 AND is_active = TRUE AND deleted_at IS NULL
ORDER BY severity DESC, created_at;

-- name: CheckDrugAllergy :one
SELECT * FROM allergies
WHERE patient_id = $1 AND category = 'drug' AND allergen ILIKE $2 AND is_active = TRUE AND deleted_at IS NULL
LIMIT 1;
