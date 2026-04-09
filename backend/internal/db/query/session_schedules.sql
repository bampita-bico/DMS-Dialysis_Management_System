-- name: CreateSessionSchedule :one
INSERT INTO session_schedules (
    hospital_id, patient_id, machine_id, shift, days_of_week, frequency_weeks,
    modality, prescribed_duration_mins, effective_from, effective_until,
    is_active, created_by, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetSessionSchedule :one
SELECT * FROM session_schedules WHERE id = $1 AND deleted_at IS NULL;

-- name: ListActiveSchedulesByPatient :many
SELECT * FROM session_schedules
WHERE patient_id = $1 AND is_active = TRUE AND deleted_at IS NULL
ORDER BY effective_from DESC;

-- name: ListActiveSchedulesByMachine :many
SELECT * FROM session_schedules
WHERE machine_id = $1 AND is_active = TRUE AND deleted_at IS NULL
ORDER BY effective_from DESC;

-- name: ListSchedulesByHospital :many
SELECT * FROM session_schedules
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateSessionSchedule :one
UPDATE session_schedules
SET machine_id = $2, shift = $3, days_of_week = $4, frequency_weeks = $5,
    modality = $6, prescribed_duration_mins = $7, effective_from = $8,
    effective_until = $9, is_active = $10, notes = $11
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeactivateSchedule :one
UPDATE session_schedules
SET is_active = FALSE, effective_until = CURRENT_DATE
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteSessionSchedule :exec
UPDATE session_schedules SET deleted_at = now() WHERE id = $1;
