-- name: CreateWaterTreatmentLog :one
INSERT INTO water_treatment_logs (
    hospital_id, test_date, test_time, tested_by, sample_location,
    bacterial_count_cfu_ml, endotoxin_level_eu_ml, chlorine_level_ppm,
    chloramine_level_ppm, ph_level, conductivity_us_cm, hardness_mg_l,
    bacteria_result, endotoxin_result, chlorine_result, overall_result,
    out_of_spec_parameters, corrective_action_taken, corrective_action_by,
    retest_required, retest_date, systems_shut_down, shutdown_time,
    resumed_time, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)
RETURNING *;

-- name: GetWaterTreatmentLog :one
SELECT * FROM water_treatment_logs WHERE id = $1 AND deleted_at IS NULL;

-- name: ListWaterTestsByDate :many
SELECT * FROM water_treatment_logs
WHERE hospital_id = $1 AND test_date = $2 AND deleted_at IS NULL
ORDER BY test_time;

-- name: ListFailedWaterTests :many
SELECT * FROM water_treatment_logs
WHERE hospital_id = $1
  AND overall_result = 'fail'
  AND deleted_at IS NULL
ORDER BY test_date DESC, test_time DESC;

-- name: ListPendingWaterTests :many
SELECT * FROM water_treatment_logs
WHERE hospital_id = $1
  AND overall_result = 'pending'
  AND deleted_at IS NULL
ORDER BY test_date DESC, test_time DESC;

-- name: ApproveWaterTest :one
UPDATE water_treatment_logs
SET approved_by = $2, approved_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateWaterTestResults :one
UPDATE water_treatment_logs
SET bacteria_result = $2, endotoxin_result = $3, chlorine_result = $4,
    overall_result = $5, out_of_spec_parameters = $6,
    corrective_action_taken = $7, corrective_action_by = $8,
    retest_required = $9, retest_date = $10
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteWaterTreatmentLog :exec
UPDATE water_treatment_logs SET deleted_at = now() WHERE id = $1;
