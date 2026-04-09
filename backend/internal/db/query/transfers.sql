-- name: CreateTransfer :one
INSERT INTO transfers (
    hospital_id, patient_id, from_hospital_id, to_hospital_id, from_department,
    to_department, reason, reason_notes, requested_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: ListTransfersByPatient :many
SELECT * FROM transfers
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY requested_at DESC;

-- name: UpdateTransferStatus :one
UPDATE transfers
SET status = $2, approved_by = $3, approved_at = $4
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;
