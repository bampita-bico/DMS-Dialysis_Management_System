-- name: CreateLabResult :one
INSERT INTO lab_results (
    hospital_id, order_item_id, value_text, value_numeric, unit, reference_range,
    is_abnormal, is_critical, status, result_date, result_time, entered_by, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetLabResult :one
SELECT * FROM lab_results WHERE id = $1 AND deleted_at IS NULL;

-- name: GetLabResultByOrderItem :one
SELECT * FROM lab_results WHERE order_item_id = $1 AND deleted_at IS NULL;

-- name: ListLabResultsByOrder :many
SELECT lr.* FROM lab_results lr
JOIN lab_order_items loi ON lr.order_item_id = loi.id
WHERE loi.order_id = $1 AND lr.deleted_at IS NULL
ORDER BY lr.result_date DESC, lr.result_time DESC;

-- name: ListPendingVerificationResults :many
SELECT * FROM lab_results
WHERE hospital_id = $1 AND status = 'preliminary' AND verified_by IS NULL AND deleted_at IS NULL
ORDER BY result_date, result_time;

-- name: ListCriticalResults :many
SELECT * FROM lab_results
WHERE hospital_id = $1 AND is_critical = TRUE AND deleted_at IS NULL
ORDER BY result_date DESC, result_time DESC;

-- name: VerifyLabResult :one
UPDATE lab_results
SET verified_by = $2, verified_at = now(), status = 'final'
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateLabResult :one
UPDATE lab_results
SET value_text = $2, value_numeric = $3, unit = $4, reference_range = $5,
    is_abnormal = $6, is_critical = $7, status = $8, notes = $9
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteLabResult :exec
UPDATE lab_results SET deleted_at = now() WHERE id = $1;
