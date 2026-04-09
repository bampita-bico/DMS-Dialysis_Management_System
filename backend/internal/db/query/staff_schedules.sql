-- name: CreateStaffSchedule :one
INSERT INTO staff_schedules (
    hospital_id, staff_id, schedule_name, schedule_type, effective_from,
    effective_until, monday_shift, tuesday_shift, wednesday_shift,
    thursday_shift, friday_shift, saturday_shift, sunday_shift, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: GetStaffSchedule :one
SELECT * FROM staff_schedules WHERE id = $1 AND deleted_at IS NULL;

-- name: ListSchedulesByStaff :many
SELECT * FROM staff_schedules
WHERE staff_id = $1 AND deleted_at IS NULL
ORDER BY effective_from DESC;

-- name: GetActiveScheduleByStaff :one
SELECT * FROM staff_schedules
WHERE staff_id = $1
  AND effective_from <= CURRENT_DATE
  AND (effective_until IS NULL OR effective_until >= CURRENT_DATE)
  AND is_active = TRUE
  AND deleted_at IS NULL
ORDER BY effective_from DESC
LIMIT 1;

-- name: DeleteStaffSchedule :exec
UPDATE staff_schedules SET deleted_at = now() WHERE id = $1;
