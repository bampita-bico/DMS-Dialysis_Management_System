-- name: CreateLabOrderItem :one
INSERT INTO lab_order_items (
    hospital_id, order_id, test_id, specimen_type, notes
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetLabOrderItem :one
SELECT * FROM lab_order_items WHERE id = $1 AND deleted_at IS NULL;

-- name: ListLabOrderItemsByOrder :many
SELECT * FROM lab_order_items
WHERE order_id = $1 AND deleted_at IS NULL
ORDER BY created_at;

-- name: ListLabOrderItemsByStatus :many
SELECT * FROM lab_order_items
WHERE hospital_id = $1 AND status = $2 AND deleted_at IS NULL
ORDER BY created_at;

-- name: UpdateLabOrderItemStatus :one
UPDATE lab_order_items
SET status = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: CollectSpecimen :one
UPDATE lab_order_items
SET specimen_collected_by = $2, specimen_collected_at = now(),
    specimen_barcode = $3, status = 'collected'
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: RejectSpecimen :one
UPDATE lab_order_items
SET specimen_rejected = TRUE, rejection_reason = $2, status = 'rejected'
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteLabOrderItem :exec
UPDATE lab_order_items SET deleted_at = now() WHERE id = $1;
