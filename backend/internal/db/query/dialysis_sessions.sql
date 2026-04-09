-- name: CreateDialysisSession :one
INSERT INTO dialysis_sessions (
    hospital_id, patient_id, schedule_id, machine_id, access_id, modality, shift,
    status, scheduled_date, scheduled_start_time, prescribed_duration_mins,
    primary_nurse_id, supervising_doctor_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetDialysisSession :one
SELECT * FROM dialysis_sessions WHERE id = $1 AND deleted_at IS NULL;

-- name: ListSessionsByPatient :many
SELECT * FROM dialysis_sessions
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY scheduled_date DESC, scheduled_start_time DESC
LIMIT $2 OFFSET $3;

-- name: ListSessionsByDate :many
SELECT * FROM dialysis_sessions
WHERE hospital_id = $1 AND scheduled_date = $2 AND deleted_at IS NULL
ORDER BY scheduled_start_time;

-- name: ListActiveSessionsByMachine :many
SELECT * FROM dialysis_sessions
WHERE machine_id = $1 AND status = 'in_progress' AND deleted_at IS NULL
ORDER BY actual_start_time DESC;

-- name: ListActiveSessions :many
SELECT * FROM dialysis_sessions
WHERE hospital_id = $1 AND status = 'in_progress' AND deleted_at IS NULL
ORDER BY actual_start_time;

-- name: UpdateSessionStatus :one
UPDATE dialysis_sessions
SET status = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: StartSession :one
UPDATE dialysis_sessions
SET status = 'in_progress', actual_start_time = now(),
    pre_weight_kg = $2, pre_bp_systolic = $3, pre_bp_diastolic = $4,
    pre_hr = $5, pre_temp = $6
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: CompleteSession :one
UPDATE dialysis_sessions
SET status = 'completed', actual_end_time = now(), actual_duration_mins = $2,
    post_weight_kg = $3, post_bp_systolic = $4, post_bp_diastolic = $5,
    post_hr = $6, was_patient_reviewed = $7, reviewed_by = $8, session_notes = $9
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: AbortSession :one
UPDATE dialysis_sessions
SET status = 'aborted', actual_end_time = now(), aborted_reason = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteDialysisSession :exec
UPDATE dialysis_sessions SET deleted_at = now() WHERE id = $1;
