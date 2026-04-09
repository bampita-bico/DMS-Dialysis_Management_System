-- name: CreateAnticoagulationRecord :one
INSERT INTO anticoagulation_records (
    hospital_id, session_id, patient_id, administered_by, anticoagulant,
    route, loading_dose_units, loading_dose_time, maintenance_dose_units,
    maintenance_rate, total_dose_units, reversal_agent_given, reversal_agent,
    reversal_dose_units, reversal_time, bleeding_complications, clotting_observed,
    clotting_location, aptt_pre, aptt_post, act_pre, act_post, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
RETURNING *;

-- name: GetAnticoagulationRecord :one
SELECT * FROM anticoagulation_records WHERE id = $1 AND deleted_at IS NULL;

-- name: GetAnticoagulationBySession :one
SELECT * FROM anticoagulation_records WHERE session_id = $1 AND deleted_at IS NULL;

-- name: ListAnticoagulationRecordsByPatient :many
SELECT * FROM anticoagulation_records
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListRecordsWithComplications :many
SELECT * FROM anticoagulation_records
WHERE hospital_id = $1
  AND (bleeding_complications = TRUE OR clotting_observed = TRUE)
  AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateAnticoagulationRecord :one
UPDATE anticoagulation_records
SET loading_dose_units = $2, loading_dose_time = $3, maintenance_dose_units = $4,
    maintenance_rate = $5, total_dose_units = $6, reversal_agent_given = $7,
    reversal_agent = $8, reversal_dose_units = $9, reversal_time = $10,
    bleeding_complications = $11, clotting_observed = $12, clotting_location = $13,
    aptt_pre = $14, aptt_post = $15, act_pre = $16, act_post = $17, notes = $18
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteAnticoagulationRecord :exec
UPDATE anticoagulation_records SET deleted_at = now() WHERE id = $1;
