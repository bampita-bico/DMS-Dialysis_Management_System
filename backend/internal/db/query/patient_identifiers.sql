-- name: CreatePatientIdentifier :one
INSERT INTO patient_identifiers (
    hospital_id, patient_id, id_type, id_value, issuing_country, issuing_authority,
    issued_date, expiry_date, is_verified, verified_by, verified_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: FindByIdentifier :one
SELECT * FROM patient_identifiers
WHERE id_type = $1 AND id_value = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: ListIdentifiersByPatient :many
SELECT * FROM patient_identifiers
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY created_at;
