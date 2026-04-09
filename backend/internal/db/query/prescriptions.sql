-- name: CreatePrescription :one
INSERT INTO prescriptions (
    hospital_id, patient_id, session_id, prescribed_by, prescribed_date,
    prescribed_time, valid_from, valid_until, diagnosis, clinical_notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetPrescription :one
SELECT * FROM prescriptions WHERE id = $1 AND deleted_at IS NULL;

-- name: ListMedicationPrescriptionsByPatient :many
SELECT * FROM prescriptions
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY prescribed_date DESC, prescribed_time DESC
LIMIT $2 OFFSET $3;

-- name: ListActivePrescriptions :many
SELECT * FROM prescriptions
WHERE hospital_id = $1 AND status = 'active' AND deleted_at IS NULL
ORDER BY prescribed_date DESC;

-- name: VerifyPrescription :one
UPDATE prescriptions
SET pharmacist_verified_by = $2, pharmacist_verified_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DispensePrescription :one
UPDATE prescriptions
SET dispensed_by = $2, dispensed_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: CancelPrescription :one
UPDATE prescriptions
SET status = 'cancelled', cancelled_by = $2, cancelled_at = now(), cancellation_reason = $3
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeletePrescription :exec
UPDATE prescriptions SET deleted_at = now() WHERE id = $1;
