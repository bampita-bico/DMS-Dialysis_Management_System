-- name: CreateShiftAssignment :one
INSERT INTO shift_assignments (
    hospital_id, staff_id, shift_date, shift_type, shift_start_time,
    shift_end_time, machine_ids, assigned_by, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetShiftAssignment :one
SELECT * FROM shift_assignments WHERE id = $1 AND deleted_at IS NULL;

-- name: ListShiftsByDate :many
SELECT * FROM shift_assignments
WHERE hospital_id = $1 AND shift_date = $2 AND deleted_at IS NULL
ORDER BY shift_start_time;

-- name: ListShiftsByStaff :many
SELECT * FROM shift_assignments
WHERE staff_id = $1 AND shift_date BETWEEN $2 AND $3 AND deleted_at IS NULL
ORDER BY shift_date, shift_start_time;

-- name: ListUnconfirmedShifts :many
SELECT * FROM shift_assignments
WHERE hospital_id = $1
  AND shift_date >= CURRENT_DATE
  AND is_confirmed = FALSE
  AND deleted_at IS NULL
ORDER BY shift_date, shift_start_time;

-- name: ConfirmShift :one
UPDATE shift_assignments
SET is_confirmed = TRUE
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: ClockIn :one
UPDATE shift_assignments
SET clock_in_time = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: ClockOut :one
UPDATE shift_assignments
SET clock_out_time = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteShiftAssignment :exec
UPDATE shift_assignments SET deleted_at = now() WHERE id = $1;
