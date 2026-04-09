-- name: CreateLabCriticalAlert :one
INSERT INTO lab_critical_alerts (
    hospital_id, result_id, patient_id, test_name, critical_value,
    reference_range, severity
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetLabCriticalAlert :one
SELECT * FROM lab_critical_alerts WHERE id = $1 AND deleted_at IS NULL;

-- name: ListLabCriticalAlertsByPatient :many
SELECT * FROM lab_critical_alerts
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY alerted_at DESC;

-- name: ListUnacknowledgedCriticalAlerts :many
SELECT * FROM lab_critical_alerts
WHERE hospital_id = $1 AND acknowledged_at IS NULL AND deleted_at IS NULL
ORDER BY alerted_at;

-- name: ListCriticalAlertsByDateRange :many
SELECT * FROM lab_critical_alerts
WHERE hospital_id = $1
  AND alerted_at >= $2
  AND alerted_at <= $3
  AND deleted_at IS NULL
ORDER BY alerted_at DESC;

-- name: AcknowledgeCriticalAlert :one
UPDATE lab_critical_alerts
SET acknowledged_by = $2, acknowledged_at = now(), action_taken = $3
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: NotifyDoctorOfCriticalAlert :one
UPDATE lab_critical_alerts
SET doctor_notified = TRUE, doctor_notified_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteLabCriticalAlert :exec
UPDATE lab_critical_alerts SET deleted_at = now() WHERE id = $1;
