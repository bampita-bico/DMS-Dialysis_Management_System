-- name: CreateInfectionControlLog :one
INSERT INTO infection_control_logs (
    hospital_id, machine_id, session_id, activity_type, performed_by,
    performed_at, disinfectant_used, disinfectant_concentration,
    contact_time_mins, rinse_cycles_count, bacterial_test_done,
    bacterial_result, cfu_count, was_hbv_patient, bleach_disinfection_done,
    external_surfaces_cleaned, chair_cleaned, machine_ready_for_use, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
RETURNING *;

-- name: GetInfectionControlLog :one
SELECT * FROM infection_control_logs WHERE id = $1 AND deleted_at IS NULL;

-- name: ListInfectionLogsByMachine :many
SELECT * FROM infection_control_logs
WHERE machine_id = $1 AND deleted_at IS NULL
ORDER BY performed_at DESC;

-- name: ListUnverifiedMachines :many
SELECT * FROM infection_control_logs
WHERE hospital_id = $1
  AND machine_ready_for_use = FALSE
  AND deleted_at IS NULL
ORDER BY performed_at;

-- name: ListHBVDisinfections :many
SELECT * FROM infection_control_logs
WHERE hospital_id = $1
  AND was_hbv_patient = TRUE
  AND deleted_at IS NULL
ORDER BY performed_at DESC;

-- name: GetLatestDisinfectionForMachine :one
SELECT * FROM infection_control_logs
WHERE machine_id = $1
  AND machine_ready_for_use = TRUE
  AND deleted_at IS NULL
ORDER BY performed_at DESC
LIMIT 1;

-- name: VerifyMachineReady :one
UPDATE infection_control_logs
SET machine_ready_for_use = TRUE, verified_by = $2, verified_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteInfectionControlLog :exec
UPDATE infection_control_logs SET deleted_at = now() WHERE id = $1;
