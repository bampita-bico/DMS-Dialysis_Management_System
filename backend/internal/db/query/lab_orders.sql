-- name: CreateLabOrder :one
INSERT INTO lab_orders (
    hospital_id, patient_id, session_id, ordered_by, order_date, order_time,
    priority, clinical_notes, diagnosis_code
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetLabOrder :one
SELECT * FROM lab_orders WHERE id = $1 AND deleted_at IS NULL;

-- name: ListLabOrdersByPatient :many
SELECT * FROM lab_orders
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY order_date DESC, order_time DESC
LIMIT $2 OFFSET $3;

-- name: ListLabOrdersBySession :many
SELECT * FROM lab_orders
WHERE session_id = $1 AND deleted_at IS NULL
ORDER BY order_date, order_time;

-- name: ListLabOrdersByStatus :many
SELECT * FROM lab_orders
WHERE hospital_id = $1 AND status = $2 AND deleted_at IS NULL
ORDER BY priority DESC, order_date, order_time;

-- name: ListPendingLabOrders :many
SELECT * FROM lab_orders
WHERE hospital_id = $1 AND status = 'pending' AND deleted_at IS NULL
ORDER BY priority DESC, order_date, order_time;

-- name: UpdateLabOrderStatus :one
UPDATE lab_orders
SET status = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: CancelLabOrder :one
UPDATE lab_orders
SET status = 'cancelled', cancelled_by = $2, cancelled_at = now(), cancellation_reason = $3
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteLabOrder :exec
UPDATE lab_orders SET deleted_at = now() WHERE id = $1;
