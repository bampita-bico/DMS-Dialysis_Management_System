-- name: CreatePatientFlag :one
INSERT INTO patient_flags (
    hospital_id, patient_id, flag_type, reason, flagged_by, expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetActiveFlags :many
SELECT * FROM patient_flags
WHERE patient_id = $1 AND is_active = TRUE AND deleted_at IS NULL
ORDER BY flagged_at DESC;

-- name: GetInfectiousFlags :many
SELECT * FROM patient_flags
WHERE patient_id = $1 AND flag_type IN ('infectious','hiv_positive','hepatitis_b','hepatitis_c')
    AND is_active = TRUE AND deleted_at IS NULL;

-- name: ResolveFlag :exec
UPDATE patient_flags
SET is_active = FALSE, resolved_by = $2, resolved_at = now()
WHERE id = $1;
