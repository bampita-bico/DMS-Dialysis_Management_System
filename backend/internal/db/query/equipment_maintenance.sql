-- name: CreateEquipmentMaintenance :one
INSERT INTO equipment_maintenance (
    hospital_id, equipment_id, maintenance_type, scheduled_date,
    performed_date, performed_by, technician_name, technician_company,
    next_due_date, cost, findings, actions_taken, parts_replaced, is_completed, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING *;

-- name: GetEquipmentMaintenance :one
SELECT * FROM equipment_maintenance WHERE id = $1 AND deleted_at IS NULL;

-- name: ListMaintenanceByEquipment :many
SELECT * FROM equipment_maintenance
WHERE equipment_id = $1 AND deleted_at IS NULL
ORDER BY performed_date DESC NULLS LAST, scheduled_date DESC;

-- name: ListUpcomingMaintenance :many
SELECT * FROM equipment_maintenance
WHERE hospital_id = $1 AND is_completed = FALSE AND deleted_at IS NULL
ORDER BY next_due_date ASC NULLS LAST, scheduled_date ASC;

-- name: ListOverdueMaintenance :many
SELECT * FROM equipment_maintenance
WHERE hospital_id = $1
  AND is_completed = FALSE
  AND next_due_date < CURRENT_DATE
  AND deleted_at IS NULL
ORDER BY next_due_date ASC;

-- name: CompleteMaintenance :one
UPDATE equipment_maintenance
SET is_completed = TRUE, performed_date = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteEquipmentMaintenance :exec
UPDATE equipment_maintenance SET deleted_at = now() WHERE id = $1;
