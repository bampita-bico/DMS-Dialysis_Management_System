-- name: CreateComorbidity :one
INSERT INTO comorbidities (
    hospital_id, patient_id, condition, icd10_code, status, diagnosed_at, diagnosed_by, notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: ListComorbiditiesByPatient :many
SELECT * FROM comorbidities
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY diagnosed_at DESC;

-- name: UpdateComorbidityStatus :one
UPDATE comorbidities
SET status = $2, resolved_at = $3
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;
