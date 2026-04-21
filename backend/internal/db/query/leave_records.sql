-- name: CreateLeaveRecord :one
INSERT INTO leave_records (
    hospital_id, staff_id, leave_type, start_date, end_date,
    days_requested, reason, status
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetLeaveRecord :one
SELECT * FROM leave_records WHERE id = $1 AND deleted_at IS NULL;

-- name: ListLeaveByStaff :many
SELECT * FROM leave_records
WHERE staff_id = $1 AND deleted_at IS NULL
ORDER BY start_date DESC;

-- name: ListPendingLeave :many
SELECT * FROM leave_records
WHERE hospital_id = $1 AND status = 'pending' AND deleted_at IS NULL
ORDER BY requested_at ASC;

-- name: ListLeaveByDateRange :many
SELECT * FROM leave_records
WHERE hospital_id = $1
  AND start_date <= $3
  AND end_date >= $2
  AND status = 'approved'
  AND deleted_at IS NULL
ORDER BY start_date;

-- name: ApproveLeave :one
UPDATE leave_records
SET status = 'approved', approved_by = $2, approved_at = now(), days_approved = $3
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: RejectLeave :one
UPDATE leave_records
SET status = 'rejected', rejection_reason = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: HasApprovedLeaveOnDate :one
SELECT EXISTS(
    SELECT 1 FROM leave_records
    WHERE staff_id = $1
      AND start_date <= $2
      AND end_date >= $2
      AND status = 'approved'
      AND deleted_at IS NULL
) AS on_leave;

-- name: DeleteLeaveRecord :exec
UPDATE leave_records SET deleted_at = now() WHERE id = $1;
