-- name: CreateEPORecord :one
INSERT INTO epo_records (
    hospital_id, patient_id, session_id, administered_by, product_name,
    dose_units, route, injection_site, hb_at_time, hb_target_min, hb_target_max,
    ferritin_at_time, tsat_at_time, dose_adjustment_reason,
    next_dose_recommendation, adverse_reaction, adverse_reaction_details, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
RETURNING *;

-- name: GetEPORecord :one
SELECT * FROM epo_records WHERE id = $1 AND deleted_at IS NULL;

-- name: ListEPORecordsByPatient :many
SELECT * FROM epo_records
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY administered_at DESC
LIMIT $2 OFFSET $3;

-- name: GetLatestEPOForPatient :one
SELECT * FROM epo_records
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY administered_at DESC
LIMIT 1;

-- name: DeleteEPORecord :exec
UPDATE epo_records SET deleted_at = now() WHERE id = $1;
