-- name: CreatePatientContact :one
INSERT INTO patient_contacts (
    hospital_id, patient_id, contact_type, value, label, is_primary, is_verified
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetPrimaryPhone :one
SELECT * FROM patient_contacts
WHERE patient_id = $1 AND contact_type = 'phone' AND is_primary = TRUE AND deleted_at IS NULL
LIMIT 1;

-- name: ListContactsByPatient :many
SELECT * FROM patient_contacts
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY is_primary DESC, created_at;

-- name: UpdateContactPrimary :exec
UPDATE patient_contacts
SET is_primary = FALSE
WHERE patient_id = $1 AND contact_type = $2 AND id != $3;
