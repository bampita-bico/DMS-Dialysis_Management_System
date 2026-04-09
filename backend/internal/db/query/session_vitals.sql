-- name: CreateSessionVital :one
INSERT INTO session_vitals (
    hospital_id, session_id, patient_id, recorded_by, recorded_at,
    time_on_dialysis_mins, bp_systolic, bp_diastolic, heart_rate, temperature,
    spo2, respiratory_rate, blood_flow_actual, dialysate_flow_actual,
    venous_pressure, arterial_pressure, tmp, uf_removed_so_far,
    conductivity_actual, has_hypotension_alert, has_hypertension_alert,
    has_tachycardia_alert, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
RETURNING *;

-- name: GetSessionVital :one
SELECT * FROM session_vitals WHERE id = $1 AND deleted_at IS NULL;

-- name: ListVitalsBySession :many
SELECT * FROM session_vitals
WHERE session_id = $1 AND deleted_at IS NULL
ORDER BY recorded_at;

-- name: ListVitalsWithAlerts :many
SELECT * FROM session_vitals
WHERE session_id = $1
  AND (has_hypotension_alert = TRUE OR has_hypertension_alert = TRUE OR has_tachycardia_alert = TRUE)
  AND deleted_at IS NULL
ORDER BY recorded_at;

-- name: AcknowledgeAlert :one
UPDATE session_vitals
SET alert_acknowledged_by = $2, alert_acknowledged_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteSessionVital :exec
UPDATE session_vitals SET deleted_at = now() WHERE id = $1;
