-- name: CreateNextOfKin :one
INSERT INTO next_of_kin (
    hospital_id, patient_id, full_name, relationship, phone_primary, phone_secondary,
    address, national_id, is_legal_guardian, is_emergency_contact, priority_order
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetPrimaryContact :one
SELECT * FROM next_of_kin
WHERE patient_id = $1 AND is_emergency_contact = TRUE AND deleted_at IS NULL
ORDER BY priority_order
LIMIT 1;

-- name: GetLegalGuardian :one
SELECT * FROM next_of_kin
WHERE patient_id = $1 AND is_legal_guardian = TRUE AND deleted_at IS NULL
LIMIT 1;

-- name: ListNextOfKinByPatient :many
SELECT * FROM next_of_kin
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY priority_order;
