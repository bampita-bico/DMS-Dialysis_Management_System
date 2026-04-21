-- name: CountSessionsByStatusForDate :many
SELECT status, COUNT(*) AS count
FROM dialysis_sessions
WHERE hospital_id = $1 AND scheduled_date = $2 AND deleted_at IS NULL
GROUP BY status;

-- name: CountActiveSessions :one
SELECT COUNT(*) AS count FROM dialysis_sessions
WHERE hospital_id = $1 AND status = 'in_progress' AND deleted_at IS NULL;

-- name: CountUnacknowledgedCriticalAlerts :one
SELECT COUNT(*) AS count FROM lab_critical_alerts
WHERE hospital_id = $1 AND acknowledged_at IS NULL AND deleted_at IS NULL;

-- name: CountOverdueInvoices :one
SELECT COUNT(*) AS count FROM invoices
WHERE hospital_id = $1 AND status = 'overdue' AND deleted_at IS NULL;

-- name: CountLowStockItems :one
SELECT COUNT(*) AS count FROM consumables_inventory
WHERE hospital_id = $1 AND is_low_stock = TRUE AND quantity_current > 0 AND deleted_at IS NULL;

-- name: CountStaffOnDutyToday :one
SELECT COUNT(DISTINCT sa.staff_id) AS count
FROM shift_assignments sa
WHERE sa.hospital_id = $1
  AND sa.shift_date = $2
  AND sa.status IN ('confirmed', 'clocked_in')
  AND sa.deleted_at IS NULL;

-- name: CountActivePatients :one
SELECT COUNT(*) AS count FROM patients
WHERE hospital_id = $1 AND is_active = TRUE AND deleted_at IS NULL;
