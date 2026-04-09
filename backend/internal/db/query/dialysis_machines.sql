-- name: CreateDialysisMachine :one
INSERT INTO dialysis_machines (
    hospital_id, machine_code, serial_number, model, manufacturer,
    manufacture_year, installation_date, location, status, is_hbv_dedicated,
    last_service_date, next_service_date, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetDialysisMachine :one
SELECT * FROM dialysis_machines WHERE id = $1 AND deleted_at IS NULL;

-- name: ListDialysisMachinesByHospital :many
SELECT * FROM dialysis_machines
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY machine_code;

-- name: ListAvailableMachines :many
SELECT * FROM dialysis_machines
WHERE hospital_id = $1 AND status = 'available' AND deleted_at IS NULL
ORDER BY machine_code;

-- name: ListHBVDedicatedMachines :many
SELECT * FROM dialysis_machines
WHERE hospital_id = $1 AND is_hbv_dedicated = TRUE AND deleted_at IS NULL
ORDER BY machine_code;

-- name: UpdateDialysisMachine :one
UPDATE dialysis_machines
SET machine_code = $2, serial_number = $3, model = $4, manufacturer = $5,
    manufacture_year = $6, installation_date = $7, location = $8, status = $9,
    is_hbv_dedicated = $10, last_service_date = $11, next_service_date = $12, notes = $13
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateMachineStatus :one
UPDATE dialysis_machines
SET status = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteDialysisMachine :exec
UPDATE dialysis_machines SET deleted_at = now() WHERE id = $1;
