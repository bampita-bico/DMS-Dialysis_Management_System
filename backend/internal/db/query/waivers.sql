-- name: CreateWaiver :one
INSERT INTO waivers (
    hospital_id, invoice_id, patient_id, waiver_number, waiver_amount,
    waiver_percentage, waiver_type, reason, requested_by, status, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetWaiver :one
SELECT * FROM waivers WHERE id = $1 AND deleted_at IS NULL;

-- name: GetWaiverByNumber :one
SELECT * FROM waivers
WHERE hospital_id = $1 AND waiver_number = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: ListWaiversByInvoice :many
SELECT * FROM waivers
WHERE invoice_id = $1 AND deleted_at IS NULL
ORDER BY requested_at DESC;

-- name: ListWaiversByPatient :many
SELECT * FROM waivers
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY requested_at DESC;

-- name: ListPendingWaivers :many
SELECT * FROM waivers
WHERE hospital_id = $1 AND status = 'pending' AND deleted_at IS NULL
ORDER BY requested_at ASC;

-- name: ApproveWaiver :one
UPDATE waivers
SET status = 'approved', approved_by = $2, approved_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: RejectWaiver :one
UPDATE waivers
SET status = 'rejected', rejection_reason = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteWaiver :exec
UPDATE waivers SET deleted_at = now() WHERE id = $1;
