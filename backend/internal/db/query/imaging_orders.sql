-- name: CreateImagingOrder :one
INSERT INTO imaging_orders (
    hospital_id, patient_id, session_id, ordered_by, order_date, order_time,
    modality, body_part, laterality, clinical_indication, priority, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetImagingOrder :one
SELECT * FROM imaging_orders WHERE id = $1 AND deleted_at IS NULL;

-- name: ListImagingOrdersByPatient :many
SELECT * FROM imaging_orders
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY order_date DESC, order_time DESC
LIMIT $2 OFFSET $3;

-- name: ListImagingOrdersBySession :many
SELECT * FROM imaging_orders
WHERE session_id = $1 AND deleted_at IS NULL
ORDER BY order_date, order_time;

-- name: ListImagingOrdersByModality :many
SELECT * FROM imaging_orders
WHERE hospital_id = $1 AND modality = $2 AND deleted_at IS NULL
ORDER BY order_date DESC, order_time DESC;

-- name: ListPendingImagingOrders :many
SELECT * FROM imaging_orders
WHERE hospital_id = $1 AND status = 'pending' AND deleted_at IS NULL
ORDER BY priority DESC, order_date, order_time;

-- name: ScheduleImagingOrder :one
UPDATE imaging_orders
SET scheduled_date = $2, scheduled_time = $3
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: PerformImagingOrder :one
UPDATE imaging_orders
SET performed_at = now(), performed_by = $2, status = 'completed'
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: CancelImagingOrder :one
UPDATE imaging_orders
SET status = 'cancelled', cancelled_by = $2, cancelled_at = now(), cancellation_reason = $3
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteImagingOrder :exec
UPDATE imaging_orders SET deleted_at = now() WHERE id = $1;
