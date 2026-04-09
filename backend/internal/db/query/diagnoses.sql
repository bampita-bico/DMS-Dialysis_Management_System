-- name: CreateDiagnosis :one
INSERT INTO diagnoses (
    hospital_id, patient_id, icd10_code, description, diagnosis_type, diagnosed_by, admission_id, notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetPrimaryDiagnosis :one
SELECT * FROM diagnoses
WHERE patient_id = $1 AND diagnosis_type = 'primary' AND deleted_at IS NULL
ORDER BY diagnosed_at DESC
LIMIT 1;

-- name: ListDiagnosesByPatient :many
SELECT * FROM diagnoses
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY diagnosed_at DESC;
