-- name: CreatePowerOutageLog :one
INSERT INTO power_outage_logs (
    hospital_id, outage_start, outage_end, duration_mins, affected_sessions,
    sessions_terminated_count, sessions_paused_count, generator_available,
    generator_start_delay_mins, backup_power_duration_mins, logged_by,
    incident_severity, patient_safety_impact, equipment_damage,
    data_loss_occurred, recovery_actions, utility_company_notified,
    incident_report_filed, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
RETURNING *;

-- name: GetPowerOutageLog :one
SELECT * FROM power_outage_logs WHERE id = $1 AND deleted_at IS NULL;

-- name: ListPowerOutagesByHospital :many
SELECT * FROM power_outage_logs
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY outage_start DESC;

-- name: ListOngoingOutages :many
SELECT * FROM power_outage_logs
WHERE hospital_id = $1
  AND outage_end IS NULL
  AND deleted_at IS NULL
ORDER BY outage_start;

-- name: ListSevereOutages :many
SELECT * FROM power_outage_logs
WHERE hospital_id = $1
  AND incident_severity IN ('severe', 'life_threatening')
  AND deleted_at IS NULL
ORDER BY outage_start DESC;

-- name: EndPowerOutage :one
UPDATE power_outage_logs
SET outage_end = $2, duration_mins = $3, recovery_actions = $4
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeletePowerOutageLog :exec
UPDATE power_outage_logs SET deleted_at = now() WHERE id = $1;
