-- name: CreateIronTherapyRecord :one
INSERT INTO iron_therapy_records (
    hospital_id, patient_id, session_id, administered_by, product, dose_mg,
    route, infusion_duration_mins, dilution, ferritin_at_time, ferritin_target_min,
    ferritin_target_max, tsat_at_time, hb_at_time, adverse_reaction,
    adverse_reaction_type, adverse_reaction_severity, treatment_given,
    test_dose_given, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
RETURNING *;

-- name: GetIronTherapyRecord :one
SELECT * FROM iron_therapy_records WHERE id = $1 AND deleted_at IS NULL;

-- name: ListIronTherapyRecordsByPatient :many
SELECT * FROM iron_therapy_records
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY administered_at DESC
LIMIT $2 OFFSET $3;

-- name: ListIronAdverseReactions :many
SELECT * FROM iron_therapy_records
WHERE hospital_id = $1 AND adverse_reaction = TRUE AND deleted_at IS NULL
ORDER BY administered_at DESC;

-- name: DeleteIronTherapyRecord :exec
UPDATE iron_therapy_records SET deleted_at = now() WHERE id = $1;
