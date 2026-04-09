-- name: CreateSessionComplication :one
INSERT INTO session_complications (
    hospital_id, session_id, patient_id, reported_by, occurred_at,
    complication_type, severity, symptoms, vital_signs_at_event,
    immediate_action_taken, outcome, required_hospitalization,
    was_session_terminated, doctor_notified, doctor_id, family_notified, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
RETURNING *;

-- name: GetSessionComplication :one
SELECT * FROM session_complications WHERE id = $1 AND deleted_at IS NULL;

-- name: ListComplicationsBySession :many
SELECT * FROM session_complications
WHERE session_id = $1 AND deleted_at IS NULL
ORDER BY occurred_at;

-- name: ListComplicationsByPatient :many
SELECT * FROM session_complications
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY occurred_at DESC;

-- name: ListSevereComplications :many
SELECT * FROM session_complications
WHERE hospital_id = $1
  AND severity IN ('severe', 'life_threatening')
  AND deleted_at IS NULL
ORDER BY occurred_at DESC;

-- name: UpdateSessionComplication :one
UPDATE session_complications
SET complication_type = $2, severity = $3, symptoms = $4,
    vital_signs_at_event = $5, immediate_action_taken = $6, outcome = $7,
    required_hospitalization = $8, was_session_terminated = $9,
    doctor_notified = $10, doctor_id = $11, family_notified = $12, notes = $13
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteSessionComplication :exec
UPDATE session_complications SET deleted_at = now() WHERE id = $1;
