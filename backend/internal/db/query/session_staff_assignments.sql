-- name: CreateSessionStaffAssignment :one
INSERT INTO session_staff_assignments (
    hospital_id, session_id, staff_id, role, assigned_by, notes
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetSessionStaffAssignment :one
SELECT * FROM session_staff_assignments WHERE id = $1 AND deleted_at IS NULL;

-- name: ListStaffAssignmentsBySession :many
SELECT * FROM session_staff_assignments
WHERE session_id = $1 AND deleted_at IS NULL
ORDER BY assigned_at;

-- name: ListActiveAssignmentsByStaff :many
SELECT * FROM session_staff_assignments
WHERE staff_id = $1
  AND started_at IS NOT NULL
  AND completed_at IS NULL
  AND deleted_at IS NULL
ORDER BY started_at;

-- name: StartAssignment :one
UPDATE session_staff_assignments
SET started_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: CompleteAssignment :one
UPDATE session_staff_assignments
SET completed_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteSessionStaffAssignment :exec
UPDATE session_staff_assignments SET deleted_at = now() WHERE id = $1;
