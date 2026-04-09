-- name: CreateEquipmentFault :one
INSERT INTO equipment_faults (
    hospital_id, equipment_id, reported_by, fault_description, severity,
    is_equipment_unusable, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetEquipmentFault :one
SELECT * FROM equipment_faults WHERE id = $1 AND deleted_at IS NULL;

-- name: ListFaultsByEquipment :many
SELECT * FROM equipment_faults
WHERE equipment_id = $1 AND deleted_at IS NULL
ORDER BY reported_at DESC;

-- name: ListUnresolvedFaults :many
SELECT * FROM equipment_faults
WHERE hospital_id = $1 AND is_resolved = FALSE AND deleted_at IS NULL
ORDER BY severity DESC, reported_at DESC;

-- name: ListCriticalFaults :many
SELECT * FROM equipment_faults
WHERE hospital_id = $1 AND severity IN ('critical', 'severe') AND is_resolved = FALSE AND deleted_at IS NULL
ORDER BY severity DESC, reported_at DESC;

-- name: ResolveFault :one
UPDATE equipment_faults
SET is_resolved = TRUE, resolved_by = $2, resolved_at = now(),
    resolution_description = $3, downtime_minutes = $4, cost = $5
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteEquipmentFault :exec
UPDATE equipment_faults SET deleted_at = now() WHERE id = $1;
